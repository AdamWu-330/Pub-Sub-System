package manage

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Topic_subscribers struct {
	Topic       string   `json:"topic"`
	Subscribers []string `json:"subscribers"`
}

func Update_topic_subscribers() {
	t1 := Topic_subscribers{"bike_stations", []string{"adam", "angela", "alex"}}
	t2 := Topic_subscribers{"bike_status", []string{"alex"}}
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Get a handle for your collection
	collection := client.Database("cvst_pubsub").Collection("topic_subscibers")

	_, err2 := collection.InsertOne(context.TODO(), t1)
	if err2 != nil {
		log.Fatal(err2)
	}

	_, err3 := collection.InsertOne(context.TODO(), t2)
	if err3 != nil {
		log.Fatal(err3)
	}

	fmt.Println("finished updating topic_subscibers")

	// Close the connection once no longer needed
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connection to MongoDB closed.")
	}
}

func Get_subscribers_of_topic(topic string) []string {
	// fetch topic-subscribers info from db
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("Connected to MongoDB!")

	// Get a handle for your collection
	collection := client.Database("cvst_pubsub").Collection("topic_subscibers")

	opt := options.FindOne()

	var res Topic_subscribers

	err2 := collection.FindOne(context.TODO(), bson.D{{"topic", topic}}, opt).Decode(&res)
	if err2 != nil {
		log.Fatal(err2)
	}
	// fmt.Println(res.Topic)
	// fmt.Println(res.Subscribers)
	return res.Subscribers
}
