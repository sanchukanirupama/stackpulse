package config

import (
	"fmt"
	"time"
)

type ServiceConfig struct {
	Host            string        `yaml:"host" json:"host"`
	Port            int           `yaml:"port" json:"port"`
	PID             int           `yaml:"pid" json:"pid"`
	InspectPort     int           `yaml:"inspectPort" json:"inspectPort"`
	PollingInterval time.Duration `yaml:"pollingInterval" json:"pollingInterval"`
	HeapLimit       string        `yaml:"heapLimit" json:"heapLimit"`
	CPUThreshold    float64       `yaml:"cpuThreshold" json:"cpuThreshold"`
}

func (sc *ServiceConfig) Validate() error {
	if sc.PID == 0 && sc.Port == 0 {
		return fmt.Errorf("must specify either PID or port")
	}
	
	if sc.CPUThreshold <= 0 || sc.CPUThreshold > 100 {
		return fmt.Errorf("CPU threshold must be between 0 and 100")
	}
	
	if sc.PollingInterval < time.Millisecond {
		return fmt.Errorf("polling interval must be at least 1ms")
	}
	
	return nil
}