package cpu

import (
	"fmt"
	"time"
)

// CPUStats represents cpu statistics for linux
type CPUStats struct {
	User      uint64 `json:"user"`
	Nice      uint64 `json:"nice"`
	System    uint64 `json:"system"`
	Idle      uint64 `json:"idle"`
	Iowait    uint64 `json:"iowait"`
	Irq       uint64 `json:"irq"`
	Softirq   uint64 `json:"softirq"`
	Steal     uint64 `json:"steal"`
	Guest     uint64 `json:"guest"`
	GuestNice uint64 `json:"guestnice"`
	Total     uint64 `json:"total"`
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

func (cs *CPUStats) String() string {
	return fmt.Sprintf(
		"(us: %d, sy: %d, ni: %d, id: %d, wa: %d, hi: %d, si: %d, st: %d, gu: %d, gn: %d)",
		cs.User, cs.System, cs.Nice, cs.Idle, cs.Iowait,
		cs.Irq, cs.Softirq, cs.Steal, cs.Guest, cs.GuestNice,
	)
}

type CPUUsage struct {
	CPUStats
	Delta time.Duration `json:"delta,omitempty"`
}

// GetCPUUsage get cpu usage between delta duration
func GetCPUUsage(delta time.Duration) (cu CPUUsage, err error) {
	var cs1, cs2 CPUStats

	cs1, err = GetCPUStats()
	if err != nil {
		return
	}

	time.Sleep(delta)

	cs2, err = GetCPUStats()
	if err != nil {
		return
	}

	cu.CPUStats = cs2
	cu.Subtract(&cs1)
	cu.Delta = delta
	return
}
