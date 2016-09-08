package cputemp

import (
	"log"
)

func ExampleCpuTemp() {
	var c CpuTemp
	c = c.Get()

	log.Printf("hostname = %s, temperature = %s \n", c.HostName, c.Temp)
	// Output:
}
