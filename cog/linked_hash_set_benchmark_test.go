package cog

import "testing"

func benchmarkLinkedHashSetGet(b *testing.B, lset *LinkedHashSet[int], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			lset.Get(n)
		}
	}
}

func benchmarkLinkedHashSetAdd(b *testing.B, lset *LinkedHashSet[int], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			lset.Add(n)
		}
	}
}

func benchmarkLinkedHashSetDel(b *testing.B, lset *LinkedHashSet[int], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			lset.Delete(n)
		}
	}
}

func BenchmarkLinkedHashSetGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	lset := NewLinkedHashSet[int]()
	for n := 0; n < size; n++ {
		lset.Add(n)
	}
	b.StartTimer()
	benchmarkLinkedHashSetGet(b, lset, size)
}

func BenchmarkLinkedHashSetGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	lset := NewLinkedHashSet[int]()
	for n := 0; n < size; n++ {
		lset.Add(n)
	}
	b.StartTimer()
	benchmarkLinkedHashSetGet(b, lset, size)
}

func BenchmarkLinkedHashSetGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	lset := NewLinkedHashSet[int]()
	for n := 0; n < size; n++ {
		lset.Add(n)
	}
	b.StartTimer()
	benchmarkLinkedHashSetGet(b, lset, size)
}

func BenchmarkLinkedHashSetAdd100(b *testing.B) {
	b.StopTimer()
	size := 100
	lset := NewLinkedHashSet[int]()
	b.StartTimer()
	benchmarkLinkedHashSetAdd(b, lset, size)
}

func BenchmarkLinkedHashSetAdd1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	lset := NewLinkedHashSet[int]()
	for n := 0; n < size; n++ {
		lset.Add(n)
	}
	b.StartTimer()
	benchmarkLinkedHashSetAdd(b, lset, size)
}

func BenchmarkLinkedHashSetAdd10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	lset := NewLinkedHashSet[int]()
	for n := 0; n < size; n++ {
		lset.Add(n)
	}
	b.StartTimer()
	benchmarkLinkedHashSetAdd(b, lset, size)
}

func BenchmarkLinkedHashSetRemove100(b *testing.B) {
	b.StopTimer()
	size := 100
	lset := NewLinkedHashSet[int]()
	for n := 0; n < size; n++ {
		lset.Add(n)
	}
	b.StartTimer()
	benchmarkLinkedHashSetDel(b, lset, size)
}

func BenchmarkLinkedHashSetRemove1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	lset := NewLinkedHashSet[int]()
	for n := 0; n < size; n++ {
		lset.Add(n)
	}
	b.StartTimer()
	benchmarkLinkedHashSetDel(b, lset, size)
}

func BenchmarkLinkedHashSetRemove10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	lset := NewLinkedHashSet[int]()
	for n := 0; n < size; n++ {
		lset.Add(n)
	}
	b.StartTimer()
	benchmarkLinkedHashSetDel(b, lset, size)
}
