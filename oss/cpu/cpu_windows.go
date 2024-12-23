package cpu

import (
	"syscall"
	"unsafe"
)

var (
	kernel         = syscall.MustLoadDLL("kernel32.dll")
	getSystemTimes = kernel.MustFindProc("GetSystemTimes")
)

// GetCPUStats get cpu statistics
func GetCPUStats() (cpu CPUStats, err error) {
	var kernel uint64
	var ret uintptr

	ret, _, err = getSystemTimes.Call(
		uintptr(unsafe.Pointer(&cpu.Idle)),
		uintptr(unsafe.Pointer(&kernel)),
		uintptr(unsafe.Pointer(&cpu.User)),
	)
	if ret == 0 {
		return
	}

	err = nil
	cpu.System = kernel - cpu.Idle
	cpu.Total = cpu.Idle + cpu.System + cpu.User
	return
}
