package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jens18/rpi_temp/cputemp"
	"golang.org/x/build/kubernetes"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const kubeMaster string = "http://carmel:8080"

var netClient = &http.Client{
	Timeout: time.Second * 10,
}

type NodeTemp struct {
	IpAddress string          `json:"ipAddress"`
	NodeTemp  cputemp.CpuTemp `json:"nodeTemp"`
}

func requestCpuTemp(hostIP string) cputemp.CpuTemp {
	var cputemp cputemp.CpuTemp

	// http client with sensible timeout
	// (https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779#.s0bop6t5g)
	response, err := netClient.Get("http://" + hostIP + ":30001")
	if err != nil {
		log.Fatal(err)
	}

	buf, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(buf, &cputemp)

	response.Body.Close()
	return cputemp
}

func Index(w http.ResponseWriter, r *http.Request) {

	var cputemp cputemp.CpuTemp

	var kubetemp []NodeTemp

	// https://godoc.org/golang.org/x/build/kubernetes
	c, err := kubernetes.NewClient(kubeMaster, http.DefaultClient)
	if err != nil {
		log.Fatalf("NewClient: %v", err)
	}

	nodes, _ := c.GetNodes(context.Background())
	if err != nil {
		log.Fatalf("GetNodes: %v", err)
	}

	for _, n := range nodes {
		for _, ip := range n.Status.Addresses {
			if strings.Compare(string(ip.Type), "LegacyHostIP") == 0 {

				cputemp = requestCpuTemp(ip.Address)

				log.Printf("ip=%s, temp=%s, hostname=%s\n",
					ip.Address,
					cputemp.Temp,
					cputemp.HostName)

				kubetemp = append(kubetemp,
					NodeTemp{ip.Address,
						cputemp})
			}
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(kubetemp)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":9999", router))
}
