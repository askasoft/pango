//go:build go1.18
// +build go1.18

package cog

// minArrayCap is smallest capacity that array may have.
const minArrayCap = 16

// doubleup grow up capcity by c *= 2 until c >= n
func doubleup(c, n int) int {
	if c < minArrayCap {
		c = minArrayCap
	}
	for c < n {
		c <<= 1
	}
	return c
}

// roundup round up size by the block size r
func roundup(n, r int) int {
	r--
	if (n & r) == 0 {
		return n
	}

	return (n + r) & (^r)
}

func contains[T any](a []T, v T) bool {
	return index(a, v) >= 0
}

func index[T any](a []T, v T) int {
	for i, e := range a {
		if any(e) == any(v) {
			return i
		}
	}
	return -1
}

func setMapPairs[K any, V any](m Map[K, V], ps ...P[K, V]) {
	for _, p := range ps {
		m.Set(p.Key, p.Value)
	}
}

// CopyMap copy map 'src' to map 'des'
func CopyMap[K any, V any](dst, src Map[K, V]) {
	if dst == src {
		return
	}

	if im, ok := src.(Iterable2[K, V]); ok {
		it := im.Iterator()
		for it.Next() {
			dst.Set(it.Key(), it.Value())
		}
		return
	}

	src.Each(func(k K, v V) bool {
		dst.Set(k, v)
		return true
	})
}
