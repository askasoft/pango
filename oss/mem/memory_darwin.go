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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Reference: man 1 vm_stat
	cmd := exec.CommandContext(ctx, "vm_stat")

	var out io.ReadCloser
	out, err = cmd.StdoutPipe()
	if err != nil {
		return
	}
	if err = cmd.Start(); err != nil {
		return
	}

	err = collectMemoryStats(out, &ms)
	if err != nil {
		// it is needed to cleanup the process, but its result is not needed.
		go cmd.Wait() //nolint:errcheck
		return
	}

	if err = cmd.Wait(); err != nil {
		return
	}

	// Reference: sys/sysctl.h, man 3 sysctl, sysctl vm.swapusage
	var ret []byte
	ret, err = unix.SysctlRaw("vm.swapusage")
	if err != nil {
		err = fmt.Errorf("failed in sysctl vm.swapusage: %w", err)
		return
	}

	err = collectSwapStats(ret, &ms)
	return
}

// References:
//   - https://support.apple.com/guide/activity-monitor/view-memory-usage-actmntr1004/10.14/mac/11.0
//   - https://opensource.apple.com/source/system_cmds/system_cmds-880.60.2/vm_stat.tproj/
func collectMemoryStats(out io.Reader, ms *MemoryStats) error {
	scanner := bufio.NewScanner(out)
	if !scanner.Scan() {
		return fmt.Errorf("failed to scan output of vm_stat")
	}

	line := scanner.Text()
	var pageSize uint64
	if _, err := fmt.Sscanf(line, "Mach Virtual Memory Statistics: (page size of %d bytes)", &pageSize); err != nil {
		return fmt.Errorf("unexpected output of vm_stat: %s", line)
	}

	var active, inactive, speculative, wired, purgeable, fileBacked, compressed uint64
	memStats := map[string]*uint64{
		"Pages free":                   &ms.Free,
		"Pages active":                 &active,
		"Pages inactive":               &inactive,
		"Pages speculative":            &speculative,
		"Pages wired down":             &wired,
		"Pages purgeable":              &purgeable,
		"File-backed pages":            &fileBacked,
		"Pages occupied by compressor": &compressed,
	}
	for scanner.Scan() {
		line := scanner.Text()
		i := strings.IndexRune(line, ':')
		if i < 0 {
			continue
		}
		if ptr := memStats[line[:i]]; ptr != nil {
			val := strings.TrimRight(strings.TrimSpace(line[i+1:]), ".")
			if v, err := strconv.ParseUint(val, 10, 64); err == nil {
				*ptr = v * pageSize
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan error for vm_stat: %s", err)
	}

	ms.Cached = purgeable + fileBacked
	used := wired + compressed + active + inactive + speculative - ms.Cached
	ms.Total = used + ms.Cached + ms.Free
	return nil
}

// xsw_usage in sys/sysctl.h
type swapUsage struct {
	Total     uint64
	Avail     uint64
	Used      uint64
	Pagesize  int32
	Encrypted bool
}

func collectSwapStats(out []byte, ms *MemoryStats) error {
	if len(out) != 32 {
		return fmt.Errorf("unexpected output of sysctl vm.swapusage: %v (len: %d)", out, len(out))
	}

	su := (*swapUsage)(unsafe.Pointer(&out[0]))

	ms.SwapTotal = su.Total
	ms.SwapFree = su.Avail
	return nil
}
