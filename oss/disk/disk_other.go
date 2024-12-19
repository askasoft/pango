//go:build !windows
// +build !windows

package disk

import "syscall"

// GetDiskStats returns an object holding the disk usage of volumePath
// or nil in case of error (invalid path, etc)
func GetDiskStats(volumePath string) (ds DiskStats, err error) {
	var sf syscall.Statfs_t

	err = syscall.Statfs(volumePath, &sf)
	if err == nil {
		ds.Free = sf.Bfree * uint64(sf.Bsize)
		ds.Available = sf.Bavail * uint64(sf.Bsize)
		ds.Total = uint64(sf.Blocks) * uint64(sf.Bsize)
	}
	return
}
