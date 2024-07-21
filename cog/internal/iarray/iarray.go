//go:build go1.18
// +build go1.18

package iarray

func Contains[T any](a []T, v T) bool {
	return Index(a, v) >= 0
}

func Index[T any](a []T, v T) int {
	for i, e := range a {
		if any(e) == any(v) {
			return i
		}
	}
	return -1
}
