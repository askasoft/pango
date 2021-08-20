package col

import (
	"fmt"
)

// LinkedSetItem is an item of a linked set.
type LinkedSetItem struct {
	// Next and previous pointers in the doubly-linked set of items.
	next, prev *LinkedSetItem

	// The set to which this item belongs.
	lset *LinkedHashSet

	// The value stored with this item.
	value interface{}
}

// root returns the root item of list.
func (li *LinkedSetItem) root() *LinkedSetItem {
	if li.lset != nil {
		return &li.lset.root
	}
	return nil
}

// isRoot returns true if this item is the root item of list.
func (li *LinkedSetItem) isRoot() bool {
	return li == li.root()
}

// Value returns the value stored with this item.
func (li *LinkedSetItem) Value() interface{} {
	return li.value
}

// SetValue set the value to the item
func (li *LinkedSetItem) SetValue(v interface{}) {
	if li.value == v {
		return
	}

	// delete old item
	delete(li.lset.hash, li.value)

	// delete duplicated item
	if di, ok := li.lset.hash[v]; ok {
		di.Remove()
	}

	// add new item
	li.value = v
	li.lset.hash[v] = li
}

// Next returns the next list item or nil.
func (li *LinkedSetItem) Next() *LinkedSetItem {
	ni := li.next
	if ni != nil && ni.isRoot() {
		return nil
	}
	return ni
}

// Prev returns the previous list item or nil.
func (li *LinkedSetItem) Prev() *LinkedSetItem {
	pi := li.prev
	if pi != nil && pi.isRoot() {
		return nil
	}
	return pi
}

// Offset returns the next +n or previous -n list item or nil if n is out of range.
func (li *LinkedSetItem) Offset(n int) *LinkedSetItem {
	if n == 0 {
		return li
	}
	if li.lset == nil {
		return nil
	}

	if n > 0 {
		ni := li
		for ; n > 0; n-- {
			ni = ni.next
			if ni.isRoot() {
				return nil
			}
		}
		return ni
	}

	pi := li
	for ; n < 0; n++ {
		pi = pi.prev
		if pi.isRoot() {
			return nil
		}
	}
	return pi
}

// Remove removes the item li from it's onwer list
func (li *LinkedSetItem) Remove() {
	if li.lset == nil || li.isRoot() {
		return
	}

	li.prev.next = li.next
	li.next.prev = li.prev

	delete(li.lset.hash, li.value)

	// remain prev/next for iterator to Prev()/Next()
	li.lset = nil
}

// insertAfter inserts item li after item at
func (li *LinkedSetItem) insertAfter(at *LinkedSetItem) {
	ni := at.next
	at.next = li
	li.prev = at
	li.next = ni
	ni.prev = li
	li.lset = at.lset

	li.lset.hash[li.value] = li
}

// moveAfter moves the item li to next to at
func (li *LinkedSetItem) moveAfter(at *LinkedSetItem) {
	li.prev.next = li.next
	li.next.prev = li.prev

	n := at.next
	at.next = li
	li.prev = at
	li.next = n
	n.prev = li
}

// String print the list item to string
func (li *LinkedSetItem) String() string {
	return fmt.Sprintf("%v", li.value)
}

//-----------------------------------------------------

// linkedSetItemIterator a iterator for LinkedSetItem
type linkedSetItemIterator struct {
	lset *LinkedHashSet
	item *LinkedSetItem
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (it *linkedSetItemIterator) Prev() bool {
	if pi := it.item.Prev(); pi != nil {
		it.item = pi
		return true
	}
	return false
}

// Next moves the iterator to the next element and returns true if there was a next element in the collection.
// If Next() returns true, then next element's value can be retrieved by Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (it *linkedSetItemIterator) Next() bool {
	if ni := it.item.Next(); ni != nil {
		it.item = ni
		return true
	}
	return false
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (it *linkedSetItemIterator) Value() interface{} {
	return it.item.value
}

// SetValue set the value to the item
func (it *linkedSetItemIterator) SetValue(v interface{}) {
	it.item.SetValue(v)
}

// Remove remove the current element
func (it *linkedSetItemIterator) Remove() {
	it.item.Remove()
}

// Reset resets the iterator to its initial state (one-before-first/one-after-last)
// Call Next()/Prev() to fetch the first/last element if any.
func (it *linkedSetItemIterator) Reset() {
	it.item = &it.lset.root
}
