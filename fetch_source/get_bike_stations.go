// reference: https://tutorialedge.net/golang/consuming-restful-api-with-go/
package fetch_source

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Response1 struct {
	Data Station1 `bson:"data"`
}

type Station1 struct {
	Stations []Detail_station `bson:"stations"`
}

type Detail_station struct {
	Station_id             string   `bson:"station_id"`
	Name                   string   `bson:"name"`
	Physical_configuration string   `bson:"physical_configuration"`
	Lat                    float32  `bson:"lat"`
	Lon                    float32  `bson:"lon"`
	Altitude               int      `bson:"altitude"`
	Address                string   `bson:"address"`
	Capacity               int      `bson:"capacity"`
	Is_charging_station    bool     `bson:"is_charging_station"`
	Post_code              string   `bson:"post_code"`
	Rental_methods         []string `bson:"rental_methods"`
	Groups                 []string `bson:"groups"`
	Obcn                   string   `bson:"obcn"`
	Nearby_distance        float32  `bson:"nearby_distance"`
	Ride_code_support      bool     `bson:"_ride_code_support"`
}

func Fetch_source_bike_station() []Detail_station {
	response, err := http.Get("https://tor.publicbikesystem.net/ube/gbfs/v1/en/station_information")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	raw_response, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	var response_obj Response1
	json.Unmarshal(raw_response, &response_obj)

	var station_objs = make([]Detail_station, len(response_obj.Data.Stations))
	station_objs = response_obj.Data.Stations

	return station_objs
}
