package mem

import "fmt"

type MemoryStats struct {
	Total     uint64 `json:"total"`      // the total accessible system memory in bytes.
	Free      uint64 `json:"free"`       // the total free system memory in bytes.
	Shared    uint64 `json:"shared"`     // the total shared system memory in bytes.
	Buffer    uint64 `json:"buffer"`     // the total buffer system memory in bytes.
	Cached    uint64 `json:"cached"`     // the total cached system memory in bytes.
	SwapTotal uint64 `json:"swap_total"` // the total swap memory in bytes.
	SwapFree  uint64 `json:"swap_free"`  // the total free swap memory in bytes.
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

func (ms *MemoryStats) String() string {
	return fmt.Sprintf("(T: %d, F: %d, S: %d, B: %d, C: %d, ST: %d, SF: %d)",
		ms.Total, ms.Free, ms.Shared, ms.Buffer, ms.Cached, ms.SwapTotal, ms.SwapFree)
}
