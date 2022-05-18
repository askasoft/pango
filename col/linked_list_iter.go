package col

// linkedListIterator a iterator for linkedListNode
type linkedListIterator struct {
	list    *LinkedList
	node    *linkedListNode
	removed bool
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (it *linkedListIterator) Prev() bool {
	if it.list.IsEmpty() {
		return false
	}

	if it.node == nil {
		it.node = it.list.tail
		it.removed = false
		return true
	}

	if pi := it.node.prev; pi != nil {
		it.node = pi
		it.removed = false
		return true
	}
	return false
}

// Next moves the iterator to the next element and returns true if there was a next element in the collection.
// If Next() returns true, then next element's value can be retrieved by Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (it *linkedListIterator) Next() bool {
	if it.list.IsEmpty() {
		return false
	}

	if it.node == nil {
		it.node = it.list.head
		it.removed = false
		return true
	}

	if ni := it.node.next; ni != nil {
		it.node = ni
		it.removed = false
		return true
	}
	return false
}

// Value returns the current element's value.
func (it *linkedListIterator) Value() T {
	if it.node == nil {
		return nil
	}
	return it.node.value
}

// SetValue set the value to the item
func (it *linkedListIterator) SetValue(v T) {
	if it.node != nil {
		it.node.value = v
	}
}

// Remove remove the current element
func (it *linkedListIterator) Remove() {
	if it.node == nil {
		return
	}

	if it.removed {
		panic("LinkedList can't remove a unlinked item")
	}

	it.list.deleteNode(it.node)
	it.removed = true
}

// Reset resets the iterator to its initial state (one-before-first/one-after-last)
// Call Next()/Prev() to fetch the first/last element if any.
func (it *linkedListIterator) Reset() {
	it.node = nil
	it.removed = false
}
