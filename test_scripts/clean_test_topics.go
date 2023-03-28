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

	filter := bson.M{"topic": bson.M{"$regex": "^test_topic\\d+"}}
	_, err = collection.DeleteMany(context.Background(), filter)

	fmt.Println("cleaned newly created topics")
	// closing connection
	err = client.Disconnect(context.TODO())

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("disconnected from MongoDB")
}