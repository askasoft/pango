package cpu

// CPUStats represents cpu statistics for linux
type CPUStats struct {
	User, Nice, System, Idle, Iowait, Irq, Softirq, Steal, Guest, GuestNice, Total uint64
}

func (cs *CPUStats) calcUsage(v uint64) float64 {
	if cs.Total == 0 {
		return 0
	}
	return float64(v) / float64(cs.Total)
}

func (cs *CPUStats) CPUUsage() float64 {
	return cs.calcUsage(cs.Total - cs.Idle)
}

func (cs *CPUStats) UserUsage() float64 {
	return cs.calcUsage(cs.User)
}

func (cs *CPUStats) NiceUsage() float64 {
	return cs.calcUsage(cs.Nice)
}

func (cs *CPUStats) SystemUsage() float64 {
	return cs.calcUsage(cs.System)
}

func (cs *CPUStats) IdleUsage() float64 {
	return cs.calcUsage(cs.Idle)
}

func (cs *CPUStats) IowaitUsage() float64 {
	return cs.calcUsage(cs.Iowait)
}

func (cs *CPUStats) IrqUsage() float64 {
	return cs.calcUsage(cs.Irq)
}

func (cs *CPUStats) SoftirqUsage() float64 {
	return cs.calcUsage(cs.Softirq)
}

func (cs *CPUStats) StealUsage() float64 {
	return cs.calcUsage(cs.Steal)
}

func (cs *CPUStats) GuestUsage() float64 {
	return cs.calcUsage(cs.Guest)
}

func (cs *CPUStats) GuestNiceUsage() float64 {
	return cs.calcUsage(cs.GuestNice)
}
