package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

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

	// query new topics from db
	new_topics := fetch_source.Get_new_topics_from_db()

	for _, new_topic := range new_topics {
		response, err := http.Get(new_topic.Url)

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		response_data, err := ioutil.ReadAll(response.Body)

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		q, err := ch.QueueDeclare(
			new_topic.Topic_name, // name
			true,                 // durable
			false,                // autoDelete
			false,                // exclusive
			false,                // noWait
			nil,                  // args
		)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// need two data objects to handle two scenarios: 1. most interested data in one field, 2. interested data as array
		var data_single fetch_source.Generic_data_single

		json.Unmarshal(response_data, &data_single)
		json.Unmarshal(response_data, &data_single.Data)

		var data_multiple fetch_source.Generic_data_multiple

		json.Unmarshal(response_data, &data_multiple)
		json.Unmarshal(response_data, &data_multiple.Data)

		var err_work_q error
		var content_work_q []byte

		if len(data_single.Data) > 0 {
			err_work_q, content_work_q = fetch_source.Encode_to_bytes(data_single.Data)
		} else {
			err_work_q, content_work_q = fetch_source.Encode_to_bytes(data_multiple.Data)
		}

		if err_work_q != nil {
			fmt.Println(err_work_q.Error())
			os.Exit(1)
		}

		// publish to work queues for saving to db
		err = ch.PublishWithContext(ctx,
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{ // messages to publish
				ContentType: "text/plain",
				Body:        content_work_q,
			})

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		//log.Printf("Successfully published: %s", content)
		log.Printf("Successfully published one message to work queue for: %s", q.Name)

		// pub-sub, use another channel
		ch2, err := conn.Channel()

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		defer ch2.Close()

		// publish to exchange
		var err_pubsub error
		var content_pubsub []byte

		if len(data_single.Data) > 0 {
			err_pubsub, content_pubsub = fetch_source.Encode_to_bytes(data_single.Data)

		} else {
			err_pubsub, content_pubsub = fetch_source.Encode_to_bytes(data_multiple.Data)
		}

		if err_pubsub != nil {
			fmt.Println(err_pubsub.Error())
			os.Exit(1)
		}

		routing_key := fmt.Sprintf("all.%s", new_topic.Topic_name)

		err = ch.PublishWithContext(ctx,
			"cvst_exchange", // exchange
			routing_key,     // routing key
			false,           // mandatory
			false,           // immediate
			amqp.Publishing{ // messages to publish
				ContentType: "text/plain",
				Body:        content_pubsub,
			})

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		log.Printf("Successfully published one message to cvst_exchange for: %s", routing_key)
	}

	// closing connection
	err0 := client.Disconnect(context.TODO())

	if err0 != nil {
		fmt.Println(err0.Error())
		os.Exit(1)
	}
	fmt.Println("disconnected from MongoDB")
}
