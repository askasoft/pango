//go:build !linux && !darwin && !freebsd && !netbsd

package network

import (
	"fmt"
	"runtime"
)

// GetNetworksStats network statistics
func GetNetworksStats() (nss NetworksStats, err error) {
	err = fmt.Errorf("network statistics for %s is not supported", runtime.GOOS)
	return
}
