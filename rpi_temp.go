//
// Name:
//
// rpi_temp
//
// Description:
//
// Get Raspberry PI CPU temperature (OS: Raspian).
//
// Example:
//
// $ curl http://carmel:9090/
// {"temp":50.4,"hostName":"carmel"}
//
// REST API notes:
//
// http://thenewstack.io/make-a-restful-json-api-go/
// https://github.com/ant0ine/go-json-rest/graphs/contributors
// https://golang.org/doc/articles/wiki/ (Writing Web Applications)
// http://jsonapi.org/
//
// jens@mesgtone.net
//

package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jens18/rpi_temp/cputemp"
	"log"
	"net/http"
)

// Identifier have to be 'exported' (the first character of the
// identifier's name is a Unicode upper case letter) to permit access
// to it from another package)
//
// https://golang.org/ref/spec#Exported_identifiers
//

func Index(w http.ResponseWriter, r *http.Request) {

	var c cputemp.CpuTemp

	c = c.Get()

	log.Printf("cpuTemp = \"%s\", hostName = \"%s\"\n",
		c.Temp,
		c.HostName)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(c)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":9090", router))
}
