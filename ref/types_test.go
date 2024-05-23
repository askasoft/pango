package ref

import (
	"testing"
)

func TestIsArrayType(t *testing.T) {
	var a [1]int

	if !IsArrayType(a) {
		t.Errorf("IsArrayType(%T) = false", a)
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
