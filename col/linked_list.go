package col

import (
	"encoding/json"
	"sort"

	"github.com/pandafw/pango/cmp"
)

// NewLinkedList returns an initialized list.
// Example: NewLinkedList(1, 2, 3)
func NewLinkedList(vs ...interface{}) *LinkedList {
	ll := &LinkedList{}
	ll.Add(vs...)
	return ll
}

// LinkedList implements a doubly linked list.
// The zero value for LinkedList is an empty list ready to use.
//
// To iterate over a list (where ll is a *LinkedList):
//	for li := ll.Front(); li != nil; li = li.Next() {
//		// do something with li.Value()
//	}
//
type LinkedList struct {
	root LinkedListItem
	len  int
}

// lazyInit lazily initializes a zero LinkedList value.
func (ll *LinkedList) lazyInit() {
	if ll.root.next == nil {
		ll.root.list = ll
		ll.root.next = &ll.root
		ll.root.prev = &ll.root
		ll.len = 0
	}
}

// insertAfter inserts item with value v after item at, increments l.len, and returns v's LinkedListItem.
func (ll *LinkedList) insertAfter(at *LinkedListItem, v interface{}) *LinkedListItem {
	li := &LinkedListItem{value: v}
	li.insertAfter(at)
	return li
}

// insertBefore inserts item with value v before item at, increments l.len, and returns v's LinkedListItem.
func (ll *LinkedList) insertBefore(at *LinkedListItem, v interface{}) *LinkedListItem {
	li := &LinkedListItem{value: v}
	li.insertAfter(at.prev)
	return li
}

//-----------------------------------------------------------
// implements Collection interface

// Len returns the length of the list.
func (ll *LinkedList) Len() int {
	return ll.len
}

// IsEmpty returns true if the list length == 0
func (ll *LinkedList) IsEmpty() bool {
	return ll.len == 0
}

// Clear clears list ll.
func (ll *LinkedList) Clear() {
	ll.root.list = nil
	ll.root.next = nil
	ll.root.prev = nil
	ll.len = 0
}

// Add adds all items of vs and returns the last added item.
func (ll *LinkedList) Add(vs ...interface{}) {
	ll.PushBack(vs...)
}

// AddAll adds all items of another collection
func (ll *LinkedList) AddAll(ac Collection) {
	ll.PushBackAll(ac)
}

func (ll *LinkedList) deleteAll(v interface{}) {
	for li := ll.Front(); li != nil; li = li.Next() {
		if li.Value() == v {
			li.Remove()
		}
	}
}

// Delete delete all items with associated value v of vs
func (ll *LinkedList) Delete(vs ...interface{}) {
	for _, v := range vs {
		ll.deleteAll(v)
	}
}

// DeleteAll delete all of this collection's elements that are also contained in the specified collection
func (ll *LinkedList) DeleteAll(ac Collection) {
	if ac.IsEmpty() {
		return
	}

	if ll == ac {
		ll.Clear()
		return
	}

	if ic, ok := ac.(Iterable); ok {
		it := ic.Iterator()
		for it.Next() {
			v := it.Value()
			ll.deleteAll(v)
		}
		return
	}

	ll.Delete(ac.Values()...)
}

// Contains Test to see if the collection contains all items of vs
func (ll *LinkedList) Contains(vs ...interface{}) bool {
	if len(vs) == 0 {
		return true
	}

	if ll.IsEmpty() {
		return false
	}

	for _, v := range vs {
		if ll.Index(v) < 0 {
			return false
		}
	}
	return true
}

// ContainsAll Test to see if the collection contains all items of another collection
func (ll *LinkedList) ContainsAll(ac Collection) bool {
	if ll == ac {
		return true
	}

	if ic, ok := ac.(Iterable); ok {
		it := ic.Iterator()
		for it.Next() {
			if ll.Index(it.Value()) < 0 {
				return false
			}
		}
		return true
	}

	return ll.Contains(ac.Values()...)
}

// Retain Retains only the elements in this collection that are contained in the argument array vs.
func (ll *LinkedList) Retain(vs ...interface{}) {
	if ll.IsEmpty() || len(vs) == 0 {
		return
	}

	ll.RetainAll(NewArrayList(vs...))
}

// RetainAll Retains only the elements in this collection that are contained in the specified collection.
func (ll *LinkedList) RetainAll(ac Collection) {
	if ll.IsEmpty() || ac.IsEmpty() || ll == ac {
		return
	}

	for li := ll.Front(); li != nil; li = li.Next() {
		if !ac.Contains(li.Value()) {
			li.Remove()
		}
	}
}

// Values returns a slice contains all the items of the list ll
func (ll *LinkedList) Values() []interface{} {
	vs := make([]interface{}, ll.Len())
	for i, li := 0, ll.Front(); li != nil; i, li = i+1, li.Next() {
		vs[i] = li.Value()
	}
	return vs
}

// Each call f for each item in the list
func (ll *LinkedList) Each(f func(interface{})) {
	for li := ll.Front(); li != nil; li = li.Next() {
		f(li.Value())
	}
}

//-----------------------------------------------------------
// implements List interface

// Get returns the element at the specified position in this list
func (ll *LinkedList) Get(index int) (interface{}, bool) {
	li := ll.Item(index)
	if li == nil {
		return nil, false
	}
	return li.Value(), true
}

// Set set the v at the specified index in this list and returns the old value.
func (ll *LinkedList) Set(index int, v interface{}) (ov interface{}) {
	li := ll.Item(index)
	if li != nil {
		ov, li.value = li.value, v
	}
	return
}

// Insert inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Does not do anything if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (ll *LinkedList) Insert(index int, vs ...interface{}) {
	n := len(vs)
	if n == 0 {
		return
	}

	len := ll.Len()
	if index < -len || index > len {
		return
	}

	if index < 0 {
		index += len
	}

	if index == len {
		// Append
		ll.Add(vs...)
		return
	}

	li := ll.Item(index)
	if li != nil {
		ll.InsertBefore(li, vs...)
	}
}

// InsertAll inserts values of another collection ac at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Does not do anything if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (ll *LinkedList) InsertAll(index int, ac Collection) {
	n := ac.Len()
	if n == 0 {
		return
	}

	len := ll.Len()
	if index < -len || index > len {
		return
	}

	if index < 0 {
		index += len
	}

	if index == len {
		// Append
		ll.AddAll(ac)
		return
	}

	li := ll.Item(index)
	if li != nil {
		ll.InsertAllBefore(li, ac)
	}
}

// Index returns the index of the first occurrence of the specified v in this list, or -1 if this list does not contain v.
func (ll *LinkedList) Index(v interface{}) int {
	i, _ := ll.Search(v)
	return i
}

// Remove removes the item li from ll if li is an item of list ll.
// It returns the item's value.
// The item li must not be nil.
func (ll *LinkedList) Remove(index int) {
	li := ll.Item(index)
	if li != nil {
		li.Remove()
	}
}

// Swap swaps values of two items at the given index.
func (ll *LinkedList) Swap(i, j int) {
	if i == j {
		return
	}

	ii := ll.Item(i)
	ij := ll.Item(j)
	if ii != nil && ij != nil && ii != ij {
		ii.value, ij.value = ij.value, ii.value
	}
}

// ReverseEach call f for each item in the list with reverse order
func (ll *LinkedList) ReverseEach(f func(interface{})) {
	for li := ll.Back(); li != nil; li = li.Prev() {
		f(li.Value())
	}
}

// Iterator returns a iterator for the list
func (ll *LinkedList) Iterator() Iterator {
	return &linkedListItemIterator{ll, &ll.root}
}

//--------------------------------------------------------------------

// Front returns the first item of list ll or nil if the list is empty.
func (ll *LinkedList) Front() *LinkedListItem {
	if ll.len == 0 {
		return nil
	}
	return ll.root.next
}

// Back returns the last item of list ll or nil if the list is empty.
func (ll *LinkedList) Back() *LinkedListItem {
	if ll.len == 0 {
		return nil
	}
	return ll.root.prev
}

// Item returns the item at the specified index i.
// if i < -ll.Len() or i >= ll.Len(), returns nil
// if i < 0, returns ll.Item(ll.Len() + i)
func (ll *LinkedList) Item(i int) *LinkedListItem {
	if i < -ll.len || i >= ll.len {
		return nil
	}

	if i < 0 {
		i += ll.len
	}
	if i >= ll.len/2 {
		return ll.Back().Offset(i + 1 - ll.len)
	}

	return ll.Front().Offset(i)
}

// Search linear search v
// returns (index, item) if it's value is v
// if not found, returns (-1, nil)
func (ll *LinkedList) Search(v interface{}) (int, *LinkedListItem) {
	for i, li := 0, ll.Front(); li != nil; li = li.Next() {
		if li.Value() == v {
			return i, li
		}
		i++
	}
	return -1, nil
}

// Sort Sorts this list according to the order induced by the specified Comparator.
func (ll *LinkedList) Sort(less cmp.Less) {
	if ll.Len() < 2 {
		return
	}
	sort.Sort(&sorter{ll, less})
}

// PushBack inserts all items of vs at the back of list ll.
// returns the last inserted item.
func (ll *LinkedList) PushBack(vs ...interface{}) *LinkedListItem {
	if len(vs) == 0 {
		return nil
	}

	ll.lazyInit()
	return ll.InsertBefore(&ll.root, vs...)
}

// PushBackAll inserts a copy of another collection at the back of list ll.
// The ll and ac may be the same. They must not be nil.
func (ll *LinkedList) PushBackAll(ac Collection) *LinkedListItem {
	if ac.IsEmpty() {
		return nil
	}

	ll.lazyInit()
	return ll.InsertAllBefore(&ll.root, ac)
}

// PushFront inserts all items of vs at the front of list ll.
// returns the last inserted item.
func (ll *LinkedList) PushFront(vs ...interface{}) *LinkedListItem {
	if len(vs) == 0 {
		return nil
	}

	ll.lazyInit()
	return ll.InsertAfter(&ll.root, vs...)
}

// PushFrontAll inserts a copy of another collection at the front of list ll.
// The ll and ac may be the same. They must not be nil.
func (ll *LinkedList) PushFrontAll(ac Collection) *LinkedListItem {
	if ac.IsEmpty() {
		return nil
	}

	ll.lazyInit()
	return ll.InsertAllAfter(&ll.root, ac)
}

// InsertAfter inserts all items of vs immediately after the item at and returns the last inserted item li.
// If at is not an item of ll, the list is not modified.
// The at must not be nil.
func (ll *LinkedList) InsertAfter(at *LinkedListItem, vs ...interface{}) *LinkedListItem {
	if at.list != ll || len(vs) == 0 {
		return nil
	}

	li := at
	for _, v := range vs {
		li = ll.insertAfter(li, v)
	}
	return li
}

// InsertAllAfter inserts all items of another collection ac immediately after the item at and returns the last inserted item li.
// If at is not an item of ll, the list is not modified.
// The at must not be nil.
func (ll *LinkedList) InsertAllAfter(at *LinkedListItem, ac Collection) *LinkedListItem {
	if at.list != ll || ac.IsEmpty() {
		return nil
	}

	if ic, ok := ac.(Iterable); ok {
		return ll.insertIterAfter(at, ic.Iterator(), ac.Len())
	}
	return ll.InsertAfter(at, ac.Values()...)
}

// insertIterAfter inserts max n items of iterator it immediately after the item at and returns the last inserted item li.
// If at is not an item of ll, the list is not modified.
// The at must not be nil.
func (ll *LinkedList) insertIterAfter(at *LinkedListItem, it Iterator, n int) (li *LinkedListItem) {
	if at.list != ll || n < 1 {
		return
	}

	for ; it.Next() && n > 0; n-- {
		li = ll.insertAfter(at, it.Value())
		at = li
	}
	return
}

// InsertBefore inserts all items of vs immediately before the item at and returns the last inserted item li.
// If at is not an item of ll, the list is not modified.
// The at must not be nil.
func (ll *LinkedList) InsertBefore(at *LinkedListItem, vs ...interface{}) (li *LinkedListItem) {
	if at.list != ll || len(vs) == 0 {
		return
	}

	for _, v := range vs {
		li = ll.insertBefore(at, v)
	}
	return
}

// InsertAllBefore inserts all items of another collection immediately before the item at and returns the last inserted item li.
// If at is not an item of ll, the list is not modified.
// The at must not be nil.
func (ll *LinkedList) InsertAllBefore(at *LinkedListItem, ac Collection) *LinkedListItem {
	if at.list != ll || ac.IsEmpty() {
		return nil
	}

	if ic, ok := ac.(Iterable); ok {
		return ll.insertIterBefore(at, ic.Iterator(), ac.Len())
	}

	return ll.InsertBefore(at, ac.Values()...)
}

// insertIterBefore inserts max n items's value before the item at and returns the last inserted item li.
// If at is not an item of ll, the list is not modified.
// The at must not be nil.
func (ll *LinkedList) insertIterBefore(at *LinkedListItem, it Iterator, n int) (li *LinkedListItem) {
	if at.list != ll || n < 1 {
		return
	}

	ll.lazyInit()

	for ; it.Next() && n > 0; n-- {
		li = ll.insertBefore(at, it.Value())
	}
	return
}

// MoveToFront moves item li to the front of list ll.
// If li is not an item of ll, the list is not modified.
// The item li must not be nil.
// Returns true if list is modified.
func (ll *LinkedList) MoveToFront(li *LinkedListItem) bool {
	if li.list != ll || ll.root.next == li {
		return false
	}

	li.moveAfter(&ll.root)
	return true
}

// MoveToBack moves item li to the back of list ll.
// If li is not an item of ll, the list is not modified.
// The item li must not be nil.
// Returns true if list is modified.
func (ll *LinkedList) MoveToBack(li *LinkedListItem) bool {
	if li.list != ll || ll.root.prev == li {
		return false
	}

	li.moveAfter(ll.root.prev)
	return true
}

// MoveBefore moves item li to its new position before at.
// If li or at is not an item of ll, or li == at, the list is not modified.
// The item li and at must not be nil.
// Returns true if list is modified.
func (ll *LinkedList) MoveBefore(at, li *LinkedListItem) bool {
	if li.list != ll || li == at || at.list != ll {
		return false
	}

	li.moveAfter(at.prev)
	return true
}

// MoveAfter moves item li to its new position after at.
// If li or at is not an item of ll, or li == at, the list is not modified.
// The item li and at must not be nil.
// Returns true if list is modified.
func (ll *LinkedList) MoveAfter(at, li *LinkedListItem) bool {
	if li.list != ll || li == at || at.list != ll {
		return false
	}

	li.moveAfter(at)
	return true
}

// SwapItem swap item's value of a, b.
// If a or b is not an item of ll, or a == b, the list is not modified.
// The item a and b must not be nil.
// Returns true if list is modified.
func (ll *LinkedList) SwapItem(a, b *LinkedListItem) bool {
	if a.list != ll || a == b || b.list != ll {
		return false
	}

	a.value, b.value = b.value, a.value
	return true
}

// String print list to string
func (ll *LinkedList) String() string {
	bs, _ := json.Marshal(ll)
	return string(bs)
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

func newJSONArrayLinkedList() jsonArray {
	return NewLinkedList()
}

func (ll *LinkedList) addJSONArrayItem(v interface{}) jsonArray {
	ll.Add(v)
	return ll
}

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(ll)
func (ll *LinkedList) MarshalJSON() (res []byte, err error) {
	return jsonMarshalList(ll)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, ll)
func (ll *LinkedList) UnmarshalJSON(data []byte) error {
	ll.Clear()
	ju := &jsonUnmarshaler{
		newArray:  newJSONArrayLinkedList,
		newObject: newJSONObject,
	}
	return ju.unmarshalJSONArray(data, ll)
}
