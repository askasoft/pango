package icap

import "testing"

func TestDoubleup(t *testing.T) {
	cs := []struct {
		s int
		n int
		w int
	}{
		{0, 1, 16},
		{16, 2, 32},
		{32, 64, 128},
	}

	for i, c := range cs {
		a := Doubleup(c.s, c.s+c.n)
		if a != c.w {
			t.Errorf("[%d] growup(%d) = %d, want %d", i, c.s, a, c.w)
		}
	}
}

func TestRoundup(t *testing.T) {
	cs := []struct {
		s int
		w int
	}{
		{0, 0},
		{10, 32},
		{20, 32},
		{31, 32},
		{32, 32},
		{33, 64},
		{63, 64},
		{65, 96},
	}

	for i, c := range cs {
		a := Roundup(c.s, 32)
		if a != c.w {
			t.Errorf("[%d] roundup(%d) = %d, want %d", i, c.s, a, c.w)
		}
	}
}
