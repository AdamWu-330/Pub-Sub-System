// reference: https://www.rabbitmq.com/tutorials/tutorial-two-go.html, https://www.rabbitmq.com/tutorials/tutorial-five-go.html
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/AdamWu-330/Pub-Sub-System/fetch_source"
	amqp "github.com/rabbitmq/amqp091-go"
)

func encode_to_bytes(obj fetch_source.Detail_station) (error, []byte) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	err := encoder.Encode(obj)
	return err, buf.Bytes()
}

func main() {
	// connecting to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	// pub-sub, use another channel
	ch2, err := conn.Channel()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer ch2.Close()

	err = ch2.ExchangeDeclare(
		"cvst_exchange", // name
		"topic",         // kind
		true,            // durable
		false,           // autoDelete
		false,           // internal
		false,           // noWait
		nil,             // args
	)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

}
