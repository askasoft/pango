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
func (le *ListItem) Next() *ListItem {
	if ni := le.next; le.list != nil && ni != &le.list.root {
		return ni
	}
	return nil
}

// Prev returns the previous list item or nil.
func (le *ListItem) Prev() *ListItem {
	if pi := le.prev; le.list != nil && pi != &le.list.root {
		return pi
	}
	return nil
}

// Offset returns the next +n or previous -n list item or nil.
func (le *ListItem) Offset(n int) *ListItem {
	if n == 0 {
		return le
	}

	if n > 0 {
		for le != nil && n > 0 {
			le = le.Next()
			n--
		}
		return le
	}

	for le != nil && n < 0 {
		le = le.Prev()
		n++
	}
	return le
}

// String print the list item to string
func (le *ListItem) String() string {
	return fmt.Sprintf("%v", le.Value)
}
