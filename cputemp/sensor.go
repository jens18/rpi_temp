package cputemp

import (
	"strings"
	"strconv"
	"os"
	"io/ioutil"
)


// CPU + GPU are inside the same SOC: BCM 2837 64bit ARMv8 Cortex A53 Quad Core
const cpuTempSysFileName string = "/sys/class/thermal/thermal_zone0/temp"

type CpuTemp struct {
	Temp string `json:"temp"`
	HostName string `json:"hostName"`
}

func (c *CpuTemp) Get() CpuTemp {
	dat, err := ioutil.ReadFile(cpuTempSysFileName)
	if err != nil {
		panic(err)
	}

	cpuTempRaw := strings.Trim(string(dat), "\n")

	temp, err := strconv.ParseInt(cpuTempRaw, 10, 32)

	// temp is a 5 digit integer
	// example: 48312
	
        temp1 := int(temp / 1000) 
        temp2 := int(temp / 100)
	tempM := int(temp2 % temp1)

	// note: this method does not work in a Docker container
	hostName, err := os.Hostname()

	return CpuTemp{strconv.Itoa(temp1) + "." + strconv.Itoa(tempM), hostName}
}
	
