package cpu

import "time"

// CPUStats represents cpu statistics for linux
type CPUStats struct {
	User      uint64
	Nice      uint64
	System    uint64
	Idle      uint64
	Iowait    uint64
	Irq       uint64
	Softirq   uint64
	Steal     uint64
	Guest     uint64
	GuestNice uint64
	Total     uint64
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

func (cs *CPUStats) Subtract(s *CPUStats) {
	cs.User -= s.User
	cs.Nice -= s.Nice
	cs.System -= s.System
	cs.Idle -= s.Idle
	cs.Iowait -= s.Iowait
	cs.Irq -= s.Irq
	cs.Softirq -= s.Softirq
	cs.Steal -= s.Steal
	cs.Guest -= s.Guest
	cs.GuestNice -= s.GuestNice
	cs.Total -= s.Total
}

type CPUStatsDelta struct {
	CPUStats
	Delta time.Duration
}

// GetCPUStatsDelta get cpu statistics between delta duration
func GetCPUStatsDelta(delta time.Duration) (csd CPUStatsDelta, err error) {
	var cs1 CPUStats

	cs1, err = GetCPUStats()
	if err != nil {
		return
	}

	time.Sleep(delta)

	csd.CPUStats, err = GetCPUStats()
	if err != nil {
		return
	}

	csd.Subtract(&cs1)
	csd.Delta = delta
	return
}
