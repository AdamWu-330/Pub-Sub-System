// declare, bing subscriber queues to the topic exchanges
// parameter: topics
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AdamWu-330/Pub-Sub-System/manage"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	for _, topic := range os.Args[1:] {
		subscribers := manage.Get_subscribers_of_topic(topic)
		//fmt.Println(subscribers)

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		exchange_name := fmt.Sprintf("%s_exchange", topic)

		err = ch.ExchangeDeclare(
			exchange_name, // name
			"topic",       // type
			true,          // durable
			false,         // auto-deleted
			false,         // internal
			false,         // no-wait
			nil,           // arguments
		)
		failOnError(err, "Failed to declare an exchange")

		for i := 0; i < len(subscribers); i++ {
			// one queue for one subscriber
			q, err := ch.QueueDeclare(
				subscribers[i], // name
				true,           // durable
				false,          // delete when unused
				false,          // exclusive
				false,          // no-wait
				nil,            // arguments
			)
			failOnError(err, "Failed to declare a queue")

			// bing the queue to the topic
			routing_key := fmt.Sprintf("%s_pubsub", topic)
			err = ch.QueueBind(
				q.Name,        // queue name
				routing_key,   // routing key
				exchange_name, // exchange
				false,
				nil)

			failOnError(err, "Failed to bind a queue")
		}
	}
}
