//TODO make this into a reusable package

package main

import (
	"github.com/PEChartrand/gpa/gpa"
	"net/http"
)

func main() {

	// Listen and Handle request
	http.HandleFunc("/api", gpa.HandleRequest)
	http.ListenAndServe(":8080", nil)

}
