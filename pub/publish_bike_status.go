// reference: https://www.rabbitmq.com/tutorials/tutorial-two-go.html, https://www.rabbitmq.com/tutorials/tutorial-five-go.html
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AdamWu-330/Pub-Sub-System/fetch_source"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	status_objs := fetch_source.Fetch_source_bike_status()

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

	q, err := ch.QueueDeclare(
		"bike_status", // name
		true,          // durable
		false,         // autoDelete
		false,         // exclusive
		false,         // noWait
		nil,           // args
	)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// publish to work queues for saving to db
	for i := 0; i < len(status_objs); i++ {
		err, content := fetch_source.Encode_to_bytes(status_objs[i])

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		err = ch.PublishWithContext(ctx,
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{ // messages to publish
				ContentType: "text/plain",
				Body:        content,
			})

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		log.Printf("Successfully published: %s", content)
	}

	// pub-sub, use another channel
	ch2, err := conn.Channel()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer ch2.Close()

	// err = ch2.ExchangeDeclare(
	// 	"bike_status_exchange", // name
	// 	"topic",                // kind
	// 	true,                   // durable
	// 	false,                  // autoDelete
	// 	false,                  // internal
	// 	false,                  // noWait
	// 	nil,                    // args
	// )

	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }

	// publish to exchange
	for i := 0; i < len(status_objs); i++ {
		var obj fetch_source.Client_message
		obj.Data = status_objs[i]
		obj.Type = "bikeStatus"

		err, content := fetch_source.Encode_to_bytes(obj)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		err = ch.PublishWithContext(ctx,
			"cvst_exchange",        // exchange
			"all.bike.bike_status", // routing key
			false,                  // mandatory
			false,                  // immediate
			amqp.Publishing{ // messages to publish
				ContentType: "text/plain",
				Body:        content,
			})

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		log.Printf("successfully sent %s", content)
	}
}
