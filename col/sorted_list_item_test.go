package col

import (
	"testing"
)

func TestSortedListItemString(t *testing.T) {
	cs := []struct {
		e string
		s interface{}
	}{
		{"a", "a"},
		{"1", 1},
	}

	for _, c := range cs {
		i := &SortedListItem{value: c.s}
		a := i.String()
		if a != c.e {
			t.Errorf("SortedListItem(%v).String() = %q, want %q", c.s, a, c.e)
		}
	}
}
