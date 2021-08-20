package col

import (
	"encoding/json"
)

// NewLinkedHashMap creates a new LinkedHashMap.
// Example: NewLinkedMap("k1", "v1", "k2", "v2")
func NewLinkedHashMap(kvs ...interface{}) *LinkedHashMap {
	lm := &LinkedHashMap{}
	lm.Set(kvs...)
	return lm
}

// LinkedHashMap implements an linked map that keeps track of the order in which keys were inserted.
// The zero value for LinkedHashMap is an empty map ready to use.
//
// To iterate over a linked map (where lm is a *LinkedHashMap):
//	for mi := lm.Front(); mi != nil; mi = mi.Next() {
//		// do something with mi.Key(), mi.Value()
//	}
//
type LinkedHashMap struct {
	hash map[interface{}]*LinkedMapItem
	root LinkedMapItem // the root item of doubly-linked list
}

// lazyInit lazily initializes a zero LinkedHashMap value.
func (lm *LinkedHashMap) lazyInit() {
	if lm.hash == nil {
		lm.hash = make(map[interface{}]*LinkedMapItem)
		lm.root.lmap = lm
		lm.root.next = &lm.root
		lm.root.prev = &lm.root
	}
}

func (lm *LinkedHashMap) add(key interface{}, value interface{}) {
	mi := &LinkedMapItem{key: key, value: value}
	mi.insertAfter(lm.root.prev)
}

//-----------------------------------------------------------
// implements Container interface

// Len returns the length of the linked map.
func (lm *LinkedHashMap) Len() int {
	return len(lm.hash)
}

// IsEmpty returns true if the map has no items
func (lm *LinkedHashMap) IsEmpty() bool {
	return len(lm.hash) == 0
}

// Clear clears the map
func (lm *LinkedHashMap) Clear() {
	lm.hash = nil
	lm.root.lmap = nil
	lm.root.next = nil
	lm.root.prev = nil
}

//-----------------------------------------------------------
// implements Map interface

// Get looks for the given key, and returns the value associated with it,
// or nil if not found. The boolean it returns says whether the key is ok in the map.
func (lm *LinkedHashMap) Get(key interface{}) (interface{}, bool) {
	if lm.hash != nil {
		if mi, ok := lm.hash[key]; ok {
			return mi.value, ok
		}
	}
	return nil, false
}

// Set sets the paired key-value items, and returns what `Get` would have returned
// on that key prior to the call to `Set`.
// Example: lm.Set("k1", "v1", "k2", "v2")
func (lm *LinkedHashMap) Set(kvs ...interface{}) (ov interface{}, ok bool) {
	if (len(kvs) % 2) != 0 {
		panic("LinkedHashMap.Set(kvs...) unpaired key-value items")
	}

	if len(kvs) < 2 {
		return
	}

	lm.lazyInit()

	var mi *LinkedMapItem
	for i := 0; i+1 < len(kvs); i += 2 {
		k := kvs[i]
		v := kvs[i+1]
		ov = nil
		if mi, ok = lm.hash[k]; ok {
			ov = mi.value
			mi.value = v
		} else {
			lm.add(k, v)
		}
	}
	return
}

// SetAll set items from another map am, override the existing items
func (lm *LinkedHashMap) SetAll(am Map) {
	setMapAll(lm, am)
}

// SetIfAbsent sets the key-value item if the key does not exists in the map,
// and returns what `Get` would have returned
// on that key prior to the call to `Set`.
// Example: lm.SetIfAbsent("k1", "v1", "k2", "v2")
func (lm *LinkedHashMap) SetIfAbsent(kvs ...interface{}) (ov interface{}, ok bool) {
	if (len(kvs) % 2) != 0 {
		panic("LinkedHashMap.SetIfAbsent(kvs...) unpaired key-value items")
	}

	if len(kvs) < 2 {
		return
	}

	lm.lazyInit()

	var mi *LinkedMapItem
	for i := 0; i+1 < len(kvs); i += 2 {
		k := kvs[i]
		v := kvs[i+1]
		ov = nil
		if mi, ok = lm.hash[k]; ok {
			ov = mi.value
		} else {
			lm.add(k, v)
		}
	}
	return
}

// Delete delete all items with key of ks,
// and returns what `Get` would have returned
// on that key prior to the call to `Set`.
func (lm *LinkedHashMap) Delete(ks ...interface{}) (ov interface{}, ok bool) {
	if lm.IsEmpty() {
		return
	}

	var mi *LinkedMapItem
	for _, k := range ks {
		if mi, ok = lm.hash[k]; ok {
			mi.Remove()
		} else {
			ov = nil
		}
	}
	return
}

// Contains looks for the given key, and returns true if the key exists in the map.
func (lm *LinkedHashMap) Contains(ks ...interface{}) bool {
	if len(ks) == 0 {
		return true
	}

	if lm.IsEmpty() {
		return false
	}

	for _, k := range ks {
		if _, ok := lm.hash[k]; !ok {
			return false
		}
	}
	return true
}

// Keys returns the key slice
func (lm *LinkedHashMap) Keys() []interface{} {
	ks := make([]interface{}, lm.Len())
	for i, mi := 0, lm.Front(); mi != nil; i, mi = i+1, mi.Next() {
		ks[i] = mi.Key()
	}
	return ks
}

// Values returns the value slice
func (lm *LinkedHashMap) Values() []interface{} {
	vs := make([]interface{}, lm.Len())
	for i, mi := 0, lm.Front(); mi != nil; i, mi = i+1, mi.Next() {
		vs[i] = mi.Value()
	}
	return vs
}

// Each call f for each item in the map
func (lm *LinkedHashMap) Each(f func(k interface{}, v interface{})) {
	for mi := lm.Front(); mi != nil; mi = mi.Next() {
		f(mi.Key(), mi.Value())
	}
}

//-----------------------------------------------------------
// implements IterableMap interface

// ReverseEach call f for each item in the map with reverse order
func (lm *LinkedHashMap) ReverseEach(f func(k interface{}, v interface{})) {
	for mi := lm.Back(); mi != nil; mi = mi.Prev() {
		f(mi.Key(), mi.Value())
	}
}

// Iterator returns a iterator for the map
func (lm *LinkedHashMap) Iterator() Iterator2 {
	return &linkedMapItemIterator{lm, &lm.root}
}

//-----------------------------------------------------------

// Front returns a pointer to the oldest item. It's meant to be used to iterate on the linked map's
// items from the oldest to the newest, e.g.:
// for item := lm.Front(); item != nil; item = item.Next() {
//    fmt.Printf("%v => %v\n", item.Key(), item.Value())
// }
func (lm *LinkedHashMap) Front() *LinkedMapItem {
	if lm.IsEmpty() {
		return nil
	}
	return lm.root.next
}

// Back returns a pointer to the newest item. It's meant to be used to iterate on the linked map's
// items from the newest to the oldest, e.g.:
// for item := lm.Back(); item != nil; item = item.Prev() {
//    fmt.Printf("%v => %v\n", item.Key(), item.Value())
// }
func (lm *LinkedHashMap) Back() *LinkedMapItem {
	if lm.IsEmpty() {
		return nil
	}
	return lm.root.prev
}

// Search looks for the given key, and returns the item associated with it,
// or nil if not found. The LinkedMapItem struct can then be used to iterate over the linked map
// from that point, either forward or backward.
func (lm *LinkedHashMap) Search(key interface{}) *LinkedMapItem {
	mi, ok := lm.hash[key]
	if ok {
		return mi
	}
	return nil
}

// Items returns the map item slice
func (lm *LinkedHashMap) Items() []*LinkedMapItem {
	mis := make([]*LinkedMapItem, lm.Len(), lm.Len())
	for i, mi := 0, lm.Front(); mi != nil; i, mi = i+1, mi.Next() {
		mis[i] = mi
	}
	return mis
}

// String print map to string
func (lm *LinkedHashMap) String() string {
	bs, _ := json.Marshal(lm)
	return string(bs)
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

func newJSONObjectAsLinkedMap() jsonObject {
	return NewLinkedHashMap()
}

func (lm *LinkedHashMap) addJSONObjectItem(k string, v interface{}) jsonObject {
	lm.Set(k, v)
	return lm
}

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(lm)
func (lm *LinkedHashMap) MarshalJSON() (res []byte, err error) {
	return jsonMarshalMap(lm)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, lm)
func (lm *LinkedHashMap) UnmarshalJSON(data []byte) error {
	lm.Clear()
	ju := &jsonUnmarshaler{
		newArray:  newJSONArray,
		newObject: newJSONObjectAsLinkedMap,
	}
	return ju.unmarshalJSONObject(data, lm)
}
