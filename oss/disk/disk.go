package disk

import "time"

type DiskUsage struct {
	// Free total free bytes on file system
	Free uint64

	// Available total available bytes on file system to an unprivileged user
	Available uint64

	// Total total size of the file system
	Total uint64
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

// DiskStats represents disk I/O statistics
type DiskStats struct {
	Name    string // device name; like "hda"
	Readed  uint64 // total number of reads completed successfully
	Written uint64 // total number of writes completed successfully
}

func (ds *DiskStats) Subtract(s *DiskStats) {
	ds.Readed -= s.Readed
	ds.Written -= s.Written
}

type DisksStats []DiskStats

func (dss DisksStats) Subtract(ss DisksStats) {
	for _, ds := range dss {
		for _, s := range ss {
			if ds.Name == s.Name {
				ds.Subtract(&s)
				break
			}
		}
	}
}

type DiskStatsDelta struct {
	DiskStats
	Delta time.Duration
}

// ReadSpeed get read speed bytes/second
func (dsd *DiskStatsDelta) ReadSpeed() float64 {
	if dsd.Delta == 0 {
		return 0
	}
	return float64(dsd.Readed) / dsd.Delta.Seconds()
}

// WriteSpeed get write speed bytes/second
func (dsd *DiskStatsDelta) WriteSpeed() float64 {
	if dsd.Delta == 0 {
		return 0
	}
	return float64(dsd.Written) / dsd.Delta.Seconds()
}

type DisksStatsDelta []DiskStatsDelta

// GetDisksStatsDelta get disk statistics between delta duration
func GetDisksStatsDelta(delta time.Duration) (dssd DisksStatsDelta, err error) {
	var dss1, dss2 DisksStats

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

	dssd = make(DisksStatsDelta, len(dss2))
	for i, ds := range dss2 {
		dssd[i].DiskStats = ds
		dssd[i].Delta = delta
	}

	return
}
