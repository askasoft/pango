package col

import (
	"testing"
)

func TestLinkedMapItemString(t *testing.T) {
	cs := []struct {
		e string
		k interface{}
		v interface{}
	}{
		{"a => b", "a", "b"},
	}

	for i, c := range cs {
		a := (&LinkedMapItem{key: c.k, value: c.v}).String()
		if a != c.e {
			t.Errorf("[%d] LinkedMapItem{%v, %v}.String() = %q, want %q", i, c.k, c.v, a, c.e)
		}
	}
}
