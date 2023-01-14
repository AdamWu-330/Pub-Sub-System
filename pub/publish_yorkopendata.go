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
	open_data := fetch_source.Fetch_source_york_open_data()

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
		"york_opendata", // name
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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// publish to work queues for saving to db
	open_data_byte := []byte(open_data)
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{ // messages to publish
			ContentType: "text/plain",
			Body:        open_data_byte,
		})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	log.Printf("Successfully published to work queue: %s", open_data[:10])

	// publish to exchange
	err = ch.PublishWithContext(ctx,
		"cvst_exchange",    // exchange
		"all.yorkopendata", // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{ // messages to publish
			ContentType: "text/plain",
			Body:        open_data_byte,
		})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	log.Printf("Successfully published to exchange: %s", open_data[:10])
}
