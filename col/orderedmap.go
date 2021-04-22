package col

// OrderedMap implements an ordered map that keeps track of the order in which keys were inserted.
type OrderedMap struct {
	entries map[interface{}]*MapEntry
	list    *List
}

// NewOrderedMap creates a new OrderedMap.
func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		entries: make(map[interface{}]*MapEntry),
		list:    NewList(),
	}
}

// GetEntry looks for the given key, and returns the entry associated with it,
// or nil if not found. The MapEntry struct can then be used to iterate over the ordered map
// from that point, either forward or backward.
func (om *OrderedMap) GetEntry(key interface{}) *MapEntry {
	return om.entries[key]
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
	for e := om.Front(); e != nil; e = e.Next() {
		ks = append(ks, e.Key())
	}
	return ks
}

// Values returns the value slice
func (om *OrderedMap) Values() []interface{} {
	vs := make([]interface{}, 0, om.Len())
	for e := om.Front(); e != nil; e = e.Next() {
		vs = append(vs, e.Value)
	}
	return vs
}

// Each Call f for each item in the map
func (om *OrderedMap) Each(f func(*MapEntry)) {
	for e := om.Front(); e != nil; e = e.Next() {
		f(e)
	}
}

// ReverseEach Call f for each item in the map with reverse order
func (om *OrderedMap) ReverseEach(f func(*MapEntry)) {
	for e := om.Back(); e != nil; e = e.Prev() {
		f(e)
	}
}
