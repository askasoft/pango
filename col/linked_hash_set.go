package col

import (
	"encoding/json"
	"sort"

	"github.com/pandafw/pango/cmp"
)

// NewLinkedHashSet returns an initialized set.
// Example: NewLinkedHashSet(1, 2, 3)
func NewLinkedHashSet(vs ...interface{}) *LinkedHashSet {
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
//	for li := ls.Front(); li != nil; li = li.Next() {
//		// do something with li.Value()
//	}
//
type LinkedHashSet struct {
	hash map[interface{}]*LinkedSetItem
	root LinkedSetItem
}

// lazyInit lazily initializes a zero LinkedHashSet value.
func (ls *LinkedHashSet) lazyInit() {
	if ls.hash == nil {
		ls.hash = make(map[interface{}]*LinkedSetItem)
		ls.root.lset = ls
		ls.root.next = &ls.root
		ls.root.prev = &ls.root
	}
}

// insertAfter inserts item with value v after item at, increments l.len, and returns v's LinkedSetItem.
func (ls *LinkedHashSet) insertAfter(at *LinkedSetItem, v interface{}) (*LinkedSetItem, bool) {
	if li, ok := ls.hash[v]; ok {
		return li, false
	}

	li := &LinkedSetItem{value: v}
	li.insertAfter(at)
	return li, true
}

// insertBefore inserts item with value v before item at, increments l.len, and returns v's LinkedSetItem.
func (ls *LinkedHashSet) insertBefore(at *LinkedSetItem, v interface{}) (*LinkedSetItem, bool) {
	if li, ok := ls.hash[v]; ok {
		return li, false
	}

	li := &LinkedSetItem{value: v}
	li.insertAfter(at.prev)
	return li, true
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
	ls.root.lset = nil
	ls.root.next = nil
	ls.root.prev = nil
}

// Add adds all items of vs and returns the last added item.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) Add(vs ...interface{}) {
	ls.PushBack(vs...)
}

// AddAll adds all items of another collection
// Note: existing item's order will not change.
func (ls *LinkedHashSet) AddAll(ac Collection) {
	ls.PushBackAll(ac)
}

// Delete delete all items with associated value v of vs
func (ls *LinkedHashSet) Delete(vs ...interface{}) {
	if ls.IsEmpty() {
		return
	}

	for _, v := range vs {
		if li, ok := ls.hash[v]; ok {
			li.Remove()
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
			if li, ok := ls.hash[it.Value()]; ok {
				li.Remove()
			}
		}
		return
	}

	ls.Delete(ac.Values()...)
}

// Contains Test to see if the collection contains all items of vs
func (ls *LinkedHashSet) Contains(vs ...interface{}) bool {
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
	if ls == ac {
		return true
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
func (ls *LinkedHashSet) Retain(vs ...interface{}) {
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

	for li := ls.Front(); li != nil; li = li.Next() {
		if !ac.Contains(li.Value()) {
			li.Remove()
		}
	}
}

// Values returns a slice contains all the items of the set ls
func (ls *LinkedHashSet) Values() []interface{} {
	vs := make([]interface{}, ls.Len())
	for i, li := 0, ls.Front(); li != nil; i, li = i+1, li.Next() {
		vs[i] = li.Value()
	}
	return vs
}

// Each call f for each item in the set
func (ls *LinkedHashSet) Each(f func(interface{})) {
	for li := ls.Front(); li != nil; li = li.Next() {
		f(li.Value())
	}
}

//-----------------------------------------------------------
// implements List interface (LinkedSet behave like List)

// Get returns the element at the specified position in this set
func (ls *LinkedHashSet) Get(index int) (interface{}, bool) {
	li := ls.Item(index)
	if li == nil {
		return nil, false
	}
	return li.Value(), true
}

// Set set the v at the specified index in this set and returns the old value.
// Old item with value v will be removed.
func (ls *LinkedHashSet) Set(index int, v interface{}) (ov interface{}) {
	li := ls.Item(index)
	if li != nil {
		ov = li.Value()
		li.SetValue(v)
	}
	return
}

// Insert inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Does not do anything if position is bigger than set's size
// Note: position equal to set's size is valid, i.e. append.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) Insert(index int, vs ...interface{}) {
	n := len(vs)
	if n == 0 {
		return
	}

	len := ls.Len()
	if index < -len || index > len {
		return
	}

	if index < 0 {
		index += len
	}

	if index == len {
		// Append
		ls.Add(vs...)
		return
	}

	li := ls.Item(index)
	if li != nil {
		ls.InsertBefore(li, vs...)
	}
}

// InsertAll inserts values of another collection ac at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Does not do anything if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) InsertAll(index int, ac Collection) {
	n := ac.Len()
	if n == 0 {
		return
	}

	len := ls.Len()
	if index < -len || index > len {
		return
	}

	if index < 0 {
		index += len
	}

	if index == len {
		// Append
		ls.AddAll(ac)
		return
	}

	li := ls.Item(index)
	if li != nil {
		ls.InsertAllBefore(li, ac)
	}
}

// Remove removes the item li from ls if li is an item of set ls.
// It returns the item's value.
// The item li must not be nil.
func (ls *LinkedHashSet) Remove(index int) {
	li := ls.Item(index)
	if li != nil {
		li.Remove()
	}
}

// Swap swaps values of two items at the given index.
func (ls *LinkedHashSet) Swap(i, j int) {
	if i == j {
		return
	}

	ii := ls.Item(i)
	ij := ls.Item(j)
	if ii != nil && ij != nil && ii != ij {
		ii.value, ij.value = ij.value, ii.value
	}
}

// ReverseEach call f for each item in the set with reverse order
func (ls *LinkedHashSet) ReverseEach(f func(interface{})) {
	for li := ls.Back(); li != nil; li = li.Prev() {
		f(li.Value())
	}
}

// Iterator returns a iterator for the set
func (ls *LinkedHashSet) Iterator() Iterator {
	return &linkedSetItemIterator{ls, &ls.root}
}

//--------------------------------------------------------------------

// Front returns the first item of set ls or nil if the set is empty.
func (ls *LinkedHashSet) Front() *LinkedSetItem {
	if ls.IsEmpty() {
		return nil
	}
	return ls.root.next
}

// Back returns the last item of set ls or nil if the set is empty.
func (ls *LinkedHashSet) Back() *LinkedSetItem {
	if ls.IsEmpty() {
		return nil
	}
	return ls.root.prev
}

// Item returns the item at the specified index i.
// if i < -ls.Len() or i >= ls.Len(), returns nil
// if i < 0, returns ls.Item(ls.Len() + i)
func (ls *LinkedHashSet) Item(i int) *LinkedSetItem {
	len := ls.Len()

	if i < -len || i >= len {
		return nil
	}

	if i < 0 {
		i += len
	}
	if i >= len/2 {
		return ls.Back().Offset(i + 1 - len)
	}

	return ls.Front().Offset(i)
}

// Search looks for the given v, and returns the item associated with it,
// or nil if not found. The LinkedSetItem struct can then be used to iterate over the linked set
// from that point, either forward or backward.
func (ls *LinkedHashSet) Search(v interface{}) *LinkedSetItem {
	li, ok := ls.hash[v]
	if ok {
		return li
	}
	return nil
}

// Sort Sorts this set according to the order induced by the specified Comparator.
func (ls *LinkedHashSet) Sort(less cmp.Less) {
	if ls.Len() < 2 {
		return
	}
	sort.Sort(&sorter{ls, less})
}

// PushBack inserts all items of vs at the back of list ll.
// returns the last inserted item.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) PushBack(vs ...interface{}) *LinkedSetItem {
	if len(vs) == 0 {
		return nil
	}

	ls.lazyInit()
	return ls.InsertBefore(&ls.root, vs...)
}

// PushBackAll inserts a copy of another collection at the back of list ll.
// The ll and ac may be the same. They must not be nil.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) PushBackAll(ac Collection) *LinkedSetItem {
	if ac.IsEmpty() {
		return nil
	}

	ls.lazyInit()
	return ls.InsertAllBefore(&ls.root, ac)
}

// PushFront inserts all items of vs at the front of list ll.
// returns the last inserted item.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) PushFront(vs ...interface{}) *LinkedSetItem {
	if len(vs) == 0 {
		return nil
	}

	ls.lazyInit()
	return ls.InsertAfter(&ls.root, vs...)
}

// PushFrontAll inserts a copy of another collection at the front of list ll.
// The ll and ac may be the same. They must not be nil.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) PushFrontAll(ac Collection) *LinkedSetItem {
	if ac.IsEmpty() {
		return nil
	}

	ls.lazyInit()
	return ls.InsertAllAfter(&ls.root, ac)
}

// InsertAfter inserts all items of vs immediately after the item at and returns the last inserted item li.
// If at is not an item of ll, the list is not modified.
// The at must not be nil.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) InsertAfter(at *LinkedSetItem, vs ...interface{}) (li *LinkedSetItem) {
	if at.lset != ls || len(vs) == 0 {
		return
	}

	var ok bool
	for _, v := range vs {
		if li, ok = ls.insertAfter(at, v); ok {
			at = li
		}
	}
	return
}

// InsertAllAfter inserts all items of another collection ac immediately after the item at and returns the last inserted item li.
// If at is not an item of ll, the list is not modified.
// The at must not be nil.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) InsertAllAfter(at *LinkedSetItem, ac Collection) *LinkedSetItem {
	if at.lset != ls || ac.IsEmpty() {
		return nil
	}

	if ic, ok := ac.(Iterable); ok {
		return ls.insertIterAfter(at, ic.Iterator(), ac.Len())
	}
	return ls.InsertAfter(at, ac.Values()...)
}

// insertIterAfter inserts max n items of iterator it immediately after the item at and returns the last inserted item li.
// If at is not an item of ll, the list is not modified.
// The at must not be nil.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) insertIterAfter(at *LinkedSetItem, it Iterator, n int) (li *LinkedSetItem) {
	if at.lset != ls || n < 1 {
		return
	}

	var ok bool
	for ; it.Next() && n > 0; n-- {
		if li, ok = ls.insertAfter(at, it.Value()); ok {
			at = li
		}
	}
	return
}

// InsertBefore inserts all items of vs immediately before the item at and returns the last inserted item li.
// If at is not an item of ll, the list is not modified.
// The at must not be nil.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) InsertBefore(at *LinkedSetItem, vs ...interface{}) (li *LinkedSetItem) {
	if at.lset != ls || len(vs) == 0 {
		return
	}

	for _, v := range vs {
		li, _ = ls.insertBefore(at, v)
	}
	return
}

// InsertAllBefore inserts all items of another collection immediately before the item at and returns the last inserted item li.
// If at is not an item of ll, the list is not modified.
// The at must not be nil.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) InsertAllBefore(at *LinkedSetItem, ac Collection) *LinkedSetItem {
	if at.lset != ls || ac.IsEmpty() {
		return nil
	}

	if ic, ok := ac.(Iterable); ok {
		return ls.insertIterBefore(at, ic.Iterator(), ac.Len())
	}

	return ls.InsertBefore(at, ac.Values()...)
}

// insertIterBefore inserts max n items's value before the item at and returns the last inserted item li.
// If at is not an item of ll, the list is not modified.
// The at must not be nil.
// Note: existing item's order will not change.
func (ls *LinkedHashSet) insertIterBefore(at *LinkedSetItem, it Iterator, n int) (li *LinkedSetItem) {
	if at.lset != ls || n < 1 {
		return
	}

	ls.lazyInit()
	for ; it.Next() && n > 0; n-- {
		li, _ = ls.insertBefore(at, it.Value())
	}
	return
}

// MoveToFront moves item li to the front of set ls.
// If li is not an item of ls, the set is not modified.
// The item li must not be nil.
// Returns true if set is modified.
func (ls *LinkedHashSet) MoveToFront(li *LinkedSetItem) bool {
	if li.lset != ls || ls.root.next == li {
		return false
	}

	li.moveAfter(&ls.root)
	return true
}

// MoveToBack moves item li to the back of set ls.
// If li is not an item of ls, the set is not modified.
// The item li must not be nil.
// Returns true if set is modified.
func (ls *LinkedHashSet) MoveToBack(li *LinkedSetItem) bool {
	if li.lset != ls || ls.root.prev == li {
		return false
	}

	li.moveAfter(ls.root.prev)
	return true
}

// MoveBefore moves item li to its new position before at.
// If li or at is not an item of ls, or li == at, the set is not modified.
// The item li and at must not be nil.
// Returns true if set is modified.
func (ls *LinkedHashSet) MoveBefore(at, li *LinkedSetItem) bool {
	if li.lset != ls || li == at || at.lset != ls {
		return false
	}

	li.moveAfter(at.prev)
	return true
}

// MoveAfter moves item li to its new position after at.
// If li or at is not an item of ls, or li == at, the set is not modified.
// The item li and at must not be nil.
// Returns true if set is modified.
func (ls *LinkedHashSet) MoveAfter(at, li *LinkedSetItem) bool {
	if li.lset != ls || li == at || at.lset != ls {
		return false
	}

	li.moveAfter(at)
	return true
}

// SwapItem swap item's value of a, b.
// If a or b is not an item of ls, or a == b, the set is not modified.
// The item a and b must not be nil.
// Returns true if set is modified.
func (ls *LinkedHashSet) SwapItem(a, b *LinkedSetItem) bool {
	if a.lset != ls || a == b || b.lset != ls {
		return false
	}

	a.value, b.value = b.value, a.value
	return true
}

// String print set to string
func (ls *LinkedHashSet) String() string {
	bs, _ := json.Marshal(ls)
	return string(bs)
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

func (ls *LinkedHashSet) addJSONArrayItem(v interface{}) jsonArray {
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
