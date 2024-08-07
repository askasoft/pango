//go:build go1.18
// +build go1.18

package treemap

import (
	"testing"

	"github.com/askasoft/pango/cmp"
)

func benchmarkTreeMapGet(b *testing.B, tree *TreeMap[int, any], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Get(n)
		}
	}
}

func benchmarkTreeMapSet(b *testing.B, tree *TreeMap[int, any], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Set(n, struct{}{})
		}
	}
}

func benchmarkTreeMapDel(b *testing.B, tree *TreeMap[int, any], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Remove(n)
		}
	}
}

func BenchmarkTreeMapGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	tree := NewTreeMap[int, any](cmp.CompareInt)
	for n := 0; n < size; n++ {
		tree.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkTreeMapGet(b, tree, size)
}

func BenchmarkTreeMapGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	tree := NewTreeMap[int, any](cmp.CompareInt)
	for n := 0; n < size; n++ {
		tree.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkTreeMapGet(b, tree, size)
}

func BenchmarkTreeMapGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	tree := NewTreeMap[int, any](cmp.CompareInt)
	for n := 0; n < size; n++ {
		tree.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkTreeMapGet(b, tree, size)
}

func BenchmarkTreeMapPut100(b *testing.B) {
	b.StopTimer()
	size := 100
	tree := NewTreeMap[int, any](cmp.CompareInt)
	b.StartTimer()
	benchmarkTreeMapSet(b, tree, size)
}

func BenchmarkTreeMapPut1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	tree := NewTreeMap[int, any](cmp.CompareInt)
	for n := 0; n < size; n++ {
		tree.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkTreeMapSet(b, tree, size)
}

func BenchmarkTreeMapPut10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	tree := NewTreeMap[int, any](cmp.CompareInt)
	for n := 0; n < size; n++ {
		tree.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkTreeMapSet(b, tree, size)
}

func BenchmarkTreeMapDelete100(b *testing.B) {
	b.StopTimer()
	size := 100
	tree := NewTreeMap[int, any](cmp.CompareInt)
	for n := 0; n < size; n++ {
		tree.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkTreeMapDel(b, tree, size)
}

func BenchmarkTreeMapDelete1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	tree := NewTreeMap[int, any](cmp.CompareInt)
	for n := 0; n < size; n++ {
		tree.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkTreeMapDel(b, tree, size)
}

func BenchmarkTreeMapDelete10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	tree := NewTreeMap[int, any](cmp.CompareInt)
	for n := 0; n < size; n++ {
		tree.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkTreeMapDel(b, tree, size)
}
