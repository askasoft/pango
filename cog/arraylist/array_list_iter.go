package arraylist

// arrayListIterator a iterator for array list
type arrayListIterator[T any] struct {
	list  *ArrayList[T]
	start int
	index int
}

// Prev moves the iterator to the previous item and returns true if there was a previous item in the container.
// If Prev() returns true, then previous item's value can be retrieved by Value().
// Modifies the state of the iterator.
func (it *arrayListIterator[T]) Prev() bool {
	if it.list.IsEmpty() {
		return false
	}

	if it.index < 0 && it.start >= 0 {
		it.index = it.start
	}

	if it.index == 0 {
		return false
	}

	if it.index < 0 {
		it.index = it.list.Len() - 1
		return true
	}
	if it.index > it.list.Len() {
		return false
	}
	it.index--
	return true
}

// Next moves the iterator to the next item and returns true if there was a next item in the collection.
// If Next() returns true, then next item's value can be retrieved by Value().
// If Next() was called for the first time, then it will point the iterator to the first item if it exists.
// Modifies the state of the iterator.
func (it *arrayListIterator[T]) Next() bool {
	if it.list.IsEmpty() {
		return false
	}

	if it.index < 0 && it.start > 0 {
		it.index = it.start - 1
	}
	if it.index < -1 || it.index >= it.list.Len()-1 {
		return false
	}
	it.index++
	return true
}

// Value returns the current item's value.
func (it *arrayListIterator[T]) Value() (v T) {
	if it.index >= 0 && it.index < it.list.Len() {
		v = it.list.data[it.index]
	}
	return
}

// SetValue set the value to the item
func (it *arrayListIterator[T]) SetValue(v T) {
	if it.index >= 0 && it.index < it.list.Len() {
		it.list.data[it.index] = v
	}
}

// Remove remove the current item
func (it *arrayListIterator[T]) Remove() {
	if it.index < 0 {
		return
	}

	it.list.DeleteAt(it.index)
	it.start = it.index
	it.index = -1
}

// Reset resets the iterator to its initial state (one-before-first/one-after-last)
// Call Next()/Prev() to fetch the first/last item if any.
func (it *arrayListIterator[T]) Reset() {
	it.start = -1
	it.index = -1
}
