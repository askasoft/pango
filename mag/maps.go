package mag

import (
	"maps"

	"github.com/askasoft/pango/asg"
)

// Keys return the slice of map's key
func Keys[K comparable, V any](m map[K]V) []K {
	ks := make([]K, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

// Values return the slice of map's value
func Values[K comparable, V any](m map[K]V) []V {
	vs := make([]V, 0, len(m))
	for _, v := range m {
		vs = append(vs, v)
	}
	return vs
}

// Equal check equal for two map a and b.
func Equal[K comparable, V comparable](a map[K]V, b map[K]V) bool {
	return maps.Equal(a, b)
}

// EqualFunc is like Equal, but compares values using eq.
// Keys are still compared with ==.
func EqualFunc[K comparable, V1, V2 any](m1 map[K]V1, m2 map[K]V2, eq func(V1, V2) bool) bool {
	return maps.EqualFunc(m1, m2, eq)
}

// Clone returns a copy of the map with additional +n cap.
// The elements are copied using assignment, so this is a shallow clone.
func Clone[K comparable, V any](m map[K]V, n ...int) map[K]V {
	// Preserve nil in case it matters.
	if m == nil {
		return nil
	}

	x := asg.First(n)
	if x == 0 {
		return maps.Clone(m)
	}

	c := make(map[K]V, len(m)+x)
	Copy(c, m)
	return c
}

// Copy copy src map to dst map.
func Copy[K comparable, V any](dst, src map[K]V) {
	maps.Copy(dst, src)
}

// Merge merge all map `ms`, return the merged map.
func Merge[K comparable, V any](ms ...map[K]V) map[K]V {
	sz := 0
	for _, m := range ms {
		sz += len(m)
	}

	mm := make(map[K]V, sz)
	for _, m := range ms {
		for k, v := range m {
			mm[k] = v
		}
	}
	return mm
}

// Reverse reverse key/value from map m, return the reversed map
func Reverse[K comparable, V comparable](m map[K]V) map[V]K {
	rm := make(map[V]K, len(m))
	for k, v := range m {
		rm[v] = k
	}
	return rm
}
