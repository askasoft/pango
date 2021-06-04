package col

import (
	"encoding/json"
	"fmt"
)

// HashSet an unordered collection of unique values.
// http://en.wikipedia.org/wiki/Set_(computer_science%29)
type HashSet struct {
	hash map[interface{}]bool
}

// NewHashSet Create a new hash set
func NewHashSet(vs ...interface{}) *HashSet {
	hs := &HashSet{make(map[interface{}]bool)}
	hs.AddAll(vs...)
	return hs
}

// Len Return the number of items in the set
func (hs *HashSet) Len() int {
	return len(hs.hash)
}

// IsEmpty returns true if the set's length == 0
func (hs *HashSet) IsEmpty() bool {
	return hs.Len() == 0
}

// Add Add an v to the set
func (hs *HashSet) Add(v interface{}) {
	hs.hash[v] = true
}

// AddAll Add values vs to the set
func (hs *HashSet) AddAll(vs ...interface{}) {
	for _, v := range vs {
		hs.hash[v] = true
	}
}

// AddSet Add values of another set a
func (hs *HashSet) AddSet(a *HashSet) {
	for k := range a.hash {
		hs.hash[k] = true
	}
}

// Clear clears the hash set.
func (hs *HashSet) Clear() {
	hs.hash = make(map[interface{}]bool)
}

// Delete an v from the set
func (hs *HashSet) Delete(v interface{}) {
	delete(hs.hash, v)
}

// Contains Test to see whether or not the v is in the set
func (hs *HashSet) Contains(v interface{}) bool {
	return hs.hash[v]
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

// Each Call f for each item in the set
func (hs *HashSet) Each(f func(interface{})) {
	for k := range hs.hash {
		f(k)
	}
}

// Values returns a slice contains all the items of the set hs
func (hs *HashSet) Values() []interface{} {
	a := make([]interface{}, 0, hs.Len())
	for k := range hs.hash {
		a = append(a, k)
	}
	return a
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

/*------------- JSON -----------------*/

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
