//go:build go1.18
// +build go1.18

package cog

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/pandafw/pango/ars"
)

// NewLinkedList returns an initialized list.
// Example: NewLinkedList(1, 2, 3)
func NewLinkedList[T any](vs ...T) *LinkedList[T] {
	ll := &LinkedList[T]{}
	ll.Add(vs...)
	return ll
}

// LinkedList implements a doubly linked list.
// The zero value for LinkedList is an empty list ready to use.
//
// To iterate over a list (where ll is a *LinkedList):
//	it := ll.Iterator()
//	for it.Next() {
//		// do something with it.Value()
//	}
//
type LinkedList[T any] struct {
	head, tail *linkedListNode[T]
	len        int
}

//-----------------------------------------------------------
// implements Collection interface

// Len returns the length of the list.
func (ll *LinkedList[T]) Len() int {
	return ll.len
}

// IsEmpty returns true if the list length == 0
func (ll *LinkedList[T]) IsEmpty() bool {
	return ll.len == 0
}

// Clear clears list ll.
func (ll *LinkedList[T]) Clear() {
	ll.head = nil
	ll.tail = nil
	ll.len = 0
}

// Add adds all items of vs and returns the last added item.
func (ll *LinkedList[T]) Add(vs ...T) {
	ll.Insert(ll.len, vs...)
}

// AddAll adds all items of another collection
func (ll *LinkedList[T]) AddAll(ac Collection[T]) {
	ll.InsertAll(ll.len, ac)
}

// Delete delete all items with associated value v of vs
func (ll *LinkedList[T]) Delete(vs ...T) {
	for _, v := range vs {
		ll.deleteAll(v)
	}
}

// DeleteAll delete all of this collection's elements that are also contained in the specified collection
func (ll *LinkedList[T]) DeleteAll(ac Collection[T]) {
	if ll.IsEmpty() || ac.IsEmpty() {
		return
	}

	if ll == ac {
		ll.Clear()
		return
	}

	if ic, ok := ac.(Iterable[T]); ok {
		it := ic.Iterator()
		for it.Next() {
			ll.deleteAll(it.Value())
		}
		return
	}

	ll.Delete(ac.Values()...)
}

// Contains Test to see if the collection contains all items of vs
func (ll *LinkedList[T]) Contains(vs ...T) bool {
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

// ContainsAll Test to see if the collection contains all items of another collection
func (ll *LinkedList[T]) ContainsAll(ac Collection[T]) bool {
	if ac.IsEmpty() || ll == ac {
		return true
	}

	if ll.IsEmpty() {
		return false
	}

	if ic, ok := ac.(Iterable[T]); ok {
		it := ic.Iterator()
		for it.Next() {
			if ll.Index(it.Value()) < 0 {
				return false
			}
		}
		return true
	}

	return ll.Contains(ac.Values()...)
}

// Retain Retains only the elements in this collection that are contained in the argument array vs.
func (ll *LinkedList[T]) Retain(vs ...T) {
	if ll.IsEmpty() {
		return
	}

	if len(vs) == 0 {
		ll.Clear()
		return
	}

	for ln := ll.head; ln != nil; ln = ln.next {
		if !ars.ContainsOf(vs, ln.value) {
			ll.deleteNode(ln)
		}
	}
}

// RetainAll Retains only the elements in this collection that are contained in the specified collection.
func (ll *LinkedList[T]) RetainAll(ac Collection[T]) {
	if ll.IsEmpty() || ll == ac {
		return
	}

	if ac.IsEmpty() {
		ll.Clear()
		return
	}

	for ln := ll.head; ln != nil; ln = ln.next {
		if !ac.Contains(ln.value) {
			ll.deleteNode(ln)
		}
	}
}

// Values returns a slice contains all the items of the list ll
func (ll *LinkedList[T]) Values() []T {
	vs := make([]T, ll.Len())
	for i, ln := 0, ll.head; ln != nil; i, ln = i+1, ln.next {
		vs[i] = ln.value
	}
	return vs
}

// Each call f for each item in the list
func (ll *LinkedList[T]) Each(f func(T)) {
	for ln := ll.head; ln != nil; ln = ln.next {
		f(ln.value)
	}
}

// ReverseEach call f for each item in the list with reverse order
func (ll *LinkedList[T]) ReverseEach(f func(T)) {
	for ln := ll.tail; ln != nil; ln = ln.prev {
		f(ln.value)
	}
}

// Iterator returns a iterator for the list
func (ll *LinkedList[T]) Iterator() Iterator[T] {
	return &linkedListIterator[T]{list: ll}
}

//-----------------------------------------------------------
// implements List interface

// Get returns the element at the specified position in this list
// if i < -ll.Len() or i >= ll.Len(), panic
// if i < 0, returns ll.Get(ll.Len() + i)
func (ll *LinkedList[T]) Get(index int) T {
	index = ll.checkItemIndex(index)

	return ll.node(index).value
}

// Set set the v at the specified index in this list and returns the old value.
func (ll *LinkedList[T]) Set(index int, v T) (ov T) {
	index = ll.checkItemIndex(index)

	ln := ll.node(index)
	ov, ln.value = ln.value, v
	return
}

// Insert inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (ll *LinkedList[T]) Insert(index int, vs ...T) {
	index = ll.checkSizeIndex(index)

	if len(vs) == 0 {
		return
	}

	var prev, next *linkedListNode[T]
	if index == ll.len {
		next = nil
		prev = ll.tail
	} else {
		next = ll.node(index)
		prev = next.prev
	}

	for _, v := range vs {
		nn := &linkedListNode[T]{prev: prev, value: v, next: nil}
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

// InsertAll inserts values of another collection ac at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (ll *LinkedList[T]) InsertAll(index int, ac Collection[T]) {
	index = ll.checkSizeIndex(index)

	if ac.IsEmpty() {
		return
	}

	ll.Insert(index, ac.Values()...)
}

// Index returns the index of the first occurrence of the specified v in this list, or -1 if this list does not contain v.
func (ll *LinkedList[T]) Index(v T) int {
	for i, ln := 0, ll.head; ln != nil; ln = ln.next {
		if any(ln.value) == any(v) {
			return i
		}
		i++
	}
	return -1
}

// LastIndex returns the index of the last occurrence of the specified v in this list, or -1 if this list does not contain v.
func (ll *LinkedList[T]) LastIndex(v T) int {
	for i, ln := 0, ll.tail; ln != nil; ln = ln.prev {
		if any(ln.value) == any(v) {
			return i
		}
		i++
	}
	return -1
}

// Remove removes the element at the specified position in this list.
func (ll *LinkedList[T]) Remove(index int) {
	index = ll.checkItemIndex(index)

	ln := ll.node(index)
	ll.deleteNode(ln)
}

// Swap swaps values of two items at the given index.
func (ll *LinkedList[T]) Swap(i, j int) {
	i = ll.checkItemIndex(i)
	j = ll.checkItemIndex(j)

	if i != j {
		ni, nj := ll.node(i), ll.node(j)
		ni.value, nj.value = nj.value, ni.value
	}
}

// Sort Sorts this list according to the order induced by the specified Comparator.
func (ll *LinkedList[T]) Sort(less Less[T]) {
	if ll.Len() < 2 {
		return
	}
	sort.Sort(&sorter[T]{ll, less})
}

// Head get the first item of list.
func (ll *LinkedList[T]) Head() (v T) {
	v, _ = ll.PeekHead()
	return
}

// Tail get the last item of list.
func (ll *LinkedList[T]) Tail() (v T) {
	v, _ = ll.PeekTail()
	return
}

//--------------------------------------------------------------------
// implements Queue interface

// Peek get the first item of list.
func (ll *LinkedList[T]) Peek() (v T, ok bool) {
	return ll.PeekHead()
}

// Poll get and remove the first item of list.
func (ll *LinkedList[T]) Poll() (T, bool) {
	return ll.PollHead()
}

// Push inserts all items of vs at the tail of list al.
func (ll *LinkedList[T]) Push(vs ...T) {
	ll.Insert(ll.Len(), vs...)
}

//--------------------------------------------------------------------
// implements Deque interface

// PeekHead get the first item of list.
func (ll *LinkedList[T]) PeekHead() (v T, ok bool) {
	if ll.head != nil {
		v, ok = ll.head.value, true
	}
	return
}

// PeekTail get the last item of list.
func (ll *LinkedList[T]) PeekTail() (v T, ok bool) {
	if ll.tail != nil {
		v, ok = ll.tail.value, true
	}
	return
}

// PollHead remove the first item of list.
func (ll *LinkedList[T]) PollHead() (v T, ok bool) {
	v, ok = ll.PeekHead()
	if ok {
		ll.deleteNode(ll.head)
	}
	return
}

// PollTail remove the last item of list.
func (ll *LinkedList[T]) PollTail() (v T, ok bool) {
	v, ok = ll.PeekTail()
	if ok {
		ll.deleteNode(ll.tail)
	}
	return
}

// PushHead inserts all items of vs at the head of list ll.
func (ll *LinkedList[T]) PushHead(vs ...T) {
	ll.Insert(0, vs...)
}

// PushHeadAll inserts a copy of another collection at the head of list ll.
// The ll and ac may be the same. They must not be nil.
func (ll *LinkedList[T]) PushHeadAll(ac Collection[T]) {
	ll.InsertAll(0, ac)
}

// PushTail inserts all items of vs at the tail of list ll.
func (ll *LinkedList[T]) PushTail(vs ...T) {
	ll.Insert(ll.len, vs...)
}

// PushTailAll inserts a copy of another collection at the tail of list ll.
// The ll and ac may be the same. They must not be nil.
func (ll *LinkedList[T]) PushTailAll(ac Collection[T]) {
	ll.InsertAll(ll.len, ac)
}

// String print list to string
func (ll *LinkedList[T]) String() string {
	bs, _ := json.Marshal(ll)
	return string(bs)
}

//-----------------------------------------------------------
func (ll *LinkedList[T]) deleteAll(v T) {
	for ln := ll.head; ln != nil; ln = ln.next {
		if any(ln.value) == any(v) {
			ll.deleteNode(ln)
		}
	}
}

func (ll *LinkedList[T]) deleteNode(ln *linkedListNode[T]) {
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
func (ll *LinkedList[T]) node(i int) *linkedListNode[T] {
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

func (ll *LinkedList[T]) checkItemIndex(index int) int {
	if index >= ll.len || index < -ll.len {
		panic(fmt.Sprintf("LinkedList out of bounds: index=%d, len=%d", index, ll.len))
	}

	if index < 0 {
		index += ll.len
	}
	return index
}

func (ll *LinkedList[T]) checkSizeIndex(index int) int {
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

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(ll)
func (ll *LinkedList[T]) MarshalJSON() (res []byte, err error) {
	return jsonMarshalCol[T](ll)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, ll)
func (ll *LinkedList[T]) UnmarshalJSON(data []byte) error {
	ll.Clear()
	return jsonUnmarshalCol[T](data, ll)
}
