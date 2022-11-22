// reference: https://www.rabbitmq.com/tutorials/tutorial-two-go.html, https://www.rabbitmq.com/tutorials/tutorial-five-go.html,
// https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AdamWu-330/Pub-Sub-System/fetch_source"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
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
	collection_history := client.Database("cvst_pubsub").Collection("camera")

	client.Database("cvst_pubsub").Collection("cameraCurrent").Drop(context.TODO())

	collection_realtime := client.Database("cvst_pubsub").Collection("cameraCurrent")

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
		"camera", // name
		true,     // durable
		false,    // autoDelete
		false,    // exclusive
		false,    // noWait
		nil,      // args
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
			var obj fetch_source.Camera
			json.Unmarshal(msg.Body, &obj)
			//fmt.Println(obj)

			// pop last update and last modified fields

			// store which indices from responseObjects have the same content from the history collection, so no need to insert
			insert := true

			filter := bson.D{{"Id", obj.Id}}
			var result fetch_source.Camera
			err := collection_history.FindOne(context.TODO(), filter).Decode(&result)

			if err != nil && err == mongo.ErrNoDocuments {
				obj.LastUpdate = time.Now().Format(time.RFC3339)
				obj.LastModified = time.Now().Format(time.RFC3339)
			} else {
				obj.LastUpdate = result.LastUpdate
				obj.LastModified = result.LastModified
				insert = false

				if result.Image != obj.Image {
					obj.LastUpdate = time.Now().Format(time.RFC3339)
					insert = true
				} else if result.Name != obj.Name || result.Url != obj.Url ||
					result.Status != obj.Status || result.Description != obj.Description {
					obj.LastModified = time.Now().Format(time.RFC3339)
					insert = true
				}
			}

			if insert {
				_, err := collection_history.InsertOne(context.TODO(), obj)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				fmt.Println("inserted one camera entry to history collection")
			}

			_, err2 := collection_realtime.InsertOne(context.TODO(), obj)
			if err2 != nil {
				fmt.Println(err2.Error())
				os.Exit(1)
			}
			fmt.Println("inserted one camera entry to realtime collection")

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
