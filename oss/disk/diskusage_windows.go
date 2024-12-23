package disk

import (
	"syscall"
	"unsafe"
)

var (
	kernel              = syscall.MustLoadDLL("kernel32.dll")
	getDiskFreeSpaceExW = kernel.MustFindProc("GetDiskFreeSpaceExW")
)

// GetDiskUsage returns an object holding the disk usage of volumePath
// or nil in case of error (invalid path, etc)
func GetDiskUsage(volumePath string) (du DiskUsage, err error) {
	var pp *uint16
	pp, err = syscall.UTF16PtrFromString(volumePath)
	if err != nil {
		return
	}

	var ret uintptr
	ret, _, err = getDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(pp)),
		uintptr(unsafe.Pointer(&du.Free)),
		uintptr(unsafe.Pointer(&du.Total)),
		uintptr(unsafe.Pointer(&du.Available)),
	)
	if ret != 0 {
		err = nil
	}

	return
}
