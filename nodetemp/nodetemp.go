package nodetemp

import (
	"github.com/jens18/rpi_temp/cputemp"
)

type NodeTemp struct {
	IpAddress string          `json:"ipAddress"`
	HostName  string          `json:"hostName"` // the real hostname
	NodeTemp  cputemp.CpuTemp `json:"nodeTemp"`
}
