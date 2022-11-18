// reference: https://tutorialedge.net/golang/consuming-restful-api-with-go/
// https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial
package fetch_source

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func Fetch_source_bike_station() []Detail {
	response, err := http.Get("https://tor.publicbikesystem.net/ube/gbfs/v1/en/station_information")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	raw_response, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response_obj Response
	json.Unmarshal(raw_response, &response_obj)

	var station_objs = make([]Detail, len(response_obj.Data.Stations))
	station_objs = response_obj.Data.Stations

	return station_objs
}
