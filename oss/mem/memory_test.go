package mem

import "testing"

func TestNotZero(t *testing.T) {
	ms, err := GetMemoryStats()

	if err != nil {
		t.Fatal(err)
	}
	if ms.Total == 0 {
		t.Fatal("TotalMemory returned 0")
	}
	if ms.Free == 0 {
		t.Fatal("FreeMemory returned 0")
	}
}
