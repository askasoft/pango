package cog

// ringBufferIterator a iterator for the RingBuffer
type ringBufferIterator[T any] struct {
	rb    *RingBuffer[T]
	start int
	index int
}

// Prev moves the iterator to the previous item and returns true if there was a previous item in the container.
// If Prev() returns true, then previous item's value can be retrieved by Value().
// Modifies the state of the iterator.
func (it *ringBufferIterator[T]) Prev() bool {
	if it.rb.IsEmpty() {
		return false
	}

	if it.index < 0 && it.start >= 0 {
		it.index = it.start
	}

	if it.index == 0 {
		return false
	}

	if it.index < 0 {
		it.index = it.rb.len - 1
		return true
	}
	if it.index > it.rb.len {
		return false
	}
	it.index--
	return true
}

// Next moves the iterator to the next item and returns true if there was a next item in the collection.
// If Next() returns true, then next item's value can be retrieved by Value().
// If Next() was called for the first time, then it will point the iterator to the first item if it exists.
// Modifies the state of the iterator.
func (it *ringBufferIterator[T]) Next() bool {
	if it.rb.IsEmpty() {
		return false
	}

	if it.index < 0 && it.start > 0 {
		it.index = it.start - 1
	}
	if it.index < -1 || it.index >= it.rb.len-1 {
		return false
	}
	it.index++
	return true
}

// Value returns the current item's value.
func (it *ringBufferIterator[T]) Value() T {
	if it.index >= 0 && it.index < it.rb.len {
		return it.rb.Get(it.index)
	}

	var v T
	return v
}

// SetValue set the value to the item
func (it *ringBufferIterator[T]) SetValue(v T) {
	if it.index >= 0 && it.index < it.rb.len {
		it.rb.Set(it.index, v)
	}
}

// Remove remove the current item
func (it *ringBufferIterator[T]) Remove() {
	if it.index < 0 {
		return
	}

	it.rb.Remove(it.index)
	it.start = it.index
	it.index = -1
}

// Reset resets the iterator to its initial state (one-before-first/one-after-last)
// Call Next()/Prev() to fetch the first/last item if any.
func (it *ringBufferIterator[T]) Reset() {
	it.start = -1
	it.index = -1
}
