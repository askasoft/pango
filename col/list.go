package col

// List implements a doubly linked list.
// The zero value for List is an empty list ready to use.
//
// To iterate over a list (where l is a *List):
//	for e := l.Front(); e != nil; e = e.Next() {
//		// do something with e.Value
//	}
//
type List struct {
	root ListEntry // sentinel list entry, only &root, root.prev, and root.next are used
	len  int       // current list length excluding (this) sentinel entry
}

// Clear clears list l.
func (l *List) Clear() {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
}

// NewList returns an initialized list.
func NewList() *List {
	l := &List{}
	l.Clear()
	return l
}

// Len returns the number of entries of list l.
// The complexity is O(1).
func (l *List) Len() int {
	return l.len
}

// IsEmpty returns true if the list length == 0
func (l *List) IsEmpty() bool {
	return l.len == 0
}

// At returns the entry at the specified index
func (l *List) At(i int) *ListEntry {
	if i < 0 || i >= l.Len() {
		return nil
	}

	if i < l.len/2 {
		return l.Front().NextAt(i)
	}

	return l.Back().PrevAt(i)
}

// Front returns the first entry of list l or nil if the list is empty.
func (l *List) Front() *ListEntry {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last entry of list l or nil if the list is empty.
func (l *List) Back() *ListEntry {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// insert inserts e after at, increments l.len, and returns e.
func (l *List) insert(e, at *ListEntry) *ListEntry {
	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&Entry{Value: v}, at).
func (l *List) insertValue(v interface{}, at *ListEntry) *ListEntry {
	return l.insert(&ListEntry{Value: v}, at)
}

// remove removes e from its list, decrements l.len, and returns e.
func (l *List) remove(e *ListEntry) *ListEntry {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
	return e
}

// move moves e to next to at and returns e.
func (l *List) move(e, at *ListEntry) *ListEntry {
	if e == at {
		return e
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e

	return e
}

// Remove removes e from l if e is an entry of list l.
// It returns the entry value e.Value.
// The entry must not be nil.
func (l *List) Remove(e *ListEntry) interface{} {
	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero Entry) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}

// PushFront inserts a new entry e with value v at the front of list l and returns e.
func (l *List) PushFront(v interface{}) *ListEntry {
	return l.insertValue(v, &l.root)
}

// PushBack inserts a new entry e with value v at the back of list l and returns e.
func (l *List) PushBack(v interface{}) *ListEntry {
	return l.insertValue(v, l.root.prev)
}

// InsertBefore inserts a new entry e with value v immediately before mark and returns e.
// If mark is not an entry of l, the list is not modified.
// The mark must not be nil.
func (l *List) InsertBefore(v interface{}, mark *ListEntry) *ListEntry {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark.prev)
}

// InsertAfter inserts a new entry e with value v immediately after mark and returns e.
// If mark is not an entry of l, the list is not modified.
// The mark must not be nil.
func (l *List) InsertAfter(v interface{}, mark *ListEntry) *ListEntry {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark)
}

// MoveToFront moves entry e to the front of list l.
// If e is not an entry of l, the list is not modified.
// The entry must not be nil.
func (l *List) MoveToFront(e *ListEntry) {
	if e.list != l || l.root.next == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, &l.root)
}

// MoveToBack moves entry e to the back of list l.
// If e is not an entry of l, the list is not modified.
// The entry must not be nil.
func (l *List) MoveToBack(e *ListEntry) {
	if e.list != l || l.root.prev == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, l.root.prev)
}

// MoveBefore moves entry e to its new position before mark.
// If e or mark is not an entry of l, or e == mark, the list is not modified.
// The entry and mark must not be nil.
func (l *List) MoveBefore(e, mark *ListEntry) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark.prev)
}

// MoveAfter moves entry e to its new position after mark.
// If e or mark is not an entry of l, or e == mark, the list is not modified.
// The entry and mark must not be nil.
func (l *List) MoveAfter(e, mark *ListEntry) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark)
}

// PushBackList inserts a copy of an other list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (l *List) PushBackList(other *List) {
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

// PushFrontList inserts a copy of an other list at the front of list l.
// The lists l and other may be the same. They must not be nil.
func (l *List) PushFrontList(other *List) {
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}

// Each Call f for each item in the set
func (l *List) Each(f func(interface{})) {
	for e := l.Front(); e != nil; e = e.Next() {
		f(e.Value)
	}
}

// ReverseEach Call f for each item in the set with reverse order
func (l *List) ReverseEach(f func(interface{})) {
	for e := l.Back(); e != nil; e = e.Prev() {
		f(e.Value)
	}
}
