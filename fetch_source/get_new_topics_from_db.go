package fetch_source

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type New_topic struct {
	Topic_name string `bson:"topic"`
	Url        string `bson:"api_url"`
}

func Get_new_topics_from_db() []New_topic {
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

	//fmt.Println("successfully connected to MongoDB!")

	// define collection
	collection := client.Database("cvst_pubsub").Collection("topics")

	opt := options.Find()

	var res []New_topic

	cur, err2 := collection.Find(context.TODO(), bson.D{{"unknown_source", "yes"}}, opt)

	if err2 != nil {
		fmt.Println(err2.Error())
		os.Exit(1)
	}

	err3 := cur.All(context.TODO(), &res)

	if err3 != nil {
		fmt.Println(err3.Error())
		os.Exit(1)
	}

	//var decoded_res []New_topic

	for i := 0; i < len(res); i++ {
		cur.Decode(&res[i])
	}
	// closing connection
	err = client.Disconnect(context.TODO())

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	//fmt.Println("disconnected from MongoDB")
	//fmt.Println(topic)
	return res
}
