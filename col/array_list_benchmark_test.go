package col

import "testing"

func benchmarkArrayListGet(b *testing.B, list *ArrayList, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			list.Get(n)
		}
	}
}

func benchmarkArrayListAdd(b *testing.B, list *ArrayList, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			list.Add(n)
		}
	}
}

func benchmarkArrayListRemove(b *testing.B, list *ArrayList, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			list.Remove(n)
		}
	}
}

func BenchmarkArrayListGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	list := NewArrayList()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkArrayListGet(b, list, size)
}

func BenchmarkArrayListGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	list := NewArrayList()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkArrayListGet(b, list, size)
}

func BenchmarkArrayListGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	list := NewArrayList()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkArrayListGet(b, list, size)
}

func BenchmarkArrayListGet100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	list := NewArrayList()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkArrayListGet(b, list, size)
}

func BenchmarkArrayListAdd100(b *testing.B) {
	b.StopTimer()
	size := 100
	list := NewArrayList()
	b.StartTimer()
	benchmarkArrayListAdd(b, list, size)
}

func BenchmarkArrayListAdd1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	list := NewArrayList()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkArrayListAdd(b, list, size)
}

func BenchmarkArrayListAdd10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	list := NewArrayList()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkArrayListAdd(b, list, size)
}

func BenchmarkArrayListAdd100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	list := NewArrayList()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkArrayListAdd(b, list, size)
}

func BenchmarkArrayListRemove100(b *testing.B) {
	b.StopTimer()
	size := 100
	list := NewArrayList()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkArrayListRemove(b, list, size)
}

func BenchmarkArrayListRemove1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	list := NewArrayList()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkArrayListRemove(b, list, size)
}

func BenchmarkArrayListRemove10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	list := NewArrayList()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkArrayListRemove(b, list, size)
}

func BenchmarkArrayListRemove100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	list := NewArrayList()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkArrayListRemove(b, list, size)
}
