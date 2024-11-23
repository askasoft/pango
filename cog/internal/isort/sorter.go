package isort

import (
	"sort"

	"github.com/askasoft/pango/cog"
)

// sorter A sortable collection sorter implements sort.Interface
type sorter[T any] struct {
	scol cog.SortIF[T]
	less cog.Less[T]
}

func Sort[T any](scol cog.SortIF[T], less cog.Less[T]) {
	if scol.Len() < 2 {
		return
	}

	sort.Sort(&sorter[T]{scol, less})
}

func (ss *sorter[T]) Len() int {
	return ss.scol.Len()
}

func (ss *sorter[T]) Swap(i, j int) {
	ss.scol.Swap(i, j)
}

func (ss *sorter[T]) Less(i, j int) bool {
	vi := ss.scol.Get(i)
	vj := ss.scol.Get(j)
	return ss.less(vi, vj)
}
