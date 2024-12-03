package du

import (
	"syscall"
	"unsafe"
)

type DiskUsage struct {
	freeBytes  int64
	totalBytes int64
	availBytes int64
}

var (
	kernel              = syscall.MustLoadDLL("kernel32.dll")
	getDiskFreeSpaceExW = kernel.MustFindProc("GetDiskFreeSpaceExW")
)

// NewDiskUsages returns an object holding the disk usage of volumePath
// or nil in case of error (invalid path, etc)
func NewDiskUsage(volumePath string) *DiskUsage {
	du := &DiskUsage{}

	pp, _ := syscall.UTF16PtrFromString(volumePath)

	getDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(pp)),
		uintptr(unsafe.Pointer(&du.freeBytes)),
		uintptr(unsafe.Pointer(&du.totalBytes)),
		uintptr(unsafe.Pointer(&du.availBytes)))

	return du
}

// Free returns total free bytes on file system
func (du *DiskUsage) Free() uint64 {
	return uint64(du.freeBytes)
}

// Available returns total available bytes on file system to an unprivileged user
func (du *DiskUsage) Available() uint64 {
	return uint64(du.availBytes)
}

// Total returns total size of the file system
func (du *DiskUsage) Total() uint64 {
	return uint64(du.totalBytes)
}

// Used returns total bytes used in file system
func (du *DiskUsage) Used() uint64 {
	return du.Total() - du.Free()
}

// Usage returns percentage of use on the file system
func (du *DiskUsage) Usage() float64 {
	return float64(du.Used()) / float64(du.Total())
}
