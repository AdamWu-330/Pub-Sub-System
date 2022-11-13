package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Camera struct {
	ID        string  `json:"ID"`
	Latitude  float32 `json:"Latitude"`
	Longitude float32 `json:"Longitude"`
	URL       string  `json:"Url"`
	CityName  string  `json:"CityName"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func fetchFromDB() []*Camera {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("traffic").Collection("cameras")

	findOptions := options.Find()
	findOptions.SetLimit(1000)

	var results []*Camera

	// Finding multiple documents returns a cursor
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through the cursor
	for cur.Next(context.TODO()) {
		var elem Camera
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	// Close the connection once no longer needed
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connection to MongoDB closed.")
	}
	return results
}

func pub(start int, end int, results []*Camera) {
	// fmt.Println("test")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"traffic_topic", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// fmt.Println("test")
	// fmt.Println(start, end)
	for i := start; i < end; i++ {
		// fmt.Println(i)
		concatenated := fmt.Sprintf("%s %v %v", results[i].ID, results[i].Latitude, results[i].Longitude)
		// fmt.Println(concatenated)
		body := concatenated
		// body := bodyFromDB(os.Args)
		err = ch.PublishWithContext(ctx,
			"traffic_topic", // exchange
			"cameras",       // routing key
			false,           // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")

		log.Printf(" [x] Sent %s", body)
	}
}

func main() {
	// RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"traffic_topic", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var results []*Camera = fetchFromDB()

	//pub(0, len(results), results, ch, ctx, err)
	//pub(0, len(results)/2, results)

	pub := func(start int, end int) {
		for i := start; i < end; i++ {
			concatenated := fmt.Sprintf("%s %v %v", results[i].ID, results[i].Latitude, results[i].Longitude)
			fmt.Println(concatenated)
			body := concatenated
			// body := bodyFromDB(os.Args)
			err = ch.PublishWithContext(ctx,
				"traffic_topic", // exchange
				"cameras",       // routing key
				false,           // mandatory
				false,           // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})
			failOnError(err, "Failed to publish a message")
			log.Printf(" [x] Sent num %d message -- %s", i, body)
		}
	}

	go pub(0, len(results)/2)
	go pub(len(results)/2, len(results))

	time.Sleep(10 * time.Second)
	// fmt.Println("----------\n")
	// pub(0, 10)
}

func bodyFromDB(args []string) string {
	var s string
	if (len(args) < 3) || os.Args[2] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[2:], " ")
	}
	return s
}

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "anonymous.info"
	} else {
		s = os.Args[1]
	}
	return s
}
