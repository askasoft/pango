package cpu

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestGetCPU(t *testing.T) {
	if runtime.GOOS == "windows" {
		time.Sleep(time.Second)
	}

	cpu, err := GetCPUStats()
	if err != nil {
		t.Fatalf("error should be nil but got: %v", err)
	}

	if cpu.Total <= 0 {
		t.Fatalf("invalid cpu value: %+v", cpu)
	}

	fmt.Printf("cpu value: %+v\n", cpu)
	fmt.Printf("cpu usage: %v\n", cpu.CPUUsage())
}
