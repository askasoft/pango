package disk

import (
	"syscall"
	"unsafe"
)

var (
	kernel              = syscall.MustLoadDLL("kernel32.dll")
	getDiskFreeSpaceExW = kernel.MustFindProc("GetDiskFreeSpaceExW")
)

// GetDiskStats returns an object holding the disk usage of volumePath
// or nil in case of error (invalid path, etc)
func GetDiskStats(volumePath string) (ds DiskStats, err error) {
	var pp *uint16
	pp, err = syscall.UTF16PtrFromString(volumePath)
	if err != nil {
		return
	}

	var ret uintptr
	ret, _, err = getDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(pp)),
		uintptr(unsafe.Pointer(&ds.Free)),
		uintptr(unsafe.Pointer(&ds.Total)),
		uintptr(unsafe.Pointer(&ds.Available)),
	)
	if ret != 0 {
		err = nil
	}

	return
}
