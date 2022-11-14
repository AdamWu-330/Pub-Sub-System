// reference: https://tutorialedge.net/golang/consuming-restful-api-with-go/
// https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial
// https://freshman.tech/snippets/go/image-to-base64/
package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Response struct {
	ID                string  `bson:"ID"`
	Organization      string  `bson:"Organization"`
	RoadwayName       string  `bson:"RoadwayName"`
	DirectionOfTravel string  `bson:"DirectionOfTravel"`
	Latitude          float32 `bson:"Latitude"`
	Longitude         float32 `bson:"Longitude"`
	Name              string  `bson:"Name"`
	Url               string  `bson:"Url"`
	Status            string  `bson:"Status"`
	Description       string  `bson:"Description"`
	CityName          string  `bson:"CityName"`
	Image             string  `bson:"Image"`
	LastUpdate        string  `bson:"LastUpdate"`
	LastModified      string  `bson:"LastModified"`
}

func main() {
	response, err := http.Get("https://511on.ca/api/v2/get/cameras")

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

	// base64 encode image from url, add multiple goroutines to speed up
	to_base64 := func(start int, end int) {
		for i := start; i < end; i++ {
			rs, err := http.Get(responseObjects[i].Url)
			if err != nil {
				log.Fatal(err)
			}

			defer rs.Body.Close()

			bytes, err := ioutil.ReadAll(rs.Body)
			if err != nil {
				log.Fatal(err)
			}

			responseObjects[i].Image = base64.StdEncoding.EncodeToString(bytes)
			// fmt.Println(responseObjects[i].Image)
		}
	}

	go to_base64(0, len(responseObjects)/4)
	go to_base64(len(responseObjects)/4, len(responseObjects)/2)
	go to_base64(len(responseObjects)/2, len(responseObjects)/4*3)
	go to_base64(len(responseObjects)/4*3, len(responseObjects))

	time.Sleep(100 * time.Second)

	fmt.Println("finished base64 encoding for images")

	for i := 0; i < len(responseObjects); i++ {
		responseObjects[i].LastUpdate = time.Now().Format(time.RFC3339)
		responseObjects[i].LastModified = time.Now().Format(time.RFC3339)
	}

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
	collection_history := client.Database("cvst_pubsub").Collection("camera")

	client.Database("cvst_pubsub").Collection("cameraCurrent").Drop(context.TODO())

	collection_realtime := client.Database("cvst_pubsub").Collection("cameraCurrent")

	for i := 0; i < len(responseObjects); i++ {
		// Insert a single document
		insertResult, err := collection_history.InsertOne(context.TODO(), responseObjects[i])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted a single document to history collection: ", insertResult.InsertedID)
	}

	for i := 0; i < len(responseObjects); i++ {
		// Insert a single document
		insertResult, err := collection_realtime.InsertOne(context.TODO(), responseObjects[i])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted a single document to realtime collection: ", insertResult.InsertedID)
	}

	// Close the connection once no longer needed
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connection to MongoDB closed.")
	}

}
