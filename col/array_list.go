package col

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/askasoft/pango/bye"
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
//
//	it := al.Iterator()
//	for it.Next() {
//		// do something with it.Value()
//	}
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

// Add add the item v
func (al *ArrayList) Add(v T) {
	al.Insert(al.Len(), v)
}

// Adds adds all items of vs
func (al *ArrayList) Adds(vs ...T) {
	al.Inserts(al.Len(), vs...)
}

// AddCol adds all items of another collection
func (al *ArrayList) AddCol(ac Collection) {
	al.InsertCol(al.Len(), ac)
}

// Remove remove all items with associated value v of vs
func (al *ArrayList) Remove(v T) {
	for i := al.Len() - 1; i >= 0; i-- {
		if al.data[i] == v {
			al.DeleteAt(i)
		}
	}
}

// Removes remove all items with associated value v of vs
func (al *ArrayList) Removes(vs ...T) {
	if al.IsEmpty() {
		return
	}

	for _, v := range vs {
		al.Remove(v)
	}
}

// RemoveFunc remove all items that function f returns true
func (al *ArrayList) RemoveFunc(f func(T) bool) {
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
func (al *ArrayList) RemoveCol(ac Collection) {
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
func (al *ArrayList) Contain(v T) bool {
	return al.Index(v) >= 0
}

// Contains Test to see if the collection contains all items of vs
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

// ContainCol Test to see if the collection contains all items of another collection
func (al *ArrayList) ContainCol(ac Collection) bool {
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

// Retains Retains only the elements in this collection that are contained in the argument array vs.
func (al *ArrayList) Retains(vs ...T) {
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
func (al *ArrayList) RetainCol(ac Collection) {
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

// Insert insert the item v at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (al *ArrayList) Insert(index int, v T) {
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
func (al *ArrayList) Inserts(index int, vs ...T) {
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
func (al *ArrayList) InsertCol(index int, ac Collection) {
	al.Inserts(index, ac.Values()...)
}

// Index returns the index of the first occurrence of the specified v in this list, or -1 if this list does not contain v.
func (al *ArrayList) Index(v T) int {
	return index(al.data, v)
}

// IndexFunc returns the index of the first true returned by function f in this list, or -1 if this list does not contain v.
func (al *ArrayList) IndexFunc(f func(T) bool) int {
	for i, v := range al.data {
		if f(v) {
			return i
		}
	}
	return -1
}

// DeleteAt remove the item at the specified position in this list.
func (al *ArrayList) DeleteAt(index int) {
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

// Push insert item v at the tail of list al.
func (al *ArrayList) Push(v T) {
	al.Insert(al.Len(), v)
}

// Push inserts all items of vs at the tail of list al.
func (al *ArrayList) Pushs(vs ...T) {
	al.Inserts(al.Len(), vs...)
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
		al.DeleteAt(0)
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

// PushHead inserts the item v at the head of list al.
func (al *ArrayList) PushHead(v T) {
	al.Insert(0, v)
}

// PushHeads inserts all items of vs at the head of list al.
func (al *ArrayList) PushHeads(vs ...T) {
	al.Inserts(0, vs...)
}

// PushHeadCol inserts a copy of another collection at the head of list al.
// The al and ac may be the same. They must not be nil.
func (al *ArrayList) PushHeadCol(ac Collection) {
	al.InsertCol(0, ac)
}

// PushTail inserts the item v at the tail of list al.
func (al *ArrayList) PushTail(v T) {
	al.Insert(al.Len(), v)
}

// PushTails inserts all items of vs at the tail of list al.
func (al *ArrayList) PushTails(vs ...T) {
	al.Inserts(al.Len(), vs...)
}

// PushTailCol inserts a copy of another collection at the tail of list al.
// The al and ac may be the same. They must not be nil.
func (al *ArrayList) PushTailCol(ac Collection) {
	al.InsertCol(al.Len(), ac)
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
	return bye.UnsafeString(bs)
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
	sz := al.Len()
	if index >= sz || index < -sz {
		panic(fmt.Sprintf("ArrayList out of bounds: index=%d, len=%d", index, sz))
	}

	if index < 0 {
		index += sz
	}
	return index
}

func (al *ArrayList) checkSizeIndex(index int) int {
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

func (al *ArrayList) addJSONArrayItem(v T) jsonArray {
	al.Add(v)
	return al
}

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(al)
func (al *ArrayList) MarshalJSON() ([]byte, error) {
	return jsonMarshalArray(al)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, al)
func (al *ArrayList) UnmarshalJSON(data []byte) error {
	al.Clear()
	return jsonUnmarshalArray(data, al)
}
