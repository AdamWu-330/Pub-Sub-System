// reference: https://tutorialedge.net/golang/consuming-restful-api-with-go/

package fetch_source

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func Fetch_source_york_open_data() string {
	response, err := http.Get("https://ww6.yorkmaps.ca/traveltime/iteris_traveltimes_out.xml")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	raw_response, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	return string(raw_response)
}
