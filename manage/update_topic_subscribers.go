// https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial
package manage

import (
	"context"
	"fmt"
	"log"
	"os"

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

	// connecting to mongodb
	client_options := options.Client().ApplyURI("mongodb://localhost:27017")

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
	collection := client.Database("cvst_pubsub").Collection("topic_subscibers")

	_, err2 := collection.InsertOne(context.TODO(), t1)
	if err2 != nil {
		fmt.Println(err2.Error())
		os.Exit(1)
	}

	_, err3 := collection.InsertOne(context.TODO(), t2)
	if err3 != nil {
		fmt.Println(err3.Error())
		os.Exit(1)
	}

	fmt.Println("finished updating topic_subscibers")

	// closing connection
	err = client.Disconnect(context.TODO())

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("disconnected from MongoDB")
}

func Get_subscribers_of_topic(topic string) []string {
	// fetch topic-subscribers info from db
	client_options := options.Client().ApplyURI("mongodb://localhost:27017")

	// connecting to mongodb
	client, err := mongo.Connect(context.TODO(), client_options)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("successfully connected to MongoDB!")

	// define collection
	collection := client.Database("cvst_pubsub").Collection("topic_subscibers")

	opt := options.FindOne()

	var res Topic_subscribers

	err2 := collection.FindOne(context.TODO(), bson.D{{"topic", topic}}, opt).Decode(&res)
	if err2 != nil {
		fmt.Println(err2.Error())
		os.Exit(1)
	}

	// closing connection
	err = client.Disconnect(context.TODO())

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("disconnected from MongoDB")

	return res.Subscribers
}
