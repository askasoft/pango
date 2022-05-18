//go:build go1.18
// +build go1.18

package cog

import (
	"testing"
)

func benchmarkTreeSetContains(b *testing.B, set *TreeSet[int], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			set.Contains(n)
		}
	}
}

func benchmarkTreeSetAdd(b *testing.B, set *TreeSet[int], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			set.Add(n)
		}
	}
}

func benchmarkTreeSetDelete(b *testing.B, set *TreeSet[int], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			set.Delete(n)
		}
	}
}

func BenchmarkTreeSetContains100(b *testing.B) {
	b.StopTimer()
	size := 100
	set := NewTreeSet(CompareInt)
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkTreeSetContains(b, set, size)
}

func BenchmarkTreeSetContains1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	set := NewTreeSet(CompareInt)
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkTreeSetContains(b, set, size)
}

func BenchmarkTreeSetContains10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	set := NewTreeSet(CompareInt)
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkTreeSetContains(b, set, size)
}

func BenchmarkTreeSetAdd100(b *testing.B) {
	b.StopTimer()
	size := 100
	set := NewTreeSet(CompareInt)
	b.StartTimer()
	benchmarkTreeSetAdd(b, set, size)
}

func BenchmarkTreeSetAdd1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	set := NewTreeSet(CompareInt)
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkTreeSetAdd(b, set, size)
}

func BenchmarkTreeSetAdd10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	set := NewTreeSet(CompareInt)
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkTreeSetAdd(b, set, size)
}

func BenchmarkTreeSetDelete100(b *testing.B) {
	b.StopTimer()
	size := 100
	set := NewTreeSet(CompareInt)
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkTreeSetDelete(b, set, size)
}

func BenchmarkTreeSetDelete1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	set := NewTreeSet(CompareInt)
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkTreeSetDelete(b, set, size)
}

func BenchmarkTreeSetDelete10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	set := NewTreeSet(CompareInt)
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkTreeSetDelete(b, set, size)
}
