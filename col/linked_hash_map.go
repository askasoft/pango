package col

import (
	"encoding/json"
	"fmt"
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
//	it := lm.Iterator()
//	for it.Next() {
//		// do something with it.Key(), it.Value()
//	}
//
type LinkedHashMap struct {
	front, back *LinkedMapNode
	hash        map[interface{}]*LinkedMapNode
}

//-----------------------------------------------------------
// implements Collection interface

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
	lm.front = nil
	lm.back = nil
}

//-----------------------------------------------------------
// implements Map interface

// Get looks for the given key, and returns the value associated with it,
// or nil if not found. The boolean it returns says whether the key is ok in the map.
func (lm *LinkedHashMap) Get(key interface{}) (interface{}, bool) {
	if lm.hash != nil {
		if ln, ok := lm.hash[key]; ok {
			return ln.value, ok
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

	var ln *LinkedMapNode
	for i := 0; i+1 < len(kvs); i += 2 {
		k := kvs[i]
		v := kvs[i+1]

		ov = nil
		if ln, ok = lm.hash[k]; ok {
			ov = ln.value
			ln.value = v
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

	var ln *LinkedMapNode
	for i := 0; i+1 < len(kvs); i += 2 {
		k := kvs[i]
		v := kvs[i+1]

		ov = nil
		if ln, ok = lm.hash[k]; ok {
			ov = ln.value
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

	var ln *LinkedMapNode
	for _, k := range ks {
		if ln, ok = lm.hash[k]; ok {
			lm.deleteNode(ln)
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
	for i, ln := 0, lm.front; ln != nil; i, ln = i+1, ln.next {
		ks[i] = ln.key
	}
	return ks
}

// Values returns the value slice
func (lm *LinkedHashMap) Values() []interface{} {
	vs := make([]interface{}, lm.Len())
	for i, ln := 0, lm.front; ln != nil; i, ln = i+1, ln.next {
		vs[i] = ln.value
	}
	return vs
}

// Each call f for each item in the map
func (lm *LinkedHashMap) Each(f func(k interface{}, v interface{})) {
	for ln := lm.front; ln != nil; ln = ln.next {
		f(ln.key, ln.value)
	}
}

// ReverseEach call f for each item in the map with reverse order
func (lm *LinkedHashMap) ReverseEach(f func(k interface{}, v interface{})) {
	for ln := lm.back; ln != nil; ln = ln.prev {
		f(ln.key, ln.value)
	}
}

// Iterator returns a iterator for the map
func (lm *LinkedHashMap) Iterator() Iterator2 {
	return &linkedHashMapIterator{lmap: lm}
}

// IteratorOf returns a iterator at the specified key
// Returns nil if the map does not contains the key
func (lm *LinkedHashMap) IteratorOf(k interface{}) Iterator2 {
	if lm.IsEmpty() {
		return nil
	}
	if ln, ok := lm.hash[k]; ok {
		return &linkedHashMapIterator{lmap: lm, node: ln}
	}
	return nil
}

//-----------------------------------------------------------

// Front returns the oldest key/value item.
func (lm *LinkedHashMap) Front() *LinkedMapNode {
	return lm.front
}

// Back returns the newest key/value item.
func (lm *LinkedHashMap) Back() *LinkedMapNode {
	return lm.back
}

// PopFront remove the first item of map.
func (lm *LinkedHashMap) PopFront() *LinkedMapNode {
	ln := lm.front
	if ln != nil {
		lm.deleteNode(ln)
	}
	return ln
}

// PopBack remove the last item of map.
func (lm *LinkedHashMap) PopBack() *LinkedMapNode {
	ln := lm.back
	if ln != nil {
		lm.deleteNode(ln)
	}
	return ln
}

// Items returns the map item slice
func (lm *LinkedHashMap) Items() []*LinkedMapNode {
	mis := make([]*LinkedMapNode, lm.Len(), lm.Len())
	for i, ln := 0, lm.front; ln != nil; i, ln = i+1, ln.next {
		mis[i] = ln
	}
	return mis
}

// String print map to string
func (lm *LinkedHashMap) String() string {
	bs, _ := json.Marshal(lm)
	return string(bs)
}

//-----------------------------------------------------

func (lm *LinkedHashMap) add(k interface{}, v interface{}) {
	ln := &LinkedMapNode{prev: lm.back, key: k, value: v}
	if ln.prev == nil {
		lm.front = ln
	} else {
		ln.prev.next = ln
	}
	lm.back = ln

	if lm.hash == nil {
		lm.hash = make(map[interface{}]*LinkedMapNode)
	}
	lm.hash[k] = ln
}

func (lm *LinkedHashMap) deleteNode(ln *LinkedMapNode) {
	if ln.prev == nil {
		lm.front = ln.next
	} else {
		ln.prev.next = ln.next
	}

	if ln.next == nil {
		lm.back = ln.prev
	} else {
		ln.next.prev = ln.prev
	}

	delete(lm.hash, ln.key)
}

//-----------------------------------------------------

// LinkedMapNode is a node of a linked hash map.
type LinkedMapNode struct {
	prev  *LinkedMapNode
	next  *LinkedMapNode
	key   interface{}
	value interface{}
}

// Key returns the key
func (ln *LinkedMapNode) Key() interface{} {
	return ln.key
}

// Value returns the key
func (ln *LinkedMapNode) Value() interface{} {
	return ln.value
}

// SetValue sets the value
func (ln *LinkedMapNode) SetValue(v interface{}) {
	ln.value = v
}

// String print the list item to string
func (ln *LinkedMapNode) String() string {
	return fmt.Sprintf("%v => %v", ln.key, ln.value)
}

// linkedHashMapIterator a iterator for LinkedMap
type linkedHashMapIterator struct {
	lmap    *LinkedHashMap
	node    *LinkedMapNode
	removed bool
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (it *linkedHashMapIterator) Prev() bool {
	if it.lmap.IsEmpty() {
		return false
	}

	if it.node == nil {
		it.node = it.lmap.back
		it.removed = false
		return true
	}

	if pi := it.node.prev; pi != nil {
		it.node = pi
		it.removed = false
		return true
	}
	return false
}

// Next moves the iterator to the next element and returns true if there was a next element in the collection.
// If Next() returns true, then next element's value can be retrieved by Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (it *linkedHashMapIterator) Next() bool {
	if it.lmap.IsEmpty() {
		return false
	}

	if it.node == nil {
		it.node = it.lmap.front
		it.removed = false
		return true
	}

	if ni := it.node.next; ni != nil {
		it.node = ni
		it.removed = false
		return true
	}
	return false
}

// Key returns the current element's key.
func (it *linkedHashMapIterator) Key() interface{} {
	if it.node == nil {
		return nil
	}
	return it.node.key
}

// Value returns the current element's value.
func (it *linkedHashMapIterator) Value() interface{} {
	if it.node == nil {
		return nil
	}
	return it.node.value
}

// SetValue set the value to the item
func (it *linkedHashMapIterator) SetValue(v interface{}) {
	if it.node != nil {
		it.node.value = v
	}
}

// Remove remove the current element
func (it *linkedHashMapIterator) Remove() {
	if it.node == nil {
		return
	}

	if it.removed {
		panic("LinkedHashMap can't remove a unlinked item")
	}

	it.lmap.deleteNode(it.node)
	it.removed = true
}

// Reset resets the iterator to its initial state (one-before-first/one-after-last)
// Call Next()/Prev() to fetch the first/last element if any.
func (it *linkedHashMapIterator) Reset() {
	it.node = nil
	it.removed = false
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
