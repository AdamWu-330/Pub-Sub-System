package main

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
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

	// define collection
	collection := client.Database("cvst_pubsub").Collection("topics")
			
	_, err = collection.InsertOne(context.TODO(),bson.D{
		{"topic", "large_message_topic"},
		{"title", "large message topic"},
		{"api_url", "https://data.cityofnewyork.us/resource/kpav-sd4t.json"},
		{"unknown_source", "yes"},
	})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("successfully inserted large message topic to topics collection")
	// closing connection
	err = client.Disconnect(context.TODO())

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("disconnected from MongoDB")
}