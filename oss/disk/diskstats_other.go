//go:build !linux

package disk

import (
	"fmt"
	"runtime"
)

// GetDisksStats get disk I/O statistics.
func GetDisksStats() (dss DisksIOStats, err error) {
	err = fmt.Errorf("disk I/O statistics for %s is not supported", runtime.GOOS)
	return
}
