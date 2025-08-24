package bol

import (
	"testing"
)

func TestAtob(t *testing.T) {
	cs := []struct {
		w bool
		s string
		n []bool
	}{
		{false, "F", []bool{}},
		{false, "False", []bool{}},
		{true, "", []bool{false, true}},
	}

	for i, c := range cs {
		a := Atob(c.s, c.n...)
		if a != c.w {
			t.Errorf("[%d] Atob(%s, %v) = %v, want %v", i, c.s, c.n, a, c.w)
		}
	}
}
