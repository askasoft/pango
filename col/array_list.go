package col

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/pandafw/pango/ars"
	"github.com/pandafw/pango/cmp"
)

// NewArrayList returns an initialized list.
// Example: NewArrayList(1, 2, 3)
func NewArrayList(vs ...T) *ArrayList {
	al := &ArrayList{data: vs}
	return al
}

// ArrayList implements a list holdes the item in a array.
// The zero value for ArrayList is an empty list ready to use.
//
// To iterate over a list (where al is a *ArrayList):
//	for li := al.Front(); li != nil; li = li.Next() {
//		// do something with li.Value()
//	}
//
type ArrayList struct {
	data []T
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
	n := len(vs)
	if n == 0 {
		return
	}

	al.grow(n)
	copy(al.data[al.Len()-n:], vs)
}

// AddAll adds all items of another collection
func (al *ArrayList) AddAll(ac Collection) {
	al.Add(ac.Values()...)
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
	return
}

// DeleteAll delete all of this collection's elements that are also contained in the specified collection
func (al *ArrayList) DeleteAll(ac Collection) {
	al.Delete(ac.Values()...)
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
	if al == ac || ac.IsEmpty() {
		return true
	}

	if al.IsEmpty() {
		return false
	}

	return al.Contains(ac.Values()...)
}

// Retain Retains only the elements in this collection that are contained in the argument array vs.
func (al *ArrayList) Retain(vs ...T) {
	if al.IsEmpty() || len(vs) == 0 {
		return
	}

	al.RetainAll(NewArrayList(vs...))
}

// RetainAll Retains only the elements in this collection that are contained in the specified collection.
func (al *ArrayList) RetainAll(ac Collection) {
	if al.IsEmpty() || ac.IsEmpty() || al == ac {
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
	if index == len {
		al.Add(vs...)
		return
	}

	al.grow(n)
	copy(al.data[index+n:], al.data[index:len-index])
	copy(al.data[index:], vs)
}

// InsertAll inserts values of another collection ac at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (al *ArrayList) InsertAll(index int, ac Collection) {
	index = al.checkSizeIndex(index)
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
func (al *ArrayList) Sort(less cmp.Less) {
	if al.Len() < 2 {
		return
	}
	sort.Sort(&sorter{al, less})
}

//--------------------------------------------------------------------

// Front returns the first item of list al or nil if the list is empty.
func (al *ArrayList) Front() T {
	if al.IsEmpty() {
		return nil
	}

	return al.data[0]
}

// Back returns the last item of list al or nil if the list is empty.
func (al *ArrayList) Back() T {
	if al.IsEmpty() {
		return nil
	}

	return al.data[al.Len()-1]
}

// PopFront remove the first item of list.
func (al *ArrayList) PopFront() (v T) {
	if al.IsEmpty() {
		return
	}

	v = al.data[0]
	al.Remove(0)
	return
}

// PopBack remove the last item of list.
func (al *ArrayList) PopBack() (v T) {
	if al.IsEmpty() {
		return
	}

	v = al.data[al.Len()-1]
	al.data = al.data[:al.Len()-1]
	return
}

// PushFront inserts all items of vs at the front of list al.
func (al *ArrayList) PushFront(vs ...T) {
	if len(vs) == 0 {
		return
	}

	al.Insert(0, vs...)
}

// PushFrontAll inserts a copy of another collection at the front of list al.
// The al and ac may be the same. They must not be nil.
func (al *ArrayList) PushFrontAll(ac Collection) {
	if ac.IsEmpty() {
		return
	}

	al.InsertAll(0, ac)
}

// PushBack inserts all items of vs at the back of list al.
func (al *ArrayList) PushBack(vs ...T) {
	if len(vs) == 0 {
		return
	}

	al.Insert(al.Len(), vs...)
}

// PushBackAll inserts a copy of another collection at the back of list al.
// The al and ac may be the same. They must not be nil.
func (al *ArrayList) PushBackAll(ac Collection) {
	if ac.IsEmpty() {
		return
	}

	al.InsertAll(al.Len(), ac)
}

//------------------------------------------------------------

// Reserve Increase the capacity of the underlying array.
func (al *ArrayList) Reserve(n int) {
	l := al.Len()
	n -= l
	if n > 0 {
		al.grow(n)
		al.data = al.data[:l]
	}
}

// String print list to string
func (al *ArrayList) String() string {
	bs, _ := json.Marshal(al)
	return string(bs)
}

//-----------------------------------------------------------
const nblock = 0x1F

// roundup round up size
func (al *ArrayList) roundup(n int) int {
	if (n & nblock) == 0 {
		return n
	}

	return (n + nblock) & (^nblock)
}

// grow grows the buffer to guarantee space for n more elements.
func (al *ArrayList) grow(n int) {
	if al.data == nil {
		c := al.roundup(n)
		al.data = make([]T, n, c)
		return
	}

	l := len(al.data)
	if n <= cap(al.data)-l {
		al.data = al.data[:l+n]
		return
	}

	c := al.roundup(l + n)
	data := make([]T, l+n, c)
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

func newJSONArrayArrayList() jsonArray {
	return NewArrayList()
}

func (al *ArrayList) addJSONArrayItem(v T) jsonArray {
	al.Add(v)
	return al
}

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(al)
func (al *ArrayList) MarshalJSON() (res []byte, err error) {
	return jsonMarshalList(al)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, al)
func (al *ArrayList) UnmarshalJSON(data []byte) error {
	al.Clear()
	ju := &jsonUnmarshaler{
		newArray:  newJSONArrayArrayList,
		newObject: newJSONObject,
	}
	return ju.unmarshalJSONArray(data, al)
}
