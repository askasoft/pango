package osm

import (
	"time"

	"github.com/askasoft/pango/cog/ringbuffer"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/oss/cpu"
	"github.com/askasoft/pango/oss/disk"
	"github.com/askasoft/pango/oss/mem"
)

type StatsMonitor struct {
	Logger    log.Logger
	Interval  time.Duration // stats collect interval
	DiskFree  int64         // disk free threshold (average by count)
	DiskCount int           // disk stats collect count
	CPUUsage  float64       // cpu usage threshold (average by count)
	CPUCount  int           // cpu stats collect count
	MemUsage  float64       // memory usage threshold (average by count)
	MemCount  int           // mem stats collect count

	timer *time.Timer

	disks ringbuffer.RingBuffer[uint64]
	cpus  ringbuffer.RingBuffer[float64]
	mems  ringbuffer.RingBuffer[float64]

	lastCPUStats cpu.CPUStats
	lastCPUUsage cpu.CPUUsage
}

func NewStatsMonitor() *StatsMonitor {
	return &StatsMonitor{
		Logger:    log.GetLogger("OSM"),
		Interval:  time.Second,
		DiskFree:  num.GB,
		DiskCount: 5,
		CPUUsage:  0.9,
		CPUCount:  5,
		MemUsage:  0.9,
		MemCount:  5,
	}
}

// Start start monitor
func (sm *StatsMonitor) Start() {
	if sm.timer != nil {
		return
	}

	sm.Logger.Info("ossmonitor: start")

	go sm.start()
}

func (sm *StatsMonitor) start() {
	sm.timer = time.NewTimer(sm.Interval)
	for {
		if _, ok := <-sm.timer.C; !ok {
			break
		}

		sm.Collect()
		sm.timer.Reset(sm.Interval)
	}
}

// Stop stop monitor
func (sm *StatsMonitor) Stop() {
	if timer := sm.timer; timer != nil {
		sm.Logger.Info("ossmonitor: stop")

		timer.Stop()
		sm.timer = nil
	}
}

// Monitoring return the monitor is running or not
func (sm *StatsMonitor) Monitoring() bool {
	return sm.timer != nil
}

// LastCPUStats get the last cpu stats
func (sm *StatsMonitor) LastCPUStats() cpu.CPUStats {
	return sm.lastCPUStats
}

// LastCPUUsage get the last cpu usage
func (sm *StatsMonitor) LastCPUUsage() cpu.CPUUsage {
	return sm.lastCPUUsage
}

func (sm *StatsMonitor) Collect() {
	sm.collectDisk()
	sm.collectCPUUsage()
	sm.collectMemUsage()
}

func (sm *StatsMonitor) collectDisk() {
	if sm.DiskFree <= 0 {
		return
	}

	du, err := disk.GetDiskUsage(".")
	if err != nil {
		sm.Logger.Error(err)
		return
	}

	sm.Logger.Infof("ossmonitor: collect disk usage %s", du.String())

	sm.disks.Push(du.Available)
	if sm.disks.Len() > sm.DiskCount {
		sm.disks.Poll()
	}

	if sm.disks.Len() == sm.DiskCount {
		daa := calcAverage(sm.disks)
		if daa < uint64(sm.DiskFree) {
			sm.Logger.Errorf("insufficient free disk space %s", num.HumanSize(du.Available))

			sm.disks.Clear()
			sm.disks.Push(du.Available)
		}
	}
}

func (sm *StatsMonitor) collectCPUUsage() {
	if sm.CPUUsage <= 0 {
		return
	}

	cs, err := cpu.GetCPUStats()
	if err != nil {
		sm.Logger.Error(err)
		return
	}

	sm.Logger.Infof("ossmonitor: collect cpu statistics %s", cs.String())

	thisCPUStats := cs
	cs.Subtract(&sm.lastCPUStats)

	sm.lastCPUStats = thisCPUStats
	sm.lastCPUUsage.CPUStats = cs
	sm.lastCPUUsage.Delta = sm.Interval

	sm.cpus.Push(cs.CPUUsage())
	if sm.cpus.Len() > sm.CPUCount {
		sm.cpus.Poll()
	}

	if sm.cpus.Len() == sm.CPUCount {
		cua := calcAverage(sm.cpus)
		if cua > sm.CPUUsage {
			sm.Logger.Errorf("cpu usage %.2f%% is too high", cua*100)

			sm.cpus.Clear()
			sm.cpus.Push(cs.CPUUsage())
		}
	}
}

func (sm *StatsMonitor) collectMemUsage() {
	if sm.MemUsage <= 0 {
		return
	}

	ms, err := mem.GetMemoryStats()
	if err != nil {
		sm.Logger.Error(err)
		return
	}

	sm.Logger.Infof("ossmonitor: collect memory statistics %s", ms.String())

	sm.mems.Push(ms.Usage())
	if sm.mems.Len() > sm.MemCount {
		sm.mems.Poll()
	}

	if sm.mems.Len() == sm.MemCount {
		mua := calcAverage(sm.mems)
		if mua > sm.MemUsage {
			sm.Logger.Errorf("memory usage %.2f%% is too high", mua*100)

			sm.mems.Clear()
			sm.mems.Push(ms.Usage())
		}
	}
}

func calcAverage[E uint64 | float64](rb ringbuffer.RingBuffer[E]) E {
	var total E
	for it := rb.Iterator(); it.Next(); {
		total += it.Value()
	}
	return total / E(rb.Len())
}
