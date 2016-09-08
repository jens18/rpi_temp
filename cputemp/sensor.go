//

package cputemp

import (
	//	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// The sys file containing CPU temperature on a Raspberry PI running the Raspian/Hypriot OS (Debian)
// CPU + GPU are inside the same SOC: BCM 2837 64bit ARMv8 Cortex A53 Quad Core
const cpuArmTempSysFileName string = "/sys/class/thermal/thermal_zone0/temp"

// The sys file containing CPU temperature on an Intel PC running the Ubuntu (Debian)
// (sys file information from 'strace sensors')
const cpuIntelTempSysFileName string = "/sys/class/hwmon/hwmon2/temp1_input"

var cpuTempSysFileName string

type CpuTemp struct {
	Temp     string `json:"temp"`
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

// Determine the correct 'sys' file at startup time.
func init() {
	unameCmd := exec.Command("uname", "-m")
	unameOut, err := unameCmd.Output()
	if err != nil {
		panic(err)
	}

	switch strings.Trim(string(unameOut), "\n") {
	case "x86_64":
		cpuTempSysFileName = cpuIntelTempSysFileName
	case "armv6l":
		cpuTempSysFileName = cpuArmTempSysFileName
	case "armv7l":
		cpuTempSysFileName = cpuArmTempSysFileName
	default:
		panic("machine architecture not supported")
	}
}
