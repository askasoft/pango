package loadavg

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

func GetLoadAvg() (la LoadAvg, err error) {
	var ret []byte

	ret, err = unix.SysctlRaw("vm.loadavg")
	if err != nil {
		err = fmt.Errorf("failed in sysctl vm.loadavg: %s", err)
		return
	}

	err = collectLoadavgStats(ret, &la)
	return
}

// loadavg in sys/sysctl.h
type loadavg struct {
	Loads  [3]uint32
	Fscale uint64
}

// Reference: sys/sysctl.h
func collectLoadavgStats(out []byte, la *LoadAvg) error {
	if len(out) != 24 {
		return fmt.Errorf("unexpected output of sysctl vm.loadavg: %v (len: %d)", out, len(out))
	}

	load := *(*loadavg)(unsafe.Pointer(&out[0]))

	la.Loadavg1 = float64(load.Loads[0]) / float64(load.Fscale)
	la.Loadavg5 = float64(load.Loads[1]) / float64(load.Fscale)
	la.Loadavg15 = float64(load.Loads[2]) / float64(load.Fscale)
	return nil
}
