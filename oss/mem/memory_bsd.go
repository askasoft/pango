//go:build freebsd || openbsd || dragonfly || netbsd

package mem

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

func GetMemoryStats() (ms MemoryStats, err error) {
	err = collectMemoryStats(&ms)
	return
}

type memStat struct {
	name  string
	ptr   *uint64
	scale *uint64
}

func collectMemoryStats(ms *MemoryStats) error {
	var pageSize uint64
	one := uint64(1)

	memStats := []memStat{
		{"vm.stats.vm.v_page_size", &pageSize, &one},
		{"hw.physmem", &ms.Total, &one},
		{"vm.stats.vm.v_cache_count", &ms.Cached, &pageSize},
		{"vm.stats.vm.v_free_count", &ms.Free, &pageSize},
		// {"vm.stats.vm.v_active_count", &ms.Active, &pageSize},
		// {"vm.stats.vm.v_inactive_count", &ms.Inactive, &pageSize},
		// {"vm.stats.vm.v_wire_count", &ms.Wired, &pageSize},
	}

	for _, stat := range memStats {
		ret, err := unix.SysctlRaw(stat.name)
		if err != nil {
			return fmt.Errorf("failed in sysctl %s: %w", stat.name, err)
		}

		if len(ret) == 8 {
			*stat.ptr = *(*uint64)(unsafe.Pointer(&ret[0])) * *stat.scale
		} else if len(ret) == 4 {
			*stat.ptr = uint64(*(*uint32)(unsafe.Pointer(&ret[0]))) * *stat.scale
		} else {
			return fmt.Errorf("failed in sysctl %s: %w", stat.name, err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// collect swap statistics from swapinfo command
	cmd := exec.CommandContext(ctx, "swapinfo", "-k")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	err = collectSwapStats(out, ms)
	if err != nil {
		go cmd.Wait()
		return err
	}
	return cmd.Wait()
}

func collectSwapStats(out io.Reader, ms *MemoryStats) error {
	scanner := bufio.NewScanner(out)
	if !scanner.Scan() {
		return fmt.Errorf("failed to scan output of swapinfo")
	}

	line := scanner.Text()
	if !strings.HasPrefix(line, "Device") {
		return fmt.Errorf("unexpected output of swapinfo: %s", line)
	}

	var total, used uint64
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}
		if v, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
			total += v * 1024
		}
		if v, err := strconv.ParseUint(fields[2], 10, 64); err == nil {
			used += v * 1024
		}
	}

	ms.SwapTotal = total
	ms.SwapFree = total - used
	return nil
}
