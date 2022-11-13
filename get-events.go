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
	ID        string  `json:"ID"`
	Latitude  float32 `json:"Latitude"`
	Longitude float32 `json:"Longitude"`
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
	// fmt.Println(string(responseData))
	// fmt.Printf("type of response: %T\n", responseData)

	responseObjects := make([]Response, 0, 0)
	json.Unmarshal(responseData, &responseObjects)
	//fmt.Println(len(responseObjects))

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
	collection := client.Database("traffic").Collection("events")

	for i := 0; i < len(responseObjects); i++ {
		// Insert a single document
		insertResult, err := collection.InsertOne(context.TODO(), responseObjects[i])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	}

	// Close the connection once no longer needed
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connection to MongoDB closed.")
	}

}
