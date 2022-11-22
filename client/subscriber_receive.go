// client program, to receive subscribed content
// to run: go run subscriber_receive.go <your name>

// reference: https://www.rabbitmq.com/tutorials/tutorial-two-go.html, https://www.rabbitmq.com/tutorials/tutorial-five-go.html

package main

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// validate input
	if len(os.Args) <= 1 {
		fmt.Println("we did not receive your name")
		os.Exit(1)
	}
	if len(os.Args) > 2 {
		fmt.Println("only one user name is allowed")
		os.Exit(2)
	}

	// connecting to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	defer ch.Close()

	// declare user queue for the current user
	q, err := ch.QueueDeclare(
		os.Args[1], // name
		true,       // durable
		false,      // autoDelete
		false,      // exclusive
		false,      // noWait
		nil,        // args
	)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	// consumer messages from queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // autoAck
		false,  // exclusive
		false,  // noLocal
		false,  // noWait
		nil,    // args
	)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	go func() {
		for msg := range msgs {
			log.Printf("Successfully received: %s", msg.Body[:1000])
		}
	}()

	log.Printf("Waiting for new messages...")

	var listen chan struct{}
	<-listen
}
