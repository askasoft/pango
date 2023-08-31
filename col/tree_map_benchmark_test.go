package col

import (
	"testing"
)

func benchmarkTreeMapGet(b *testing.B, tree *TreeMap, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Get(n)
		}
	}
}

func benchmarkTreeMapSet(b *testing.B, tree *TreeMap, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Set(n, struct{}{})
		}
	}
}

func benchmarkTreeMapDel(b *testing.B, tree *TreeMap, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Remove(n)
		}
	}
}

func BenchmarkTreeMapGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	tree := NewTreeMap(CompareInt)
	for n := 0; n < size; n++ {
		tree.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkTreeMapGet(b, tree, size)
}

func BenchmarkTreeMapGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	tree := NewTreeMap(CompareInt)
	for n := 0; n < size; n++ {
		tree.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkTreeMapGet(b, tree, size)
}

func BenchmarkTreeMapGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	tree := NewTreeMap(CompareInt)
	for n := 0; n < size; n++ {
		tree.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkTreeMapGet(b, tree, size)
}

func BenchmarkTreeMapPut100(b *testing.B) {
	b.StopTimer()
	size := 100
	tree := NewTreeMap(CompareInt)
	b.StartTimer()
	benchmarkTreeMapSet(b, tree, size)
}

func BenchmarkTreeMapPut1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	tree := NewTreeMap(CompareInt)
	for n := 0; n < size; n++ {
		tree.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkTreeMapSet(b, tree, size)
}

func BenchmarkTreeMapPut10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	tree := NewTreeMap(CompareInt)
	for n := 0; n < size; n++ {
		tree.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkTreeMapSet(b, tree, size)
}

func BenchmarkTreeMapDelete100(b *testing.B) {
	b.StopTimer()
	size := 100
	tree := NewTreeMap(CompareInt)
	for n := 0; n < size; n++ {
		tree.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkTreeMapDel(b, tree, size)
}

func BenchmarkTreeMapDelete1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	tree := NewTreeMap(CompareInt)
	for n := 0; n < size; n++ {
		tree.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkTreeMapDel(b, tree, size)
}

func BenchmarkTreeMapDelete10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	tree := NewTreeMap(CompareInt)
	for n := 0; n < size; n++ {
		tree.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkTreeMapDel(b, tree, size)
}
