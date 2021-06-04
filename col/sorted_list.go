package col

// SortedList implements an sorted list.
type SortedList struct {
	list *List
	less func(a, b interface{}) bool
}

// NewSortedList returns an initialized list.
func NewSortedList(less func(a, b interface{}) bool, vs ...interface{}) *SortedList {
	sl := &SortedList{
		list: NewList(),
		less: less,
	}
	sl.AddAll(vs...)
	return sl
}

// Len returns the length of the list.
// The complexity is O(1).
func (sl *SortedList) Len() int {
	return sl.list.Len()
}

// IsEmpty checks if the list is empty.
func (sl *SortedList) IsEmpty() bool {
	return sl.list.IsEmpty()
}

// Item returns the item at the specified index
func (sl *SortedList) Item(i int) *ListItem {
	return sl.list.Item(i)
}

// Front returns the first item of list l or nil if the list is empty.
func (sl *SortedList) Front() *ListItem {
	return sl.list.Front()
}

// Back returns the last item of list l or nil if the list is empty.
func (sl *SortedList) Back() *ListItem {
	return sl.list.Back()
}

// Contains Test to see whether or not the v is in the list
func (sl *SortedList) Contains(v interface{}) bool {
	_, li := sl.Search(v)
	return li != nil
}

// Search binary search v
// returns (index, item) if it's value is v
// if not found, returns (-1, nil)
func (sl *SortedList) Search(v interface{}) (int, *ListItem) {
	n, li := sl.binarySearch(v)
	if li != nil && li.Value == v {
		return n, li
	}

	return -1, nil
}

// binarySearch binary search v
// returns (index, item) if it's value is >= v
// if not found, returns (-1, nil)
func (sl *SortedList) binarySearch(v interface{}) (int, *ListItem) {
	if sl.IsEmpty() {
		return -1, nil
	}

	li := sl.Front()
	p, i, j := 0, 0, sl.Len()
	for i < j && li != nil {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		li = li.Offset(h - p)
		p = h
		// i ≤ h < j
		if sl.less(li.Value, v) {
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

// Add inserts a new item li with value v and returns li.
func (sl *SortedList) Add(v interface{}) *ListItem {
	if sl.IsEmpty() {
		return sl.list.PushBack(v)
	}

	_, li := sl.binarySearch(v)
	if li != nil {
		return sl.list.InsertBefore(v, li)
	}

	return sl.list.PushBack(v)
}

// AddAll adds all items of vs.
func (sl *SortedList) AddAll(vs ...interface{}) {
	for _, v := range vs {
		sl.Add(v)
	}
}

// AddList adds a copy of another list.
// The lists l and other may be the same. They must not be nil.
func (sl *SortedList) AddList(other *List) {
	for li := other.Front(); li != nil; li = li.Next() {
		sl.Add(li.Value)
	}
}

// Delete delete the first item with associated value v
// returns true if v is in the list
// returns false if the the list is not changed
func (sl *SortedList) Delete(v interface{}) bool {
	_, li := sl.Search(v)
	if li != nil {
		sl.list.remove(li)
		return true
	}

	return false
}

// DeleteAll delete all items with associated value v
// returns the deleted count
func (sl *SortedList) DeleteAll(v interface{}) int {
	n := 0

	_, li := sl.binarySearch(v)
	for li != nil && li.Value == v {
		ni := li.Next()
		sl.list.remove(li)
		n++
		li = ni
	}

	return n
}

// Remove The item must not be nil.
func (sl *SortedList) Remove(li *ListItem) {
	sl.list.Remove(li)
}

// Values returns a slice contains all the items of the list l
func (sl *SortedList) Values() []interface{} {
	return sl.list.Values()
}

// Each Call f for each item in the set
func (sl *SortedList) Each(f func(interface{})) {
	sl.list.Each(f)
}

// ReverseEach Call f for each item in the set with reverse order
func (sl *SortedList) ReverseEach(f func(interface{})) {
	sl.list.ReverseEach(f)
}

// String print list to string
func (sl *SortedList) String() string {
	return sl.list.String()
}

/*------------- JSON -----------------*/

func (sl *SortedList) addJSONArrayItem(v interface{}) jsonArray {
	sl.Add(v)
	return sl
}

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(l)
func (sl *SortedList) MarshalJSON() (res []byte, err error) {
	return sl.list.MarshalJSON()
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, l)
func (sl *SortedList) UnmarshalJSON(data []byte) error {
	ju := &jsonUnmarshaler{
		newArray:  newJSONArray,
		newObject: newJSONObject,
	}
	return ju.unmarshalJSONArray(data, sl)
}
