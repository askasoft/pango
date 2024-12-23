package disk

import (
	"fmt"
	"time"
)

type DiskUsage struct {
	Free      uint64 `json:"free"`      // total free bytes on file system
	Available uint64 `json:"available"` // total available bytes on file system to an unprivileged user
	Total     uint64 `json:"total"`     // total size of the file system
}

// Used returns total bytes used in file system
func (du *DiskUsage) Used() uint64 {
	return du.Total - du.Free
}

// Usage returns percentage of use on the file system
func (du *DiskUsage) Usage() float64 {
	if du.Total == 0 {
		return 0
	}
	return float64(du.Total-du.Free) / float64(du.Total)
}

func (du *DiskUsage) String() string {
	return fmt.Sprintf("(F: %d, A: %d, T: %d)", du.Free, du.Available, du.Total)
}

// DiskIOStats represents disk I/O statistics
type DiskIOStats struct {
	Name    string `json:"name"`    // device name; like "hda"
	Readed  uint64 `json:"readed"`  // total number of reads completed successfully
	Written uint64 `json:"written"` // total number of writes completed successfully
}

func (ds *DiskIOStats) Subtract(s *DiskIOStats) {
	ds.Readed -= s.Readed
	ds.Written -= s.Written
}

func (ds *DiskIOStats) String() string {
	return fmt.Sprintf("(%q R: %d, W: %d)", ds.Name, ds.Readed, ds.Written)
}

type DisksIOStats []DiskIOStats

func (dss DisksIOStats) Subtract(ss DisksIOStats) {
	for _, ds := range dss {
		for _, s := range ss {
			if ds.Name == s.Name {
				ds.Subtract(&s)
				break
			}
		}
	}
}

type DiskIOUsage struct {
	DiskIOStats
	Delta time.Duration `json:"delta,omitempty"`
}

// ReadSpeed get read speed bytes/second
func (du *DiskIOUsage) ReadSpeed() float64 {
	if du.Delta == 0 {
		return 0
	}
	return float64(du.Readed) / du.Delta.Seconds()
}

// WriteSpeed get write speed bytes/second
func (du *DiskIOUsage) WriteSpeed() float64 {
	if du.Delta == 0 {
		return 0
	}
	return float64(du.Written) / du.Delta.Seconds()
}

func (du *DiskIOUsage) String() string {
	return fmt.Sprintf("(%q R: %d, W: %d, D: %s)", du.Name, du.Readed, du.Written, du.Delta)
}

type DisksIOUsage []DiskIOUsage

// GetDisksIOUsage get disks I/O usages between delta duration
func GetDisksIOUsage(delta time.Duration) (dsu DisksIOUsage, err error) {
	var dss1, dss2 DisksIOStats

	dss1, err = GetDisksStats()
	if err != nil {
		return
	}

	time.Sleep(delta)

	dss2, err = GetDisksStats()
	if err != nil {
		return
	}

	dss2.Subtract(dss1)

	dsu = make(DisksIOUsage, len(dss2))
	for i, ds := range dss2 {
		dsu[i].DiskIOStats = ds
		dsu[i].Delta = delta
	}

	return
}
