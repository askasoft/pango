package col

import (
	"encoding/json"
	"fmt"

	"github.com/pandafw/pango/cmp"
)

// NewSortedList returns an initialized sorted list.
// Example: NewSortedList(cmp.LessInt, 1, 2, 3)
func NewSortedList(less cmp.Less, vs ...interface{}) *SortedList {
	sl := &SortedList{
		less: less,
	}
	sl.Add(vs...)
	return sl
}

// SortedList implements a doubly linked sorted list.
type SortedList struct {
	front, back *sortedListNode
	less        cmp.Less
	len         int
}

//-----------------------------------------------------------
// implements Collection interface

// Len returns the length of the list.
func (sl *SortedList) Len() int {
	return sl.len
}

// IsEmpty checks if the list is empty.
func (sl *SortedList) IsEmpty() bool {
	return sl.len == 0
}

// Clear clears list sl.
func (sl *SortedList) Clear() {
	sl.front = nil
	sl.back = nil
	sl.len = 0
}

// Add adds all items of vs and returns the last added item.
func (sl *SortedList) Add(vs ...interface{}) {
	if len(vs) == 0 {
		return
	}

	for _, v := range vs {
		sl.add(v)
	}
}

// AddAll adds all items of another collection
func (sl *SortedList) AddAll(ac Collection) {
	sl.Add(ac.Values()...)
}

// Delete delete all items with associated value v of vs
func (sl *SortedList) Delete(vs ...interface{}) {
	if sl.IsEmpty() {
		return
	}

	for _, v := range vs {
		sl.deleteAll(v)
	}
}

// DeleteAll delete all of this collection's elements that are also contained in the specified collection
func (sl *SortedList) DeleteAll(ac Collection) {
	if sl.IsEmpty() || ac.IsEmpty() {
		return
	}

	if sl == ac {
		sl.Clear()
		return
	}

	if ic, ok := ac.(Iterable); ok {
		it := ic.Iterator()
		for it.Next() {
			sl.deleteAll(it.Value())
		}
		return
	}

	sl.Delete(ac.Values()...)
}

// Contains Test to see if the collection contains all items of vs
func (sl *SortedList) Contains(vs ...interface{}) bool {
	if len(vs) == 0 {
		return true
	}

	if sl.IsEmpty() {
		return false
	}

	for _, v := range vs {
		if sl.Index(v) < 0 {
			return false
		}
	}
	return true
}

// ContainsAll Test to see if the collection contains all items of another collection
func (sl *SortedList) ContainsAll(ac Collection) bool {
	if sl == ac || ac.IsEmpty() {
		return true
	}

	if sl.IsEmpty() {
		return false
	}

	if ic, ok := ac.(Iterable); ok {
		it := ic.Iterator()
		for it.Next() {
			if sl.Index(it.Value()) < 0 {
				return false
			}
		}
		return true
	}

	return sl.Contains(ac.Values()...)
}

// Retain Retains only the elements in this collection that are contained in the argument array vs.
func (sl *SortedList) Retain(vs ...interface{}) {
	if sl.IsEmpty() || len(vs) == 0 {
		return
	}

	sl.RetainAll(NewArrayList(vs...))
}

// RetainAll Retains only the elements in this collection that are contained in the specified collection.
func (sl *SortedList) RetainAll(ac Collection) {
	if sl.IsEmpty() || ac.IsEmpty() || sl == ac {
		return
	}

	for ln := sl.front; ln != nil; ln = ln.next {
		if !ac.Contains(ln.value) {
			sl.deleteNode(ln)
		}
	}
}

// Values returns a slice contains all the items of the list l
func (sl *SortedList) Values() []interface{} {
	vs := make([]interface{}, sl.Len())
	for i, ln := 0, sl.front; ln != nil; i, ln = i+1, ln.next {
		vs[i] = ln.value
	}
	return vs
}

// Each call f for each item in the list
func (sl *SortedList) Each(f func(interface{})) {
	for ln := sl.front; ln != nil; ln = ln.next {
		f(ln.value)
	}
}

//-----------------------------------------------------------
// implements List interface

// Get returns the element at the specified position in this list
// if i < -sl.Len() or i >= sl.Len(), panic
// if i < 0, returns sl.Get(sl.Len() + i)
func (sl *SortedList) Get(index int) interface{} {
	index = sl.checkItemIndex(index)

	return sl.node(index).value
}

// Set set the v at the specified index in this list and returns the old value.
// Just implements the List.Set() method
// Same as sl.Delete(sl.Get(index)), sl.Add(v)
func (sl *SortedList) Set(index int, v interface{}) (ov interface{}) {
	index = sl.checkItemIndex(index)

	ln := sl.node(index)
	ov = ln.value
	sl.setValue(ln, v)
	return
}

// Insert inserts values, same as Add(...)
// Just implements the List.Insert() method
// Panic if position is bigger than list's size
func (sl *SortedList) Insert(index int, vs ...interface{}) {
	sl.checkSizeIndex(index)
	sl.Add(vs...)
}

// InsertAll inserts values, same as AddAll(...)
// Just implements the List.InsertAll() method
// Panic if position is bigger than list's size
func (sl *SortedList) InsertAll(index int, ac Collection) {
	sl.checkSizeIndex(index)
	sl.AddAll(ac)
}

// Index returns the index of the first occurrence of the specified v in this list, or -1 if this list does not contain v.
func (sl *SortedList) Index(v interface{}) int {
	i, ln := sl.binarySearch(v)
	if ln != nil && ln.value == v {
		return i
	}
	return -1
}

// Remove removes the item li from l if li is an item of list l.
// It returns the item's value.
// The item li must not be nil.
func (sl *SortedList) Remove(index int) {
	index = sl.checkItemIndex(index)

	ln := sl.node(index)
	sl.deleteNode(ln)
}

// Swap swaps values of two items at the given index.
// Do nothing because all items are sorted.
func (sl *SortedList) Swap(i, j int) {
	// do nothing
}

// ReverseEach Call f for each item in the list with reverse order
func (sl *SortedList) ReverseEach(f func(interface{})) {
	for ln := sl.back; ln != nil; ln = ln.prev {
		f(ln.value)
	}
}

// Iterator returns a iterator for the list
func (sl *SortedList) Iterator() Iterator {
	return &sortedListIterator{list: sl}
}

//--------------------------------------------------------------------

// Front returns the first item of list l or nil if the list is empty.
func (sl *SortedList) Front() interface{} {
	if sl.front == nil {
		return nil
	}
	return sl.front.value
}

// Back returns the last item of list l or nil if the list is empty.
func (sl *SortedList) Back() interface{} {
	if sl.back == nil {
		return nil
	}
	return sl.back.value
}

// PopFront remove the first item of list.
func (sl *SortedList) PopFront() (v interface{}) {
	if sl.front != nil {
		v = sl.front.value
		sl.deleteNode(sl.front)
	}
	return
}

// PopBack remove the last item of list.
func (sl *SortedList) PopBack() (v interface{}) {
	if sl.back != nil {
		v = sl.back.value
		sl.deleteNode(sl.back)
	}
	return
}

// String print list to string
func (sl *SortedList) String() string {
	bs, _ := json.Marshal(sl)
	return string(bs)
}

//-----------------------------------------------------------

func (sl *SortedList) checkItemIndex(index int) int {
	if index >= sl.len || index < -sl.len {
		panic(fmt.Sprintf("SortedList out of bounds: index=%d, len=%d", index, sl.len))
	}

	if index < 0 {
		index += sl.len
	}
	return index
}

func (sl *SortedList) checkSizeIndex(index int) int {
	if index > sl.len || index < -sl.len {
		panic(fmt.Sprintf("SortedList out of bounds: index=%d, len=%d", index, sl.len))
	}

	if index < 0 {
		index += sl.len
	}
	return index
}

// node returns the node at the specified index i.
func (sl *SortedList) node(i int) *sortedListNode {
	if i < (sl.len >> 1) {
		ln := sl.front
		for ; i > 0; i-- {
			ln = ln.next
		}
		return ln
	}

	ln := sl.back
	for i = sl.len - i - 1; i > 0; i-- {
		ln = ln.prev
	}
	return ln
}

// add inserts a new item with value v.
func (sl *SortedList) add(v interface{}) (ln *sortedListNode) {
	ln = &sortedListNode{value: v}

	if sl.IsEmpty() {
		// Assert key is of less's type for initial list
		sl.less(v, v)

		sl.front, sl.back = ln, ln
		sl.len++
		return
	}

	_, at := sl.binarySearch(v)
	if at != nil {
		ln.next = at
		ln.prev = at.prev
		at.prev = ln
		if ln.prev == nil {
			sl.front = ln
		} else {
			ln.prev.next = ln
		}
		sl.len++
		return
	}

	ln.prev = sl.back
	ln.prev.next = ln
	sl.back = ln
	sl.len++
	return
}

func (sl *SortedList) setValue(ln *sortedListNode, v interface{}) *sortedListNode {
	if ln.value == v {
		return ln
	}

	// delete and insert again
	sl.deleteNode(ln)
	return sl.add(v)
}

func (sl *SortedList) deleteAll(v interface{}) {
	_, ln := sl.binarySearch(v)
	for ln != nil && ln.value == v {
		sl.deleteNode(ln)
		ln = ln.next
	}
}

func (sl *SortedList) deleteNode(ln *sortedListNode) {
	if ln.prev == nil {
		sl.front = ln.next
	} else {
		ln.prev.next = ln.next
	}

	if ln.next == nil {
		sl.back = ln.prev
	} else {
		ln.next.prev = ln.prev
	}

	sl.len--
}

// binarySearch binary search v
// returns (index, item) if it's value is >= v
// if not found, returns (-1, nil)
func (sl *SortedList) binarySearch(v interface{}) (int, *sortedListNode) {
	if sl.IsEmpty() {
		return -1, nil
	}

	ln := sl.front
	p, i, j := 0, 0, sl.Len()
	for i < j && ln != nil {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		ln = ln.offset(h - p)
		p = h
		// i ≤ h < j
		if sl.less(ln.value, v) {
			i = h + 1
		} else {
			j = h
		}
	}

	if i < sl.Len() {
		ln = ln.offset(i - p)
		return i, ln
	}
	return -1, nil
}

//-----------------------------------------------------
// sortedListNode is a node of a doublly-linked list.
type sortedListNode struct {
	prev  *sortedListNode
	next  *sortedListNode
	value interface{}
}

// String print the list item to string
func (ln *sortedListNode) String() string {
	return fmt.Sprintf("%v", ln.value)
}

// offset returns the next +n or previous -n list item.
func (ln *sortedListNode) offset(n int) *sortedListNode {
	if n == 0 {
		return ln
	}

	if n > 0 {
		ni := ln
		for ; n > 0; n-- {
			ni = ni.next
		}
		return ni
	}

	pi := ln
	for ; n < 0; n++ {
		pi = pi.prev
	}
	return pi
}

//-----------------------------------------------------
// sortedListIterator a iterator for sortedListNode
type sortedListIterator struct {
	list    *SortedList
	node    *sortedListNode
	removed bool
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (it *sortedListIterator) Prev() bool {
	if it.list.IsEmpty() {
		return false
	}

	if it.node == nil {
		it.node = it.list.back
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
func (it *sortedListIterator) Next() bool {
	if it.list.IsEmpty() {
		return false
	}

	if it.node == nil {
		it.node = it.list.front
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
func (it *sortedListIterator) Value() interface{} {
	if it.node == nil {
		return nil
	}
	return it.node.value
}

// SetValue set the value to the item
// NOTE: Prev()/Next() will change
func (it *sortedListIterator) SetValue(v interface{}) {
	if it.node == nil {
		return
	}

	if it.removed {
		// unlinked item
		it.node.value = v
		return
	}

	it.node = it.list.setValue(it.node, v)
}

// Remove remove the current element
func (it *sortedListIterator) Remove() {
	if it.node == nil {
		return
	}

	if it.removed {
		panic("SortedList can't remove a unlinked item")
	}

	it.list.deleteNode(it.node)
	it.removed = true
}

// Reset resets the iterator to its initial state (one-before-first/one-after-last)
// Call Next()/Prev() to fetch the first/last element if any.
func (it *sortedListIterator) Reset() {
	it.node = nil
	it.removed = false
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

func (sl *SortedList) addJSONArrayItem(v interface{}) jsonArray {
	sl.Add(v)
	return sl
}

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(l)
func (sl *SortedList) MarshalJSON() (res []byte, err error) {
	return jsonMarshalList(sl)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, l)
func (sl *SortedList) UnmarshalJSON(data []byte) error {
	sl.Clear()
	ju := &jsonUnmarshaler{
		newArray:  newJSONArray,
		newObject: newJSONObject,
	}
	return ju.unmarshalJSONArray(data, sl)
}
