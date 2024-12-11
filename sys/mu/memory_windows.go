package mu

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

// omitting a few fields for brevity...
// https://msdn.microsoft.com/en-us/library/windows/desktop/aa366589(v=vs.85).aspx
type memStatusEx struct {
	dwLength     uint32
	dwMemoryLoad uint32
	ullTotalPhys uint64
	ullAvailPhys uint64
	unused       [5]uint64
}

func sysTotalMemory() uint64 {
	msx := &memStatusEx{
		dwLength: 64,
	}

	r, _, _ := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(msx)))
	if r == 0 {
		return 0
	}
	return msx.ullTotalPhys
}

func sysFreeMemory() uint64 {
	msx := &memStatusEx{
		dwLength: 64,
	}
	r, _, _ := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(msx)))
	if r == 0 {
		return 0
	}
	return msx.ullAvailPhys
}
