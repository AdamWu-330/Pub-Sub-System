// reference: https://tutorialedge.net/golang/consuming-restful-api-with-go/
// https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Response struct {
	ID                 string  `bson:"ID"`
	Organization       string  `bson:"Organization"`
	RoadwayName        string  `bson:"RoadwayName"`
	DirectionOfTravel  string  `bson:"DirectionOfTravel"`
	Reported           float32 `bson:"Reported"`
	LastUpdated        float32 `bson:"LastUpdated"`
	StartDate          float32 `bson:"StartDate"`
	PlannedEndDate     float32 `bson:"PlannedEndDate"`
	LanesAffected      string  `bson:"LanesAffected"`
	Latitude           float32 `bson:"Latitude"`
	Longitude          float32 `bson:"Longitude"`
	LatitudeSecondary  float32 `bson:"LatitudeSecondary"`
	LongitudeSecondary float32 `bson:"LongitudeSecondary"`
	EventType          string  `bson:"EventType"`
	IsFullClosure      bool    `bson:"IsFullClosure"`
	LinkId             string  `bson:"LinkId"`
	Comment            string  `bson:"Comment"`
}

func main() {
	response, err := http.Get("https://511on.ca/api/v2/get/event")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseObjects := make([]Response, 0, 0)
	json.Unmarshal(responseData, &responseObjects)

	// save to mongodb
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
	collection := client.Database("cvst_pubsub").Collection("event")

	for i := 0; i < len(responseObjects); i++ {
		// Insert a single document
		_, err := collection.InsertOne(context.TODO(), responseObjects[i])
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("finished inserting to event collection")

	// Close the connection once no longer needed
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connection to MongoDB closed.")
	}

}
