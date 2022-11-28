// reference: https://www.rabbitmq.com/tutorials/tutorial-two-go.html, https://www.rabbitmq.com/tutorials/tutorial-five-go.html,
// https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/AdamWu-330/Pub-Sub-System/fetch_source"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// query new topics from db
	new_topics := fetch_source.Get_new_topics_from_db()

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

	fmt.Println(new_topics)
	// declare work queue

	//for _, new_topic := range new_topics[:2] {

	save_one_topic := func(topic_name string) {
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

		//fmt.Println(new_topic.Topic_name)
		q, err := ch.QueueDeclare(
			topic_name, // name
			true,       // durable
			false,      // autoDelete
			false,      // exclusive
			false,      // noWait
			nil,        // args
		)
		fmt.Printf("q.name: %s\n", q.Name)
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

		// define collection
		collection := client.Database("cvst_pubsub").Collection(topic_name)
		fmt.Printf("collection name: %s\n", collection.Name())

		for msg := range msgs {
			//fmt.Println(msg.Body)

			// multiple becomes single after publishing
			//var obj fetch_source.Generic_data_single
			var obj fetch_source.Generic_data_single

			json.Unmarshal(msg.Body, &obj)
			json.Unmarshal(msg.Body, &obj.Data)

			var obj2 fetch_source.Generic_data_multiple

			json.Unmarshal(msg.Body, &obj2)
			json.Unmarshal(msg.Body, &obj2.Data)

			if len(obj.Data) > 0 {
				_, err = collection.InsertOne(context.TODO(), obj.Data)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			} else {
				for i := 0; i < len(obj2.Data); i++ {
					_, err = collection.InsertOne(context.TODO(), obj2.Data[i])
					if err != nil {
						fmt.Println(err.Error())
						os.Exit(1)
					}
				}
			}

			// for i := 0; i < len(obj2.Data); i++ {
			// 	_, err = collection.InsertOne(context.TODO(), obj2.Data[i])
			// 	if err != nil {
			// 		fmt.Println(err.Error())
			// 		os.Exit(1)
			// 	}
			// }

			msg.Ack(false)
		}
	}

	for i := 0; i < len(new_topics); i++ {
		go save_one_topic(new_topics[i].Topic_name)
	}
	// go save_one_topic(new_topics[0].Topic_name)
	// go save_one_topic(new_topics[1].Topic_name)

	time.Sleep(20 * time.Second)

	// closing connection
	err = client.Disconnect(context.TODO())

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("disconnected from MongoDB")
}
