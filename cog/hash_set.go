//go:build go1.18
// +build go1.18

package cog

import (
	"fmt"

	"github.com/askasoft/pango/asg"
)

// NewHashSet Create a new hash set
func NewHashSet[T comparable](vs ...T) *HashSet[T] {
	hs := &HashSet[T]{}
	hs.Adds(vs...)
	return hs
}

// HashSet an unordered collection of unique values.
// The zero value for HashSet is an empty set ready to use.
// http://en.wikipedia.org/wiki/Set_(computer_science%29)
type HashSet[T comparable] struct {
	hash map[T]struct{}
}

// lazyInit lazily initializes a zero List value.
func (hs *HashSet[T]) lazyInit() {
	if hs.hash == nil {
		hs.hash = make(map[T]struct{})
	}
}

//-----------------------------------------------------------
// implements Collection interface

// Len Return the number of items in the set
func (hs *HashSet[T]) Len() int {
	return len(hs.hash)
}

// IsEmpty returns true if the set's length == 0
func (hs *HashSet[T]) IsEmpty() bool {
	return hs.Len() == 0
}

// Clear clears the hash set.
func (hs *HashSet[T]) Clear() {
	hs.hash = nil
}

// Add Add the item v to the set
func (hs *HashSet[T]) Add(v T) {
	hs.lazyInit()
	hs.hash[v] = struct{}{}
}

// Adds Adds all items of vs to the set
func (hs *HashSet[T]) Adds(vs ...T) {
	if len(vs) == 0 {
		return
	}

	hs.lazyInit()
	for _, v := range vs {
		hs.hash[v] = struct{}{}
	}
}

// AddCol adds all items of another collection
func (hs *HashSet[T]) AddCol(ac Collection[T]) {
	if ac.IsEmpty() || hs == ac {
		return
	}

	hs.lazyInit()
	if ic, ok := ac.(Iterable[T]); ok {
		it := ic.Iterator()
		for it.Next() {
			hs.hash[it.Value()] = struct{}{}
		}
		return
	}

	hs.Adds(ac.Values()...)
}

// Remove remove all items with associated value v of vs
func (hs *HashSet[T]) Remove(v T) {
	delete(hs.hash, v)
}

// Removes delete items of vs
func (hs *HashSet[T]) Removes(vs ...T) {
	if hs.IsEmpty() {
		return
	}

	for _, v := range vs {
		delete(hs.hash, v)
	}
}

// RemoveCol remove all of this collection's elements that are also contained in the specified collection
func (hs *HashSet[T]) RemoveCol(ac Collection[T]) {
	if hs == ac {
		hs.Clear()
		return
	}

	if ic, ok := ac.(Iterable[T]); ok {
		hs.RemoveIter(ic.Iterator())
		return
	}

	hs.Removes(ac.Values()...)
}

// RemoveIter remove all items in the iterator it
func (hs *HashSet[T]) RemoveIter(it Iterator[T]) {
	for it.Next() {
		delete(hs.hash, it.Value())
	}
}

// RemoveFunc remove all items that function f returns true
func (hs *HashSet[T]) RemoveFunc(f func(T) bool) {
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
func (hs *HashSet[T]) Contain(v T) bool {
	if hs.IsEmpty() {
		return false
	}
	if _, ok := hs.hash[v]; ok {
		return true
	}
	return false
}

// Contains Test to see if the collection contains all items of vs
func (hs *HashSet[T]) Contains(vs ...T) bool {
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
func (hs *HashSet[T]) ContainCol(ac Collection[T]) bool {
	if ac.IsEmpty() || hs == ac {
		return true
	}

	if hs.IsEmpty() {
		return false
	}

	if ic, ok := ac.(Iterable[T]); ok {
		return hs.ContainIter(ic.Iterator())
	}

	return hs.Contains(ac.Values()...)
}

// ContainIter Test to see if the collection contains all items of iterator 'it'
func (hs *HashSet[T]) ContainIter(it Iterator[T]) bool {
	for it.Next() {
		if _, ok := hs.hash[it.Value()]; !ok {
			return false
		}
	}
	return true
}

// Retains Retains only the elements in this collection that are contained in the argument array vs.
func (hs *HashSet[T]) Retains(vs ...T) {
	if hs.IsEmpty() {
		return
	}

	if len(vs) == 0 {
		hs.Clear()
		return
	}

	for k := range hs.hash {
		if !asg.Contains(vs, k) {
			delete(hs.hash, k)
		}
	}
}

// RetainCol Retains only the elements in this collection that are contained in the specified collection.
func (hs *HashSet[T]) RetainCol(ac Collection[T]) {
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
func (hs *HashSet[T]) RetainFunc(f func(T) bool) {
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
func (hs *HashSet[T]) Values() []T {
	vs := make([]T, hs.Len())
	i := 0
	for k := range hs.hash {
		vs[i] = k
		i++
	}
	return vs
}

// Each Call f for each item in the set
func (hs *HashSet[T]) Each(f func(T)) {
	for k := range hs.hash {
		f(k)
	}
}

//-----------------------------------------------------------

// Difference Find the difference btween two sets
func (hs *HashSet[T]) Difference(a *HashSet[T]) *HashSet[T] {
	b := make(map[T]struct{})

	for k := range hs.hash {
		if _, ok := a.hash[k]; !ok {
			b[k] = struct{}{}
		}
	}

	return &HashSet[T]{b}
}

// Intersection Find the intersection of two sets
func (hs *HashSet[T]) Intersection(a *HashSet[T]) *HashSet[T] {
	b := make(map[T]struct{})

	for k := range hs.hash {
		if _, ok := a.hash[k]; ok {
			b[k] = struct{}{}
		}
	}

	return &HashSet[T]{b}
}

// String print the set to string
func (hs *HashSet[T]) String() string {
	return fmt.Sprintf("%v", hs.hash)
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(hs)
func (hs *HashSet[T]) MarshalJSON() ([]byte, error) {
	return jsonMarshalCol[T](hs)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, hs)
func (hs *HashSet[T]) UnmarshalJSON(data []byte) error {
	hs.Clear()
	return jsonUnmarshalCol[T](data, hs)
}
