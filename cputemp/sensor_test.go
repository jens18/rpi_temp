package cputemp

import (
	"testing"
)

func TestTemp(t *testing.T) {
	var c CpuTemp
	
	if !(c.Get().HostName != "") {
		t.Error(`CpuTemp.Get().Hostname not set`)
	}

	if !(c.Get().Temp[len(c.Get().Temp) - 2] == '.') {
		t.Error(`CpuTemp.Get().Temp format should have 1 position after the decimal point`)
	}

	

}
