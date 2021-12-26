package col

// roundup round up size by the block size r
func roundup(n, r int) int {
	r--
	if (n & r) == 0 {
		return n
	}

	return (n + r) & (^r)
}

func setMapPairs(m Map, ps ...P) {
	for _, p := range ps {
		m.Set(p.Key, p.Value)
	}
}

func setMapAll(dst, src Map) {
	if dst == src {
		return
	}

	if im, ok := src.(Iterable2); ok {
		it := im.Iterator()
		for it.Next() {
			dst.Set(it.Key(), it.Value())
		}
		return
	}

	src.Each(func(k interface{}, v interface{}) {
		dst.Set(k, v)
	})
}
