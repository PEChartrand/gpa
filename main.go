//TODO make this into a reusable package

package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/PEChartrand/gpa/gpa"
)

func main() {

	qrfj := new(gpa.QueryResult)
	tsr := gpa.TextSearchRequest{
		ApiKey:             "AIzaSyBeAFL1pXObTDrziojUEJyx6x9b2FWUSvA",
		LocationName:       "Auberge du dragon rouge",
		LocationAreaName:   "Montreal",
		ResponseType:       "json",
		ResultLimit:        1,
		OptionalParameters: map[string]string{"language": "fr", "type": "restaurant"},
	}
	tsr.BuildUrlForQuery("")
	tsr.Query(qrfj)

	if tsr.ResponseType == "xml" {
		output, err := xml.MarshalIndent(qrfj, "  ", "  ")
		if err == nil {
			fmt.Printf("%s\n", output)
		} else {
			fmt.Printf("%s\n", err.Error())
		}

	} else {
		output, err := json.MarshalIndent(qrfj, "  ", "  ")

		if err == nil {
			fmt.Printf("%s\n", output)
		} else {
			fmt.Printf("%s\n", err.Error())
		}
	}

	fmt.Printf("\n" + qrfj.Results[0].Name + "\n")

}
