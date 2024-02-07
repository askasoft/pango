package col

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/askasoft/pango/bye"
)

// NewLinkedList returns an initialized list.
// Example: col.NewLinkedList(1, 2, 3)
func NewLinkedList(vs ...T) *LinkedList {
	ll := &LinkedList{}
	ll.Adds(vs...)
	return ll
}

// LinkedList implements a doubly linked list.
// The zero value for LinkedList is an empty list ready to use.
//
// To iterate over a list (where ll is a *LinkedList):
//
//	it := ll.Iterator()
//	for it.Next() {
//		// do something with it.Value()
//	}
type LinkedList struct {
	head, tail *linkedListNode
	len        int
}

//-----------------------------------------------------------
// implements Collection interface

// Len returns the length of the list.
func (ll *LinkedList) Len() int {
	return ll.len
}

// IsEmpty returns true if the list length == 0
func (ll *LinkedList) IsEmpty() bool {
	return ll.len == 0
}

// Clear clears list ll.
func (ll *LinkedList) Clear() {
	ll.head = nil
	ll.tail = nil
	ll.len = 0
}

// Add add the item v.
func (ll *LinkedList) Add(v T) {
	ll.Insert(ll.len, v)
}

// Adds adds all items of vs.
func (ll *LinkedList) Adds(vs ...T) {
	ll.Inserts(ll.len, vs...)
}

// AddCol adds all items of another collection
func (ll *LinkedList) AddCol(ac Collection) {
	ll.InsertCol(ll.len, ac)
}

// Remove remove all items with associated value v of vs
func (ll *LinkedList) Remove(v T) {
	for ln := ll.head; ln != nil; ln = ln.next {
		if ln.value == v {
			ll.deleteNode(ln)
		}
	}
}

// Removes remove all items in the array vs
func (ll *LinkedList) Removes(vs ...T) {
	if ll.IsEmpty() {
		return
	}

	for _, v := range vs {
		ll.Remove(v)
	}
}

// RemoveCol remove all of this collection's elements that are also contained in the specified collection
func (ll *LinkedList) RemoveCol(ac Collection) {
	if ll.IsEmpty() || ac.IsEmpty() {
		return
	}

	if ll == ac {
		ll.Clear()
		return
	}

	if ic, ok := ac.(Iterable); ok {
		ll.RemoveIter(ic.Iterator())
		return
	}

	ll.Removes(ac.Values()...)
}

// RemoveIter remove all items in the iterator it
func (ll *LinkedList) RemoveIter(it Iterator) {
	for it.Next() {
		ll.Remove(it.Value())
	}
}

// RemoveFunc remove all items that function f returns true
func (ll *LinkedList) RemoveFunc(f func(T) bool) {
	if ll.IsEmpty() {
		return
	}

	for ln := ll.head; ln != nil; ln = ln.next {
		if f(ln.value) {
			ll.deleteNode(ln)
		}
	}
}

// Contain Test to see if the list contains the value v
func (ll *LinkedList) Contain(v T) bool {
	return ll.Index(v) >= 0
}

// Contains Test to see if the collection contains all items of vs
func (ll *LinkedList) Contains(vs ...T) bool {
	if len(vs) == 0 {
		return true
	}

	if ll.IsEmpty() {
		return false
	}

	for _, v := range vs {
		if ll.Index(v) < 0 {
			return false
		}
	}
	return true
}

// ContainCol Test to see if the collection contains all items of another collection
func (ll *LinkedList) ContainCol(ac Collection) bool {
	if ac.IsEmpty() || ll == ac {
		return true
	}

	if ll.IsEmpty() {
		return false
	}

	if ic, ok := ac.(Iterable); ok {
		return ll.ContainIter(ic.Iterator())
	}

	return ll.Contains(ac.Values()...)
}

// ContainIter Test to see if the collection contains all items of iterator 'it'
func (ll *LinkedList) ContainIter(it Iterator) bool {
	for it.Next() {
		if ll.Index(it.Value()) < 0 {
			return false
		}
	}
	return true
}

// Retains Retains only the elements in this collection that are contained in the argument array vs.
func (ll *LinkedList) Retains(vs ...T) {
	if ll.IsEmpty() {
		return
	}

	if len(vs) == 0 {
		ll.Clear()
		return
	}

	for ln := ll.head; ln != nil; ln = ln.next {
		if !contains(vs, ln.value) {
			ll.deleteNode(ln)
		}
	}
}

// RetainCol Retains only the elements in this collection that are contained in the specified collection.
func (ll *LinkedList) RetainCol(ac Collection) {
	if ll.IsEmpty() || ll == ac {
		return
	}

	if ac.IsEmpty() {
		ll.Clear()
		return
	}

	ll.RetainFunc(ac.Contain)
}

// RetainFunc Retains all items that function f returns true
func (ll *LinkedList) RetainFunc(f func(T) bool) {
	if ll.IsEmpty() {
		return
	}

	for ln := ll.head; ln != nil; ln = ln.next {
		if !f(ln.value) {
			ll.deleteNode(ln)
		}
	}
}

// Values returns a slice contains all the items of the list ll
func (ll *LinkedList) Values() []T {
	vs := make([]T, ll.Len())
	for i, ln := 0, ll.head; ln != nil; i, ln = i+1, ln.next {
		vs[i] = ln.value
	}
	return vs
}

// Each call f for each item in the list
func (ll *LinkedList) Each(f func(T)) {
	for ln := ll.head; ln != nil; ln = ln.next {
		f(ln.value)
	}
}

// ReverseEach call f for each item in the list with reverse order
func (ll *LinkedList) ReverseEach(f func(T)) {
	for ln := ll.tail; ln != nil; ln = ln.prev {
		f(ln.value)
	}
}

// Iterator returns a iterator for the list
func (ll *LinkedList) Iterator() Iterator {
	return &linkedListIterator{list: ll}
}

//-----------------------------------------------------------
// implements List interface

// Get returns the element at the specified position in this list
// if i < -ll.Len() or i >= ll.Len(), panic
// if i < 0, returns ll.Get(ll.Len() + i)
func (ll *LinkedList) Get(index int) T {
	index = ll.checkItemIndex(index)

	return ll.node(index).value
}

// Set set the v at the specified index in this list and returns the old value.
func (ll *LinkedList) Set(index int, v T) (ov T) {
	index = ll.checkItemIndex(index)

	ln := ll.node(index)
	ov, ln.value = ln.value, v
	return
}

// Insert insert value v at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (ll *LinkedList) Insert(index int, v T) {
	index = ll.checkSizeIndex(index)

	var prev, next *linkedListNode
	if index == ll.len {
		next = nil
		prev = ll.tail
	} else {
		next = ll.node(index)
		prev = next.prev
	}

	nn := &linkedListNode{prev: prev, value: v, next: nil}
	if prev == nil {
		ll.head = nn
	} else {
		prev.next = nn
	}
	prev = nn

	if next == nil {
		ll.tail = prev
	} else {
		prev.next = next
		next.prev = prev
	}

	ll.len++
}

// Inserts inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (ll *LinkedList) Inserts(index int, vs ...T) {
	index = ll.checkSizeIndex(index)

	if len(vs) == 0 {
		return
	}

	var prev, next *linkedListNode
	if index == ll.len {
		next = nil
		prev = ll.tail
	} else {
		next = ll.node(index)
		prev = next.prev
	}

	for _, v := range vs {
		nn := &linkedListNode{prev: prev, value: v, next: nil}
		if prev == nil {
			ll.head = nn
		} else {
			prev.next = nn
		}
		prev = nn
	}

	if next == nil {
		ll.tail = prev
	} else {
		prev.next = next
		next.prev = prev
	}

	ll.len += len(vs)
}

// InsertCol inserts values of another collection ac at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (ll *LinkedList) InsertCol(index int, ac Collection) {
	index = ll.checkSizeIndex(index)

	if ac.IsEmpty() {
		return
	}

	ll.Inserts(index, ac.Values()...)
}

// Index returns the index of the first occurrence of the specified v in this list, or -1 if this list does not contain v.
func (ll *LinkedList) Index(v T) int {
	for i, ln := 0, ll.head; ln != nil; ln = ln.next {
		if ln.value == v {
			return i
		}
		i++
	}
	return -1
}

// IndexFunc returns the index of the first true returned by function f in this list, or -1 if this list does not contain v.
func (ll *LinkedList) IndexFunc(f func(T) bool) int {
	for i, ln := 0, ll.head; ln != nil; ln = ln.next {
		if f(ln.value) {
			return i
		}
		i++
	}
	return -1
}

// LastIndex returns the index of the last occurrence of the specified v in this list, or -1 if this list does not contain v.
func (ll *LinkedList) LastIndex(v T) int {
	for i, ln := 0, ll.tail; ln != nil; ln = ln.prev {
		if ln.value == v {
			return i
		}
		i++
	}
	return -1
}

// DeleteAt remove the element at the specified position in this list.
func (ll *LinkedList) DeleteAt(index int) {
	index = ll.checkItemIndex(index)

	ln := ll.node(index)
	ll.deleteNode(ln)
}

// Swap swaps values of two items at the given index.
func (ll *LinkedList) Swap(i, j int) {
	i = ll.checkItemIndex(i)
	j = ll.checkItemIndex(j)

	if i != j {
		ni, nj := ll.node(i), ll.node(j)
		ni.value, nj.value = nj.value, ni.value
	}
}

// Sort Sorts this list according to the order induced by the specified Comparator.
func (ll *LinkedList) Sort(less Less) {
	if ll.Len() < 2 {
		return
	}
	sort.Sort(&sorter{ll, less})
}

// Head get the first item of list.
func (ll *LinkedList) Head() (v T) {
	v, _ = ll.PeekHead()
	return
}

// Tail get the last item of list.
func (ll *LinkedList) Tail() (v T) {
	v, _ = ll.PeekTail()
	return
}

//--------------------------------------------------------------------
// implements Queue interface

// Peek get the first item of list.
func (ll *LinkedList) Peek() (v T, ok bool) {
	return ll.PeekHead()
}

// Poll get and remove the first item of list.
func (ll *LinkedList) Poll() (T, bool) {
	return ll.PollHead()
}

// Push insert the item v at the tail of list al.
func (ll *LinkedList) Push(v T) {
	ll.Insert(ll.Len(), v)
}

// Pushs inserts all items of vs at the tail of list al.
func (ll *LinkedList) Pushs(vs ...T) {
	ll.Inserts(ll.Len(), vs...)
}

//--------------------------------------------------------------------
// implements Deque interface

// PeekHead get the first item of list.
func (ll *LinkedList) PeekHead() (v T, ok bool) {
	if ll.head != nil {
		v, ok = ll.head.value, true
	}
	return
}

// PeekTail get the last item of list.
func (ll *LinkedList) PeekTail() (v T, ok bool) {
	if ll.tail != nil {
		v, ok = ll.tail.value, true
	}
	return
}

// PollHead remove the first item of list.
func (ll *LinkedList) PollHead() (v T, ok bool) {
	v, ok = ll.PeekHead()
	if ok {
		ll.deleteNode(ll.head)
	}
	return
}

// PollTail remove the last item of list.
func (ll *LinkedList) PollTail() (v T, ok bool) {
	v, ok = ll.PeekTail()
	if ok {
		ll.deleteNode(ll.tail)
	}
	return
}

// PushHead insert the item v at the head of list ll.
func (ll *LinkedList) PushHead(v T) {
	ll.Insert(0, v)
}

// PushHeads inserts all items of vs at the head of list ll.
func (ll *LinkedList) PushHeads(vs ...T) {
	ll.Inserts(0, vs...)
}

// PushHeadCol inserts a copy of another collection at the head of list ll.
// The ll and ac may be the same. They must not be nil.
func (ll *LinkedList) PushHeadCol(ac Collection) {
	ll.InsertCol(0, ac)
}

// PushTail insert the item v at the tail of list ll.
func (ll *LinkedList) PushTail(v T) {
	ll.Insert(ll.len, v)
}

// PushTails inserts all items of vs at the tail of list ll.
func (ll *LinkedList) PushTails(vs ...T) {
	ll.Inserts(ll.len, vs...)
}

// PushTailCol inserts a copy of another collection at the tail of list ll.
// The ll and ac may be the same. They must not be nil.
func (ll *LinkedList) PushTailCol(ac Collection) {
	ll.InsertCol(ll.len, ac)
}

// String print list to string
func (ll *LinkedList) String() string {
	bs, _ := json.Marshal(ll)
	return bye.UnsafeString(bs)
}

// -----------------------------------------------------------
func (ll *LinkedList) deleteNode(ln *linkedListNode) {
	if ln.prev == nil {
		ll.head = ln.next
	} else {
		ln.prev.next = ln.next
	}

	if ln.next == nil {
		ll.tail = ln.prev
	} else {
		ln.next.prev = ln.prev
	}

	ll.len--
}

// node returns the node at the specified index i.
func (ll *LinkedList) node(i int) *linkedListNode {
	if i < (ll.len >> 1) {
		ln := ll.head
		for ; i > 0; i-- {
			ln = ln.next
		}
		return ln
	}

	ln := ll.tail
	for i = ll.len - i - 1; i > 0; i-- {
		ln = ln.prev
	}
	return ln
}

func (ll *LinkedList) checkItemIndex(index int) int {
	if index >= ll.len || index < -ll.len {
		panic(fmt.Sprintf("LinkedList out of bounds: index=%d, len=%d", index, ll.len))
	}

	if index < 0 {
		index += ll.len
	}
	return index
}

func (ll *LinkedList) checkSizeIndex(index int) int {
	if index > ll.len || index < -ll.len {
		panic(fmt.Sprintf("LinkedList out of bounds: index=%d, len=%d", index, ll.len))
	}

	if index < 0 {
		index += ll.len
	}
	return index
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

func (ll *LinkedList) addJSONArrayItem(v T) jsonArray {
	ll.Add(v)
	return ll
}

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(ll)
func (ll *LinkedList) MarshalJSON() ([]byte, error) {
	return jsonMarshalArray(ll)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, ll)
func (ll *LinkedList) UnmarshalJSON(data []byte) error {
	ll.Clear()
	return jsonUnmarshalArray(data, ll)
}
