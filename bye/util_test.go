package bye

import "testing"

func TestCountByte(t *testing.T) {
	cs := []struct {
		w int
		s string
		b byte
	}{
		{0, "", 'a'},
		{1, " a ", 'a'},
		{0, " b ", 'a'},
		{2, " a a ", 'a'},
	}

	for i, c := range cs {
		a := CountByte([]byte(c.s), c.b)
		if a != c.w {
			t.Errorf("[%d] CountByte(%q, %v) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}
