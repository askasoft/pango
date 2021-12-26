package col

import "testing"

func assertPanics(t *testing.T, name string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("%s: didn't panic as expected", name)
		}
	}()

	f()
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
		a := roundup(c.s, 32)
		if a != c.w {
			t.Errorf("[%d] roundup(%d) = %d, want %d", i, c.s, a, c.w)
		}
	}
}
