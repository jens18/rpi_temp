// inspired by: http://reprage.com/post/using-golang-to-connect-raspberrypi-and-arduino

package main

import (
	"encoding/json"
	"github.com/huin/goserial"
	"github.com/jens18/rpi_temp/nodetemp"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

/*
type NodeTemp struct {
	IpAddress string          `json:"ipAddress"`
	HostName  string          `json:"hostName"` // the real hostname
	NodeTemp  cputemp.CpuTemp `json:"nodeTemp"`
}
*/

var netClient = &http.Client{
	Timeout: time.Second * 10,
}

func requestKubeTemp(clusterTemp *[]string) {
	var nodeTemp []nodetemp.NodeTemp

	// http client with sensible timeout
	// (https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779#.s0bop6t5g)
	response, err := netClient.Get("http://192.168.29.13:9999/kubetemp")
	if err != nil {
		log.Fatal(err)
	}

	buf, _ := ioutil.ReadAll(response.Body)
	log.Printf(string(buf))
	json.Unmarshal(buf, &nodeTemp)
	for _, val := range nodeTemp {
		n := strings.TrimSuffix(val.HostName, ".mesgtone.lan.") + " | " + val.NodeTemp.Temp + " C"

		*clusterTemp = append(*clusterTemp, n)
		log.Printf("%s\n", n)
	}
}

func main() {
	// Find the device that represents the arduino serial
	// connection.
	c := &goserial.Config{Name: "/dev/ttyACM1", Baud: 9600}
	s, err := goserial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	// When connecting to an older revision Arduino, you need to wait
	// a little while it resets.
	time.Sleep(1 * time.Second)

	/*
		clusterTemp := []string{
			"carmel  | 49 C\n",
			"aptos   | 47 C\n",
			"venice  | 49 C\n",
			"salinas | 56 C\n"}
	*/

	var clusterTemp []string

	for {
		log.Printf("%q\n", clusterTemp)

		clusterTemp = clusterTemp[:0]

		requestKubeTemp(&clusterTemp)

		for _, t := range clusterTemp {
			_, err := s.Write([]byte(t))
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("%s\n", t)
			time.Sleep(5000 * time.Millisecond)

		}
	}
	/*
		buf := make([]byte, 128)
		for {
			n, err = s.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			// log.Print("%q", buf[:n])
			log.Printf("%s\n", string(buf[:n]))
		}
	*/
}
