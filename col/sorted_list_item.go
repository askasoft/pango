package col

import (
	"fmt"
)

// SortedListItem is an item of a sorted list.
type SortedListItem struct {
	// Next and previous pointers in the doubly-linked list of items.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next item of the last
	// list item (l.Back()) and the previous item of the first list
	// item (l.Front()).
	next, prev *SortedListItem

	// The list to which this item belongs.
	list *SortedList

	// The value stored with this item.
	value interface{}
}

// Value returns the value stored with this item.
func (li *SortedListItem) Value() interface{} {
	return li.value
}

// SetValue set the value to the item
func (li *SortedListItem) SetValue(v interface{}) {
	if li.isRoot() {
		return
	}

	if li.value == v {
		return
	}

	li.value = v
	if li.list != nil {
		li.list.onValueChanged(li)
	}
}

// root returns the root item of list.
func (li *SortedListItem) root() *SortedListItem {
	if li.list != nil {
		return &li.list.root
	}
	return nil
}

// isRoot returns true if this item is the root item of list.
func (li *SortedListItem) isRoot() bool {
	return li == li.root()
}

// Next returns the next list item or nil.
func (li *SortedListItem) Next() *SortedListItem {
	ni := li.next
	if ni != nil && ni.isRoot() {
		return nil
	}
	return ni
}

// Prev returns the previous list item or nil.
func (li *SortedListItem) Prev() *SortedListItem {
	pi := li.prev
	if pi != nil && pi.isRoot() {
		return nil
	}
	return pi
}

// Offset returns the next +n or previous -n list item or nil if n is out of range.
func (li *SortedListItem) Offset(n int) *SortedListItem {
	if n == 0 {
		return li
	}
	if li.list == nil {
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
func (li *SortedListItem) Remove() {
	if li.list == nil || li.isRoot() {
		return
	}

	li.prev.next = li.next
	li.next.prev = li.prev

	li.list.onItemRemoved(li)

	li.list = nil

	// remain prev/next for iterator to Prev()/Next()
}

// insertAfter inserts item li after item at
func (li *SortedListItem) insertAfter(at *SortedListItem) {
	ni := at.next
	at.next = li
	li.prev = at
	li.next = ni
	ni.prev = li
	li.list = at.list

	li.list.onItemInserted(li)
}

// moveAfter moves the item li to next to at
func (li *SortedListItem) moveAfter(at *SortedListItem) {
	li.prev.next = li.next
	li.next.prev = li.prev

	n := at.next
	at.next = li
	li.prev = at
	li.next = n
	n.prev = li
}

// String print the list item to string
func (li *SortedListItem) String() string {
	return fmt.Sprintf("%v", li.value)
}

//-----------------------------------------------------

// sortedListItemIterator a iterator for SortedListItem
type sortedListItemIterator struct {
	list *SortedList
	item *SortedListItem
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (it *sortedListItemIterator) Prev() bool {
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
func (it *sortedListItemIterator) Next() bool {
	if ni := it.item.Next(); ni != nil {
		it.item = ni
		return true
	}
	return false
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (it *sortedListItemIterator) Value() interface{} {
	return it.item.value
}

// SetValue set the value to the item
func (it *sortedListItemIterator) SetValue(v interface{}) {
	it.item.SetValue(v)
}

// Remove remove the current element
func (it *sortedListItemIterator) Remove() {
	it.item.Remove()
}

// Reset resets the iterator to its initial state (one-before-first/one-after-last)
// Call Next()/Prev() to fetch the first/last element if any.
func (it *sortedListItemIterator) Reset() {
	it.item = &it.list.root
}
