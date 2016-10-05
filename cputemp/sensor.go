// Package cputemp provides Arm (Raspbian) and Intel (Ubuntu) CPU temperature
// measurements. AMD CPU temperature is not supported.
package cputemp

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// cpuArmTempSysFileName is the 'sys' file name containing CPU temperature
// on a Raspberry PI running Raspian/Hypriot (Debian).
// CPU + GPU are inside the same SOC. In the case of an RPI3: BCM 2837
// 64bit ARMv8 Cortex A53 Quad Core.
const cpuArmTempSysFileName string = "/sys/class/thermal/thermal_zone0/temp"

// cpuIntelTempSysFileName is the 'sys' file name containing CPU temperature
// on an Intel PC running Ubuntu (Debian).
// https://www.kernel.org/doc/Documentation/hwmon/coretemp
const cpuIntelTempSysFileName string = "/sys/devices/platform/coretemp.0/hwmon/hwmon2/temp1_input"

// cpuTempSysFileName is initialized with the 'sys' file name for the current
// machine architecture (arm or x86_64)
var cpuTempSysFileName string

// cpuArch is initialized with the CPU architecture (as returned by 'uname -m')
var cpuArch string

// CpuTemp stores temperature and hostname information. The temperature unit
// is celsius (Example: 46.3).
type CpuTemp struct {
	Temp     string `json:"temp"`
	CpuArch  string `json:"cpuArch"`
	HostName string `json:"hostName"`
}

// Get determines the current CPU temperature and returns a new cputemp instance.
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

	return CpuTemp{strconv.Itoa(temp1) + "." + strconv.Itoa(tempM), cpuArch, hostName}
}

// init determines the correct 'sys' file at startup time based on the
// output of the 'uname -m' command.
func init() {
	unameCmd := exec.Command("uname", "-m")
	unameOut, err := unameCmd.Output()
	if err != nil {
		panic(err)
	}

	cpuArch = strings.Trim(string(unameOut), "\n")

	switch cpuArch {
	case "x86_64":
		cpuTempSysFileName = cpuIntelTempSysFileName
	case "armv6l":
		cpuTempSysFileName = cpuArmTempSysFileName
	case "armv7l":
		cpuTempSysFileName = cpuArmTempSysFileName
	default:
		panic("machine architecture not supported")
	}

	// test if cpuTempSysFileName exists
	if _, err := os.Stat(cpuTempSysFileName); os.IsNotExist(err) {
		log.Printf("sensor.go: No sensors found at %s.\n",
			cpuTempSysFileName)
		//os.Exit(1)
	}
}
