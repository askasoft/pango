//go:build !linux
// +build !linux

package disk

import (
	"fmt"
	"runtime"
)

// GetDisksStats get disk I/O statistics.
func GetDisksStats() (dss DisksStats, err error) {
	err = fmt.Errorf("disk I/O statistics for %s is not supported", runtime.GOOS)
	return
}
