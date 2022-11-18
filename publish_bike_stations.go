package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/AdamWu-330/Pub-Sub-System/fetch_source"
	amqp "github.com/rabbitmq/amqp091-go"
)

func encode_to_bytes(obj fetch_source.Detail_station) (error, []byte) {
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
	station_objs := fetch_source.Fetch_source_bike_station()

	// RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"bike_stations", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// publish to work queues for saving to db
	for i := 0; i < len(station_objs); i++ {
		err, content := encode_to_bytes(station_objs[i])
		failOnError(err, "Failed to convert to bytes")

		err = ch.PublishWithContext(ctx,
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        content,
			})
		failOnError(err, "Failed to publish a message")

		log.Printf(" [x] Sent %s", content)
	}

	// pub-sub, use another channel
	ch2, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch2.ExchangeDeclare(
		"bike_stations_exchange", // name
		"topic",                  // type
		true,                     // durable
		false,                    // auto-deleted
		false,                    // internal
		false,                    // no-wait
		nil,                      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// publish to exchange
	for i := 0; i < len(station_objs); i++ {
		err, content := encode_to_bytes(station_objs[i])
		failOnError(err, "Failed to convert to bytes")

		err = ch.PublishWithContext(ctx,
			"bike_stations_exchange", // exchange
			"bike_stations_pubsub",   // routing key
			false,                    // mandatory
			false,                    // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        content,
			})
		failOnError(err, "Failed to publish a message")

		log.Printf(" [x] Sent %s", content)
	}

}
