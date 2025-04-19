package npid

import (
	"testing"
)

func TestBuild(t *testing.T) {
	cs := []struct {
		ip4     uint32
		pid     int
		ip4Bits int
		pidBits int
		w       int64
	}{
		{0xA1, 0xB2, 4, 0, 0x01},
		{0xA1, 0xB2, 0, 4, 0x02},
		{0xF1, 0xF1, 4, 4, 0x11},
		{0xFF, 0xF1, 4, 4, 0xF1},
		{0xE1, 0xE1, 4, 4, 0x11},
		{0xEF, 0xE1, 4, 4, 0xF1},
	}

	for i, c := range cs {
		a := Build(c.ip4, c.pid, c.ip4Bits, c.pidBits)
		if a != c.w {
			t.Errorf("#%d Build(%x, %x, %d, %d) = %x, want %x", i, c.ip4, c.pid, c.ip4Bits, c.pidBits, a, c.w)
		}
	}
}
