package col

import (
	"encoding/json"
	"fmt"
)

// OrderedMap implements an ordered map that keeps track of the order in which keys were inserted.
type OrderedMap struct {
	hash map[interface{}]*OrderedMapItem
	list *List
}

// NewOrderedMap creates a new OrderedMap.
// Example: NewOrderedMap("k1", "v1", "k2", "v2")
func NewOrderedMap(kvs ...interface{}) *OrderedMap {
	om := &OrderedMap{
		hash: make(map[interface{}]*OrderedMapItem),
		list: NewList(),
	}
	for i := 0; i+1 < len(kvs); i += 2 {
		om.Set(kvs[i], kvs[i+1])
	}
	return om
}

// Len returns the length of the ordered map.
func (om *OrderedMap) Len() int {
	return len(om.hash)
}

// IsEmpty returns true if the map has no items
func (om *OrderedMap) IsEmpty() bool {
	return len(om.hash) == 0
}

// Item looks for the given key, and returns the item associated with it,
// or nil if not found. The OrderedMapItem struct can then be used to iterate over the ordered map
// from that point, either forward or backward.
func (om *OrderedMap) Item(key interface{}) *OrderedMapItem {
	return om.hash[key]
}

// Has looks for the given key, and returns true if the key exists in the map.
func (om *OrderedMap) Has(key interface{}) bool {
	_, ok := om.hash[key]
	return ok
}

// Get looks for the given key, and returns the value associated with it,
// or nil if not found. The boolean it returns says whether the key is ok in the map.
func (om *OrderedMap) Get(key interface{}) (interface{}, bool) {
	if mi, ok := om.hash[key]; ok {
		return mi.Value, ok
	}
	return nil, false
}

func (om *OrderedMap) put(key interface{}, value interface{}) {
	mi := &OrderedMapItem{}
	mi.key = key
	mi.Value = value
	mi.item = om.list.PushBack(mi)
	om.hash[key] = mi
}

// Set sets the key-value item, and returns what `Get` would have returned
// on that key prior to the call to `Set`.
func (om *OrderedMap) Set(key interface{}, value interface{}) (interface{}, bool) {
	if mi, ok := om.hash[key]; ok {
		ov := mi.Value
		mi.Value = value
		return ov, true
	}

	om.put(key, value)
	return nil, false
}

// SetIfAbsent sets the key-value item if the key does not exists in the map,
// and returns what `Get` would have returned
// on that key prior to the call to `Set`.
func (om *OrderedMap) SetIfAbsent(key interface{}, value interface{}) (interface{}, bool) {
	if mi, ok := om.hash[key]; ok {
		return mi.Value, true
	}

	om.put(key, value)
	return nil, false
}

// Copy copy items from another map am, override the existing items
func (om *OrderedMap) Copy(am *OrderedMap) {
	for mi := am.Front(); mi != nil; mi = mi.Next() {
		om.Set(mi.key, mi.Value)
	}
}

// Delete delete the item with key, and returns what `Get` would have returned
// on that key prior to the call to `Delete`.
func (om *OrderedMap) Delete(key interface{}) (interface{}, bool) {
	if mi, ok := om.hash[key]; ok {
		om.list.Remove(mi.item)
		delete(om.hash, key)
		return mi.Value, true
	}

	return nil, false
}

// Clear clears the map
func (om *OrderedMap) Clear() {
	om.hash = make(map[interface{}]*OrderedMapItem)
	om.list.Clear()
}

// Front returns a pointer to the oldest item. It's meant to be used to iterate on the ordered map's
// items from the oldest to the newest, e.g.:
// for item := orderedMap.Front(); item != nil; item = item.Next() { fmt.Printf("%v => %v\n", item.Key(), item.Value()) }
func (om *OrderedMap) Front() *OrderedMapItem {
	return toOrderedMapItem(om.list.Front())
}

// Back returns a pointer to the newest item. It's meant to be used to iterate on the ordered map's
// items from the newest to the oldest, e.g.:
// for item := orderedMap.Back(); item != nil; item = item.Prev() { fmt.Printf("%v => %v\n", item.Key(), item.Value()) }
func (om *OrderedMap) Back() *OrderedMapItem {
	return toOrderedMapItem(om.list.Back())
}

// Keys returns the key slice
func (om *OrderedMap) Keys() []interface{} {
	ks := make([]interface{}, 0, om.Len())
	for mi := om.Front(); mi != nil; mi = mi.Next() {
		ks = append(ks, mi.Key())
	}
	return ks
}

// Values returns the value slice
func (om *OrderedMap) Values() []interface{} {
	vs := make([]interface{}, 0, om.Len())
	for mi := om.Front(); mi != nil; mi = mi.Next() {
		vs = append(vs, mi.Value)
	}
	return vs
}

// Items returns the mep item slice
func (om *OrderedMap) Items() []*OrderedMapItem {
	vs := make([]*OrderedMapItem, 0, om.Len())
	for mi := om.Front(); mi != nil; mi = mi.Next() {
		vs = append(vs, mi)
	}
	return vs
}

// Each Call f for each item in the map
func (om *OrderedMap) Each(f func(*OrderedMapItem)) {
	for mi := om.Front(); mi != nil; mi = mi.Next() {
		f(mi)
	}
}

// ReverseEach Call f for each item in the map with reverse order
func (om *OrderedMap) ReverseEach(f func(*OrderedMapItem)) {
	for mi := om.Back(); mi != nil; mi = mi.Prev() {
		f(mi)
	}
}

// String print map to string
func (om *OrderedMap) String() string {
	bs, _ := json.Marshal(om)
	return string(bs)
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
	for mi := om.Front(); mi != nil; mi = mi.Next() {
		k, ok := mi.key.(string)
		if !ok {
			err = fmt.Errorf("expecting JSON key should be always a string: %T: %v", mi.key, mi.key)
			return
		}

		res = append(res, fmt.Sprintf("%q:", k)...)
		var b []byte
		b, err = json.Marshal(mi.Value)
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
