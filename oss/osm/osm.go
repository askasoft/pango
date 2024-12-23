package osm

import "github.com/askasoft/pango/oss/cpu"

var _osm = NewStatsMonitor()

// Default returns the default StatsMonitor instance used by the package-level functions.
func Default() *StatsMonitor {
	return _osm
}

// SetDefault set the default StatsMonitor instance used by the package-level functions.
func SetDefault(sm *StatsMonitor) {
	_osm = sm
}

// Start start monitor
func Start() {
	_osm.Start()
}

// Stop stop monitor
func Stop() {
	_osm.Stop()
}

// Monitoring return the monitor is running or not
func Monitoring() bool {
	return _osm.Monitoring()
}

// LastCPUStats get the last cpu stats
func LastCPUStats() cpu.CPUStats {
	return _osm.LastCPUStats()
}

// LastCPUUsage get the last cpu usage
func LastCPUUsage() cpu.CPUUsage {
	return _osm.LastCPUUsage()
}
