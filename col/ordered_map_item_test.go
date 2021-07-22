package col

import (
	"testing"
)

func TestOrderedMapItemString(t *testing.T) {
	cs := []struct {
		e string
		k interface{}
		v interface{}
	}{
		{"a => b", "a", "b"},
	}

	for i, c := range cs {
		a := (&OrderedMapItem{key: c.k, Value: c.v}).String()
		if a != c.e {
			t.Errorf("[%d] OrderedMapItem{%v, %v}.String() = %q, want %q", i, c.k, c.v, a, c.e)
		}
	}
}
