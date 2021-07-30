package col

// ArrayListIterator a iterator for array list
type ArrayListIterator struct {
	list  *ArrayList
	start int
	index int
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's value can be retrieved by Value().
// Modifies the state of the iterator.
func (it *ArrayListIterator) Prev() bool {
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

// Next moves the iterator to the next element and returns true if there was a next element in the collection.
// If Next() returns true, then next element's value can be retrieved by Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (it *ArrayListIterator) Next() bool {
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

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (it *ArrayListIterator) Value() interface{} {
	if it.index >= 0 && it.index < it.list.Len() {
		return it.list.data[it.index]
	}
	return nil
}

// SetValue set the value to the item
func (it *ArrayListIterator) SetValue(v interface{}) {
	if it.index >= 0 && it.index < it.list.Len() {
		it.list.data[it.index] = v
	}
}

// Remove remove the current element
func (it *ArrayListIterator) Remove() {
	if it.index < 0 {
		return
	}

	it.list.Remove(it.index)
	it.start = it.index
	it.index = -1
}

// Reset resets the iterator to its initial state (one-before-first/one-after-last)
// Call Next()/Prev() to fetch the first/last element if any.
func (it *ArrayListIterator) Reset() {
	it.start = -1
	it.index = -1
}
