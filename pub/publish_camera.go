// reference: https://www.rabbitmq.com/tutorials/tutorial-two-go.html, https://www.rabbitmq.com/tutorials/tutorial-five-go.html
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AdamWu-330/Pub-Sub-System/fetch_source"
	amqp "github.com/rabbitmq/amqp091-go"
)

func encode_to_bytes(obj fetch_source.Camera) (error, []byte) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	err := encoder.Encode(obj)
	return err, buf.Bytes()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	camera_objs := fetch_source.Fetch_source_camera()

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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// publish to work queues for saving to db
	for i := 0; i < len(camera_objs); i++ {
		err, content := encode_to_bytes(camera_objs[i])

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
	}
	log.Printf("successfully published %d camera entries to work queue", len(camera_objs))

	// pub-sub, use another channel
	ch2, err := conn.Channel()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer ch2.Close()

	err = ch2.ExchangeDeclare(
		"camera_exchange", // name
		"topic",           // kind
		true,              // durable
		false,             // autoDelete
		false,             // internal
		false,             // noWait
		nil,               // args
	)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// publish to exchange
	for i := 0; i < len(camera_objs); i++ {
		err, content := encode_to_bytes(camera_objs[i])

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		err = ch.PublishWithContext(ctx,
			"camera_exchange", // exchange
			"camera_pubsub",   // routing key
			false,             // mandatory
			false,             // immediate
			amqp.Publishing{ // messages to publish
				ContentType: "text/plain",
				Body:        content,
			})

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
	log.Printf("successfully published %d camera entries to exchange", len(camera_objs))
}
