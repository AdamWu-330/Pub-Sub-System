// reference: https://www.rabbitmq.com/tutorials/tutorial-two-go.html, https://www.rabbitmq.com/tutorials/tutorial-five-go.html,
// https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/AdamWu-330/Pub-Sub-System/fetch_source"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// connecting to mongodb
	client_options := options.Client().ApplyURI("mongodb://localhost:27017")

	// connecting to mongodb
	client, err := mongo.Connect(context.TODO(), client_options)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("successfully connected to MongoDB!")

	// define collection
	collection := client.Database("cvst_pubsub").Collection("bikeStation")

	// connecting to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer ch.Close()

	// declare work queue
	q, err := ch.QueueDeclare(
		"bike_stations", // name
		true,            // durable
		false,           // autoDelete
		false,           // exclusive
		false,           // noWait
		nil,             // args
	)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// consume from queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // autoAck
		false,  // exclusive
		false,  // noLocal
		false,  // noWait
		nil,    // args
	)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	go func() {
		for msg := range msgs {
			var obj fetch_source.Detail_station
			json.Unmarshal(msg.Body, &obj)
			fmt.Println(obj)

			_, err := collection.InsertOne(context.TODO(), obj)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			msg.Ack(false)
		}
	}()

	log.Printf("Waiting for new messages...")
	var listen chan struct{}
	<-listen

	// closing connection
	err = client.Disconnect(context.TODO())

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("disconnected from MongoDB")
}
