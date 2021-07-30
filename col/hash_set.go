package col

import (
	"encoding/json"
	"fmt"
)

// NewHashSet Create a new hash set
func NewHashSet(vs ...interface{}) *HashSet {
	hs := &HashSet{}
	hs.Add(vs...)
	return hs
}

// NewStringHashSet Create a new hash set
func NewStringHashSet(vs ...string) *HashSet {
	hs := &HashSet{}
	if len(vs) > 0 {
		hs.lazyInit()
		for _, v := range vs {
			hs.hash[v] = true
		}
	}
	return hs
}

// HashSet an unordered collection of unique values.
// The zero value for HashSet is an empty set ready to use.
// http://en.wikipedia.org/wiki/Set_(computer_science%29)
type HashSet struct {
	hash map[interface{}]bool
}

// lazyInit lazily initializes a zero List value.
func (hs *HashSet) lazyInit() {
	if hs.hash == nil {
		hs.hash = make(map[interface{}]bool)
	}
}

//-----------------------------------------------------------
// implements Collection interface

// Len Return the number of items in the set
func (hs *HashSet) Len() int {
	return len(hs.hash)
}

// IsEmpty returns true if the set's length == 0
func (hs *HashSet) IsEmpty() bool {
	return hs.Len() == 0
}

// Clear clears the hash set.
func (hs *HashSet) Clear() {
	hs.hash = nil
}

// Add Adds all items of vs to the set
func (hs *HashSet) Add(vs ...interface{}) {
	if len(vs) == 0 {
		return
	}

	hs.lazyInit()
	for _, v := range vs {
		hs.hash[v] = true
	}
}

// AddAll adds all items of another collection
func (hs *HashSet) AddAll(ac Collection) {
	hs.Add(ac.Values()...)
}

// Delete delete items of vs
func (hs *HashSet) Delete(vs ...interface{}) {
	if len(hs.hash) == 0 {
		return
	}

	for _, v := range vs {
		delete(hs.hash, v)
	}
}

// DeleteAll delete all of this collection's elements that are also contained in the specified collection
func (hs *HashSet) DeleteAll(ac Collection) {
	hs.Delete(ac.Values()...)
}

// Contains Test to see if the collection contains all items of vs
func (hs *HashSet) Contains(vs ...interface{}) bool {
	if len(vs) == 0 {
		return true
	}

	if hs.IsEmpty() {
		return false
	}

	for _, v := range vs {
		if _, ok := hs.hash[v]; !ok {
			return false
		}
	}
	return true
}

// ContainsAll Test to see if the collection contains all items of another collection
func (hs *HashSet) ContainsAll(ac Collection) bool {
	if hs == ac {
		return true
	}
	return hs.Contains(ac.Values()...)
}

// Retain Retains only the elements in this collection that are contained in the argument array vs.
func (hs *HashSet) Retain(vs ...interface{}) {
	if hs.IsEmpty() || len(vs) == 0 {
		return
	}

	hs.RetainAll(NewArrayList(vs...))
}

// RetainAll Retains only the elements in this collection that are contained in the specified collection.
func (hs *HashSet) RetainAll(ac Collection) {
	if hs.IsEmpty() || ac.IsEmpty() || hs == ac {
		return
	}

	for k := range hs.hash {
		if !ac.Contains(k) {
			delete(hs.hash, k)
		}
	}
}

// Values returns a slice contains all the items of the set hs
func (hs *HashSet) Values() []interface{} {
	vs := make([]interface{}, hs.Len())
	i := 0
	for k := range hs.hash {
		vs[i] = k
		i++
	}
	return vs
}

// Each Call f for each item in the set
func (hs *HashSet) Each(f func(interface{})) {
	for k := range hs.hash {
		f(k)
	}
}

//-----------------------------------------------------------

// AddSet Add values of another set a
func (hs *HashSet) AddSet(a *HashSet) {
	if a.IsEmpty() {
		return
	}

	hs.lazyInit()
	for k := range a.hash {
		hs.hash[k] = true
	}
}

// ContainsSet returns true if HashSet hs contains the HashSet a.
func (hs *HashSet) ContainsSet(a *HashSet) bool {
	if hs.Len() < a.Len() {
		return false
	}

	for k := range a.hash {
		if !a.hash[k] {
			return false
		}
	}
	return true
}

// Difference Find the difference btween two sets
func (hs *HashSet) Difference(a *HashSet) *HashSet {
	b := make(map[interface{}]bool)

	for k := range hs.hash {
		if _, ok := a.hash[k]; !ok {
			b[k] = true
		}
	}

	return &HashSet{b}
}

// Intersection Find the intersection of two sets
func (hs *HashSet) Intersection(a *HashSet) *HashSet {
	b := make(map[interface{}]bool)

	for k := range hs.hash {
		if _, ok := a.hash[k]; ok {
			b[k] = true
		}
	}

	return &HashSet{b}
}

// String print the set to string
func (hs *HashSet) String() string {
	return fmt.Sprintf("%v", hs.hash)
}

//-----------------------------------------------------------
// implements JSON Marshaller/Unmarshaller interface

func (hs *HashSet) addJSONArrayItem(v interface{}) jsonArray {
	hs.Add(v)
	return hs
}

// MarshalJSON implements type json.Marshaler interface, so can be called in json.Marshal(hs)
func (hs *HashSet) MarshalJSON() (res []byte, err error) {
	if hs.IsEmpty() {
		return []byte("[]"), nil
	}

	res = append(res, '[')
	for v := range hs.hash {
		var b []byte
		b, err = json.Marshal(v)
		if err != nil {
			return
		}
		res = append(res, b...)
		res = append(res, ',')
	}
	res[len(res)-1] = ']'
	return
}

// UnmarshalJSON implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, hs)
func (hs *HashSet) UnmarshalJSON(data []byte) error {
	ju := &jsonUnmarshaler{
		newArray:  newJSONArray,
		newObject: newJSONObject,
	}
	return ju.unmarshalJSONArray(data, hs)
}
