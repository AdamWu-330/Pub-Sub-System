// reference: https://techoverflow.net/2019/11/16/how-to-download-and-parse-html-page-in-go/
package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Generic_data_xml struct {
	Data map[string]interface{} `xml:"match_summary_data>match_summary"`
}

func main() {
	response, err := http.Get("https://ww6.yorkmaps.ca/traveltime/iteris_traveltimes_out.xml")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	// doc, err := goquery.NewDocumentFromReader(response.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// doc.Find("body").Each(func(i int, s *goquery.Selection) {
	// 	fmt.Printf("body of the page: %s\n", s.Text())
	// })
	raw_response, err := ioutil.ReadAll(response.Body)

	fmt.Println(string(raw_response))

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	var response_obj Generic_data_xml
	xml.Unmarshal(raw_response, &response_obj)

	fmt.Println(response_obj.Data)
	// // for i := 0; i < len(response_obj.Data); i++ {
	// // 	for k, v := range response_obj.Data[i] {
	// // 		fmt.Printf("key: %s , val: %s\n", k, v)
	// // 	}
	// // }

	for k, v := range response_obj.Data {
		fmt.Printf("key: %s , val: %s\n", k, v)
	}
}
