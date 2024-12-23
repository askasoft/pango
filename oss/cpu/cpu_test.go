package cpu

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestGetCPUUsage(t *testing.T) {
	for i := 0; i < 8; i++ {
		time.Sleep(time.Millisecond * 250)

		cpu, err := GetCPUUsage(time.Millisecond * 250)
		if err != nil {
			t.Fatalf("error should be nil but got: %v", err)
		}

		if cpu.Total <= 0 {
			t.Fatalf("invalid cpu value: %+v", cpu)
		}

		bs, _ := json.Marshal(cpu)
		fmt.Printf("[%d] cpu value: %s\n", i, string(bs))
		fmt.Printf("[%d] cpu value: %+v\n", i, cpu)
		fmt.Printf("[%d] cpu usage: %.2f%%\n", i, cpu.CPUUsage()*100)
	}
}
