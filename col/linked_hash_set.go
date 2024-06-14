package col

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/askasoft/pango/str"
)

// NewLinkedHashSet returns an initialized set.
// Example: col.NewLinkedHashSet(1, 2, 3)
func NewLinkedHashSet(vs ...T) *LinkedHashSet {
	ls := &LinkedHashSet{}
	ls.Adds(vs...)
	return ls
}

// LinkedHashSet implements a doubly linked set.
// The zero value for LinkedHashSet is an empty set ready to use.
// Note that insertion order is not affected if an element is re-inserted into the set.
// (An element e is reinserted into a set s if s.Add(e) is invoked when s.Contains(e) would return true immediately prior to the invocation.)
//
// To iterate over a set (where ls is a *LinkedHashSet):
//
//	it := ls.Iterator()
//	for it.Next() {
//		// do something with it.Value()
//	}
type LinkedHashSet struct {
	head, tail *linkedSetNode
	hash       map[T]*linkedSetNode
}

//-----------------------------------------------------------
// implements Collection interface

// Len returns the length of the set.
func (ls *LinkedHashSet) Len() int {
	return len(ls.hash)
}

// IsEmpty returns true if the set length == 0
func (ls *LinkedHashSet) IsEmpty() bool {
	return len(ls.hash) == 0
}

// Clear clears set ls.
func (ls *LinkedHashSet) Clear() {
	ls.hash = nil
	ls.head = nil
	ls.tail = nil
}

// Add add the item v.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) Add(v T) {
	ls.Insert(ls.Len(), v)
}

// Adds adds all items of vs.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) Adds(vs ...T) {
	ls.Inserts(ls.Len(), vs...)
}

// AddCol adds all items of another collection
// Note: existing item's order will not change.
func (ls *LinkedHashSet) AddCol(ac Collection) {
	ls.InsertCol(ls.Len(), ac)
}

// Remove remove all items with associated value v
func (ls *LinkedHashSet) Remove(v T) {
	if ln, ok := ls.hash[v]; ok {
		ls.deleteNode(ln)
	}
}

// Removes remove all items in the array vs
func (ls *LinkedHashSet) Removes(vs ...T) {
	if ls.IsEmpty() {
		return
	}

	for _, v := range vs {
		ls.Remove(v)
	}
}

// RemoveCol remove all of this collection's elements that are also contained in the specified collection
func (ls *LinkedHashSet) RemoveCol(ac Collection) {
	if ls.IsEmpty() || ac.IsEmpty() {
		return
	}

	if ls == ac {
		ls.Clear()
		return
	}

	if ic, ok := ac.(Iterable); ok {
		ls.RemoveIter(ic.Iterator())
		return
	}

	ls.Removes(ac.Values()...)
}

// RemoveIter remove all items in the iterator it
func (ls *LinkedHashSet) RemoveIter(it Iterator) {
	for it.Next() {
		ls.Remove(it.Value())
	}
}

// RemoveFunc remove all items that function f returns true
func (ls *LinkedHashSet) RemoveFunc(f func(T) bool) {
	if ls.IsEmpty() {
		return
	}

	for ln := ls.head; ln != nil; ln = ln.next {
		if f(ln.value) {
			ls.deleteNode(ln)
		}
	}
}

// Contain Test to see if the list contains the value v
func (ls *LinkedHashSet) Contain(v T) bool {
	if ls.IsEmpty() {
		return false
	}
	if _, ok := ls.hash[v]; ok {
		return true
	}
	return false
}

// Contains Test to see if the collection contains all items of vs
func (ls *LinkedHashSet) Contains(vs ...T) bool {
	if len(vs) == 0 {
		return true
	}

	if ls.IsEmpty() {
		return false
	}

	for _, v := range vs {
		if _, ok := ls.hash[v]; !ok {
			return false
		}
	}
	return true
}

// ContainCol Test to see if the collection contains all items of another collection
func (ls *LinkedHashSet) ContainCol(ac Collection) bool {
	if ac.IsEmpty() || ls == ac {
		return true
	}

	if ls.IsEmpty() {
		return false
	}

	if ic, ok := ac.(Iterable); ok {
		return ls.ContainIter(ic.Iterator())
	}

	return ls.Contains(ac.Values()...)
}

// ContainIter Test to see if the collection contains all items of iterator 'it'
func (ls *LinkedHashSet) ContainIter(it Iterator) bool {
	for it.Next() {
		if _, ok := ls.hash[it.Value()]; !ok {
			return false
		}
	}
	return true
}

// Retains Retains only the elements in this collection that are contained in the argument array vs.
func (ls *LinkedHashSet) Retains(vs ...T) {
	if ls.IsEmpty() {
		return
	}

	if len(vs) == 0 {
		ls.Clear()
		return
	}

	for ln := ls.head; ln != nil; ln = ln.next {
		if !contains(vs, ln.value) {
			ls.deleteNode(ln)
		}
	}
}

// RetainCol Retains only the elements in this collection that are contained in the specified collection.
func (ls *LinkedHashSet) RetainCol(ac Collection) {
	if ls.IsEmpty() || ls == ac {
		return
	}

	if ac.IsEmpty() {
		ls.Clear()
		return
	}

	ls.RetainFunc(ac.Contain)
}

// RetainFunc Retains all items that function f returns true
func (ls *LinkedHashSet) RetainFunc(f func(T) bool) {
	if ls.IsEmpty() {
		return
	}

	for ln := ls.head; ln != nil; ln = ln.next {
		if !f(ln.value) {
			ls.deleteNode(ln)
		}
	}
}

// Values returns a slice contains all the items of the set ls
func (ls *LinkedHashSet) Values() []T {
	vs := make([]T, ls.Len())
	for i, ln := 0, ls.head; ln != nil; i, ln = i+1, ln.next {
		vs[i] = ln.value
	}
	return vs
}

// Each call f for each item in the set
func (ls *LinkedHashSet) Each(f func(int, T) bool) {
	i := 0
	for ln := ls.head; ln != nil; ln = ln.next {
		if !f(i, ln.value) {
			return
		}
		i++
	}
}

// ReverseEach call f for each item in the set with reverse order
func (ls *LinkedHashSet) ReverseEach(f func(int, T) bool) {
	i := ls.Len() - 1
	for ln := ls.tail; ln != nil; ln = ln.prev {
		if !f(i, ln.value) {
			return
		}
		i--
	}
}

// Iterator returns a iterator for the set
func (ls *LinkedHashSet) Iterator() Iterator {
	return &linkedHashSetIterator{lset: ls}
}

//-----------------------------------------------------------
// implements List interface (LinkedSet behave like List)

// Get returns the element at the specified position in this set
// if i < -ls.Len() or i >= ls.Len(), panic
// if i < 0, returns ls.Get(ls.Len() + i)
func (ls *LinkedHashSet) Get(index int) T {
	index = ls.checkItemIndex(index)

	return ls.node(index).value
}

// Set set the v at the specified index in this set and returns the old value.
// Old item at index will be removed.
func (ls *LinkedHashSet) Set(index int, v T) (ov T) {
	index = ls.checkItemIndex(index)

	ln := ls.node(index)
	ov = ln.value
	ls.setValue(ln, v)
	return
}

// SetValue set the value to the node
func (ls *LinkedHashSet) setValue(ln *linkedSetNode, v T) {
	if ln.value == v {
		return
	}

	// delete old item
	delete(ls.hash, ln.value)

	// delete duplicated item
	if dn, ok := ls.hash[v]; ok {
		ls.deleteNode(dn)
	}

	// add new item
	ln.value = v
	ls.hash[v] = ln
}

// Insert insert value v at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than set's size
// Note: position equal to set's size is valid, i.e. append.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) Insert(index int, v T) {
	index = ls.checkSizeIndex(index)

	if ls.hash == nil {
		ls.hash = make(map[T]*linkedSetNode)
	}

	if _, ok := ls.hash[v]; ok {
		return
	}

	var prev, next *linkedSetNode
	if index == ls.Len() {
		next = nil
		prev = ls.tail
	} else {
		next = ls.node(index)
		prev = next.prev
	}

	nn := &linkedSetNode{prev: prev, value: v, next: nil}
	if prev == nil {
		ls.head = nn
	} else {
		prev.next = nn
	}
	prev = nn
	ls.hash[v] = nn

	if next == nil {
		ls.tail = prev
	} else {
		prev.next = next
		next.prev = prev
	}
}

// Insert inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than set's size
// Note: position equal to set's size is valid, i.e. append.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) Inserts(index int, vs ...T) {
	index = ls.checkSizeIndex(index)

	n := len(vs)
	if n == 0 {
		return
	}

	if ls.hash == nil {
		ls.hash = make(map[T]*linkedSetNode)
	}

	var prev, next *linkedSetNode
	if index == ls.Len() {
		next = nil
		prev = ls.tail
	} else {
		next = ls.node(index)
		prev = next.prev
	}

	for _, v := range vs {
		if _, ok := ls.hash[v]; ok {
			continue
		}

		nn := &linkedSetNode{prev: prev, value: v, next: nil}
		if prev == nil {
			ls.head = nn
		} else {
			prev.next = nn
		}
		prev = nn
		ls.hash[v] = nn
	}

	if next == nil {
		ls.tail = prev
	} else if prev != nil {
		prev.next = next
		next.prev = prev
	}
}

// InsertCol inserts values of another collection ac at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) InsertCol(index int, ac Collection) {
	index = ls.checkSizeIndex(index)

	if ac.IsEmpty() || ls == ac {
		return
	}

	if ls.hash == nil {
		ls.hash = make(map[T]*linkedSetNode)
	}

	if ic, ok := ac.(Iterable); ok {
		var prev, next *linkedSetNode
		if index == ls.Len() {
			next = nil
			prev = ls.tail
		} else {
			next = ls.node(index)
			prev = next.prev
		}

		it := ic.Iterator()
		for it.Next() {
			v := it.Value()
			if _, ok := ls.hash[v]; ok {
				continue
			}

			nn := &linkedSetNode{prev: prev, value: v, next: nil}
			if prev == nil {
				ls.head = nn
			} else {
				prev.next = nn
			}
			prev = nn
			ls.hash[v] = nn
		}

		if next == nil {
			ls.tail = prev
		} else if prev != nil {
			prev.next = next
			next.prev = prev
		}
		return
	}

	ls.Inserts(index, ac.Values()...)
}

// Index returns the index of the specified v in this set, or -1 if this set does not contain v.
func (ls *LinkedHashSet) Index(v T) int {
	for i, ln := 0, ls.head; ln != nil; ln = ln.next {
		if ln.value == v {
			return i
		}
		i++
	}
	return -1
}

// DeleteAt removes the element at the specified position in this set.
func (ls *LinkedHashSet) DeleteAt(index int) {
	index = ls.checkItemIndex(index)

	ln := ls.node(index)
	ls.deleteNode(ln)
}

// Swap swaps values of two items at the given index.
func (ls *LinkedHashSet) Swap(i, j int) {
	i = ls.checkItemIndex(i)
	j = ls.checkItemIndex(j)

	if i != j {
		ni, nj := ls.node(i), ls.node(j)
		ni.value, nj.value = nj.value, ni.value
	}
}

// Sort Sorts this set according to the order induced by the specified Comparator.
func (ls *LinkedHashSet) Sort(less Less) {
	if ls.Len() < 2 {
		return
	}
	sort.Sort(&sorter{ls, less})
}

// Head get the first item of set.
func (ls *LinkedHashSet) Head() (v T) {
	v, _ = ls.PeekHead()
	return
}

// Tail get the last item of set.
func (ls *LinkedHashSet) Tail() (v T) {
	v, _ = ls.PeekTail()
	return
}

//--------------------------------------------------------------------
// implements Queue interface

// Peek get the first item of set.
func (ls *LinkedHashSet) Peek() (v T, ok bool) {
	return ls.PeekHead()
}

// Poll get and remove the first item of set.
func (ls *LinkedHashSet) Poll() (T, bool) {
	return ls.PollHead()
}

// Push insert the item v at the tail of set al.
func (ls *LinkedHashSet) Push(v T) {
	ls.Insert(ls.Len(), v)
}

// Pushs inserts all items of vs at the tail of set al.
func (ls *LinkedHashSet) Pushs(vs ...T) {
	ls.Inserts(ls.Len(), vs...)
}

//--------------------------------------------------------------------
// implements Deque interface

// PeekHead get the first item of set.
func (ls *LinkedHashSet) PeekHead() (v T, ok bool) {
	if ls.head != nil {
		v, ok = ls.head.value, true
	}
	return
}

// PeekTail get the last item of set.
func (ls *LinkedHashSet) PeekTail() (v T, ok bool) {
	if ls.tail != nil {
		v, ok = ls.tail.value, true
	}
	return
}

// PollHead remove the first item of set.
func (ls *LinkedHashSet) PollHead() (v T, ok bool) {
	v, ok = ls.PeekHead()
	if ok {
		ls.deleteNode(ls.head)
	}
	return
}

// PollTail remove the last item of set.
func (ls *LinkedHashSet) PollTail() (v T, ok bool) {
	v, ok = ls.PeekTail()
	if ok {
		ls.deleteNode(ls.tail)
	}
	return
}

// PushHead insert the item v at the head of set ls.
func (ls *LinkedHashSet) PushHead(v T) {
	ls.Insert(0, v)
}

// PushHeads inserts all items of vs at the head of set ls.
func (ls *LinkedHashSet) PushHeads(vs ...T) {
	ls.Inserts(0, vs...)
}

// PushHeadCol inserts a copy of another collection at the head of set ls.
// The ls and ac may be the same. They must not be nil.
func (ls *LinkedHashSet) PushHeadCol(ac Collection) {
	ls.InsertCol(0, ac)
}

// PushTail insert the item v at the tail of set ls.
func (ls *LinkedHashSet) PushTail(v T) {
	ls.Insert(ls.Len(), v)
}

// PushTails inserts all items of vs at the tail of set ls.
func (ls *LinkedHashSet) PushTails(vs ...T) {
	ls.Inserts(ls.Len(), vs...)
}

// PushTailCol inserts a copy of another collection at the tail of set ls.
// The ls and ac may be the same. They must not be nil.
func (ls *LinkedHashSet) PushTailCol(ac Collection) {
	ls.InsertCol(ls.Len(), ac)
}

// String print list to string
func (ls *LinkedHashSet) String() string {
	bs, _ := json.Marshal(ls)
	return str.UnsafeString(bs)
}

// -----------------------------------------------------------
func (ls *LinkedHashSet) deleteNode(ln *linkedSetNode) {
	if ln.prev == nil {
		ls.head = ln.next
	} else {
		ln.prev.next = ln.next
	}

	if ln.next == nil {
		ls.tail = ln.prev
	} else {
		ln.next.prev = ln.prev
	}

	delete(ls.hash, ln.value)
}

// node returns the node at the specified index i.
func (ls *LinkedHashSet) node(i int) *linkedSetNode {
	if i < (ls.Len() >> 1) {
		ln := ls.head
		for ; i > 0; i-- {
			ln = ln.next
		}
		return ln
	}

	ln := ls.tail
	for i = ls.Len() - i - 1; i > 0; i-- {
		ln = ln.prev
	}
	return ln
}

func (ls *LinkedHashSet) checkItemIndex(index int) int {
	sz := ls.Len()
	if index >= sz || index < -sz {
		panic(fmt.Sprintf("LinkedHashSet out of bounds: index=%d, len=%d", index, sz))
	}

	if index < 0 {
		index += sz
	}
	return index
}

func (ls *LinkedHashSet) checkSizeIndex(index int) int {
	sz := ls.Len()
	if index > sz || index < -sz {
		panic(fmt.Sprintf("LinkedHashSet out of bounds: index=%d, len=%d", index, sz))
	}

	if index < 0 {
		index += sz
	}
	return index
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

func (ls *LinkedHashSet) addJSONArrayItem(v T) jsonArray {
	ls.Add(v)
	return ls
}

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(ls)
func (ls *LinkedHashSet) MarshalJSON() ([]byte, error) {
	return jsonMarshalArray(ls)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, ls)
func (ls *LinkedHashSet) UnmarshalJSON(data []byte) error {
	ls.Clear()
	return jsonUnmarshalArray(data, ls)
}
