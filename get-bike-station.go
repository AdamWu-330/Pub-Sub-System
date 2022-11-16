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
	Data Station `bson:"data"`
}

type Station struct {
	Stations []Detail `bson:"stations"`
}

type Detail struct {
	Station_id             string   `bson:"station_id"`
	Name                   string   `bson:"name"`
	Physical_configuration string   `bson:"physical_configuration"`
	Lat                    float32  `bson:"lat"`
	Lon                    float32  `bson:"lon"`
	Altitude               float32  `bson:"altitude"`
	Address                string   `bson:"address"`
	Capacity               int      `bson:"capacity"`
	Is_charging_station    bool     `bson:"is_charging_station"`
	Rental_methods         []string `bson:"rental_methods"`
	Groups                 []string `bson:"groups"`
	Obcn                   string   `bson:"obcn"`
	Nearby_distance        float32  `bson:"nearby_distance"`
	Ride_code_support      bool     `bson:"_ride_code_support"`
}

func main() {
	response, err := http.Get("https://tor.publicbikesystem.net/ube/gbfs/v1/en/station_information")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	var stationObjects = make([]Detail, len(responseObject.Data.Stations))
	stationObjects = responseObject.Data.Stations

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
	collection := client.Database("cvst_pubsub").Collection("bikeStation")

	for i := 0; i < len(stationObjects); i++ {
		// Insert a single document
		_, err := collection.InsertOne(context.TODO(), stationObjects[i])
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("finished inserting to bikeStatus collection")

	// Close the connection once no longer needed
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connection to MongoDB closed.")
	}

}
