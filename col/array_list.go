package col

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/pandafw/pango/ars"
)

// NewArrayList returns an initialized list.
// Example: NewArrayList(1, 2, 3)
func NewArrayList(vs ...T) *ArrayList {
	al := &ArrayList{data: vs}
	return al
}

// AsArrayList returns an initialized list.
// Example: AsArrayList([]T{1, 2, 3})
func AsArrayList(vs []T) *ArrayList {
	al := &ArrayList{data: vs}
	return al
}

// ArrayList implements a list holdes the item in a array.
// The zero value for ArrayList is an empty list ready to use.
//
// To iterate over a list (where al is a *ArrayList):
//	it := al.Iterator()
//	for it.Next() {
//		// do something with it.Value()
//	}
//
type ArrayList struct {
	data []T
}

// Cap returns the capcity of the list.
func (al *ArrayList) Cap() int {
	return cap(al.data)
}

//-----------------------------------------------------------
// implements Collection interface

// Len returns the length of the list.
func (al *ArrayList) Len() int {
	return len(al.data)
}

// IsEmpty returns true if the list length == 0
func (al *ArrayList) IsEmpty() bool {
	return al.Len() == 0
}

// Clear clears list al.
func (al *ArrayList) Clear() {
	if al.data != nil {
		al.data = al.data[:0]
	}
}

// Add adds all items of vs and returns the last added item.
func (al *ArrayList) Add(vs ...T) {
	al.PushTail(vs...)
}

// AddAll adds all items of another collection
func (al *ArrayList) AddAll(ac Collection) {
	al.PushTailAll(ac)
}

// Delete delete all items with associated value v of vs
func (al *ArrayList) Delete(vs ...T) {
	if len(vs) == 0 {
		return
	}

	if len(vs) == 1 {
		for i := al.Len() - 1; i >= 0; i-- {
			if al.data[i] == vs[0] {
				al.Remove(i)
			}
		}
		return
	}

	for i := al.Len() - 1; i >= 0; i-- {
		if ars.Contains(vs, al.data[i]) {
			al.Remove(i)
		}
	}
}

// DeleteAll delete all of this collection's elements that are also contained in the specified collection
func (al *ArrayList) DeleteAll(ac Collection) {
	if ac.IsEmpty() {
		return
	}

	if al == ac {
		al.Clear()
		return
	}

	for i := al.Len() - 1; i >= 0; i-- {
		if ac.Contains(al.data[i]) {
			al.Remove(i)
		}
	}
}

// Contains Test to see if the list contains the value v
func (al *ArrayList) Contains(vs ...T) bool {
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

// ContainsAll Test to see if the collection contains all items of another collection
func (al *ArrayList) ContainsAll(ac Collection) bool {
	if ac.IsEmpty() || al == ac {
		return true
	}

	if al.IsEmpty() {
		return false
	}

	if ic, ok := ac.(Iterable); ok {
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

// Retain Retains only the elements in this collection that are contained in the argument array vs.
func (al *ArrayList) Retain(vs ...T) {
	if al.IsEmpty() {
		return
	}

	if len(vs) == 0 {
		al.Clear()
		return
	}

	for i := al.Len() - 1; i >= 0; i-- {
		if !ars.Contains(vs, al.data[i]) {
			al.Remove(i)
		}
	}
}

// RetainAll Retains only the elements in this collection that are contained in the specified collection.
func (al *ArrayList) RetainAll(ac Collection) {
	if al.IsEmpty() || al == ac {
		return
	}

	if ac.IsEmpty() {
		al.Clear()
		return
	}

	for i := al.Len() - 1; i >= 0; i-- {
		if !ac.Contains(al.data[i]) {
			al.Remove(i)
		}
	}
}

// Values returns a slice contains all the items of the list al
func (al *ArrayList) Values() []T {
	return al.data
}

// Each call f for each item in the list
func (al *ArrayList) Each(f func(T)) {
	for _, v := range al.data {
		f(v)
	}
}

// ReverseEach call f for each item in the list with reverse order
func (al *ArrayList) ReverseEach(f func(T)) {
	for i := al.Len() - 1; i >= 0; i-- {
		f(al.data[i])
	}
}

// Iterator returns a iterator for the list
func (al *ArrayList) Iterator() Iterator {
	return &arrayListIterator{al, -1, -1}
}

//-----------------------------------------------------------
// implements List interface

// Get returns the item at the specified position in this list
// if i < -al.Len() or i >= al.Len(), panic
// if i < 0, returns al.Get(al.Len() + i)
func (al *ArrayList) Get(index int) T {
	index = al.checkItemIndex(index)

	return al.data[index]
}

// Set set the v at the specified index in this list and returns the old value.
func (al *ArrayList) Set(index int, v T) (ov T) {
	index = al.checkItemIndex(index)

	ov = al.data[index]
	al.data[index] = v
	return
}

// Insert inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (al *ArrayList) Insert(index int, vs ...T) {
	index = al.checkSizeIndex(index)

	n := len(vs)
	if n == 0 {
		return
	}

	len := al.Len()

	al.expand(n)
	if index < len {
		copy(al.data[index+n:], al.data[index:len-index])
	}
	copy(al.data[index:], vs)
}

// InsertAll inserts values of another collection ac at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (al *ArrayList) InsertAll(index int, ac Collection) {
	al.Insert(index, ac.Values()...)
}

// Index returns the index of the first occurrence of the specified v in this list, or -1 if this list does not contain v.
func (al *ArrayList) Index(v T) int {
	for i, d := range al.data {
		if d == v {
			return i
		}
	}
	return -1
}

// Remove removes the item at the specified position in this list.
func (al *ArrayList) Remove(index int) {
	index = al.checkItemIndex(index)

	al.data[index] = nil
	copy(al.data[index:], al.data[index+1:])
	al.data = al.data[:al.Len()-1]
}

// Swap swaps values of two items at the given index.
func (al *ArrayList) Swap(i, j int) {
	i = al.checkItemIndex(i)
	j = al.checkItemIndex(j)

	if i != j {
		al.data[i], al.data[j] = al.data[j], al.data[i]
	}
}

// Sort Sorts this list according to the order induced by the specified Comparator.
func (al *ArrayList) Sort(less Less) {
	if al.Len() < 2 {
		return
	}
	sort.Sort(&sorter{al, less})
}

// Head get the first item of list.
func (al *ArrayList) Head() (v T) {
	v, _ = al.PeekHead()
	return
}

// Tail get the last item of list.
func (al *ArrayList) Tail() (v T) {
	v, _ = al.PeekTail()
	return
}

//--------------------------------------------------------------------
// implements Queue interface

// Peek get the first item of list.
func (al *ArrayList) Peek() (v T, ok bool) {
	return al.PeekHead()
}

// Poll get and remove the first item of list.
func (al *ArrayList) Poll() (T, bool) {
	return al.PollHead()
}

// Push inserts all items of vs at the tail of list al.
func (al *ArrayList) Push(vs ...T) {
	al.Insert(al.Len(), vs...)
}

//--------------------------------------------------------------------
// implements Deque interface

// PeekHead get the first item of list.
func (al *ArrayList) PeekHead() (v T, ok bool) {
	if al.IsEmpty() {
		return
	}

	v, ok = al.data[0], true
	return
}

// PeekTail get the last item of list.
func (al *ArrayList) PeekTail() (v T, ok bool) {
	if al.IsEmpty() {
		return
	}

	v, ok = al.data[al.Len()-1], true
	return
}

// PollHead get and remove the first item of list.
func (al *ArrayList) PollHead() (v T, ok bool) {
	v, ok = al.PeekHead()
	if ok {
		al.Remove(0)
	}
	return
}

// PollTail get and remove the last item of list.
func (al *ArrayList) PollTail() (v T, ok bool) {
	v, ok = al.PeekTail()
	if ok {
		al.data = al.data[:al.Len()-1]
	}
	return
}

// PushHead inserts all items of vs at the head of list al.
func (al *ArrayList) PushHead(vs ...T) {
	al.Insert(0, vs...)
}

// PushHeadAll inserts a copy of another collection at the head of list al.
// The al and ac may be the same. They must not be nil.
func (al *ArrayList) PushHeadAll(ac Collection) {
	al.InsertAll(0, ac)
}

// PushTail inserts all items of vs at the tail of list al.
func (al *ArrayList) PushTail(vs ...T) {
	al.Insert(al.Len(), vs...)
}

// PushTailAll inserts a copy of another collection at the tail of list al.
// The al and ac may be the same. They must not be nil.
func (al *ArrayList) PushTailAll(ac Collection) {
	al.InsertAll(al.Len(), ac)
}

//------------------------------------------------------------

// Reserve Increase the capacity of the underlying array.
func (al *ArrayList) Reserve(n int) {
	l := al.Len()
	n -= l
	if n > 0 {
		al.expand(n)
		al.data = al.data[:l]
	}
}

// String print list to string
func (al *ArrayList) String() string {
	bs, _ := json.Marshal(al)
	return string(bs)
}

//-----------------------------------------------------------

// expand expand the buffer to guarantee space for n more elements.
func (al *ArrayList) expand(x int) {
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

func (al *ArrayList) checkItemIndex(index int) int {
	len := al.Len()
	if index >= len || index < -len {
		panic(fmt.Sprintf("ArrayList out of bounds: index=%d, len=%d", index, len))
	}

	if index < 0 {
		index += len
	}
	return index
}

func (al *ArrayList) checkSizeIndex(index int) int {
	len := al.Len()
	if index > len || index < -len {
		panic(fmt.Sprintf("ArrayList out of bounds: index=%d, len=%d", index, len))
	}

	if index < 0 {
		index += len
	}
	return index
}

//-----------------------------------------------------

// arrayListIterator a iterator for array list
type arrayListIterator struct {
	list  *ArrayList
	start int
	index int
}

// Prev moves the iterator to the previous item and returns true if there was a previous item in the container.
// If Prev() returns true, then previous item's value can be retrieved by Value().
// Modifies the state of the iterator.
func (it *arrayListIterator) Prev() bool {
	if it.list.IsEmpty() {
		return false
	}

	if it.index < 0 && it.start >= 0 {
		it.index = it.start
	}

	if it.index == 0 {
		return false
	}

	if it.index < 0 {
		it.index = it.list.Len() - 1
		return true
	}
	if it.index > it.list.Len() {
		return false
	}
	it.index--
	return true
}

// Next moves the iterator to the next item and returns true if there was a next item in the collection.
// If Next() returns true, then next item's value can be retrieved by Value().
// If Next() was called for the first time, then it will point the iterator to the first item if it exists.
// Modifies the state of the iterator.
func (it *arrayListIterator) Next() bool {
	if it.list.IsEmpty() {
		return false
	}

	if it.index < 0 && it.start > 0 {
		it.index = it.start - 1
	}
	if it.index < -1 || it.index >= it.list.Len()-1 {
		return false
	}
	it.index++
	return true
}

// Value returns the current item's value.
func (it *arrayListIterator) Value() T {
	if it.index >= 0 && it.index < it.list.Len() {
		return it.list.data[it.index]
	}
	return nil
}

// SetValue set the value to the item
func (it *arrayListIterator) SetValue(v T) {
	if it.index >= 0 && it.index < it.list.Len() {
		it.list.data[it.index] = v
	}
}

// Remove remove the current item
func (it *arrayListIterator) Remove() {
	if it.index < 0 {
		return
	}

	it.list.Remove(it.index)
	it.start = it.index
	it.index = -1
}

// Reset resets the iterator to its initial state (one-before-first/one-after-last)
// Call Next()/Prev() to fetch the first/last item if any.
func (it *arrayListIterator) Reset() {
	it.start = -1
	it.index = -1
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

func (al *ArrayList) addJSONArrayItem(v T) jsonArray {
	al.Add(v)
	return al
}

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(al)
func (al *ArrayList) MarshalJSON() (res []byte, err error) {
	return jsonMarshalArray(al)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, al)
func (al *ArrayList) UnmarshalJSON(data []byte) error {
	al.Clear()
	return jsonUnmarshalArray(data, al)
}
