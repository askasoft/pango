package cpu

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// GetCPUStats get cpu statistics
func GetCPUStats() (cpu CPUStats, err error) {
	var file *os.File

	// Reference: man 5 proc, Documentation/filesystems/proc.txt in Linux source code
	file, err = os.Open("/proc/stat")
	if err != nil {
		return
	}
	defer file.Close()

	err = collectCPUStats(file, &cpu)
	return
}

type cpuStat struct {
	name string
	ptr  *uint64
}

func collectCPUStats(out io.Reader, cpu *CPUStats) error {
	scanner := bufio.NewScanner(out)

	cpuStats := []cpuStat{
		{"user", &cpu.User},
		{"nice", &cpu.Nice},
		{"system", &cpu.System},
		{"idle", &cpu.Idle},
		{"iowait", &cpu.Iowait},
		{"irq", &cpu.Irq},
		{"softirq", &cpu.Softirq},
		{"steal", &cpu.Steal},
		{"guest", &cpu.Guest},
		{"guest_nice", &cpu.GuestNice},
	}

	if !scanner.Scan() {
		return errors.New("failed to scan /proc/stat")
	}

	valStrs := strings.Fields(scanner.Text())[1:]
	for i, valStr := range valStrs {
		val, err := strconv.ParseUint(valStr, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to scan %q from /proc/stat", cpuStats[i].name)
		}

		*cpuStats[i].ptr = val
		cpu.Total += val
	}

	// Since cpustat[CPUTIME_USER] includes cpustat[CPUTIME_GUEST], subtract the duplicated values from total.
	// https://github.com/torvalds/linux/blob/4ec9f7a18/kernel/sched/cputime.c#L151-L158
	cpu.Total -= cpu.Guest

	// cpustat[CPUTIME_NICE] includes cpustat[CPUTIME_GUEST_NICE]
	cpu.Total -= cpu.GuestNice

	return nil
}
