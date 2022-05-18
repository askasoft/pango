package cog

import "testing"

func benchmarkArrayListGet(b *testing.B, list *ArrayList[int], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			list.Get(n)
		}
	}
}

func benchmarkArrayListAdd(b *testing.B, list *ArrayList[int], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			list.Add(n)
		}
	}
}

func benchmarkArrayListDelete(b *testing.B, list *ArrayList[int], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			list.Delete(n)
		}
	}
}

func BenchmarkArrayListGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	list := NewArrayList[int]()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkArrayListGet(b, list, size)
}

func BenchmarkArrayListGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	list := NewArrayList[int]()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkArrayListGet(b, list, size)
}

func BenchmarkArrayListGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	list := NewArrayList[int]()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkArrayListGet(b, list, size)
}

func BenchmarkArrayListAdd100(b *testing.B) {
	b.StopTimer()
	size := 100
	list := NewArrayList[int]()
	b.StartTimer()
	benchmarkArrayListAdd(b, list, size)
}

func BenchmarkArrayListAdd1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	list := NewArrayList[int]()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkArrayListAdd(b, list, size)
}

func BenchmarkArrayListAdd10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	list := NewArrayList[int]()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkArrayListAdd(b, list, size)
}

func BenchmarkArrayListDelete100(b *testing.B) {
	b.StopTimer()
	size := 100
	list := NewArrayList[int]()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkArrayListDelete(b, list, size)
}

func BenchmarkArrayListDelete1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	list := NewArrayList[int]()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkArrayListDelete(b, list, size)
}

func BenchmarkArrayListDelete10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	list := NewArrayList[int]()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkArrayListDelete(b, list, size)
}
