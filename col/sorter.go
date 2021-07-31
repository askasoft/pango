package col

import (
	"github.com/pandafw/pango/cmp"
)

type sortable interface {
	// Len returns the length of the collection.
	Len() int

	// Get returns the value at the specified index in this list
	Get(index int) (interface{}, bool)

	// Swap swaps values of two items at the given index.
	Swap(i, j int)
}

// sorter A sortable collection sorter implements sort.Interface
type sorter struct {
	scol sortable
	less cmp.Less
}

func (ss *sorter) Len() int {
	return ss.scol.Len()
}

func (ss *sorter) Swap(i, j int) {
	ss.scol.Swap(i, j)
}

func (ss *sorter) Less(i, j int) bool {
	vi, _ := ss.scol.Get(i)
	vj, _ := ss.scol.Get(j)
	return ss.less(vi, vj)
}
