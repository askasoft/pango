//go:build linux
// +build linux

package disk

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// GetDisksStats get disk I/O statistics.
func GetDisksStats() (DisksStats, error) {
	// Reference: Documentation/iostats.txt in the source of Linux
	file, err := os.Open("/proc/diskstats")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return collectDiskStats(file)
}

func collectDiskStats(out io.Reader) (stats DisksStats, err error) {
	var readed, written uint64

	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 14 {
			continue
		}

		name := fields[2]
		readed, err = strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			err = fmt.Errorf("failed to parse reads completed of %s", name)
			return
		}

		written, err = strconv.ParseUint(fields[7], 10, 64)
		if err != nil {
			err = fmt.Errorf("failed to parse writes completed of %s", name)
			return
		}

		stats = append(stats, DiskStats{
			Name:    name,
			Readed:  readed,
			Written: written,
		})
	}

	if err = scanner.Err(); err != nil {
		err = fmt.Errorf("scan error for /proc/diskstats: %w", err)
	}

	return
}
