package cog

type sortable[T any] interface {
	// Len returns the length of the collection.
	Len() int

	// Get returns the value at the specified index in this list
	Get(index int) T

	// Swap swaps values of two items at the given index.
	Swap(i, j int)
}

// sorter A sortable collection sorter implements sort.Interface
type sorter[T any] struct {
	scol sortable[T]
	less Less[T]
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
