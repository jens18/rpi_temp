package cputemp

import (
	"log"
)

func ExampleCpuTemp() {
	var c CpuTemp
	c = c.Get()

	log.Printf("hostname = %s, cpu = %s, temperature = %s \n", c.HostName, c.CpuArch, c.Temp)
	// Output:
}
