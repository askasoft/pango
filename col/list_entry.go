package col

// ListEntry is an entry of a linked list.
type ListEntry struct {
	// Next and previous pointers in the doubly-linked list of entries.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next entry of the last
	// list entry (l.Back()) and the previous entry of the first list
	// entry (l.Front()).
	next, prev *ListEntry

	// The list to which this entry belongs.
	list *List

	// The value stored with this entry.
	Value interface{}
}

// Next returns the next list entry or nil.
func (e *ListEntry) Next() *ListEntry {
	if n := e.next; e.list != nil && n != &e.list.root {
		return n
	}
	return nil
}

// Move returns the next +n or previous -n list entry or nil.
func (e *ListEntry) Move(n int) *ListEntry {
	if n == 0 {
		return e
	}

	if n > 0 {
		for e != nil && n > 0 {
			e = e.Next()
			n--
		}
		return e
	}

	for e != nil && n < 0 {
		e = e.Prev()
		n++
	}
	return e
}

// Prev returns the previous list entry or nil.
func (e *ListEntry) Prev() *ListEntry {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}
