package treemap

//-----------------------------------------------------

// treeMapIterator a iterator for TreeMap
type treeMapIterator[K any, V any] struct {
	tree    *TreeMap[K, V]
	node    *TreeMapNode[K, V]
	removed bool
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (it *treeMapIterator[K, V]) Prev() bool {
	if it.tree.IsEmpty() {
		return false
	}

	if it.node == nil {
		it.node = it.tree.tail()
		it.removed = false
		return true
	}

	if it.removed {
		if it.node.left == nil {
			return false
		}
		it.node = it.node.left
		it.removed = false
		return true
	}

	if pi := it.node.prev(); pi != nil {
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
func (it *treeMapIterator[K, V]) Next() bool {
	if it.tree.IsEmpty() {
		return false
	}

	if it.node == nil {
		it.node = it.tree.head()
		it.removed = false
		return true
	}

	if it.removed {
		if it.node.right == nil {
			return false
		}
		it.node = it.node.right
		it.removed = false
		return true
	}

	if ni := it.node.next(); ni != nil {
		it.node = ni
		it.removed = false
		return true
	}
	return false
}

// Key returns the current element's key.
func (it *treeMapIterator[K, V]) Key() K {
	if it.node == nil {
		var k K
		return k
	}
	return it.node.key
}

// Value returns the current element's value.
func (it *treeMapIterator[K, V]) Value() V {
	if it.node == nil {
		var v V
		return v
	}
	return it.node.value
}

// SetValue set the value to the item
func (it *treeMapIterator[K, V]) SetValue(v V) {
	if it.node != nil {
		it.node.value = v
	}
}

// Remove remove the current element
func (it *treeMapIterator[K, V]) Remove() {
	if it.node == nil {
		return
	}

	if it.removed {
		panic("TreeMap can't remove a unlinked item")
	}

	p, n := it.node.prev(), it.node.next()
	d := it.tree.deleteNode(it.node)
	if d != it.node {
		n, it.node = it.node, d
	}

	// save prev/next for iterator
	it.node.left, it.node.right = p, n
	it.removed = true
}

// Reset resets the iterator to its initial state (one-before-first/one-after-last)
// Call Next()/Prev() to fetch the first/last element if any.
func (it *treeMapIterator[K, V]) Reset() {
	it.node = nil
	it.removed = false
}
