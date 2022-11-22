// fetch user subscription information from db, declare user queues, bind subscriber queues to the topic exchanges
// to run:
// go run update_subscriber_queues.go <topic 1> ... <topic n>

// reference: https://www.rabbitmq.com/tutorials/tutorial-two-go.html, https://www.rabbitmq.com/tutorials/tutorial-five-go.html

package main

import (
	"fmt"
	"os"

	"github.com/AdamWu-330/Pub-Sub-System/manage"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// connecting to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	for _, topic := range os.Args[1:] {
		subscribers := manage.Get_subscribers_of_topic(topic)

		ch, err := conn.Channel()

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		defer ch.Close()

		exchange_name := fmt.Sprintf("%s_exchange", topic)

		// decalre exchange for the current topic
		err = ch.ExchangeDeclare(
			exchange_name, // name
			"topic",       // kind
			true,          // durable
			false,         // autoDelete
			false,         // internal
			false,         // noWait
			nil,           // args
		)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// declare subscriber queues, bind the topic exchange
		for i := 0; i < len(subscribers); i++ {
			// one queue for one subscriber
			q, err := ch.QueueDeclare(
				subscribers[i], // name
				true,           // durable
				false,          // autoDelete
				false,          // exclusive
				false,          // noWait
				nil,            // args
			)

			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			// bind the queue to the topic
			routing_key := fmt.Sprintf("%s_pubsub", topic)

			err = ch.QueueBind(
				q.Name,        // name
				routing_key,   // routing key
				exchange_name, // exchange
				false,         // noWait
				nil,           // args
			)

			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
	}
}
