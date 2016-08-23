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
	"strings"
	"strconv"
	"os"
	"log"
	"io/ioutil"
	"net/http"
	"github.com/gorilla/mux"
)

// CPU + GPU are inside the same SOC: BCM 2837 64bit ARMv8 Cortex A53 Quad Core
const cpuTempSysFileName string = "/sys/class/thermal/thermal_zone0/temp"

type CpuTemp struct {
	Temp string `json:"temp"`
	HostName string `json:"hostName"`
}

// Identifier have to be 'exported' (the first character of the
// identifier's name is a Unicode upper case letter) to permit access
// to it from another package)
//
// https://golang.org/ref/spec#Exported_identifiers
//

func Index(w http.ResponseWriter, r *http.Request) {

	dat, err := ioutil.ReadFile(cpuTempSysFileName)
	check(err)
	cpuTempRaw := strings.Trim(string(dat), "\n")

	temp, err := strconv.ParseInt(cpuTempRaw, 10, 32)
	// temp is a 5 digit integer
	// example: 48312
	
        temp1 := int(temp / 1000) 
        temp2 := int(temp / 100)
	tempM := int(temp2 % temp1)
	hostName, err := os.Hostname()

	cpuTemp := CpuTemp{strconv.Itoa(temp1) + "." + strconv.Itoa(tempM), hostName}

	log.Printf("cpuTemp = \"%s\", hostName = \"%s\"\n", 
		cpuTemp.Temp, 
		cpuTemp.HostName)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(cpuTemp)
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
    
