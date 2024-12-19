package mem

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func GetMemoryStats() (ms MemoryStats, err error) {
	var file *os.File

	// Reference: man 5 proc, Documentation/filesystems/proc.txt in Linux source code
	file, err = os.Open("/proc/meminfo")
	if err != nil {
		return
	}
	defer file.Close()

	err = collectMemoryStats(file, &ms)
	return
}

func collectMemoryStats(out io.Reader, ms *MemoryStats) error {
	scanner := bufio.NewScanner(out)

	memStats := map[string]*uint64{
		"MemTotal": &ms.Total,
		"MemFree":  &ms.Free,
		"Buffers":  &ms.Buffer,
		"Cached":   &ms.Cached,
		// "Active":       &ms.Active,
		// "Inactive":     &ms.Inactive,
		// "SwapCached":   &ms.SwapCached,
		"SwapTotal": &ms.SwapTotal,
		"SwapFree":  &ms.SwapFree,
		// "Mapped":       &ms.Mapped,
		// "Shmem":        &ms.Shmem,
		// "Slab":         &ms.Slab,
		// "PageTables":   &ms.PageTables,
		// "Committed_AS": &ms.Committed,
		// "VmallocUsed":  &ms.VmallocUsed,
	}

	for scanner.Scan() {
		line := scanner.Text()
		i := strings.IndexRune(line, ':')
		if i < 0 {
			continue
		}

		fld := line[:i]
		if ptr := memStats[fld]; ptr != nil {
			val := strings.TrimSpace(strings.TrimRight(line[i+1:], "kB"))
			if v, err := strconv.ParseUint(val, 10, 64); err == nil {
				*ptr = v * 1024
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan error for /proc/meminfo: %w", err)
	}

	return nil
}
