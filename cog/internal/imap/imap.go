//go:build go1.18
// +build go1.18

package imap

import "github.com/askasoft/pango/cog"

func SetMapPairs[K any, V any](m cog.Map[K, V], ps ...cog.P[K, V]) {
	for _, p := range ps {
		m.Set(p.Key, p.Value)
	}
}

// CopyMap copy map 'src' to map 'des'
func CopyMap[K any, V any](dst, src cog.Map[K, V]) {
	if dst == src {
		return
	}

	if im, ok := src.(cog.Iterable2[K, V]); ok {
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
