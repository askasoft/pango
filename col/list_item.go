package col

import "fmt"

// ListItem is an item of a linked list.
type ListItem struct {
	// Next and previous pointers in the doubly-linked list of items.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next item of the last
	// list item (l.Back()) and the previous item of the first list
	// item (l.Front()).
	next, prev *ListItem

	// The list to which this item belongs.
	list *List

	// The value stored with this item.
	Value interface{}
}

// Next returns the next list item or nil.
func (li *ListItem) Next() *ListItem {
	if ni := li.next; li.list != nil && ni != &li.list.root {
		return ni
	}
	return nil
}

// Prev returns the previous list item or nil.
func (li *ListItem) Prev() *ListItem {
	if pi := li.prev; li.list != nil && pi != &li.list.root {
		return pi
	}
	return nil
}

// Offset returns the next +n or previous -n list item or nil.
func (li *ListItem) Offset(n int) *ListItem {
	if n == 0 {
		return li
	}

	if n > 0 {
		for li != nil && n > 0 {
			li = li.Next()
			n--
		}
		return li
	}

	for li != nil && n < 0 {
		li = li.Prev()
		n++
	}
	return li
}

// Remove removes the item li from it's onwer list
func (li *ListItem) Remove() {
	if li.list == nil {
		return
	}

	li.list.len--

	li.prev.next = li.next
	li.next.prev = li.prev

	// avoid memory leaks
	li.next = nil
	li.prev = nil
	li.list = nil
}

// String print the list item to string
func (li *ListItem) String() string {
	return fmt.Sprintf("%v", li.Value)
}
