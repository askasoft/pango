package col

// LessString string less function
func LessString(a interface{}, b interface{}) bool {
	return a.(string) < b.(string)
}

// LessInt int less function
func LessInt(a interface{}, b interface{}) bool {
	return a.(int) < b.(int)
}

// LessInt32 int32 less function
func LessInt32(a interface{}, b interface{}) bool {
	return a.(int32) < b.(int32)
}

// LessInt64 int64 less function
func LessInt64(a interface{}, b interface{}) bool {
	return a.(int64) < b.(int64)
}

// LessFloat32 float32 less function
func LessFloat32(a interface{}, b interface{}) bool {
	return a.(float32) < b.(float32)
}

// LessFloat64 float64 less function
func LessFloat64(a interface{}, b interface{}) bool {
	return a.(float64) < b.(float64)
}

// SortedList implements an sorted entry list.
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

// Len returns the number of entries of list l.
// The complexity is O(1).
func (sl *SortedList) Len() int {
	return sl.list.Len()
}

// IsEmpty checks if the list is empty.
func (sl *SortedList) IsEmpty() bool {
	return sl.list.IsEmpty()
}

// At returns the entry at the specified index
func (sl *SortedList) At(i int) *ListEntry {
	return sl.list.At(i)
}

// Front returns the first entry of list l or nil if the list is empty.
func (sl *SortedList) Front() *ListEntry {
	return sl.list.Front()
}

// Back returns the last entry of list l or nil if the list is empty.
func (sl *SortedList) Back() *ListEntry {
	return sl.list.Back()
}

// Contains Test to see whether or not the v is in the list
func (sl *SortedList) Contains(v interface{}) bool {
	_, e := sl.Search(v)
	return e != nil
}

// Search binary search v
// returns index, entry if it's value is v
// if not found, returns -1, nil
func (sl *SortedList) Search(v interface{}) (int, *ListEntry) {
	n, e := sl.binarySearch(v)
	if e != nil && e.Value == v {
		return n, e
	}

	return -1, nil
}

// binarySearch binary search v
// returns index, entry if it's value is >= v
// if not found, returns -1, nil
func (sl *SortedList) binarySearch(v interface{}) (int, *ListEntry) {
	if sl.IsEmpty() {
		return -1, nil
	}

	e := sl.Front()
	p, i, j := 0, 0, sl.Len()
	for i < j && e != nil {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		e = e.Move(h - p)
		p = h
		// i ≤ h < j
		if sl.less(e.Value, v) {
			i = h + 1
		} else {
			j = h
		}
	}

	if i < sl.Len() {
		e = e.Move(i - p)
		return i, e
	}
	return -1, nil
}

// Add inserts a new entry e with value v and returns e.
func (sl *SortedList) Add(v interface{}) *ListEntry {
	if sl.IsEmpty() {
		return sl.list.PushBack(v)
	}

	_, e := sl.binarySearch(v)
	if e != nil {
		return sl.list.InsertBefore(v, e)
	}

	return sl.list.PushBack(v)
}

// AddList adds a copy of another list.
// The lists l and other may be the same. They must not be nil.
func (sl *SortedList) AddList(other *List) {
	for e := other.Front(); e != nil; e = e.Next() {
		sl.Add(e.Value)
	}
}

// AddAll adds all items of vs.
func (sl *SortedList) AddAll(vs ...interface{}) {
	for _, v := range vs {
		sl.Add(v)
	}
}

// Remove The entry must not be nil.
func (sl *SortedList) Remove(e *ListEntry) interface{} {
	return sl.list.Remove(e)
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

/*------------- JSON -----------------*/

func newJSONArrayAsSortedList() jsonArray {
	return NewSet()
}

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
