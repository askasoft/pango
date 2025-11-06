//go:build darwin && !cgo

package cpu

import "fmt"

// GetCPUStats get cpu statistics
// CPU counters for darwin is unavailable without cgo.
func GetCPUStats() (cpu CPUStats, err error) {
	err = fmt.Errorf("CPUStats for darwin is not supported")
	return
}
