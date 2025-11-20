package mag

import "github.com/askasoft/pango/asg"

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
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if w, ok := b[k]; !ok || w != v {
			return false
		}
	}

	return true
}

// Clone returns a copy of the map with additional +n cap.
// The elements are copied using assignment, so this is a shallow clone.
func Clone[K comparable, V any](src map[K]V, n ...int) map[K]V {
	dst := make(map[K]V, len(src)+asg.First(n))
	Copy(dst, src)
	return dst
}

// Copy copy src map to dst map.
func Copy[K comparable, V any](dst, src map[K]V) {
	for k, v := range src {
		dst[k] = v
	}
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
