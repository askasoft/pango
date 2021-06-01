package col

import (
	"encoding/json"
	"fmt"
)

// OrderedMapEntry key/value entry
type OrderedMapEntry struct {
	MapEntry
	entry *ListEntry
}

// Next returns a pointer to the next entry.
func (e *OrderedMapEntry) Next() *OrderedMapEntry {
	return toOrderedMapEntry(e.entry.Next())
}

// Prev returns a pointer to the previous entry.
func (e *OrderedMapEntry) Prev() *OrderedMapEntry {
	return toOrderedMapEntry(e.entry.Prev())
}

func toOrderedMapEntry(le *ListEntry) *OrderedMapEntry {
	if le == nil {
		return nil
	}
	return le.Value.(*OrderedMapEntry)
}

// OrderedMap implements an ordered map that keeps track of the order in which keys were inserted.
type OrderedMap struct {
	entries map[interface{}]*OrderedMapEntry
	list    *List
}

// NewOrderedMap creates a new OrderedMap.
// Example: NewOrderedMap("k1", "v1", "k2", "v2")
func NewOrderedMap(kvs ...interface{}) *OrderedMap {
	om := &OrderedMap{
		entries: make(map[interface{}]*OrderedMapEntry),
		list:    NewList(),
	}
	for i := 0; i+1 < len(kvs); i += 2 {
		om.Set(kvs[i], kvs[i+1])
	}
	return om
}

// GetEntry looks for the given key, and returns the entry associated with it,
// or nil if not found. The OrderedMapEntry struct can then be used to iterate over the ordered map
// from that point, either forward or backward.
func (om *OrderedMap) GetEntry(key interface{}) *OrderedMapEntry {
	return om.entries[key]
}

// Has looks for the given key, and returns true if the key exists in the map.
func (om *OrderedMap) Has(key interface{}) bool {
	_, ok := om.entries[key]
	return ok
}

// Get looks for the given key, and returns the value associated with it,
// or nil if not found. The boolean it returns says whether the key is ok in the map.
func (om *OrderedMap) Get(key interface{}) (interface{}, bool) {
	if me, ok := om.entries[key]; ok {
		return me.Value, ok
	}
	return nil, false
}

// Set sets the key-value entry, and returns what `Get` would have returned
// on that key prior to the call to `Set`.
func (om *OrderedMap) Set(key interface{}, value interface{}) (interface{}, bool) {
	if me, ok := om.entries[key]; ok {
		old := me.Value
		me.Value = value
		return old, true
	}

	me := &OrderedMapEntry{}
	me.key = key
	me.Value = value
	me.entry = om.list.PushBack(me)
	om.entries[key] = me

	return nil, false
}

// Copy copy entries from another map am, override the existing entries
func (om *OrderedMap) Copy(am *OrderedMap) {
	for e := am.Front(); e != nil; e = e.Next() {
		om.Set(e.key, e.Value)
	}
}

// Remove removes the key-value entry, and returns what `Get` would have returned
// on that key prior to the call to `Remove`.
func (om *OrderedMap) Remove(key interface{}) (interface{}, bool) {
	if me, ok := om.entries[key]; ok {
		om.list.Remove(me.entry)
		delete(om.entries, key)
		return me.Value, true
	}

	return nil, false
}

// Len returns the length of the ordered map.
func (om *OrderedMap) Len() int {
	return len(om.entries)
}

// IsEmpty returns true if the map has no entries
func (om *OrderedMap) IsEmpty() bool {
	return len(om.entries) == 0
}

// Clear clears the map
func (om *OrderedMap) Clear() {
	om.entries = make(map[interface{}]*OrderedMapEntry)
	om.list.Clear()
}

// Front returns a pointer to the oldest entry. It's meant to be used to iterate on the ordered map's
// entries from the oldest to the newest, e.g.:
// for entry := orderedMap.Front(); entry != nil; entry = entry.Next() { fmt.Printf("%v => %v\n", entry.Key(), entry.Value()) }
func (om *OrderedMap) Front() *OrderedMapEntry {
	return toOrderedMapEntry(om.list.Front())
}

// Back returns a pointer to the newest entry. It's meant to be used to iterate on the ordered map's
// entries from the newest to the oldest, e.g.:
// for entry := orderedMap.Back(); entry != nil; entry = entry.Prev() { fmt.Printf("%v => %v\n", entry.Key(), entry.Value()) }
func (om *OrderedMap) Back() *OrderedMapEntry {
	return toOrderedMapEntry(om.list.Back())
}

// Keys returns the key slice
func (om *OrderedMap) Keys() []interface{} {
	ks := make([]interface{}, 0, om.Len())
	for me := om.Front(); me != nil; me = me.Next() {
		ks = append(ks, me.Key())
	}
	return ks
}

// Values returns the value slice
func (om *OrderedMap) Values() []interface{} {
	vs := make([]interface{}, 0, om.Len())
	for me := om.Front(); me != nil; me = me.Next() {
		vs = append(vs, me.Value)
	}
	return vs
}

// Entries returns the mep entry slice
func (om *OrderedMap) Entries() []*OrderedMapEntry {
	vs := make([]*OrderedMapEntry, 0, om.Len())
	for me := om.Front(); me != nil; me = me.Next() {
		vs = append(vs, me)
	}
	return vs
}

// Each Call f for each item in the map
func (om *OrderedMap) Each(f func(*OrderedMapEntry)) {
	for me := om.Front(); me != nil; me = me.Next() {
		f(me)
	}
}

// ReverseEach Call f for each item in the map with reverse order
func (om *OrderedMap) ReverseEach(f func(*OrderedMapEntry)) {
	for me := om.Back(); me != nil; me = me.Prev() {
		f(me)
	}
}

/*------------- JSON -----------------*/
func newJSONObjectAsOrderedMap() jsonObject {
	return NewOrderedMap()
}

func (om *OrderedMap) addJSONObjectItem(k string, v interface{}) jsonObject {
	om.Set(k, v)
	return om
}

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(om)
func (om *OrderedMap) MarshalJSON() (res []byte, err error) {
	if om.IsEmpty() {
		return []byte("{}"), nil
	}

	res = append(res, '{')
	for me := om.Front(); me != nil; me = me.Next() {
		k, ok := me.key.(string)
		if !ok {
			err = fmt.Errorf("expecting JSON key should be always a string: %T: %v", me.key, me.key)
			return
		}

		res = append(res, fmt.Sprintf("%q:", k)...)
		var b []byte
		b, err = json.Marshal(me.Value)
		if err != nil {
			return
		}
		res = append(res, b...)
		res = append(res, ',')
	}
	res[len(res)-1] = '}'
	return
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, om)
func (om *OrderedMap) UnmarshalJSON(data []byte) error {
	ju := &jsonUnmarshaler{
		newArray:  newJSONArray,
		newObject: newJSONObjectAsOrderedMap,
	}
	return ju.unmarshalJSONObject(data, om)
}
