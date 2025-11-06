//go:build !windows

package disk

import "syscall"

// GetDiskUsage returns an object holding the disk usage of volumePath
// or nil in case of error (invalid path, etc)
func GetDiskUsage(volumePath string) (du DiskUsage, err error) {
	var sf syscall.Statfs_t

	err = syscall.Statfs(volumePath, &sf)
	if err == nil {
		du.Free = sf.Bfree * uint64(sf.Bsize)
		du.Available = sf.Bavail * uint64(sf.Bsize)
		du.Total = uint64(sf.Blocks) * uint64(sf.Bsize)
	}
	return
}
