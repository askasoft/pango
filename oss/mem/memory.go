package mem

type MemoryStats struct {
	// Total the total accessible system memory in bytes.
	Total uint64

	// Free the total free system memory in bytes.
	Free uint64

	// Shared the total shared system memory in bytes.
	Shared uint64

	// Buffer the total buffer system memory in bytes.
	Buffer uint64

	// Cached the total cached system memory in bytes.
	Cached uint64

	SwapTotal uint64
	SwapFree  uint64
}

func (ms *MemoryStats) Used() uint64 {
	return ms.Total - ms.Free - ms.Buffer - ms.Cached
}

func (ms *MemoryStats) Available() uint64 {
	return ms.Free + ms.Buffer + ms.Cached
}

func (ms *MemoryStats) SwapUsed() uint64 {
	return ms.SwapTotal - ms.SwapFree
}

func (ms *MemoryStats) Usage() float64 {
	if ms.Total == 0 {
		return 0
	}
	return float64(ms.Used()) / float64(ms.Total)
}
