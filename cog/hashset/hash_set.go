package hashset

import (
	"iter"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/cog"
	"github.com/askasoft/pango/cog/internal/icol"
	"github.com/askasoft/pango/doc/jsonx"
)

// NewHashSet Create a new hash set
func NewHashSet[T comparable](vs ...T) *HashSet[T] {
	hs := &HashSet[T]{}
	hs.AddAll(vs...)
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
	clear(hs.hash)
}

// Add Add the item v to the set
func (hs *HashSet[T]) Add(v T) {
	hs.lazyInit()
	hs.hash[v] = struct{}{}
}

// AddAll AddAll all items of vs to the set
func (hs *HashSet[T]) AddAll(vs ...T) {
	if len(vs) == 0 {
		return
	}

	hs.lazyInit()
	for _, v := range vs {
		hs.hash[v] = struct{}{}
	}
}

// AddCol adds all items of another collection
func (hs *HashSet[T]) AddCol(ac cog.Collection[T]) {
	if ac.IsEmpty() || hs == ac {
		return
	}

	hs.lazyInit()
	if ic, ok := ac.(cog.Iterable[T]); ok {
		it := ic.Iterator()
		for it.Next() {
			hs.hash[it.Value()] = struct{}{}
		}
		return
	}

	hs.AddAll(ac.Values()...)
}

// Remove remove all items with associated value v
func (hs *HashSet[T]) Remove(v T) {
	delete(hs.hash, v)
}

// RemoveAll delete items of vs
func (hs *HashSet[T]) RemoveAll(vs ...T) {
	if hs.IsEmpty() {
		return
	}

	for _, v := range vs {
		delete(hs.hash, v)
	}
}

// RemoveCol remove all of this collection's elements that are also contained in the specified collection
func (hs *HashSet[T]) RemoveCol(ac cog.Collection[T]) {
	if hs == ac {
		hs.Clear()
		return
	}

	if ic, ok := ac.(cog.Iterable[T]); ok {
		hs.RemoveIter(ic.Iterator())
		return
	}

	hs.RemoveAll(ac.Values()...)
}

// RemoveIter remove all items in the iterator it
func (hs *HashSet[T]) RemoveIter(it cog.Iterator[T]) {
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

// Contains Test to see if the list contains the value v
func (hs *HashSet[T]) Contains(v T) bool {
	if hs.IsEmpty() {
		return false
	}
	_, ok := hs.hash[v]
	return ok
}

// ContainsAll Test to see if the collection contains any item of vs
func (hs *HashSet[T]) ContainsAny(vs ...T) bool {
	if len(vs) == 0 {
		return true
	}

	if hs.IsEmpty() {
		return false
	}

	for _, v := range vs {
		if _, ok := hs.hash[v]; ok {
			return true
		}
	}
	return false
}

// ContainsAll Test to see if the collection contains all items of vs
func (hs *HashSet[T]) ContainsAll(vs ...T) bool {
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

// ContainsCol Test to see if the collection contains all items of another collection
func (hs *HashSet[T]) ContainsCol(ac cog.Collection[T]) bool {
	return icol.ContainsCol(hs, ac)
}

// ContainsIter Test to see if the collection contains all items of iterator 'it'
func (hs *HashSet[T]) ContainsIter(it cog.Iterator[T]) bool {
	for it.Next() {
		if _, ok := hs.hash[it.Value()]; !ok {
			return false
		}
	}
	return true
}

// RetainAll Retains only the elements in this collection that are contained in the argument array vs.
func (hs *HashSet[T]) RetainAll(vs ...T) {
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
func (hs *HashSet[T]) RetainCol(ac cog.Collection[T]) {
	if hs.IsEmpty() || hs == ac {
		return
	}

	if ac.IsEmpty() {
		hs.Clear()
		return
	}

	hs.RetainFunc(ac.Contains)
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
func (hs *HashSet[T]) Each(f func(int, T) bool) {
	i := 0
	for k := range hs.hash {
		if !f(i, k) {
			return
		}
		i++
	}
}

// Seq returns a iter.Seq[T] for range
func (hs *HashSet[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		for k := range hs.hash {
			if !yield(k) {
				return
			}
		}
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
	return jsonx.Stringify(hs)
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(hs)
func (hs *HashSet[T]) MarshalJSON() ([]byte, error) {
	return icol.JsonMarshalCol(hs)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, hs)
func (hs *HashSet[T]) UnmarshalJSON(data []byte) error {
	hs.Clear()
	return icol.JsonUnmarshalCol(data, hs)
}
