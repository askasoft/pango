package col

import (
	"encoding/json"
)

// NewHashMap creates a new HashMap.
// Example: NewHashMap("k1", "v1", "k2", "v2")
func NewHashMap(kvs ...interface{}) *HashMap {
	hm := &HashMap{}
	hm.Set(kvs...)
	return hm
}

// AsHashMap creates a new HashMap from a map.
// Example: ToHashMap(map[interface{}]interface{}{"k1": "v1", "k2": "v2"})
func AsHashMap(m map[interface{}]interface{}) *HashMap {
	hm := &HashMap{m}
	return hm
}

// HashMap hash map type
type HashMap struct {
	hash map[interface{}]interface{}
}

// lazyInit lazily initializes a zero HashMap value.
func (hm *HashMap) lazyInit() {
	if hm.hash == nil {
		hm.hash = make(map[interface{}]interface{})
	}
}

//-----------------------------------------------------------
// implements Container interface

// Len returns the length of the linked map.
func (hm *HashMap) Len() int {
	return len(hm.hash)
}

// IsEmpty returns true if the map has no items
func (hm *HashMap) IsEmpty() bool {
	return len(hm.hash) == 0
}

// Clear clears the map
func (hm *HashMap) Clear() {
	// for AsHashMap()
	for k := range hm.hash {
		delete(hm.hash, k)
	}
}

//-----------------------------------------------------------
// implements Map interface

// Get looks for the given key, and returns the value associated with it,
// or nil if not found. The boolean it returns says whether the key is ok in the map.
func (hm *HashMap) Get(key interface{}) (v interface{}, ok bool) {
	if hm.hash == nil {
		return
	}

	v, ok = hm.hash[key]
	return
}

// Set sets the paired key-value items, and returns what `Get` would have returned
// on that key prior to the call to `Set`.
// Example: lm.Set("k1", "v1", "k2", "v2")
func (hm *HashMap) Set(kvs ...interface{}) (ov interface{}, ok bool) {
	if (len(kvs) % 2) != 0 {
		panic("HashMap.Set(kvs...) unpaired key-value items")
	}

	if len(kvs) < 2 {
		return
	}

	hm.lazyInit()

	for i := 0; i+1 < len(kvs); i += 2 {
		k := kvs[i]
		v := kvs[i+1]
		ov, ok = hm.hash[k]
		hm.hash[k] = v
	}
	return
}

// SetAll set items from another map am, override the existing items
func (hm *HashMap) SetAll(am Map) {
	setMapAll(hm, am)
}

// SetIfAbsent sets the key-value item if the key does not exists in the map,
// and returns what `Get` would have returned
// on that key prior to the call to `Set`.
// Example: lm.SetIfAbsent("k1", "v1", "k2", "v2")
func (hm *HashMap) SetIfAbsent(kvs ...interface{}) (ov interface{}, ok bool) {
	if (len(kvs) % 2) != 0 {
		panic("HashMap.SetIfAbsent(kvs...) unpaired key-value items")
	}

	if len(kvs) < 2 {
		return
	}

	hm.lazyInit()

	for i := 0; i+1 < len(kvs); i += 2 {
		k := kvs[i]
		v := kvs[i+1]
		if ov, ok = hm.hash[k]; !ok {
			hm.hash[k] = v
		}
	}
	return
}

// Delete delete all items with key of ks,
// and returns what `Get` would have returned
// on that key prior to the call to `Set`.
func (hm *HashMap) Delete(ks ...interface{}) (ov interface{}, ok bool) {
	if hm.IsEmpty() {
		return
	}

	for _, k := range ks {
		ov, ok = hm.hash[k]
		delete(hm.hash, k)
	}
	return
}

// Contains looks for the given key, and returns true if the key exists in the map.
func (hm *HashMap) Contains(ks ...interface{}) bool {
	if len(ks) == 0 {
		return true
	}

	if hm.IsEmpty() {
		return false
	}

	for _, k := range ks {
		if _, ok := hm.hash[k]; !ok {
			return false
		}
	}
	return true
}

// Keys returns the key slice
func (hm *HashMap) Keys() []interface{} {
	ks := make([]interface{}, hm.Len())
	i := 0
	for k := range hm.hash {
		ks[i] = k
		i++
	}
	return ks
}

// Values returns the value slice
func (hm *HashMap) Values() []interface{} {
	vs := make([]interface{}, hm.Len())
	i := 0
	for _, v := range hm.hash {
		vs[i] = v
		i++
	}
	return vs
}

// Each call f for each item(k,v) in the map
func (hm *HashMap) Each(f func(k interface{}, v interface{})) {
	for k, v := range hm.hash {
		f(k, v)
	}
}

// HashMap returns underlying hash map
func (hm *HashMap) HashMap() map[interface{}]interface{} {
	return hm.hash
}

// String print map to string
func (hm *HashMap) String() string {
	bs, _ := json.Marshal(hm)
	return string(bs)
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

func newJSONObjectAsHashMap() jsonObject {
	return NewHashMap()
}

func (hm *HashMap) addJSONObjectItem(k string, v interface{}) jsonObject {
	hm.Set(k, v)
	return hm
}

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(lm)
func (hm *HashMap) MarshalJSON() (res []byte, err error) {
	return jsonMarshalHashMap(hm.hash)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, lm)
func (hm *HashMap) UnmarshalJSON(data []byte) error {
	hm.Clear()
	ju := &jsonUnmarshaler{
		newArray:  newJSONArray,
		newObject: newJSONObjectAsHashMap,
	}
	return ju.unmarshalJSONObject(data, hm)
}
