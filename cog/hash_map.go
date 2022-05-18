//go:build go1.18
// +build go1.18

package cog

import (
	"encoding/json"
)

// NewHashMap creates a new HashMap.
func NewHashMap[K comparable, V any](kvs ...P[K, V]) *HashMap[K, V] {
	hm := &HashMap[K, V]{}
	hm.SetPairs(kvs...)
	return hm
}

// AsHashMap creates a new HashMap from a map.
// Example: AsHashMap(map[K]V{"k1": "v1", "k2": "v2"})
func AsHashMap[K comparable, V any](m map[K]V) *HashMap[K, V] {
	hm := &HashMap[K, V]{m}
	return hm
}

// HashMap hash map type
type HashMap[K comparable, V any] struct {
	hash map[K]V
}

// lazyInit lazily initializes a zero HashMap value.
func (hm *HashMap[K, V]) lazyInit() {
	if hm.hash == nil {
		hm.hash = make(map[K]V)
	}
}

//-----------------------------------------------------------
// implements Container interface

// Len returns the length of the linked map.
func (hm *HashMap[K, V]) Len() int {
	return len(hm.hash)
}

// IsEmpty returns true if the map has no items
func (hm *HashMap[K, V]) IsEmpty() bool {
	return len(hm.hash) == 0
}

// Clear clears the map
func (hm *HashMap[K, V]) Clear() {
	// for AsHashMap()
	for k := range hm.hash {
		delete(hm.hash, k)
	}
}

//-----------------------------------------------------------
// implements Map interface

// Get looks for the given key, and returns the value associated with it,
// or nil if not found. The boolean it returns says whether the key is ok in the map.
func (hm *HashMap[K, V]) Get(key K) (v V, ok bool) {
	if hm.hash == nil {
		return
	}

	v, ok = hm.hash[key]
	return
}

// Set sets the paired key-value items, and returns what `Get` would have returned
// on that key prior to the call to `Set`.
func (hm *HashMap[K, V]) Set(key K, value V) (ov V, ok bool) {
	hm.lazyInit()

	ov, ok = hm.hash[key]
	hm.hash[key] = value
	return
}

// SetPairs set items from key-value items array, override the existing items
func (hm *HashMap[K, V]) SetPairs(pairs ...P[K, V]) {
	for _, p := range pairs {
		hm.Set(p.Key, p.Value)
	}
}

// SetAll set items from another map am, override the existing items
func (hm *HashMap[K, V]) SetAll(am Map[K, V]) {
	if hm == am {
		return
	}

	if im, ok := am.(Iterable2[K, V]); ok {
		it := im.Iterator()
		for it.Next() {
			hm.Set(it.Key(), it.Value())
		}
		return
	}

	am.Each(func(k K, v V) {
		hm.Set(k, v)
	})
}

// SetIfAbsent sets the key-value item if the key does not exists in the map,
// and returns what `Get` would have returned
// on that key prior to the call to `Set`.
func (hm *HashMap[K, V]) SetIfAbsent(key K, value V) (ov V, ok bool) {
	hm.lazyInit()

	if ov, ok = hm.hash[key]; !ok {
		hm.hash[key] = value
	}
	return
}

// Delete delete all items with key of ks,
// and returns what `Get` would have returned
// on that key prior to the call to `Set`.
func (hm *HashMap[K, V]) Delete(ks ...K) (ov V, ok bool) {
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
func (hm *HashMap[K, V]) Contains(ks ...K) bool {
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
func (hm *HashMap[K, V]) Keys() []K {
	ks := make([]K, hm.Len())
	i := 0
	for k := range hm.hash {
		ks[i] = k
		i++
	}
	return ks
}

// Values returns the value slice
func (hm *HashMap[K, V]) Values() []V {
	vs := make([]V, hm.Len())
	i := 0
	for _, v := range hm.hash {
		vs[i] = v
		i++
	}
	return vs
}

// Each call f for each item(k,v) in the map
func (hm *HashMap[K, V]) Each(f func(k K, v V)) {
	for k, v := range hm.hash {
		f(k, v)
	}
}

// HashMap returns underlying hash map
func (hm *HashMap[K, V]) HashMap() map[K]V {
	return hm.hash
}

// String print map to string
func (hm *HashMap[K, V]) String() string {
	bs, _ := json.Marshal(hm)
	return string(bs)
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(lm)
func (hm *HashMap[K, V]) MarshalJSON() (res []byte, err error) {
	return jsonMarshalMap[K, V](hm)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, lm)
func (hm *HashMap[K, V]) UnmarshalJSON(data []byte) error {
	hm.Clear()
	hm.lazyInit()
	return jsonUnmarshalMap[K, V](data, hm)
}
