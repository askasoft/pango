package col

import (
	"encoding/json"
	"sort"

	"github.com/pandafw/pango/ars"
	"github.com/pandafw/pango/cmp"
)

// NewArrayList returns an initialized list.
// Example: NewArrayList(1, 2, 3)
func NewArrayList(vs ...interface{}) *ArrayList {
	al := &ArrayList{data: vs}
	return al
}

// ArrayList implements a list holdes the element in a array.
// The zero value for ArrayList is an empty list ready to use.
//
// To iterate over a list (where al is a *ArrayList):
//	for li := al.Front(); li != nil; li = li.Next() {
//		// do something with li.Value()
//	}
//
type ArrayList struct {
	data []interface{}
}

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
		al.data = make([]interface{}, n, c)
		return
	}

	l := len(al.data)
	if n <= cap(al.data)-l {
		al.data = al.data[:l+n]
		return
	}

	c := al.roundup(l + n)
	data := make([]interface{}, l+n, c)
	copy(data, al.data)
	al.data = data
}

func (al *ArrayList) checkIndex(index int) int {
	len := al.Len()

	if index < -len || index >= len {
		return -1
	}

	if index < 0 {
		index += len
	}

	return index
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
func (al *ArrayList) Add(vs ...interface{}) {
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
func (al *ArrayList) Delete(vs ...interface{}) {
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
func (al *ArrayList) Contains(vs ...interface{}) bool {
	for _, v := range vs {
		if al.Index(v) < 0 {
			return false
		}
	}
	return true
}

// ContainsAll Test to see if the collection contains all items of another collection
func (al *ArrayList) ContainsAll(ac Collection) bool {
	if al == ac {
		return true
	}
	return al.Contains(ac.Values()...)
}

// Retain Retains only the elements in this collection that are contained in the argument array vs.
func (al *ArrayList) Retain(vs ...interface{}) {
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
func (al *ArrayList) Values() []interface{} {
	return al.data
}

// Each call f for each item in the list
func (al *ArrayList) Each(f func(interface{})) {
	for _, v := range al.data {
		f(v)
	}
}

//-----------------------------------------------------------
// implements List interface

// Get returns the element at the specified position in this list
func (al *ArrayList) Get(index int) (interface{}, bool) {
	index = al.checkIndex(index)
	if index < 0 {
		return nil, false
	}

	return al.data[index], true
}

// Set set the v at the specified index in this list and returns the old value.
func (al *ArrayList) Set(index int, v interface{}) (ov interface{}) {
	index = al.checkIndex(index)
	if index < 0 {
		return nil
	}

	ov = al.data[index]
	al.data[index] = v
	return
}

// Insert inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Does not do anything if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (al *ArrayList) Insert(index int, vs ...interface{}) {
	n := len(vs)
	if n == 0 {
		return
	}

	len := al.Len()
	if index < -len || index > len {
		return
	}

	if index < 0 {
		index += len
	}

	if index == len {
		// Append
		al.Add(vs...)
		return
	}

	al.grow(n)
	copy(al.data[index+n:], al.data[index:len-index])
	copy(al.data[index:], vs)
}

// InsertAll inserts values of another collection ac at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Does not do anything if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (al *ArrayList) InsertAll(index int, ac Collection) {
	al.Insert(index, ac.Values()...)
}

// Index returns the index of the first occurrence of the specified v in this list, or -1 if this list does not contain v.
func (al *ArrayList) Index(v interface{}) int {
	for i, d := range al.data {
		if d == v {
			return i
		}
	}
	return -1
}

// Remove removes the item li from al if li is an item of list al.
// It returns the item's value.
// The item li must not be nil.
func (al *ArrayList) Remove(index int) {
	index = al.checkIndex(index)
	if index < 0 {
		return
	}

	al.data[index] = nil
	copy(al.data[index:], al.data[index+1:])
	al.data = al.data[:al.Len()-1]
}

// Swap swaps values of two items at the given index.
func (al *ArrayList) Swap(i, j int) {
	i = al.checkIndex(i)
	j = al.checkIndex(j)
	if i < 0 || j < 0 || i == j {
		return
	}

	al.data[i], al.data[j] = al.data[j], al.data[i]
}

// ReverseEach call f for each item in the list with reverse order
func (al *ArrayList) ReverseEach(f func(interface{})) {
	for i := al.Len() - 1; i >= 0; i-- {
		f(al.data[i])
	}
}

// Iterator returns a iterator for the list
func (al *ArrayList) Iterator() Iterator {
	return &ArrayListIterator{al, -1, -1}
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

// Sort Sorts this list according to the order induced by the specified Comparator.
func (al *ArrayList) Sort(less cmp.Less) {
	if al.Len() < 2 {
		return
	}
	sort.Sort(&sorter{al, less})
}

// String print list to string
func (al *ArrayList) String() string {
	bs, _ := json.Marshal(al)
	return string(bs)
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

func newJSONArrayArrayList() jsonArray {
	return NewArrayList()
}

func (al *ArrayList) addJSONArrayItem(v interface{}) jsonArray {
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
