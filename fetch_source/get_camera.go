// reference: https://tutorialedge.net/golang/consuming-restful-api-with-go/
package fetch_source

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Camera struct {
	Id                string  `bson:"Id"`
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

func Fetch_source_camera() []Camera {
	response, err := http.Get("https://511on.ca/api/v2/get/cameras")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	raw_response, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	camera_objs := make([]Camera, 0, 0)
	json.Unmarshal(raw_response, &camera_objs)

	// base64 encode image from url, add multiple goroutines to speed up
	to_base64 := func(start int, end int) {
		for i := start; i < end; i++ {
			rs, err := http.Get(camera_objs[i].Url)

			if err != nil {
				fmt.Print(err.Error())
				os.Exit(1)
			}

			defer rs.Body.Close()

			bytes, err := ioutil.ReadAll(rs.Body)

			if err != nil {
				fmt.Print(err.Error())
				os.Exit(1)
			}

			camera_objs[i].Image = base64.StdEncoding.EncodeToString(bytes)
		}
	}

	go to_base64(0, len(camera_objs)/4)
	go to_base64(len(camera_objs)/4, len(camera_objs)/2)
	go to_base64(len(camera_objs)/2, len(camera_objs)/4*3)
	go to_base64(len(camera_objs)/4*3, len(camera_objs))

	time.Sleep(100 * time.Second)

	fmt.Println("finished base64 encoding for images")

	return camera_objs
}
