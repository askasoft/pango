package disk

type DiskStats struct {
	// Free total free bytes on file system
	Free uint64

	// Available total available bytes on file system to an unprivileged user
	Available uint64

	// Total total size of the file system
	Total uint64
}

// Used returns total bytes used in file system
func (ds *DiskStats) Used() uint64 {
	return ds.Total - ds.Free
}

// Usage returns percentage of use on the file system
func (ds *DiskStats) Usage() float64 {
	if ds.Total == 0 {
		return 0
	}
	return float64(ds.Total-ds.Free) / float64(ds.Total)
}
