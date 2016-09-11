package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jens18/rpi_temp/cputemp"
	"golang.org/x/build/kubernetes"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net"
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
	HostName  string          `json:"hostName"` // the real hostname
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

func KubeTemp(w http.ResponseWriter, r *http.Request) {

	var cputemp cputemp.CpuTemp

	var kubetemp []NodeTemp
	var ipToName = make(map[string]string)

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
				// cache the symbolic name
				name, ok := ipToName[ip.Address]
				if ok {
					name = ipToName[ip.Address]
				} else {
					names, _ := net.LookupAddr(ip.Address)
					ipToName[ip.Address] = names[0]
					name = names[0]
				}

				cputemp = requestCpuTemp(ip.Address)

				log.Printf("%v\n", cputemp)

				log.Printf("ip=%s, name=%s, cpuarch=%s, temp=%s, hostname=%s\n",
					ip.Address,
					name,
					cputemp.CpuArch,
					cputemp.Temp,
					cputemp.HostName)

				kubetemp = append(kubetemp,
					NodeTemp{ip.Address,
						name,
						cputemp})
			}
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(kubetemp)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/kubetemp", KubeTemp)
	// static pages
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	log.Fatal(http.ListenAndServe(":9999", router))
}
