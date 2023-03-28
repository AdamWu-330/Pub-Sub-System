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

	for i:=0; i<250; i++ {
		topic_name := fmt.Sprintf("test_topic%d", i) 
		title_name := fmt.Sprintf("test topic %d", i) 
		// api_urls := [16]string{"https://511on.ca/api/v2/get/constructionprojects", "https://511on.ca/api/v2/get/cameras",
		// 	"https://511on.ca/api/v2/get/groupedcameras", "https://511on.ca/api/v2/get/roadconditions", "https://511on.ca/api/v2/get/transithub", 
		// 	"https://511on.ca/api/v2/get/carpoollots", "https://511on.ca/api/v2/get/ferryterminals", "https://511on.ca/api/v2/get/servicecentres",
		// 	"https://511on.ca/api/v2/get/informationcenter", "https://511on.ca/api/v2/get/hovlanes", "https://511on.ca/api/v2/get/truckrestareas",
		// 	"https://511on.ca/api/v2/get/inspectionstations", "https://511on.ca/api/v2/get/roundabouts", "https://511on.ca/api/v2/get/seasonalloadapi",
		// 	"https://511on.ca/api/v2/get/allrestareas", "https://511on.ca/api/v2/get/alerts"}
		api_urls := []string{
			"https://myttc.ca/finch_station.json", "https://myttc.ca/spadina_station.json", 
			"https://api.ontario.ca/api/data/81527?count=0&download=1", "https://api.ontario.ca/api/data/81522?count=0&download=1"}
			
		_, err := collection.InsertOne(context.TODO(),bson.D{
			{"topic", topic_name},
			{"title", title_name},
			{"api_url", api_urls[i / 63]},
			{"unknown_source", "yes"},
		})

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	fmt.Println("successfully inserted 250 test channels to topics collection")
	// closing connection
	err = client.Disconnect(context.TODO())

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("disconnected from MongoDB")
}
