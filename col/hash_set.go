package col

import (
	"fmt"
)

// NewHashSet Create a new hash set
func NewHashSet(vs ...T) *HashSet {
	hs := &HashSet{}
	hs.Adds(vs...)
	return hs
}

// NewStringHashSet Create a new hash set
func NewStringHashSet(vs ...string) *HashSet {
	hs := &HashSet{}
	if len(vs) > 0 {
		hs.lazyInit()
		for _, v := range vs {
			hs.hash[v] = true
		}
	}
	return hs
}

// HashSet an unordered collection of unique values.
// The zero value for HashSet is an empty set ready to use.
// http://en.wikipedia.org/wiki/Set_(computer_science%29)
type HashSet struct {
	hash map[T]bool
}

// lazyInit lazily initializes a zero List value.
func (hs *HashSet) lazyInit() {
	if hs.hash == nil {
		hs.hash = make(map[T]bool)
	}
}

//-----------------------------------------------------------
// implements Collection interface

// Len Return the number of items in the set
func (hs *HashSet) Len() int {
	return len(hs.hash)
}

// IsEmpty returns true if the set's length == 0
func (hs *HashSet) IsEmpty() bool {
	return hs.Len() == 0
}

// Clear clears the hash set.
func (hs *HashSet) Clear() {
	hs.hash = nil
}

// Add Add the item v to the set
func (hs *HashSet) Add(v T) {
	hs.lazyInit()
	hs.hash[v] = true
}

// Adds Adds all items of vs to the set
func (hs *HashSet) Adds(vs ...T) {
	if len(vs) == 0 {
		return
	}

	hs.lazyInit()
	for _, v := range vs {
		hs.hash[v] = true
	}
}

// AddCol adds all items of another collection
func (hs *HashSet) AddCol(ac Collection) {
	if ac.IsEmpty() || hs == ac {
		return
	}

	hs.lazyInit()
	if ic, ok := ac.(Iterable); ok {
		it := ic.Iterator()
		for it.Next() {
			hs.hash[it.Value()] = true
		}
		return
	}

	hs.Adds(ac.Values()...)
}

// Remove remove all items with associated value v of vs
func (hs *HashSet) Remove(v T) {
	delete(hs.hash, v)
}

// Removes delete items of vs
func (hs *HashSet) Removes(vs ...T) {
	if hs.IsEmpty() {
		return
	}

	for _, v := range vs {
		delete(hs.hash, v)
	}
}

// RemoveCol remove all of this collection's elements that are also contained in the specified collection
func (hs *HashSet) RemoveCol(ac Collection) {
	if hs == ac {
		hs.Clear()
		return
	}

	if ic, ok := ac.(Iterable); ok {
		hs.RemoveIter(ic.Iterator())
		return
	}

	hs.RemoveFunc(ac.Contain)
}

// RemoveIter remove all items in the iterator it
func (hs *HashSet) RemoveIter(it Iterator) {
	for it.Next() {
		delete(hs.hash, it.Value())
	}
}

// RemoveFunc remove all items that function f returns true
func (hs *HashSet) RemoveFunc(f func(T) bool) {
	if hs.IsEmpty() {
		return
	}

	for k := range hs.hash {
		if f(k) {
			delete(hs.hash, k)
		}
	}
}

// Contain Test to see if the list contains the value v
func (hs *HashSet) Contain(v T) bool {
	if hs.IsEmpty() {
		return false
	}
	if _, ok := hs.hash[v]; ok {
		return true
	}
	return false
}

// Contains Test to see if the collection contains all items of vs
func (hs *HashSet) Contains(vs ...T) bool {
	if len(vs) == 0 {
		return true
	}

	if hs.IsEmpty() {
		return false
	}

	for _, v := range vs {
		if _, ok := hs.hash[v]; !ok {
			return false
		}
	}
	return true
}

// ContainCol Test to see if the collection contains all items of another collection
func (hs *HashSet) ContainCol(ac Collection) bool {
	if ac.IsEmpty() || hs == ac {
		return true
	}

	if hs.IsEmpty() {
		return false
	}

	if ic, ok := ac.(Iterable); ok {
		it := ic.Iterator()
		for it.Next() {
			if _, ok := hs.hash[it.Value()]; !ok {
				return false
			}
		}
		return true
	}

	return hs.Contains(ac.Values()...)
}

// Retains Retains only the elements in this collection that are contained in the argument array vs.
func (hs *HashSet) Retains(vs ...T) {
	if hs.IsEmpty() {
		return
	}

	if len(vs) == 0 {
		hs.Clear()
		return
	}

	for k := range hs.hash {
		if !contains(vs, k) {
			delete(hs.hash, k)
		}
	}
}

// RetainCol Retains only the elements in this collection that are contained in the specified collection.
func (hs *HashSet) RetainCol(ac Collection) {
	if hs.IsEmpty() || hs == ac {
		return
	}

	if ac.IsEmpty() {
		hs.Clear()
		return
	}

	hs.RetainFunc(ac.Contain)
}

// RetainFunc Retains all items that function f returns true
func (hs *HashSet) RetainFunc(f func(T) bool) {
	if hs.IsEmpty() {
		return
	}

	for k := range hs.hash {
		if !f(k) {
			delete(hs.hash, k)
		}
	}
}

// Values returns a slice contains all the items of the set hs
func (hs *HashSet) Values() []T {
	vs := make([]T, hs.Len())
	i := 0
	for k := range hs.hash {
		vs[i] = k
		i++
	}
	return vs
}

// Each Call f for each item in the set
func (hs *HashSet) Each(f func(T)) {
	for k := range hs.hash {
		f(k)
	}
}

//-----------------------------------------------------------

// Difference Find the difference btween two sets
func (hs *HashSet) Difference(a *HashSet) *HashSet {
	b := make(map[T]bool)

	for k := range hs.hash {
		if _, ok := a.hash[k]; !ok {
			b[k] = true
		}
	}

	return &HashSet{b}
}

// Intersection Find the intersection of two sets
func (hs *HashSet) Intersection(a *HashSet) *HashSet {
	b := make(map[T]bool)

	for k := range hs.hash {
		if _, ok := a.hash[k]; ok {
			b[k] = true
		}
	}

	return &HashSet{b}
}

// String print the set to string
func (hs *HashSet) String() string {
	return fmt.Sprintf("%v", hs.hash)
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

func (hs *HashSet) addJSONArrayItem(v T) jsonArray {
	hs.Add(v)
	return hs
}

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(hs)
func (hs *HashSet) MarshalJSON() ([]byte, error) {
	return jsonMarshalArray(hs)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, hs)
func (hs *HashSet) UnmarshalJSON(data []byte) error {
	hs.Clear()
	return jsonUnmarshalArray(data, hs)
}
