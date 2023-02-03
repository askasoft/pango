//go:build go1.18
// +build go1.18

package ars

// ContainsOf reports whether the c is contained in the slice a.
func ContainsOf[T any](a []T, c T) bool {
	return IndexOf(a, c) >= 0
}

// EqualOf reports whether a and b
// are the same length and contain the same element.
// A nil argument is equivalent to an empty slice.
func EqualOf[T any](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if any(v) != any(b[i]) {
			return false
		}
	}
	return true
}

// GetOf get element at the specified index i.
func GetOf[T any](a []T, i int) (v T, ok bool) {
	if i >= 0 && i < len(a) {
		v, ok = a[i], true
	}
	return
}

// IndexOf returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexOf[T any](a []T, c T) int {
	for i, e := range a {
		if any(e) == any(c) {
			return i
		}
	}
	return -1
}
