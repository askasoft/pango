package arraylist

type inta []int

// Len is the number of elements in the collection.
func (a inta) Len() int {
	return len(a)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (a inta) Less(i, j int) bool {
	return a[i] < a[j]
}

// Swap swaps the elements with indexes i and j.
func (a inta) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
