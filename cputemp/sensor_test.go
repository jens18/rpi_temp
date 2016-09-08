package cputemp

import (
	"encoding/json"
	"testing"
)

func TestTemp(t *testing.T) {
	var c CpuTemp

	if !(c.Get().HostName != "") {
		t.Error(`CpuTemp.Get().Hostname not set`)
	}

	if !(c.Get().Temp[len(c.Get().Temp)-2] == '.') {
		t.Error(`CpuTemp.Get().Temp format should have 1 position after the decimal point`)
	}
}

func TestJsonEncoding(t *testing.T) {
	var e, d CpuTemp
	var b []byte
	var err error

	// convert into JSON
	if b, err = json.Marshal(e.Get()); err != nil {
		t.Error(`JSON encoding error`)
	}

	// convert back into struct
	if err = json.Unmarshal(b, &d); err != nil {
		t.Error(`JSON decoding error`)
	}
}
