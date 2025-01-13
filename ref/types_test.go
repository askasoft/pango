package ref

import (
	"testing"
)

func TestIsArrayType(t *testing.T) {
	cs := []struct {
		w bool
		x any
	}{
		{true, [1]int{}},
		{false, nil},
	}

	for i, c := range cs {
		a := IsArrayType(c.x)
		if a != c.w {
			t.Errorf("[%d] IsArrayType(%T) = %v, want %v", i, c.x, a, c.w)
		}
	}
}

func TestIsSliceType(t *testing.T) {
	var a []int

	if !IsSliceType(a) {
		t.Errorf("IsSliceType(%T) = false", a)
	}
}

func TestIsMapType(t *testing.T) {
	var a map[int]int

	if !IsMapType(a) {
		t.Errorf("IsMapType(%T) = false", a)
	}
}

func TestIsPtrType(t *testing.T) {
	var a *string

	if !IsPtrType(a) {
		t.Errorf("IsPtrType(%T) = false", a)
	}
}
