//go:build go1.18
// +build go1.18

package linkedhashmap

import (
	"encoding/json"
	"fmt"

	"github.com/askasoft/pango/cog"
	"github.com/askasoft/pango/cog/internal/imap"
	"github.com/askasoft/pango/cog/internal/jsonmap"
	"github.com/askasoft/pango/str"
)

// NewLinkedHashMap creates a new LinkedHashMap.
func NewLinkedHashMap[K comparable, V any](kvs ...cog.P[K, V]) *LinkedHashMap[K, V] {
	lm := &LinkedHashMap[K, V]{}
	lm.SetEntries(kvs...)
	return lm
}

// LinkedHashMap implements an linked map that keeps track of the order in which keys were inserted.
// The zero value for LinkedHashMap is an empty map ready to use.
//
// To iterate over a linked map (where lm is a *LinkedHashMap):
//
//	it := lm.Iterator()
//	for it.Next() {
//		// do something with it.Key(), it.Value()
//	}
type LinkedHashMap[K comparable, V any] struct {
	head, tail *LinkedMapNode[K, V]
	hash       map[K]*LinkedMapNode[K, V]
}

//-----------------------------------------------------------
// implements Collection interface

// Len returns the length of the linked map.
func (lm *LinkedHashMap[K, V]) Len() int {
	return len(lm.hash)
}

// IsEmpty returns true if the map has no items
func (lm *LinkedHashMap[K, V]) IsEmpty() bool {
	return len(lm.hash) == 0
}

// Clear clears the map
func (lm *LinkedHashMap[K, V]) Clear() {
	lm.hash = nil
	lm.head = nil
	lm.tail = nil
}

//-----------------------------------------------------------
// implements Map interface

// Get looks for the given key, and returns the value associated with it,
// or nil if not found. The boolean it returns says whether the key is ok in the map.
func (lm *LinkedHashMap[K, V]) Get(key K) (V, bool) {
	if lm.hash != nil {
		if ln, ok := lm.hash[key]; ok {
			return ln.value, ok
		}
	}
	var v V
	return v, false
}

// MustGet looks for the given key, and returns the value associated with it.
// Panic if not found.
func (lm *LinkedHashMap[K, V]) MustGet(key K) V {
	if v, ok := lm.Get(key); ok {
		return v
	}
	panic(fmt.Errorf("LinkedHashMap key '%v' does not exist", key))
}

// SafeGet looks for the given key, and returns the value associated with it.
// If not found, return defaults[0] or zero V.
func (lm *LinkedHashMap[K, V]) SafeGet(key K, defaults ...V) V {
	if v, ok := lm.Get(key); ok {
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
func (lm *LinkedHashMap[K, V]) Set(key K, value V) (ov V, ok bool) {
	var ln *LinkedMapNode[K, V]
	if ln, ok = lm.hash[key]; ok {
		ov = ln.value
		ln.value = value
	} else {
		lm.add(key, value)
	}
	return
}

// SetIfAbsent sets the key-value item if the key does not exists in the map,
// and returns what `Get` would have returned
// on that key prior to the call to `Set`.
// Example: lm.SetIfAbsent("k1", "v1", "k2", "v2")
func (lm *LinkedHashMap[K, V]) SetIfAbsent(key K, value V) (ov V, ok bool) {
	var ln *LinkedMapNode[K, V]

	if ln, ok = lm.hash[key]; ok {
		ov = ln.value
	} else {
		lm.add(key, value)
	}

	return
}

// SetEntries set items from key-value items array, override the existing items
func (lm *LinkedHashMap[K, V]) SetEntries(pairs ...cog.P[K, V]) {
	imap.SetMapPairs[K, V](lm, pairs...)
}

// Copy copy items from another map am, override the existing items
func (lm *LinkedHashMap[K, V]) Copy(am cog.Map[K, V]) {
	imap.CopyMap[K, V](lm, am)
}

// Remove remove the item with key k,
// and returns what `Get` would have returned
// on that key prior to the call to `Set`.
func (lm *LinkedHashMap[K, V]) Remove(k K) (ov V, ok bool) {
	if lm.IsEmpty() {
		return
	}

	var ln *LinkedMapNode[K, V]
	ln, ok = lm.hash[k]
	if ok {
		ov = ln.value
		lm.deleteNode(ln)
	}
	return
}

// Remove remove all items with key of ks.
func (lm *LinkedHashMap[K, V]) Removes(ks ...K) {
	if lm.IsEmpty() {
		return
	}

	for _, k := range ks {
		if ln, ok := lm.hash[k]; ok {
			lm.deleteNode(ln)
		}
	}
}

// Contain Test to see if the list contains the key k
func (lm *LinkedHashMap[K, V]) Contain(k K) bool {
	if lm.IsEmpty() {
		return false
	}

	if _, ok := lm.hash[k]; ok {
		return true
	}
	return false
}

// Contains looks for the given key, and returns true if the key exists in the map.
func (lm *LinkedHashMap[K, V]) Contains(ks ...K) bool {
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
func (lm *LinkedHashMap[K, V]) Keys() []K {
	ks := make([]K, lm.Len())
	for i, ln := 0, lm.head; ln != nil; i, ln = i+1, ln.next {
		ks[i] = ln.key
	}
	return ks
}

// Values returns the value slice
func (lm *LinkedHashMap[K, V]) Values() []V {
	vs := make([]V, lm.Len())
	for i, ln := 0, lm.head; ln != nil; i, ln = i+1, ln.next {
		vs[i] = ln.value
	}
	return vs
}

// Entries returns the key-value pair slice
func (lm *LinkedHashMap[K, V]) Entries() []cog.P[K, V] {
	ps := make([]cog.P[K, V], lm.Len())
	for i, ln := 0, lm.head; ln != nil; i, ln = i+1, ln.next {
		ps[i] = cog.KV(ln.key, ln.value)
	}
	return ps
}

// Each call f for each item in the map
func (lm *LinkedHashMap[K, V]) Each(f func(K, V) bool) {
	for ln := lm.head; ln != nil; ln = ln.next {
		if !f(ln.key, ln.value) {
			return
		}
	}
}

// ReverseEach call f for each item in the map with reverse order
func (lm *LinkedHashMap[K, V]) ReverseEach(f func(K, V) bool) {
	for ln := lm.tail; ln != nil; ln = ln.prev {
		if !f(ln.key, ln.value) {
			return
		}
	}
}

// Iterator returns a iterator for the map
func (lm *LinkedHashMap[K, V]) Iterator() cog.Iterator2[K, V] {
	return &linkedHashMapIterator[K, V]{lmap: lm}
}

// IteratorOf returns a iterator at the specified key
// Returns nil if the map does not contains the key
func (lm *LinkedHashMap[K, V]) IteratorOf(k K) cog.Iterator2[K, V] {
	if lm.IsEmpty() {
		return nil
	}
	if ln, ok := lm.hash[k]; ok {
		return &linkedHashMapIterator[K, V]{lmap: lm, node: ln}
	}
	return nil
}

//-----------------------------------------------------------

// Head returns the oldest key/value item.
func (lm *LinkedHashMap[K, V]) Head() *LinkedMapNode[K, V] {
	return lm.head
}

// Tail returns the newest key/value item.
func (lm *LinkedHashMap[K, V]) Tail() *LinkedMapNode[K, V] {
	return lm.tail
}

// PollHead remove the first item of map.
func (lm *LinkedHashMap[K, V]) PollHead() *LinkedMapNode[K, V] {
	ln := lm.head
	if ln != nil {
		lm.deleteNode(ln)
	}
	return ln
}

// PollTail remove the last item of map.
func (lm *LinkedHashMap[K, V]) PollTail() *LinkedMapNode[K, V] {
	ln := lm.tail
	if ln != nil {
		lm.deleteNode(ln)
	}
	return ln
}

// Items returns the map item slice
func (lm *LinkedHashMap[K, V]) Items() []*LinkedMapNode[K, V] {
	mis := make([]*LinkedMapNode[K, V], lm.Len())
	for i, ln := 0, lm.head; ln != nil; i, ln = i+1, ln.next {
		mis[i] = ln
	}
	return mis
}

// String print map to string
func (lm *LinkedHashMap[K, V]) String() string {
	bs, _ := json.Marshal(lm)
	return str.UnsafeString(bs)
}

//-----------------------------------------------------

func (lm *LinkedHashMap[K, V]) add(k K, v V) {
	ln := &LinkedMapNode[K, V]{prev: lm.tail, key: k, value: v}
	if ln.prev == nil {
		lm.head = ln
	} else {
		ln.prev.next = ln
	}
	lm.tail = ln

	if lm.hash == nil {
		lm.hash = make(map[K]*LinkedMapNode[K, V])
	}
	lm.hash[k] = ln
}

func (lm *LinkedHashMap[K, V]) deleteNode(ln *LinkedMapNode[K, V]) {
	if ln.prev == nil {
		lm.head = ln.next
	} else {
		ln.prev.next = ln.next
	}

	if ln.next == nil {
		lm.tail = ln.prev
	} else {
		ln.next.prev = ln.prev
	}

	delete(lm.hash, ln.key)
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(lm)
func (lm *LinkedHashMap[K, V]) MarshalJSON() ([]byte, error) {
	return jsonmap.JsonMarshalMap[K, V](lm)
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, lm)
func (lm *LinkedHashMap[K, V]) UnmarshalJSON(data []byte) error {
	lm.Clear()
	return jsonmap.JsonUnmarshalMap[K, V](data, lm)
}
