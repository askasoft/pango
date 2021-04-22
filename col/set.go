package col

// Set an unordered collection of unique values.
// http://en.wikipedia.org/wiki/Set_(computer_science%29)
type Set struct {
	hash map[interface{}]bool
}

// NewSet Create a new set
func NewSet(vs ...interface{}) *Set {
	s := &Set{make(map[interface{}]bool)}

	s.AddAll(vs...)

	return s
}

// Len Return the number of items in the set
func (s *Set) Len() int {
	return len(s.hash)
}

// Add Add an v to the set
func (s *Set) Add(v interface{}) {
	s.hash[v] = true
}

// AddAll Add values vs to the set
func (s *Set) AddAll(vs ...interface{}) {
	for _, v := range vs {
		s.hash[v] = true
	}
}

// AddSet Add values of another set a
func (s *Set) AddSet(a *Set) {
	for k := range a.hash {
		s.hash[k] = true
	}
}

// Remove an v from the set
func (s *Set) Remove(v interface{}) {
	delete(s.hash, v)
}

// Contains Test to see whether or not the v is in the set
func (s *Set) Contains(v interface{}) bool {
	return s.hash[v]
}

// ContainsSet returns true if Set s contains the Set a.
func (s *Set) ContainsSet(a *Set) bool {
	if s.Len() < a.Len() {
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
func (s *Set) Each(f func(interface{})) {
	for k := range s.hash {
		f(k)
	}
}

// Values returns a slice contains all the items of the set s
func (s *Set) Values() []interface{} {
	a := make([]interface{}, 0, s.Len())
	for k := range s.hash {
		a = append(a, k)
	}
	return a
}

// Difference Find the difference btween two sets
func (s *Set) Difference(a *Set) *Set {
	b := make(map[interface{}]bool)

	for k := range s.hash {
		if _, ok := a.hash[k]; !ok {
			b[k] = true
		}
	}

	return &Set{b}
}

// Intersection Find the intersection of two sets
func (s *Set) Intersection(a *Set) *Set {
	b := make(map[interface{}]bool)

	for k := range s.hash {
		if _, ok := a.hash[k]; ok {
			b[k] = true
		}
	}

	return &Set{b}
}