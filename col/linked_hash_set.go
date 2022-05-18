package col

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/pandafw/pango/ars"
)

// NewLinkedHashSet returns an initialized set.
// Example: NewLinkedHashSet(1, 2, 3)
func NewLinkedHashSet(vs ...T) *LinkedHashSet {
	ls := &LinkedHashSet{}
	ls.Add(vs...)
	return ls
}

// LinkedHashSet implements a doubly linked set.
// The zero value for LinkedHashSet is an empty set ready to use.
// Note that insertion order is not affected if an element is re-inserted into the set.
// (An element e is reinserted into a set s if s.Add(e) is invoked when s.Contains(e) would return true immediately prior to the invocation.)
//
// To iterate over a set (where ls is a *LinkedHashSet):
//	it := ls.Iterator()
//	for it.Next() {
//		// do something with it.Value()
//	}
//
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

// Add adds all items of vs and returns the last added item.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) Add(vs ...T) {
	ls.Insert(ls.Len(), vs...)
}

// AddAll adds all items of another collection
// Note: existing item's order will not change.
func (ls *LinkedHashSet) AddAll(ac Collection) {
	ls.InsertAll(ls.Len(), ac)
}

// Delete delete all items with associated value v of vs
func (ls *LinkedHashSet) Delete(vs ...T) {
	if ls.IsEmpty() {
		return
	}

	for _, v := range vs {
		if ln, ok := ls.hash[v]; ok {
			ls.deleteNode(ln)
		}
	}
}

// DeleteAll delete all of this collection's elements that are also contained in the specified collection
func (ls *LinkedHashSet) DeleteAll(ac Collection) {
	if ls.IsEmpty() || ac.IsEmpty() {
		return
	}

	if ls == ac {
		ls.Clear()
		return
	}

	if ic, ok := ac.(Iterable); ok {
		it := ic.Iterator()
		for it.Next() {
			if ln, ok := ls.hash[it.Value()]; ok {
				ls.deleteNode(ln)
			}
		}
		return
	}

	ls.Delete(ac.Values()...)
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

// ContainsAll Test to see if the collection contains all items of another collection
func (ls *LinkedHashSet) ContainsAll(ac Collection) bool {
	if ac.IsEmpty() || ls == ac {
		return true
	}

	if ls.IsEmpty() {
		return false
	}

	if ic, ok := ac.(Iterable); ok {
		it := ic.Iterator()
		for it.Next() {
			if _, ok := ls.hash[it.Value()]; !ok {
				return false
			}
		}
		return true
	}

	return ls.Contains(ac.Values()...)
}

// Retain Retains only the elements in this collection that are contained in the argument array vs.
func (ls *LinkedHashSet) Retain(vs ...T) {
	if ls.IsEmpty() {
		return
	}

	if len(vs) == 0 {
		ls.Clear()
		return
	}

	for ln := ls.head; ln != nil; ln = ln.next {
		if !ars.Contains(vs, ln.value) {
			ls.deleteNode(ln)
		}
	}
}

// RetainAll Retains only the elements in this collection that are contained in the specified collection.
func (ls *LinkedHashSet) RetainAll(ac Collection) {
	if ls.IsEmpty() || ls == ac {
		return
	}

	if ac.IsEmpty() {
		ls.Clear()
		return
	}

	for ln := ls.head; ln != nil; ln = ln.next {
		if !ac.Contains(ln.value) {
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
func (ls *LinkedHashSet) Each(f func(T)) {
	for ln := ls.head; ln != nil; ln = ln.next {
		f(ln.value)
	}
}

// ReverseEach call f for each item in the set with reverse order
func (ls *LinkedHashSet) ReverseEach(f func(T)) {
	for ln := ls.tail; ln != nil; ln = ln.prev {
		f(ln.value)
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

// Insert inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than set's size
// Note: position equal to set's size is valid, i.e. append.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) Insert(index int, vs ...T) {
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

// InsertAll inserts values of another collection ac at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) InsertAll(index int, ac Collection) {
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

	ls.Insert(index, ac.Values()...)
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

// Remove removes the element at the specified position in this set.
func (ls *LinkedHashSet) Remove(index int) {
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

// Push inserts all items of vs at the tail of set al.
func (ls *LinkedHashSet) Push(vs ...T) {
	ls.Insert(ls.Len(), vs...)
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

// PushHead inserts all items of vs at the head of set ls.
func (ls *LinkedHashSet) PushHead(vs ...T) {
	ls.Insert(0, vs...)
}

// PushHeadAll inserts a copy of another collection at the head of set ls.
// The ls and ac may be the same. They must not be nil.
func (ls *LinkedHashSet) PushHeadAll(ac Collection) {
	ls.InsertAll(0, ac)
}

// PushTail inserts all items of vs at the tail of set ls.
func (ls *LinkedHashSet) PushTail(vs ...T) {
	ls.Insert(ls.Len(), vs...)
}

// PushTailAll inserts a copy of another collection at the tail of set ls.
// The ls and ac may be the same. They must not be nil.
func (ls *LinkedHashSet) PushTailAll(ac Collection) {
	ls.InsertAll(ls.Len(), ac)
}

// String print list to string
func (ls *LinkedHashSet) String() string {
	bs, _ := json.Marshal(ls)
	return string(bs)
}

//-----------------------------------------------------------
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
	len := ls.Len()
	if index >= len || index < -len {
		panic(fmt.Sprintf("LinkedHashSet out of bounds: index=%d, len=%d", index, len))
	}

	if index < 0 {
		index += len
	}
	return index
}

func (ls *LinkedHashSet) checkSizeIndex(index int) int {
	len := ls.Len()
	if index > len || index < -len {
		panic(fmt.Sprintf("LinkedHashSet out of bounds: index=%d, len=%d", index, len))
	}

	if index < 0 {
		index += len
	}
	return index
}

//-----------------------------------------------------
// linkedSetNode is a node of a doublly-linked list.
type linkedSetNode struct {
	prev  *linkedSetNode
	next  *linkedSetNode
	value T
}

// String print the list item to string
func (ln *linkedSetNode) String() string {
	return fmt.Sprintf("%v", ln.value)
}

// linkedHashSetIterator a iterator for linkedSetNode
type linkedHashSetIterator struct {
	lset    *LinkedHashSet
	node    *linkedSetNode
	removed bool
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (it *linkedHashSetIterator) Prev() bool {
	if it.lset.IsEmpty() {
		return false
	}

	if it.node == nil {
		it.node = it.lset.tail
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
func (it *linkedHashSetIterator) Next() bool {
	if it.lset.IsEmpty() {
		return false
	}

	if it.node == nil {
		it.node = it.lset.head
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
func (it *linkedHashSetIterator) Value() T {
	if it.node == nil {
		return nil
	}
	return it.node.value
}

// SetValue set the value to the item
func (it *linkedHashSetIterator) SetValue(v T) {
	if it.node == nil {
		return
	}

	if it.removed {
		// unlinked item
		it.node.value = v
		return
	}

	it.lset.setValue(it.node, v)
}

// Remove remove the current element
func (it *linkedHashSetIterator) Remove() {
	if it.node == nil {
		return
	}

	if it.removed {
		panic("LinkedHashSet can't remove a unlinked item")
	}

	it.lset.deleteNode(it.node)
	it.removed = true
}

// Reset resets the iterator to its initial state (one-before-first/one-after-last)
// Call Next()/Prev() to fetch the first/last element if any.
func (it *linkedHashSetIterator) Reset() {
	it.node = nil
	it.removed = false
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

func (ls *LinkedHashSet) addJSONArrayItem(v T) jsonArray {
	ls.Add(v)
	return ls
}

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(ls)
func (ls *LinkedHashSet) MarshalJSON() (res []byte, err error) {
	return jsonMarshalArray(ls)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, ls)
func (ls *LinkedHashSet) UnmarshalJSON(data []byte) error {
	ls.Clear()
	return jsonUnmarshalArray(data, ls)
}
