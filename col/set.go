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

// Contains Test to see whether or not the v is in the set
func (s *Set) Contains(v interface{}) bool {
	_, ok := s.hash[v]
	return ok
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

// Len Return the number of items in the set
func (s *Set) Len() int {
	return len(s.hash)
}

// Remove an v from the set
func (s *Set) Remove(v interface{}) {
	delete(s.hash, v)
}

// Each Call f for each item in the set
func (s *Set) Each(f func(interface{})) {
	for k := range s.hash {
		f(k)
	}
}

// SubsetOf Test whether or not s set is a subset of "set"
func (s *Set) SubsetOf(a *Set) bool {
	if s.Len() > a.Len() {
		return false
	}
	for k := range s.hash {
		if _, ok := a.hash[k]; !ok {
			return false
		}
	}
	return true
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
