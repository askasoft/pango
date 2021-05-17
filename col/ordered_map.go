package col

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// OrderedMap implements an ordered map that keeps track of the order in which keys were inserted.
type OrderedMap struct {
	entries map[interface{}]*MapEntry
	list    *List
}

// NewOrderedMap creates a new OrderedMap.
// Example: NewOrderedMap("k1", "v1", "k2", "v2")
func NewOrderedMap(kvs ...interface{}) *OrderedMap {
	om := &OrderedMap{
		entries: make(map[interface{}]*MapEntry),
		list:    NewList(),
	}
	for i := 0; i+1 < len(kvs); i += 2 {
		om.Set(kvs[i], kvs[i+1])
	}
	return om
}

// GetEntry looks for the given key, and returns the entry associated with it,
// or nil if not found. The MapEntry struct can then be used to iterate over the ordered map
// from that point, either forward or backward.
func (om *OrderedMap) GetEntry(key interface{}) *MapEntry {
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

	me := &MapEntry{
		key:   key,
		Value: value,
	}
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
	om.entries = make(map[interface{}]*MapEntry)
	om.list.Clear()
}

// Front returns a pointer to the oldest entry. It's meant to be used to iterate on the ordered map's
// entries from the oldest to the newest, e.g.:
// for entry := orderedMap.Front(); entry != nil; entry = entry.Next() { fmt.Printf("%v => %v\n", entry.Key(), entry.Value()) }
func (om *OrderedMap) Front() *MapEntry {
	return toMapEntry(om.list.Front())
}

// Back returns a pointer to the newest entry. It's meant to be used to iterate on the ordered map's
// entries from the newest to the oldest, e.g.:
// for entry := orderedMap.Back(); entry != nil; entry = entry.Prev() { fmt.Printf("%v => %v\n", entry.Key(), entry.Value()) }
func (om *OrderedMap) Back() *MapEntry {
	return toMapEntry(om.list.Back())
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
func (om *OrderedMap) Entries() []*MapEntry {
	vs := make([]*MapEntry, 0, om.Len())
	for me := om.Front(); me != nil; me = me.Next() {
		vs = append(vs, me)
	}
	return vs
}

// Each Call f for each item in the map
func (om *OrderedMap) Each(f func(*MapEntry)) {
	for me := om.Front(); me != nil; me = me.Next() {
		f(me)
	}
}

// ReverseEach Call f for each item in the map with reverse order
func (om *OrderedMap) ReverseEach(f func(*MapEntry)) {
	for me := om.Back(); me != nil; me = me.Prev() {
		f(me)
	}
}

/*------------- JSON -----------------*/

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
	dec := json.NewDecoder(bytes.NewReader(data))

	// must open with a delim token '{'
	t, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expect JSON object open with '{'")
	}

	err = om.parseJSONObject(dec)
	if err != nil {
		return err
	}

	t, err = dec.Token()
	if err != io.EOF {
		return fmt.Errorf("expect end of JSON object but got more token: %T: %v or err: %v", t, t, err)
	}

	return nil
}

func (om *OrderedMap) parseJSONObject(dec *json.Decoder) (err error) {
	var t json.Token
	for dec.More() {
		t, err = dec.Token()
		if err != nil {
			return err
		}

		key, ok := t.(string)
		if !ok {
			return fmt.Errorf("expecting JSON key should be always a string: %T: %v", t, t)
		}

		t, err = dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		var value interface{}
		value, err = handleDelim(t, dec)
		if err != nil {
			return err
		}

		om.Set(key, value)
	}

	t, err = dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '}' {
		return fmt.Errorf("expect JSON object close with '}'")
	}

	return nil
}

func parseJSONArray(dec *json.Decoder) (arr []interface{}, err error) {
	var t json.Token
	arr = make([]interface{}, 0)
	for dec.More() {
		t, err = dec.Token()
		if err != nil {
			return
		}

		var value interface{}
		value, err = handleDelim(t, dec)
		if err != nil {
			return
		}
		arr = append(arr, value)
	}
	t, err = dec.Token()
	if err != nil {
		return
	}
	if delim, ok := t.(json.Delim); !ok || delim != ']' {
		err = fmt.Errorf("expect JSON array close with ']'")
		return
	}

	return
}

func handleDelim(t json.Token, dec *json.Decoder) (res interface{}, err error) {
	if delim, ok := t.(json.Delim); ok {
		switch delim {
		case '{':
			om2 := NewOrderedMap()
			err = om2.parseJSONObject(dec)
			if err != nil {
				return
			}
			return om2, nil
		case '[':
			var value []interface{}
			value, err = parseJSONArray(dec)
			if err != nil {
				return
			}
			return value, nil
		default:
			return nil, fmt.Errorf("Unexpected delimiter: %q", delim)
		}
	}
	return t, nil
}
