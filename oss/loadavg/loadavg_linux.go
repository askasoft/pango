package loadavg

import "syscall"

// func get() (*Stats, error) {
// 	// Reference: man 5 proc, loadavg_proc_show in Linux source code
// 	file, err := os.Open("/proc/loadavg")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	return collectLoadavgStats(file)
// }

// func collectLoadavgStats(out io.Reader) (*Stats, error) {
// 	var loadavg Stats
// 	ret, err := fmt.Fscanf(out, "%f %f %f", &loadavg.Loadavg1, &loadavg.Loadavg5, &loadavg.Loadavg15)
// 	if err != nil || ret != 3 {
// 		return nil, fmt.Errorf("unexpected format of /proc/loadavg")
// 	}
// 	return &loadavg, nil
// }

const si_load_shift = float64(1 << 16)

func GetLoadAvg() (la LoadAvg, err error) {
	var si syscall.Sysinfo_t

	if err = syscall.Sysinfo(&si); err != nil {
		return
	}

	la.Loadavg1 = float64(si.Loads[0]) / si_load_shift
	la.Loadavg5 = float64(si.Loads[1]) / si_load_shift
	la.Loadavg15 = float64(si.Loads[2]) / si_load_shift
	return
}
