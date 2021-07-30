package col

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
