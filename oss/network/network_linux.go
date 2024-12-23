package network

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// GetNetworksStats get network statistics
func GetNetworksStats() (NetworksStats, error) {
	// Reference: man 5 proc, Documentation/filesystems/proc.txt in Linux source code
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return collectNetworksStats(file)
}

func collectNetworksStats(out io.Reader) (nss NetworksStats, err error) {
	var rxBytes, txBytes uint64

	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		// Reference: dev_seq_printf_stats in Linux source code
		kv := strings.SplitN(scanner.Text(), ":", 2)
		if len(kv) != 2 {
			continue
		}

		fields := strings.Fields(kv[1])
		if len(fields) < 16 {
			continue
		}

		name := strings.TrimSpace(kv[0])
		if name == "lo" {
			continue
		}

		if rxBytes, err = strconv.ParseUint(fields[0], 10, 64); err != nil {
			err = fmt.Errorf("failed to parse rxBytes of %s", name)
			return
		}

		if txBytes, err = strconv.ParseUint(fields[8], 10, 64); err != nil {
			err = fmt.Errorf("failed to parse txBytes of %s", name)
			return
		}

		nss = append(nss, NetworkStats{Name: name, Received: rxBytes, Transmitted: txBytes})
	}

	if err = scanner.Err(); err != nil {
		err = fmt.Errorf("scan error for /proc/net/dev: %w", err)
	}

	return
}
