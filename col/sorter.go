package col

import (
	"github.com/pandafw/pango/cmp"
)

// listSorter A list sorter implements sort.Interface
type listSorter struct {
	list List
	less cmp.Less
}

func (ls *listSorter) Len() int {
	return ls.list.Len()
}

func (ls *listSorter) Swap(i, j int) {
	ls.list.Swap(i, j)
}

func (ls *listSorter) Less(i, j int) bool {
	vi, _ := ls.list.Get(i)
	vj, _ := ls.list.Get(j)
	return ls.less(vi, vj)
}
