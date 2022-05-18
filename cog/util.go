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

func setMapPairs[K any, V any](m Map[K, V], ps ...P[K, V]) {
	for _, p := range ps {
		m.Set(p.Key, p.Value)
	}
}

func setMapAll[K any, V any](dst, src Map[K, V]) {
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

	src.Each(func(k K, v V) {
		dst.Set(k, v)
	})
}