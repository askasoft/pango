//go:build go1.18
// +build go1.18

package cog

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/askasoft/pango/bye"
)

// NewArrayList returns an initialized list.
// Example: NewArrayList(1, 2, 3)
func NewArrayList[T any](vs ...T) *ArrayList[T] {
	al := &ArrayList[T]{data: vs}
	return al
}

// AsArrayList returns an initialized list.
// Example: AsArrayList([]T{1, 2, 3})
func AsArrayList[T any](vs []T) *ArrayList[T] {
	al := &ArrayList[T]{data: vs}
	return al
}

// ArrayList implements a list holdes the item in a array.
// The zero value for ArrayList is an empty list ready to use.
//
// To iterate over a list (where al is a *ArrayList):
//
//	it := al.Iterator()
//	for it.Next() {
//		// do something with it.Value()
//	}
type ArrayList[T any] struct {
	data []T
}

// Cap returns the capcity of the list.
func (al *ArrayList[T]) Cap() int {
	return cap(al.data)
}

//-----------------------------------------------------------
// implements Collection interface

// Len returns the length of the list.
func (al *ArrayList[T]) Len() int {
	return len(al.data)
}

// IsEmpty returns true if the list length == 0
func (al *ArrayList[T]) IsEmpty() bool {
	return al.Len() == 0
}

// Clear clears list al.
func (al *ArrayList[T]) Clear() {
	if al.data != nil {
		al.data = al.data[:0]
	}
}

// Add add the item v
func (al *ArrayList[T]) Add(v T) {
	al.Insert(al.Len(), v)
}

// Adds adds all items of vs
func (al *ArrayList[T]) Adds(vs ...T) {
	al.Inserts(al.Len(), vs...)
}

// AddCol adds all items of another collection
func (al *ArrayList[T]) AddCol(ac Collection[T]) {
	al.InsertCol(al.Len(), ac)
}

// Remove remove all items with associated value v of vs
func (al *ArrayList[T]) Remove(v T) {
	i := al.Index(v)
	if i < 0 {
		return
	}

	// Don't start copying elements until we find one to delete.
	a := al.data
	for j := i + 1; j < len(a); j++ {
		if e := a[j]; any(e) != any(v) {
			a[i] = e
			i++
		}
	}
	al.data = a[:i]
}

// Removes remove all items with associated value v of vs
func (al *ArrayList[T]) Removes(vs ...T) {
	if al.IsEmpty() {
		return
	}
	for _, v := range vs {
		al.Remove(v)
	}
}

// RemoveIf remove all items that function f returns true
func (al *ArrayList[T]) RemoveIf(f func(T) bool) {
	if al.IsEmpty() {
		return
	}

	for i := al.Len() - 1; i >= 0; i-- {
		if f(al.data[i]) {
			al.DeleteAt(i)
		}
	}
}

// RemoveCol remove all of this collection's elements that are also contained in the specified collection
func (al *ArrayList[T]) RemoveCol(ac Collection[T]) {
	if ac.IsEmpty() {
		return
	}

	if al == ac {
		al.Clear()
		return
	}

	for i := al.Len() - 1; i >= 0; i-- {
		if ac.Contain(al.data[i]) {
			al.DeleteAt(i)
		}
	}
}

// Contain Test to see if the list contains the value v
func (al *ArrayList[T]) Contain(v T) bool {
	return al.Index(v) >= 0
}

// Contains Test to see if the collection contains all items of vs
func (al *ArrayList[T]) Contains(vs ...T) bool {
	if len(vs) == 0 {
		return true
	}

	if al.IsEmpty() {
		return false
	}

	for _, v := range vs {
		if al.Index(v) < 0 {
			return false
		}
	}
	return true
}

// ContainCol Test to see if the collection contains all items of another collection
func (al *ArrayList[T]) ContainCol(ac Collection[T]) bool {
	if ac.IsEmpty() || al == ac {
		return true
	}

	if al.IsEmpty() {
		return false
	}

	if ic, ok := ac.(Iterable[T]); ok {
		it := ic.Iterator()
		for it.Next() {
			if al.Index(it.Value()) < 0 {
				return false
			}
		}
		return true
	}

	return al.Contains(ac.Values()...)
}

// Retains Retains only the elements in this collection that are contained in the argument array vs.
func (al *ArrayList[T]) Retains(vs ...T) {
	if al.IsEmpty() {
		return
	}

	if len(vs) == 0 {
		al.Clear()
		return
	}

	for i := al.Len() - 1; i >= 0; i-- {
		if !contains(vs, al.data[i]) {
			al.DeleteAt(i)
		}
	}
}

// RetainCol Retains only the elements in this collection that are contained in the specified collection.
func (al *ArrayList[T]) RetainCol(ac Collection[T]) {
	if al.IsEmpty() || al == ac {
		return
	}

	if ac.IsEmpty() {
		al.Clear()
		return
	}

	for i := al.Len() - 1; i >= 0; i-- {
		if !ac.Contain(al.data[i]) {
			al.DeleteAt(i)
		}
	}
}

// Values returns a slice contains all the items of the list al
func (al *ArrayList[T]) Values() []T {
	return al.data
}

// Each call f for each item in the list
func (al *ArrayList[T]) Each(f func(T)) {
	for _, v := range al.data {
		f(v)
	}
}

// ReverseEach call f for each item in the list with reverse order
func (al *ArrayList[T]) ReverseEach(f func(T)) {
	for i := al.Len() - 1; i >= 0; i-- {
		f(al.data[i])
	}
}

// Iterator returns a iterator for the list
func (al *ArrayList[T]) Iterator() Iterator[T] {
	return &arrayListIterator[T]{al, -1, -1}
}

//-----------------------------------------------------------
// implements List interface

// Get returns the item at the specified position in this list
// if i < -al.Len() or i >= al.Len(), panic
// if i < 0, returns al.Get(al.Len() + i)
func (al *ArrayList[T]) Get(index int) T {
	index = al.checkItemIndex(index)

	return al.data[index]
}

// Set set the v at the specified index in this list and returns the old value.
func (al *ArrayList[T]) Set(index int, v T) (ov T) {
	index = al.checkItemIndex(index)

	ov = al.data[index]
	al.data[index] = v
	return
}

// Insert insert the item v at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (al *ArrayList[T]) Insert(index int, v T) {
	index = al.checkSizeIndex(index)

	z := al.Len()

	al.expand(1)
	if index < z {
		copy(al.data[index+1:], al.data[index:z-index])
	}
	al.data[index] = v
}

// Insert inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (al *ArrayList[T]) Inserts(index int, vs ...T) {
	index = al.checkSizeIndex(index)

	n := len(vs)
	if n == 0 {
		return
	}

	z := al.Len()

	al.expand(n)
	if index < z {
		copy(al.data[index+n:], al.data[index:z-index])
	}
	copy(al.data[index:], vs)
}

// InsertCol inserts values of another collection ac at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (al *ArrayList[T]) InsertCol(index int, ac Collection[T]) {
	al.Inserts(index, ac.Values()...)
}

// Index returns the index of the first occurrence of the specified v in this list, or -1 if this list does not contain v.
func (al *ArrayList[T]) Index(v T) int {
	return index(al.data, v)
}

// IndexFunc returns the index of the first true returned by function f in this list, or -1 if this list does not contain v.
func (al *ArrayList[T]) IndexFunc(f func(T) bool) int {
	for i, v := range al.data {
		if f(v) {
			return i
		}
	}
	return -1
}

// DeleteAt remove the item at the specified position in this list.
func (al *ArrayList[T]) DeleteAt(index int) {
	index = al.checkItemIndex(index)

	var v T
	al.data[index] = v
	copy(al.data[index:], al.data[index+1:])
	al.data = al.data[:al.Len()-1]
}

// Swap swaps values of two items at the given index.
func (al *ArrayList[T]) Swap(i, j int) {
	i = al.checkItemIndex(i)
	j = al.checkItemIndex(j)

	if i != j {
		al.data[i], al.data[j] = al.data[j], al.data[i]
	}
}

// Sort Sorts this list according to the order induced by the specified Comparator.
func (al *ArrayList[T]) Sort(less Less[T]) {
	if al.Len() < 2 {
		return
	}
	sort.Sort(&sorter[T]{al, less})
}

// Head get the first item of list.
func (al *ArrayList[T]) Head() (v T) {
	v, _ = al.PeekHead()
	return
}

// Tail get the last item of list.
func (al *ArrayList[T]) Tail() (v T) {
	v, _ = al.PeekTail()
	return
}

//--------------------------------------------------------------------
// implements Queue interface

// Peek get the first item of list.
func (al *ArrayList[T]) Peek() (v T, ok bool) {
	return al.PeekHead()
}

// Poll get and remove the first item of list.
func (al *ArrayList[T]) Poll() (T, bool) {
	return al.PollHead()
}

// Push insert item v at the tail of list al.
func (al *ArrayList[T]) Push(v T) {
	al.Insert(al.Len(), v)
}

// Push inserts all items of vs at the tail of list al.
func (al *ArrayList[T]) Pushs(vs ...T) {
	al.Inserts(al.Len(), vs...)
}

//--------------------------------------------------------------------
// implements Deque interface

// PeekHead get the first item of list.
func (al *ArrayList[T]) PeekHead() (v T, ok bool) {
	if al.IsEmpty() {
		return
	}

	v, ok = al.data[0], true
	return
}

// PeekTail get the last item of list.
func (al *ArrayList[T]) PeekTail() (v T, ok bool) {
	if al.IsEmpty() {
		return
	}

	v, ok = al.data[al.Len()-1], true
	return
}

// PollHead get and remove the first item of list.
func (al *ArrayList[T]) PollHead() (v T, ok bool) {
	v, ok = al.PeekHead()
	if ok {
		al.DeleteAt(0)
	}
	return
}

// PollTail get and remove the last item of list.
func (al *ArrayList[T]) PollTail() (v T, ok bool) {
	v, ok = al.PeekTail()
	if ok {
		al.data = al.data[:al.Len()-1]
	}
	return
}

// PushHead inserts the item v at the head of list al.
func (al *ArrayList[T]) PushHead(v T) {
	al.Insert(0, v)
}

// PushHeads inserts all items of vs at the head of list al.
func (al *ArrayList[T]) PushHeads(vs ...T) {
	al.Inserts(0, vs...)
}

// PushHeadCol inserts a copy of another collection at the head of list al.
// The al and ac may be the same. They must not be nil.
func (al *ArrayList[T]) PushHeadCol(ac Collection[T]) {
	al.InsertCol(0, ac)
}

// PushTail inserts the item v at the tail of list al.
func (al *ArrayList[T]) PushTail(v T) {
	al.Insert(al.Len(), v)
}

// PushTails inserts all items of vs at the tail of list al.
func (al *ArrayList[T]) PushTails(vs ...T) {
	al.Inserts(al.Len(), vs...)
}

// PushTailCol inserts a copy of another collection at the tail of list al.
// The al and ac may be the same. They must not be nil.
func (al *ArrayList[T]) PushTailCol(ac Collection[T]) {
	al.InsertCol(al.Len(), ac)
}

//------------------------------------------------------------

// Reserve Increase the capacity of the underlying array.
func (al *ArrayList[T]) Reserve(n int) {
	l := al.Len()
	n -= l
	if n > 0 {
		al.expand(n)
		al.data = al.data[:l]
	}
}

// String print list to string
func (al *ArrayList[T]) String() string {
	bs, _ := json.Marshal(al)
	return bye.UnsafeString(bs)
}

//-----------------------------------------------------------

// expand expand the buffer to guarantee space for n more elements.
func (al *ArrayList[T]) expand(x int) {
	l := len(al.data)
	c := cap(al.data)
	if l+x <= c {
		al.data = al.data[:l+x]
		return
	}

	c = doubleup(c, c+x)
	data := make([]T, l+x, c)
	copy(data, al.data)
	al.data = data
}

func (al *ArrayList[T]) checkItemIndex(index int) int {
	sz := al.Len()
	if index >= sz || index < -sz {
		panic(fmt.Sprintf("ArrayList out of bounds: index=%d, len=%d", index, sz))
	}

	if index < 0 {
		index += sz
	}
	return index
}

func (al *ArrayList[T]) checkSizeIndex(index int) int {
	sz := al.Len()
	if index > sz || index < -sz {
		panic(fmt.Sprintf("ArrayList out of bounds: index=%d, len=%d", index, sz))
	}

	if index < 0 {
		index += sz
	}
	return index
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(al)
func (al *ArrayList[T]) MarshalJSON() ([]byte, error) {
	return jsonMarshalCol[T](al)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, al)
func (al *ArrayList[T]) UnmarshalJSON(data []byte) error {
	al.Clear()
	return jsonUnmarshalCol[T](data, al)
}
