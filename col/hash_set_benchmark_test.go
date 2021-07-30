package col

import "testing"

func benchmarkHashSetContains(b *testing.B, set *HashSet, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			set.Contains(n)
		}
	}
}

func benchmarkHashSetAdd(b *testing.B, set *HashSet, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			set.Add(n)
		}
	}
}

func benchmarkHashSetDel(b *testing.B, set *HashSet, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			set.Delete(n)
		}
	}
}

func BenchmarkHashSetContains100(b *testing.B) {
	b.StopTimer()
	size := 100
	set := NewHashSet()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkHashSetContains(b, set, size)
}

func BenchmarkHashSetContains1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	set := NewHashSet()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkHashSetContains(b, set, size)
}

func BenchmarkHashSetContains10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	set := NewHashSet()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkHashSetContains(b, set, size)
}

func BenchmarkHashSetContains100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	set := NewHashSet()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkHashSetContains(b, set, size)
}

func BenchmarkHashSetAdd100(b *testing.B) {
	b.StopTimer()
	size := 100
	set := NewHashSet()
	b.StartTimer()
	benchmarkHashSetAdd(b, set, size)
}

func BenchmarkHashSetAdd1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	set := NewHashSet()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkHashSetAdd(b, set, size)
}

func BenchmarkHashSetAdd10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	set := NewHashSet()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkHashSetAdd(b, set, size)
}

func BenchmarkHashSetAdd100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	set := NewHashSet()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkHashSetAdd(b, set, size)
}

func BenchmarkHashSetRemove100(b *testing.B) {
	b.StopTimer()
	size := 100
	set := NewHashSet()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkHashSetDel(b, set, size)
}

func BenchmarkHashSetRemove1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	set := NewHashSet()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkHashSetDel(b, set, size)
}

func BenchmarkHashSetRemove10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	set := NewHashSet()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkHashSetDel(b, set, size)
}

func BenchmarkHashSetRemove100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	set := NewHashSet()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkHashSetDel(b, set, size)
}
