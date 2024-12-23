//go:build darwin && cgo
// +build darwin,cgo

package cpu

import (
	"fmt"
	"unsafe"
)

// #include <mach/mach_host.h>
// #include <mach/host_info.h>
import "C"

// GetCPUStats cpu statistics
func GetCPUStats() (cs CPUStats, err error) {
	err = collectCPUStats(&cs)
	return
}

func collectCPUStats(cs *CPUStats) error {
	var cpuLoad C.host_cpu_load_info_data_t
	var count C.mach_msg_type_number_t = C.HOST_CPU_LOAD_INFO_COUNT

	ret := C.host_statistics(C.host_t(C.mach_host_self()), C.HOST_CPU_LOAD_INFO, C.host_info_t(unsafe.Pointer(&cpuLoad)), &count)
	if ret != C.KERN_SUCCESS {
		return fmt.Errorf("host_statistics failed: %d", ret)
	}

	cs.User = uint64(cpuLoad.cpu_ticks[C.CPU_STATE_USER])
	cs.System = uint64(cpuLoad.cpu_ticks[C.CPU_STATE_SYSTEM])
	cs.Idle = uint64(cpuLoad.cpu_ticks[C.CPU_STATE_IDLE])
	cs.Nice = uint64(cpuLoad.cpu_ticks[C.CPU_STATE_NICE])
	cs.Total = cs.User + cs.System + cs.Idle + cs.Nice
	return nil
}
