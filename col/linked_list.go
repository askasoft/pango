package col

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/pandafw/pango/ars"
	"github.com/pandafw/pango/cmp"
)

// NewLinkedList returns an initialized list.
// Example: NewLinkedList(1, 2, 3)
func NewLinkedList(vs ...T) *LinkedList {
	ll := &LinkedList{}
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

// Add adds all items of vs and returns the last added item.
func (ll *LinkedList) Add(vs ...T) {
	ll.Insert(ll.len, vs...)
}

// AddAll adds all items of another collection
func (ll *LinkedList) AddAll(ac Collection) {
	ll.InsertAll(ll.len, ac)
}

// Delete delete all items with associated value v of vs
func (ll *LinkedList) Delete(vs ...T) {
	for _, v := range vs {
		ll.deleteAll(v)
	}
}

// DeleteAll delete all of this collection's elements that are also contained in the specified collection
func (ll *LinkedList) DeleteAll(ac Collection) {
	if ll.IsEmpty() || ac.IsEmpty() {
		return
	}

	if ll == ac {
		ll.Clear()
		return
	}

	if ic, ok := ac.(Iterable); ok {
		it := ic.Iterator()
		for it.Next() {
			ll.deleteAll(it.Value())
		}
		return
	}

	ll.Delete(ac.Values()...)
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

// ContainsAll Test to see if the collection contains all items of another collection
func (ll *LinkedList) ContainsAll(ac Collection) bool {
	if ll == ac || ac.IsEmpty() {
		return true
	}

	if ll.IsEmpty() {
		return false
	}

	if ic, ok := ac.(Iterable); ok {
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
func (ll *LinkedList) Retain(vs ...T) {
	if ll.IsEmpty() {
		return
	}

	if len(vs) == 0 {
		ll.Clear()
		return
	}

	for ln := ll.head; ln != nil; ln = ln.next {
		if !ars.Contains(vs, ln.value) {
			ll.deleteNode(ln)
		}
	}
}

// RetainAll Retains only the elements in this collection that are contained in the specified collection.
func (ll *LinkedList) RetainAll(ac Collection) {
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

// Insert inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (ll *LinkedList) Insert(index int, vs ...T) {
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

// InsertAll inserts values of another collection ac at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (ll *LinkedList) InsertAll(index int, ac Collection) {
	index = ll.checkSizeIndex(index)

	if ac.IsEmpty() {
		return
	}

	ll.Insert(index, ac.Values()...)
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

// Remove removes the element at the specified position in this list.
func (ll *LinkedList) Remove(index int) {
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
func (ll *LinkedList) Sort(less cmp.Less) {
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

// Push inserts all items of vs at the tail of list al.
func (ll *LinkedList) Push(vs ...T) {
	ll.Insert(ll.Len(), vs...)
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

// PushHead inserts all items of vs at the head of list ll.
func (ll *LinkedList) PushHead(vs ...T) {
	ll.Insert(0, vs...)
}

// PushHeadAll inserts a copy of another collection at the head of list ll.
// The ll and ac may be the same. They must not be nil.
func (ll *LinkedList) PushHeadAll(ac Collection) {
	ll.InsertAll(0, ac)
}

// PushTail inserts all items of vs at the tail of list ll.
func (ll *LinkedList) PushTail(vs ...T) {
	ll.Insert(ll.len, vs...)
}

// PushTailAll inserts a copy of another collection at the tail of list ll.
// The ll and ac may be the same. They must not be nil.
func (ll *LinkedList) PushTailAll(ac Collection) {
	ll.InsertAll(ll.len, ac)
}

// String print list to string
func (ll *LinkedList) String() string {
	bs, _ := json.Marshal(ll)
	return string(bs)
}

//-----------------------------------------------------------
func (ll *LinkedList) deleteAll(v T) {
	for ln := ll.head; ln != nil; ln = ln.next {
		if ln.value == v {
			ll.deleteNode(ln)
		}
	}
}

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

//-----------------------------------------------------
// linkedListNode is a node of a doublly-linked list.
type linkedListNode struct {
	prev  *linkedListNode
	next  *linkedListNode
	value T
}

// String print the list item to string
func (ln *linkedListNode) String() string {
	return fmt.Sprintf("%v", ln.value)
}

// linkedListIterator a iterator for linkedListNode
type linkedListIterator struct {
	list    *LinkedList
	node    *linkedListNode
	removed bool
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (it *linkedListIterator) Prev() bool {
	if it.list.IsEmpty() {
		return false
	}

	if it.node == nil {
		it.node = it.list.tail
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
func (it *linkedListIterator) Next() bool {
	if it.list.IsEmpty() {
		return false
	}

	if it.node == nil {
		it.node = it.list.head
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
func (it *linkedListIterator) Value() T {
	if it.node == nil {
		return nil
	}
	return it.node.value
}

// SetValue set the value to the item
func (it *linkedListIterator) SetValue(v T) {
	if it.node != nil {
		it.node.value = v
	}
}

// Remove remove the current element
func (it *linkedListIterator) Remove() {
	if it.node == nil {
		return
	}

	if it.removed {
		panic("LinkedList can't remove a unlinked item")
	}

	it.list.deleteNode(it.node)
	it.removed = true
}

// Reset resets the iterator to its initial state (one-before-first/one-after-last)
// Call Next()/Prev() to fetch the first/last element if any.
func (it *linkedListIterator) Reset() {
	it.node = nil
	it.removed = false
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

func newJSONArrayLinkedList() jsonArray {
	return NewLinkedList()
}

func (ll *LinkedList) addJSONArrayItem(v T) jsonArray {
	ll.Add(v)
	return ll
}

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(ll)
func (ll *LinkedList) MarshalJSON() (res []byte, err error) {
	return jsonMarshalList(ll)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, ll)
func (ll *LinkedList) UnmarshalJSON(data []byte) error {
	ll.Clear()
	ju := &jsonUnmarshaler{
		newArray:  newJSONArrayLinkedList,
		newObject: newJSONObject,
	}
	return ju.unmarshalJSONArray(data, ll)
}
