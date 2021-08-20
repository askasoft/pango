package col

import (
	"encoding/json"

	"github.com/pandafw/pango/cmp"
)

// NewSortedList returns an initialized sorted list.
func NewSortedList(less cmp.Less, vs ...interface{}) *SortedList {
	sl := &SortedList{
		less: less,
	}
	sl.Add(vs...)
	return sl
}

// SortedList implements a sorted list.
type SortedList struct {
	less cmp.Less
	root SortedListItem
	len  int
}

// lazyInit lazily initializes a zero LinkedList value.
func (sl *SortedList) lazyInit() {
	if sl.root.next == nil {
		sl.root.list = sl
		sl.root.next = &sl.root
		sl.root.prev = &sl.root
		sl.len = 0
	}
}

// binarySearch binary search v
// returns (index, item) if it's value is >= v
// if not found, returns (-1, nil)
func (sl *SortedList) binarySearch(v interface{}) (int, *SortedListItem) {
	if sl.IsEmpty() {
		return -1, nil
	}

	li := sl.root.next
	p, i, j := 0, 0, sl.Len()
	for i < j && li != nil {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		li = li.Offset(h - p)
		p = h
		// i ≤ h < j
		if sl.less(li.value, v) {
			i = h + 1
		} else {
			j = h
		}
	}

	if i < sl.Len() {
		li = li.Offset(i - p)
		return i, li
	}
	return -1, nil
}

// add inserts a new item li with value v and returns li.
func (sl *SortedList) add(v interface{}) (li *SortedListItem) {
	li = &SortedListItem{value: v}

	if sl.IsEmpty() {
		li.insertAfter(sl.root.prev)
		return
	}

	_, at := sl.binarySearch(v)
	if at != nil {
		li.insertAfter(at.prev)
		return
	}

	li.insertAfter(sl.root.prev)
	return
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
	sl.root.list = nil
	sl.root.next = nil
	sl.root.prev = nil
	sl.len = 0
}

// Add adds all items of vs and returns the last added item.
func (sl *SortedList) Add(vs ...interface{}) {
	if len(vs) == 0 {
		return
	}

	sl.lazyInit()
	for _, v := range vs {
		sl.add(v)
	}
}

// AddAll adds all items of another collection
func (sl *SortedList) AddAll(ac Collection) {
	if sl == ac {
		sl.Add(ac.Values()...)
		return
	}

	if ic, ok := ac.(Iterable); ok {
		it := ic.Iterator()
		for it.Next() {
			sl.Add(it.Value())
		}
		return
	}

	sl.Add(ac.Values()...)
}

func (sl *SortedList) deleteAll(v interface{}) {
	_, li := sl.binarySearch(v)
	for li != nil && li.value == v {
		li.Remove()
		li = li.Next()
	}
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
	for _, v := range vs {
		if sl.Index(v) < 0 {
			return false
		}
	}
	return true
}

// ContainsAll Test to see if the collection contains all items of another collection
func (sl *SortedList) ContainsAll(ac Collection) bool {
	if sl == ac {
		return true
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

	for li := sl.Front(); li != nil; li = li.Next() {
		if !ac.Contains(li.Value()) {
			li.Remove()
		}
	}
}

// Values returns a slice contains all the items of the list l
func (sl *SortedList) Values() []interface{} {
	vs := make([]interface{}, sl.Len())
	for i, li := 0, sl.Front(); li != nil; i, li = i+1, li.Next() {
		vs[i] = li.Value()
	}
	return vs
}

// Each call f for each item in the list
func (sl *SortedList) Each(f func(interface{})) {
	for li := sl.Front(); li != nil; li = li.Next() {
		f(li.Value())
	}
}

//-----------------------------------------------------------
// implements List interface

// Get returns the element at the specified position in this list
func (sl *SortedList) Get(index int) (interface{}, bool) {
	li := sl.Item(index)
	if li == nil {
		return nil, false
	}
	return li.Value(), true
}

// Set set the v at the specified index in this list and returns the old value.
func (sl *SortedList) Set(index int, v interface{}) (ov interface{}) {
	li := sl.Item(index)
	if li != nil {
		ov = li.Value()
		li.SetValue(v)
	}
	return
}

// Insert inserts values, same as: Adds(...)
// Just implements the List.Insert() method
// Does not do anything if position is bigger than list's size
func (sl *SortedList) Insert(index int, vs ...interface{}) {
	n := len(vs)
	if n == 0 {
		return
	}

	len := sl.Len()
	if index < -len || index > len {
		return
	}

	sl.Add(vs...)
}

// InsertAll inserts values of another collection ac at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Does not do anything if position is bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (sl *SortedList) InsertAll(index int, ac Collection) {
	n := ac.Len()
	if n == 0 {
		return
	}

	len := sl.Len()
	if index < -len || index > len {
		return
	}

	if ic, ok := ac.(Iterable); ok {
		it := ic.Iterator()
		for it.Next() {
			sl.Add(it.Value())
		}
		return
	}

	sl.Add(ac.Values()...)
}

// Index returns the index of the first occurrence of the specified v in this list, or -1 if this list does not contain v.
func (sl *SortedList) Index(v interface{}) int {
	i, _ := sl.Search(v)
	return i
}

// Remove removes the item li from l if li is an item of list l.
// It returns the item's value.
// The item li must not be nil.
func (sl *SortedList) Remove(index int) {
	li := sl.Item(index)
	if li != nil {
		li.Remove()
	}
}

// Swap swaps values of two items at the given index.
// Do nothing because all items are sorted.
func (sl *SortedList) Swap(i, j int) {
	// do nothing
}

// ReverseEach Call f for each item in the list with reverse order
func (sl *SortedList) ReverseEach(f func(interface{})) {
	for li := sl.Back(); li != nil; li = li.Prev() {
		f(li.Value())
	}
}

// Iterator returns a iterator for the list
func (sl *SortedList) Iterator() Iterator {
	return &sortedListItemIterator{sl, &sl.root}
}

//--------------------------------------------------------------------

// Front returns the first item of list l or nil if the list is empty.
func (sl *SortedList) Front() *SortedListItem {
	if sl.len == 0 {
		return nil
	}
	return sl.root.next
}

// Back returns the last item of list l or nil if the list is empty.
func (sl *SortedList) Back() *SortedListItem {
	if sl.len == 0 {
		return nil
	}
	return sl.root.prev
}

// Item returns the item at the specified index
// if i < -l.Len() or i >= l.Len(), returns nil
// if i < 0, returns l.Item(l.Len() + i)
func (sl *SortedList) Item(i int) *SortedListItem {
	if i < -sl.len || i >= sl.len {
		return nil
	}

	if i < 0 {
		i += sl.len
	}
	if i >= sl.len/2 {
		return sl.Back().Offset(i + 1 - sl.len)
	}

	return sl.Front().Offset(i)
}

// Search binary search v
// returns (index, item) if it's value is v
// if not found, returns (-1, nil)
func (sl *SortedList) Search(v interface{}) (int, *SortedListItem) {
	n, li := sl.binarySearch(v)
	if li != nil && li.Value() == v {
		return n, li
	}

	return -1, nil
}

// String print list to string
func (sl *SortedList) String() string {
	bs, _ := json.Marshal(sl)
	return string(bs)
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
