package main

import (
	"encoding/json"
	"fmt"
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

func requestCpuTemp(hostIP string) {
	var cputemp cputemp.CpuTemp

	// http client with sensible timeout
	// (https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779#.s0bop6t5g)

	response, err := netClient.Get("http://" + hostIP + ":30001")
	if err != nil {
		log.Fatal(err)
	}

	buf, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(buf, &cputemp)

	fmt.Printf("temp=%s, hostname=%s\n", cputemp.Temp, cputemp.HostName)

	response.Body.Close()
}

func main() {

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
				fmt.Printf("ip=%s, ", ip.Address)
				requestCpuTemp(ip.Address)
				//time.Sleep(1000 * time.Millisecond)
			}
		}
	}
}
