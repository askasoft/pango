package cpu

import (
	"syscall"
	"time"
	"unsafe"
)

var (
	kernel         = syscall.MustLoadDLL("kernel32.dll")
	getSystemTimes = kernel.MustFindProc("GetSystemTimes")

	// Delta time duration between 2 GetSystemTimes() call
	Delta = time.Millisecond * 250
)

// GetCPUStats get cpu statistics
func GetCPUStats() (cpu CPUStats, err error) {
	var idle1, kernel1, user1 uint64
	var idle2, kernel2, user2 uint64
	var ret uintptr

	ret, _, err = getSystemTimes.Call(
		uintptr(unsafe.Pointer(&idle1)),
		uintptr(unsafe.Pointer(&kernel1)),
		uintptr(unsafe.Pointer(&user1)),
	)
	if ret == 0 {
		return
	}

	time.Sleep(Delta)

	ret, _, err = getSystemTimes.Call(
		uintptr(unsafe.Pointer(&idle2)),
		uintptr(unsafe.Pointer(&kernel2)),
		uintptr(unsafe.Pointer(&user2)),
	)
	if ret == 0 {
		return
	}

	err = nil
	cpu.Idle = idle2 - idle1
	cpu.System = kernel2 - kernel1 - cpu.Idle
	cpu.User = user2 - user1
	cpu.Total = cpu.Idle + cpu.System + cpu.User

	return
}
