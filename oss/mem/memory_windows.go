package mem

import (
	"syscall"
	"unsafe"
)

var (
	kernel = syscall.MustLoadDLL("kernel32.dll")

	// GetPhysicallyInstalledSystemMemory is simpler, but broken on
	// older versions of windows (and uses this under the hood anyway).
	globalMemoryStatusEx = kernel.MustFindProc("GlobalMemoryStatusEx")
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa366589(v=vs.85).aspx
type memStatusEx struct {
	dwLength                uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

func GetMemoryStats() (ms MemoryStats, err error) {
	mse := memStatusEx{dwLength: 64}

	var ret uintptr
	ret, _, err = globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&mse)))
	if ret != 0 {
		err = nil
		ms.Total = mse.ullTotalPhys
		ms.Free = mse.ullAvailPhys
		ms.SwapTotal = mse.ullTotalPageFile
		ms.SwapFree = mse.ullAvailPageFile
	}

	return
}
