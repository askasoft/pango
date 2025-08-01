package hashmap

import (
	"encoding/json"
	"fmt"
	"iter"

	"github.com/askasoft/pango/cog"
	"github.com/askasoft/pango/cog/internal/imap"
	"github.com/askasoft/pango/cog/internal/jsonmap"
	"github.com/askasoft/pango/str"
)

// NewHashMap creates a new HashMap.
func NewHashMap[K comparable, V any](kvs ...cog.P[K, V]) *HashMap[K, V] {
	hm := &HashMap[K, V]{}
	hm.SetEntries(kvs...)
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
	clear(hm.hash)
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

// MustGet looks for the given key, and returns the value associated with it.
// Panic if not found.
func (hm *HashMap[K, V]) MustGet(key K) V {
	if v, ok := hm.Get(key); ok {
		return v
	}
	panic(fmt.Errorf("HashMap key '%v' does not exist", key))
}

// SafeGet looks for the given key, and returns the value associated with it.
// If not found, return defaults[0] or zero V.
func (hm *HashMap[K, V]) SafeGet(key K, defaults ...V) V {
	if v, ok := hm.Get(key); ok {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}

	var v V
	return v
}

// Set sets the paired key-value items, and returns what `Get` would have returned
// on that key prior to the call to `Set`.
func (hm *HashMap[K, V]) Set(key K, value V) (ov V, ok bool) {
	hm.lazyInit()

	ov, ok = hm.hash[key]
	hm.hash[key] = value
	return
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

// SetEntries set items from key-value items array, override the existing items
func (hm *HashMap[K, V]) SetEntries(pairs ...cog.P[K, V]) {
	imap.SetMapPairs(hm, pairs...)
}

// Copy copy items from another map am, override the existing items
func (hm *HashMap[K, V]) Copy(am cog.Map[K, V]) {
	imap.CopyMap(hm, am)
}

// Remove remove the item with key k,
// and returns what `Get` would have returned
// on that key prior to the call to `Set`.
func (hm *HashMap[K, V]) Remove(k K) (ov V, ok bool) {
	if hm.IsEmpty() {
		return
	}

	ov, ok = hm.hash[k]
	if ok {
		delete(hm.hash, k)
	}
	return
}

// RemoveAll remove all items with key of ks.
func (hm *HashMap[K, V]) RemoveAll(ks ...K) {
	if hm.IsEmpty() {
		return
	}

	for _, k := range ks {
		delete(hm.hash, k)
	}
}

// Contains Test to see if the map contains the given key
func (hm *HashMap[K, V]) Contains(k K) bool {
	if hm.IsEmpty() {
		return false
	}

	if _, ok := hm.hash[k]; ok {
		return true
	}
	return false
}

// ContainsAll Test to see if the map contains all keys of ks
func (hm *HashMap[K, V]) ContainsAll(ks ...K) bool {
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

// Entries returns the key-value pair slice
func (hm *HashMap[K, V]) Entries() []cog.P[K, V] {
	ps := make([]cog.P[K, V], hm.Len())
	i := 0
	for k, v := range hm.hash {
		ps[i] = cog.KV(k, v)
		i++
	}
	return ps
}

// Each call f for each item(k,v) in the map
func (hm *HashMap[K, V]) Each(f func(K, V) bool) {
	for k, v := range hm.hash {
		if !f(k, v) {
			return
		}
	}
}

// Seq returns a iter.Seq[K, V] for range
func (hm *HashMap[K, V]) Seq() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range hm.hash {
			if !yield(k, v) {
				return
			}
		}
	}
}

// HashMap returns underlying hash map
func (hm *HashMap[K, V]) HashMap() map[K]V {
	return hm.hash
}

// String print map to string
func (hm *HashMap[K, V]) String() string {
	bs, _ := json.Marshal(hm)
	return str.UnsafeString(bs)
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(lm)
func (hm *HashMap[K, V]) MarshalJSON() ([]byte, error) {
	return jsonmap.JsonMarshalMap(hm)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, lm)
func (hm *HashMap[K, V]) UnmarshalJSON(data []byte) error {
	hm.Clear()
	hm.lazyInit()
	return jsonmap.JsonUnmarshalMap(data, hm)
}
