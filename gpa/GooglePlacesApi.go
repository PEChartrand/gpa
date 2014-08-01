package gpa

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// TextSearchRequest stores all the parameters of query to be sent to the Goolge Places API
// The "ApiKey" and "Query" parameters are required by the api. The "Query" parameter will be built with the "LocationName" and "LocationAreaName" properties.
// <dl>
// 	<dt>ApiKey</dt><dd>Your api key</dd>
// 	<dt>LocationName</dt><dd>The name of the location (ex.: "Restaurant Le Food")</dd>
//	<dt>LocationAreaName (can be empty)</dt><dd>The name of the area you are filtering by. (ex.: a city name to avoid getting all restaurant Le Food in the world)
// Other required choices are:
// <dl>
// 	<dt>dataType string</dl><dt>"xml" for get an XML response<br/>"json" for a json response</dd>
// </dl>
// Any optional parameter should be stored in the OptinoalParameters map.
// Here's a list of all the optional parameters:
// <dl>
//  <dt>location</dt><dd>The latitude/longitude around which to retrieve place information. This must be specified as latitude,longitude. If you specify a location parameter, you must also specify a radius parameter.</dd>
//  <dt>radius</dt><dd>Defines the distance (in meters) within which to bias place results. The maximum allowed radius is 50 000 meters. Results inside of this region will be ranked higher than results outside of the search circle; however, prominent results from outside of the search radius may be included.</dd>
//  <dt>language</dt><dd>The language code, indicating in which language the results should be returned, if possible. See the list of supported languages and their codes. Note that we often update supported languages so this list may not be exhaustive.</dd>
//  <dt>minprice and maxprice (optional)</dt><dd>Restricts results to only those places within the specified price level. Valid values are in the range from 0 (most affordable) to 4 (most expensive), inclusive. The exact amount indicated by a specific value will vary from region to region.</dd>
//  <dt>opennow</dt><dd>Returns only those places that are open for business at the time the query is sent. places that do not specify opening hours in the Google Places database will not be returned if you include this parameter in your query.</dd>
//  <dt>types</dt><dd>Restricts the results to places matching at least one of the specified types. Types should be separated with a pipe symbol (type1|type2|etc). See the list of supported types.</dd>
//  <dt>zagatselected</dt><dd>Add this parameter (just the parameter name, with no associated value) to restrict your search to locations that are Zagat selected businesses. This parameter must not include a true or false value. The zagatselected parameter is experimental, and is only available to Places API enterprise customers.</dd>
// </dl>
type TextSearchRequest struct {
	ApiKey             string // your api key
	LocationName       string // the title of the place you are looking for
	LocationAreaName   string // An area name (ex: a city name)
	ResponseType       string // json or xml
	Url                string // url sent to api
	ResultLimit        int    // a many places to be included in the result. (1 will give you only the most relevant result [the first one found by google] according to your rankby parameter)
	OptionalParameters map[string]string
}

// Type representing the response receive from Google Places (json or xml)
type QueryResult struct {

	// data about the response
	HtmlAttributions []string `json:"html_attributions" xml:"html_attribution"`
	Status           string   `json:"status" xml:"status"`
	Results          []Result `json:"results" xml:"result"`
}

// One place 
type Result struct {
	PlaceId          string  `json:"place_id" xml:"place_id"`
	FormattedAddress string  `json:"formatted_address" xml:"formatted_address"`
	Icon             string  `json:"icon" xml:"icon"`
	Name             string  `json:"name" xml:"name"`
	Rating           float32 `json:"rating" xml:"rating"`
	Reference        string  `json:"reference" xml:"reference"`
	Geometry         struct {
		Location struct {
			Lat float32 `json:"lat" xml:"lat"`
			Lng float32 `json:"lng" xml:"lng"`
		} `json:"location" xml:"location"`
	} `json:"geometry" xml:"geometry"`
	OpeningHours struct {
		OpenNow bool `json:"open_now" xml:"open_now"`
	} `json:"opening_hours" xml:"opening_hours"`
	Photo struct {
		Height           int    `json:"height" xml:"height"`
		HtmlAttributions string `json:"html_attributions" xml:"html_attribution"`
		PhotoReference   string `json:"photo_reference" xml:"photo_reference"`
		Width            int    `json:"width" xml:"width"`
	} `json:"photo" xml:"photo"`
}

// buildUrlForQuery builds the full url that will posted to the api
// path string is the comman path for the api. If nil then it defaults to: "https://maps.googleapis.com/maps/api/place/textsearch/"
func (tsr *TextSearchRequest) BuildUrlForQuery(path string) {

	if path == "" {
		path = "https://maps.googleapis.com/maps/api/place/textsearch/" // Defualt path
	}

	u, err := url.Parse(path + tsr.ResponseType) // append response type parameter
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}

	v := url.Values{}
	queryPar := tsr.LocationName + " " + tsr.LocationAreaName
	v.Set("query", url.QueryEscape(queryPar))
	v.Set("key", url.QueryEscape(tsr.ApiKey))
	u.RawQuery = v.Encode()
	// fmt.Println(u.String())
	tsr.Url = u.String()
}

// Query queries the Google Place API and returns the data. If the response is error response then an error is returned
func (tsr *TextSearchRequest) Query(qrfj *QueryResult) (string, error) {

	// Get data
	res, err := http.Get(tsr.Url) // post

	if err != nil {
		return "ERROR: ", err
	}

	// Check resposne code
	if res.StatusCode != 200 {
		return "ERROR: ", errors.New(res.Status)
	}

	//DEBUG
	// // Output response headers
	// for key, ele := range res.Header {
	// 	fmt.Print(key + ": \n")
	// 	for key1, ele1 := range ele {
	// 		fmt.Print("  " + string(key1) + ":" + ele1 + "\n")
	// 	}
	// }
	//DEBUG^^^

	// Store data in it's type
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(tsr.ResponseType)
	if tsr.ResponseType == "json" {
		json.Unmarshal(body, qrfj)
	} else if tsr.ResponseType == "xml" {
		xml.Unmarshal(body, qrfj)
	}
	res.Body.Close()
	if err != nil {
		return string(body), err
	}

	// store placaes in memory	

	// return json
	return string(body), nil

}
