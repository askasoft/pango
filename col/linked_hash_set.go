package col

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/pandafw/pango/cmp"
)

// NewLinkedHashSet returns an initialized set.
// Example: NewLinkedHashSet(1, 2, 3)
func NewLinkedHashSet(vs ...T) *LinkedHashSet {
	ls := &LinkedHashSet{}
	ls.Add(vs...)
	return ls
}

// LinkedHashSet implements a doubly linked set.
// The zero value for LinkedHashSet is an empty set ready to use.
// Note that insertion order is not affected if an element is re-inserted into the set.
// (An element e is reinserted into a set s if s.Add(e) is invoked when s.Contains(e) would return true immediately prior to the invocation.)
//
// To iterate over a set (where ls is a *LinkedHashSet):
//	it := ls.Iterator()
//	for it.Next() {
//		// do something with it.Value()
//	}
//
type LinkedHashSet struct {
	front, back *linkedSetNode
	hash        map[T]*linkedSetNode
}

//-----------------------------------------------------------
// implements Collection interface

// Len returns the length of the set.
func (ls *LinkedHashSet) Len() int {
	return len(ls.hash)
}

// IsEmpty returns true if the set length == 0
func (ls *LinkedHashSet) IsEmpty() bool {
	return len(ls.hash) == 0
}

// Clear clears set ls.
func (ls *LinkedHashSet) Clear() {
	ls.hash = nil
	ls.front = nil
	ls.back = nil
}

// Add adds all items of vs and returns the last added item.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) Add(vs ...T) {
	ls.Insert(ls.Len(), vs...)
}

// AddAll adds all items of another collection
// Note: existing item's order will not change.
func (ls *LinkedHashSet) AddAll(ac Collection) {
	ls.InsertAll(ls.Len(), ac)
}

// Delete delete all items with associated value v of vs
func (ls *LinkedHashSet) Delete(vs ...T) {
	if ls.IsEmpty() {
		return
	}

	for _, v := range vs {
		if ln, ok := ls.hash[v]; ok {
			ls.deleteNode(ln)
		}
	}
}

// DeleteAll delete all of this collection's elements that are also contained in the specified collection
func (ls *LinkedHashSet) DeleteAll(ac Collection) {
	if ls.IsEmpty() || ac.IsEmpty() {
		return
	}

	if ls == ac {
		ls.Clear()
		return
	}

	if ic, ok := ac.(Iterable); ok {
		it := ic.Iterator()
		for it.Next() {
			if ln, ok := ls.hash[it.Value()]; ok {
				ls.deleteNode(ln)
			}
		}
		return
	}

	ls.Delete(ac.Values()...)
}

// Contains Test to see if the collection contains all items of vs
func (ls *LinkedHashSet) Contains(vs ...T) bool {
	if len(vs) == 0 {
		return true
	}

	if ls.IsEmpty() {
		return false
	}

	for _, v := range vs {
		if _, ok := ls.hash[v]; !ok {
			return false
		}
	}
	return true
}

// ContainsAll Test to see if the collection contains all items of another collection
func (ls *LinkedHashSet) ContainsAll(ac Collection) bool {
	if ls == ac || ac.IsEmpty() {
		return true
	}

	if ls.IsEmpty() {
		return false
	}

	if ic, ok := ac.(Iterable); ok {
		it := ic.Iterator()
		for it.Next() {
			if _, ok := ls.hash[it.Value()]; !ok {
				return false
			}
		}
		return true
	}

	return ls.Contains(ac.Values()...)
}

// Retain Retains only the elements in this collection that are contained in the argument array vs.
func (ls *LinkedHashSet) Retain(vs ...T) {
	if ls.IsEmpty() || len(vs) == 0 {
		return
	}

	ls.RetainAll(NewArrayList(vs...))
}

// RetainAll Retains only the elements in this collection that are contained in the specified collection.
func (ls *LinkedHashSet) RetainAll(ac Collection) {
	if ls.IsEmpty() || ac.IsEmpty() || ls == ac {
		return
	}

	for ln := ls.front; ln != nil; ln = ln.next {
		if !ac.Contains(ln.value) {
			ls.deleteNode(ln)
		}
	}
}

// Values returns a slice contains all the items of the set ls
func (ls *LinkedHashSet) Values() []T {
	vs := make([]T, ls.Len())
	for i, ln := 0, ls.front; ln != nil; i, ln = i+1, ln.next {
		vs[i] = ln.value
	}
	return vs
}

// Each call f for each item in the set
func (ls *LinkedHashSet) Each(f func(T)) {
	for ln := ls.front; ln != nil; ln = ln.next {
		f(ln.value)
	}
}

// ReverseEach call f for each item in the set with reverse order
func (ls *LinkedHashSet) ReverseEach(f func(T)) {
	for ln := ls.back; ln != nil; ln = ln.prev {
		f(ln.value)
	}
}

// Iterator returns a iterator for the set
func (ls *LinkedHashSet) Iterator() Iterator {
	return &linkedHashSetIterator{lset: ls}
}

//-----------------------------------------------------------
// implements List interface (LinkedSet behave like List)

// Get returns the element at the specified position in this set
// if i < -ls.Len() or i >= ls.Len(), panic
// if i < 0, returns ls.Get(ls.Len() + i)
func (ls *LinkedHashSet) Get(index int) T {
	index = ls.checkItemIndex(index)

	return ls.node(index).value
}

// Set set the v at the specified index in this set and returns the old value.
// Old item at index will be removed.
func (ls *LinkedHashSet) Set(index int, v T) (ov T) {
	index = ls.checkItemIndex(index)

	ln := ls.node(index)
	ov = ln.value
	ls.setValue(ln, v)
	return
}

// SetValue set the value to the node
func (ls *LinkedHashSet) setValue(ln *linkedSetNode, v T) {
	if ln.value == v {
		return
	}

	// delete old item
	delete(ls.hash, ln.value)

	// delete duplicated item
	if dn, ok := ls.hash[v]; ok {
		ls.deleteNode(dn)
	}

	// add new item
	ln.value = v
	ls.hash[v] = ln
}

// Insert inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than set's size
// Note: position equal to set's size is valid, i.e. append.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) Insert(index int, vs ...T) {
	index = ls.checkSizeIndex(index)

	n := len(vs)
	if n == 0 {
		return
	}

	if ls.hash == nil {
		ls.hash = make(map[T]*linkedSetNode)
	}

	var prev, next *linkedSetNode
	if index == ls.Len() {
		next = nil
		prev = ls.back
	} else {
		next = ls.node(index)
		prev = next.prev
	}

	for _, v := range vs {
		if _, ok := ls.hash[v]; ok {
			continue
		}

		nn := &linkedSetNode{prev: prev, value: v, next: nil}
		if prev == nil {
			ls.front = nn
		} else {
			prev.next = nn
		}
		prev = nn
		ls.hash[v] = nn
	}

	if next == nil {
		ls.back = prev
	} else if prev != nil {
		prev.next = next
		next.prev = prev
	}
}

// InsertAll inserts values of another collection ac at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Panic if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) InsertAll(index int, ac Collection) {
	index = ls.checkSizeIndex(index)

	if ac.IsEmpty() || ls == ac {
		return
	}

	if ls.hash == nil {
		ls.hash = make(map[T]*linkedSetNode)
	}

	if ic, ok := ac.(Iterable); ok {
		var prev, next *linkedSetNode
		if index == ls.Len() {
			next = nil
			prev = ls.back
		} else {
			next = ls.node(index)
			prev = next.prev
		}

		it := ic.Iterator()
		for it.Next() {
			v := it.Value()
			if _, ok := ls.hash[v]; ok {
				continue
			}

			nn := &linkedSetNode{prev: prev, value: v, next: nil}
			if prev == nil {
				ls.front = nn
			} else {
				prev.next = nn
			}
			prev = nn
			ls.hash[v] = nn
		}

		if next == nil {
			ls.back = prev
		} else if prev != nil {
			prev.next = next
			next.prev = prev
		}
		return
	}

	ls.Insert(index, ac.Values()...)
}

// Index returns the index of the specified v in this set, or -1 if this set does not contain v.
func (ls *LinkedHashSet) Index(v T) int {
	for i, ln := 0, ls.front; ln != nil; ln = ln.next {
		if ln.value == v {
			return i
		}
		i++
	}
	return -1
}

// Remove removes the element at the specified position in this set.
func (ls *LinkedHashSet) Remove(index int) {
	index = ls.checkItemIndex(index)

	ln := ls.node(index)
	ls.deleteNode(ln)
}

// Swap swaps values of two items at the given index.
func (ls *LinkedHashSet) Swap(i, j int) {
	i = ls.checkItemIndex(i)
	j = ls.checkItemIndex(j)

	if i != j {
		ni, nj := ls.node(i), ls.node(j)
		ni.value, nj.value = nj.value, ni.value
	}
}

// Sort Sorts this set according to the order induced by the specified Comparator.
func (ls *LinkedHashSet) Sort(less cmp.Less) {
	if ls.Len() < 2 {
		return
	}
	sort.Sort(&sorter{ls, less})
}

//--------------------------------------------------------------------

// Front returns the first item of list ls or nil if the list is empty.
func (ls *LinkedHashSet) Front() T {
	if ls.front == nil {
		return nil
	}
	return ls.front.value
}

// Back returns the last item of list ls or nil if the list is empty.
func (ls *LinkedHashSet) Back() T {
	if ls.back == nil {
		return nil
	}
	return ls.back.value
}

// PopFront remove the first item of list.
func (ls *LinkedHashSet) PopFront() (v T) {
	if ls.front != nil {
		v = ls.front.value
		ls.deleteNode(ls.front)
	}
	return
}

// PopBack remove the last item of list.
func (ls *LinkedHashSet) PopBack() (v T) {
	if ls.back != nil {
		v = ls.back.value
		ls.deleteNode(ls.back)
	}
	return
}

// PushFront inserts all items of vs at the front of list ls.
func (ls *LinkedHashSet) PushFront(vs ...T) {
	if len(vs) == 0 {
		return
	}

	ls.Insert(0, vs...)
}

// PushFrontAll inserts a copy of another collection at the front of list ls.
// The ls and ac may be the same. They must not be nil.
func (ls *LinkedHashSet) PushFrontAll(ac Collection) {
	if ac.IsEmpty() || ls == ac {
		return
	}

	ls.InsertAll(0, ac)
}

// PushBack inserts all items of vs at the back of list ls.
func (ls *LinkedHashSet) PushBack(vs ...T) {
	if len(vs) == 0 {
		return
	}

	ls.Insert(ls.Len(), vs...)
}

// PushBackAll inserts a copy of another collection at the back of list ls.
// The ls and ac may be the same. They must not be nil.
func (ls *LinkedHashSet) PushBackAll(ac Collection) {
	if ac.IsEmpty() || ls == ac {
		return
	}

	ls.InsertAll(ls.Len(), ac)
}

// String print list to string
func (ls *LinkedHashSet) String() string {
	bs, _ := json.Marshal(ls)
	return string(bs)
}

//-----------------------------------------------------------
func (ls *LinkedHashSet) deleteNode(ln *linkedSetNode) {
	if ln.prev == nil {
		ls.front = ln.next
	} else {
		ln.prev.next = ln.next
	}

	if ln.next == nil {
		ls.back = ln.prev
	} else {
		ln.next.prev = ln.prev
	}

	delete(ls.hash, ln.value)
}

// node returns the node at the specified index i.
func (ls *LinkedHashSet) node(i int) *linkedSetNode {
	if i < (ls.Len() >> 1) {
		ln := ls.front
		for ; i > 0; i-- {
			ln = ln.next
		}
		return ln
	}

	ln := ls.back
	for i = ls.Len() - i - 1; i > 0; i-- {
		ln = ln.prev
	}
	return ln
}

func (ls *LinkedHashSet) checkItemIndex(index int) int {
	len := ls.Len()
	if index >= len || index < -len {
		panic(fmt.Sprintf("LinkedHashSet out of bounds: index=%d, len=%d", index, len))
	}

	if index < 0 {
		index += len
	}
	return index
}

func (ls *LinkedHashSet) checkSizeIndex(index int) int {
	len := ls.Len()
	if index > len || index < -len {
		panic(fmt.Sprintf("LinkedHashSet out of bounds: index=%d, len=%d", index, len))
	}

	if index < 0 {
		index += len
	}
	return index
}

//-----------------------------------------------------
// linkedSetNode is a node of a doublly-linked list.
type linkedSetNode struct {
	prev  *linkedSetNode
	next  *linkedSetNode
	value T
}

// String print the list item to string
func (ln *linkedSetNode) String() string {
	return fmt.Sprintf("%v", ln.value)
}

// linkedHashSetIterator a iterator for linkedSetNode
type linkedHashSetIterator struct {
	lset    *LinkedHashSet
	node    *linkedSetNode
	removed bool
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (it *linkedHashSetIterator) Prev() bool {
	if it.lset.IsEmpty() {
		return false
	}

	if it.node == nil {
		it.node = it.lset.back
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
func (it *linkedHashSetIterator) Next() bool {
	if it.lset.IsEmpty() {
		return false
	}

	if it.node == nil {
		it.node = it.lset.front
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
func (it *linkedHashSetIterator) Value() T {
	if it.node == nil {
		return nil
	}
	return it.node.value
}

// SetValue set the value to the item
func (it *linkedHashSetIterator) SetValue(v T) {
	if it.node == nil {
		return
	}

	if it.removed {
		// unlinked item
		it.node.value = v
		return
	}

	it.lset.setValue(it.node, v)
}

// Remove remove the current element
func (it *linkedHashSetIterator) Remove() {
	if it.node == nil {
		return
	}

	if it.removed {
		panic("LinkedHashSet can't remove a unlinked item")
	}

	it.lset.deleteNode(it.node)
	it.removed = true
}

// Reset resets the iterator to its initial state (one-before-first/one-after-last)
// Call Next()/Prev() to fetch the first/last element if any.
func (it *linkedHashSetIterator) Reset() {
	it.node = nil
	it.removed = false
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

func (ls *LinkedHashSet) addJSONArrayItem(v T) jsonArray {
	ls.Add(v)
	return ls
}

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(ls)
func (ls *LinkedHashSet) MarshalJSON() (res []byte, err error) {
	return jsonMarshalSet(ls)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, ls)
func (ls *LinkedHashSet) UnmarshalJSON(data []byte) error {
	ls.Clear()
	ju := &jsonUnmarshaler{
		newArray:  newJSONArray,
		newObject: newJSONObject,
	}
	return ju.unmarshalJSONArray(data, ls)
}
