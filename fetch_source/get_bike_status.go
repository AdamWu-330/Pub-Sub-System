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

type Response2 struct {
	Data Station2 `bson:"data"`
}

type Station2 struct {
	Stations []Detail_status `bson:"stations"`
}

type Detail_status struct {
	Station_id                string `bson:"station_id"`
	Num_bikes_available       int    `bson:"num_bikes_available"`
	Num_bikes_available_types struct {
		Mechanical int `bson:"mechanical"`
		Ebike      int `bson:"ebike"`
	} `bson:"num_bikes_available_types"`
	Num_bikes_disabled  int     `bson:"num_bikes_disabled"`
	Num_docks_available int     `bson:"num_docks_available"`
	Num_docks_disabled  int     `bson:"num_docks_disabled"`
	Last_reported       float32 `bson:"last_reported"`
	Is_charging_station bool    `bson:"is_charging_station"`
	Status              string  `bson:"status"`
	Is_installed        int     `bson:"is_installed"`
	Is_renting          int     `bson:"is_renting"`
	Is_returning        int     `bson:"is_returning"`
	Traffic             string  `bson:"traffic"`
}

func Fetch_source_bike_status() []Detail_status {
	response, err := http.Get("https://tor.publicbikesystem.net/ube/gbfs/v1/en/station_status")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	raw_response, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response_obj Response2
	json.Unmarshal(raw_response, &response_obj)

	var status_objs = make([]Detail_status, len(response_obj.Data.Stations))
	status_objs = response_obj.Data.Stations

	return status_objs
}
