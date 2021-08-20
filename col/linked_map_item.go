package col

import (
	"fmt"
)

// LinkedMapItem key/value item
type LinkedMapItem struct {
	// Next and previous pointers in the doubly-linked list of items.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next item of the last
	// list item (l.Back()) and the previous item of the first list
	// item (l.Front()).
	next, prev *LinkedMapItem

	// The map to which this item belongs.
	lmap *LinkedHashMap

	key   interface{}
	value interface{}
}

// root returns the root item of list.
func (mi *LinkedMapItem) root() *LinkedMapItem {
	if mi.lmap != nil {
		return &mi.lmap.root
	}
	return nil
}

// isRoot returns true if this item is the root item of list.
func (mi *LinkedMapItem) isRoot() bool {
	return mi == mi.root()
}

// Key returns the item's key
func (mi *LinkedMapItem) Key() interface{} {
	return mi.key
}

// Value returns the value stored with this item.
func (mi *LinkedMapItem) Value() interface{} {
	return mi.value
}

// SetValue set the value to the item
func (mi *LinkedMapItem) SetValue(v interface{}) {
	mi.value = v
}

// Next returns the next item or nil.
func (mi *LinkedMapItem) Next() *LinkedMapItem {
	ni := mi.next
	if ni != nil && ni.isRoot() {
		return nil
	}
	return ni
}

// Prev returns the previous item or nil.
func (mi *LinkedMapItem) Prev() *LinkedMapItem {
	pi := mi.prev
	if pi != nil && pi.isRoot() {
		return nil
	}
	return pi
}

// Offset returns the next +n or previous -n list item or nil.
func (mi *LinkedMapItem) Offset(n int) *LinkedMapItem {
	if n == 0 {
		return mi
	}
	if mi.lmap == nil {
		return nil
	}

	if n > 0 {
		ni := mi
		for ; n > 0; n-- {
			ni = ni.next
			if ni.isRoot() {
				return nil
			}
		}
		return ni
	}

	pi := mi
	for ; n < 0; n++ {
		pi = pi.prev
		if pi.isRoot() {
			return nil
		}
	}
	return pi
}

// Remove remove this item from the map
func (mi *LinkedMapItem) Remove() {
	if mi.lmap == nil || mi.isRoot() {
		return
	}

	mi.prev.next = mi.next
	mi.next.prev = mi.prev

	delete(mi.lmap.hash, mi.key)

	// remain prev/next for iterator to Prev()/Next()
	mi.lmap = nil
}

// insertAfter inserts item mi after item at
func (mi *LinkedMapItem) insertAfter(at *LinkedMapItem) {
	ni := at.next
	at.next = mi
	mi.prev = at
	mi.next = ni
	ni.prev = mi
	mi.lmap = at.lmap

	mi.lmap.hash[mi.key] = mi
}

// moveAfter moves the item mi to next to at
func (mi *LinkedMapItem) moveAfter(at *LinkedMapItem) {
	mi.prev.next = mi.next
	mi.next.prev = mi.prev

	n := at.next
	at.next = mi
	mi.prev = at
	mi.next = n
	n.prev = mi
}

// String print the item to string
func (mi *LinkedMapItem) String() string {
	return fmt.Sprintf("%v => %v", mi.key, mi.value)
}

//-----------------------------------------------------

// linkedMapItemIterator a iterator for LinkedMapItem
type linkedMapItemIterator struct {
	lmap *LinkedHashMap
	item *LinkedMapItem
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's key/value can be retrieved by Key()/Value().
// Modifies the state of the iterator.
func (it *linkedMapItemIterator) Prev() bool {
	if pi := it.item.Prev(); pi != nil {
		it.item = pi
		return true
	}
	return false
}

// Next moves the iterator to the next element and returns true if there was a next element in the collection.
// If Next() returns true, then next element's key/value can be retrieved by Key()/Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (it *linkedMapItemIterator) Next() bool {
	if ni := it.item.Next(); ni != nil {
		it.item = ni
		return true
	}
	return false
}

// Key returns the current element's key.
func (it *linkedMapItemIterator) Key() interface{} {
	return it.item.key
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (it *linkedMapItemIterator) Value() interface{} {
	return it.item.value
}

// SetValue set the value to the item
func (it *linkedMapItemIterator) SetValue(v interface{}) {
	it.item.SetValue(v)
}

// Remove remove the current element
func (it *linkedMapItemIterator) Remove() {
	it.item.Remove()
}

// Reset resets the iterator to its initial state (one-before-first/one-after-last)
// Call Next()/Prev() to fetch the first/last element if any.
func (it *linkedMapItemIterator) Reset() {
	it.item = &it.lmap.root
}
