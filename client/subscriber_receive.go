// client program, to receive subscribed content
// to run: go run subscriber_receive.go <your name>

// reference: https://www.rabbitmq.com/tutorials/tutorial-two-go.html, https://www.rabbitmq.com/tutorials/tutorial-five-go.html

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// validate input
	if len(os.Args) <= 1 {
		fmt.Println("we did not receive your name")
		os.Exit(1)
	}
	if len(os.Args) > 3 {
		fmt.Println("only one user name and one saving directory is allowed")
		os.Exit(2)
	}

	user_name := os.Args[1]

	saving := false
	var saving_dir string

	if len(os.Args) > 2 {
		saving = true
		saving_dir = os.Args[2]

		_, err := os.Open(saving_dir)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
	}

	// connecting to RabbitMQ, use 5672 if publishing server is on local, use 5673 is publishing server is remote
	//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5673/")

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
		user_name, // name
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,       // args
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
			log.Printf("Successfully received from: %s", msg.RoutingKey)

			// save to local by topic
			if saving {
				saving_file := filepath.Join(saving_dir, msg.RoutingKey+".txt")
				f, err := os.OpenFile(saving_file, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)

				if err != nil {
					fmt.Print(err.Error())
					os.Exit(1)
				}

				_, err_w := f.Write(msg.Body)

				if err_w != nil {
					fmt.Print(err_w.Error())
					os.Exit(1)
				}

				f.Close()
			}

		}
	}()

	log.Printf("Waiting for new messages...")

	var listen chan struct{}
	<-listen
}
