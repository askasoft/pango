package col

import "testing"

func TestSimple(t *testing.T) {
	s := NewSet()

	s.Add(5)

	if s.Len() != 1 {
		t.Errorf("Length should be 1")
	}

	if !s.Contains(5) {
		t.Errorf("Membership test failed")
	}

	s.Remove(5)

	if s.Len() != 0 {
		t.Errorf("Length should be 0")
	}

	if s.Contains(5) {
		t.Errorf("The set should be empty")
	}
}

func TestContainsSet(t *testing.T) {
	s1 := NewSet()
	s2 := NewSet()

	if !s1.ContainsSet(s1) {
		t.Errorf("set should be a subset of itself")
	}

	if !s1.ContainsSet(s2) {
		t.Errorf("empty set should contains another empty set")
	}

	s1.Add(1)
	if !s1.ContainsSet(s2) {
		t.Errorf("set should contains another empty set")
	}

	s2.Add(1)
	if !s1.ContainsSet(s2) {
		t.Errorf("set should contains another same set")
	}

	s1.Add(2)
	if !s1.ContainsSet(s2) {
		t.Errorf("set should contains another small set")
	}

	s2.Add(3)
	if !s1.ContainsSet(s2) {
		t.Errorf("set should not contains another different set")
	}
}

func TestDifference(t *testing.T) {
	// Difference
	s1 := NewSet(1, 2, 3, 4, 5, 6)
	s2 := NewSet(4, 5, 6)
	s3 := s1.Difference(s2)

	if s3.Len() != 3 {
		t.Errorf("Length should be 3")
	}

	if !(s3.Contains(1) && s3.Contains(2) && s3.Contains(3)) {
		t.Errorf("Set should only contain 1, 2, 3")
	}
}

func TestIntersection(t *testing.T) {
	s1 := NewSet(1, 2, 3, 4, 5, 6)
	s2 := NewSet(4, 5, 6)

	// Intersection
	s3 := s1.Intersection(s2)
	if s3.Len() != 3 {
		t.Errorf("Length should be 3 after intersection")
	}

	if !(s3.Contains(4) && s3.Contains(5) && s3.Contains(6)) {
		t.Errorf("Set should contain 4, 5, 6")
	}
}

func TestAddSet(t *testing.T) {
	// AddSet
	s1 := NewSet(4, 5, 6)
	s2 := NewSet(7, 8, 9)
	s1.AddSet(s2)

	if s1.Len() != 6 {
		t.Errorf("Length should be 6 after union")
	}

	for i := 4; i <= 9; i++ {
		if !(s1.Contains(i)) {
			t.Errorf("Set should contains %d", i)
		}
	}
}
