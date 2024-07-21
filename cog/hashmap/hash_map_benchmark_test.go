//go:build go1.18
// +build go1.18

package hashmap

import "testing"

func benchmarkHashMapGet(b *testing.B, m *HashMap[int, any], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			m.Get(n)
		}
	}
}

func benchmarkHashMapSet(b *testing.B, m *HashMap[int, any], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			m.Set(n, struct{}{})
		}
	}
}

func benchmarkHashMapRemove(b *testing.B, m *HashMap[int, any], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			m.Remove(n)
		}
	}
}

func BenchmarkHashMapGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	m := NewHashMap[int, any]()
	for n := 0; n < size; n++ {
		m.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkHashMapGet(b, m, size)
}

func BenchmarkHashMapGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	m := NewHashMap[int, any]()
	for n := 0; n < size; n++ {
		m.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkHashMapGet(b, m, size)
}

func BenchmarkHashMapGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	m := NewHashMap[int, any]()
	for n := 0; n < size; n++ {
		m.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkHashMapGet(b, m, size)
}

func BenchmarkHashMapSet100(b *testing.B) {
	b.StopTimer()
	size := 100
	m := NewHashMap[int, any]()
	b.StartTimer()
	benchmarkHashMapSet(b, m, size)
}

func BenchmarkHashMapSet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	m := NewHashMap[int, any]()
	for n := 0; n < size; n++ {
		m.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkHashMapSet(b, m, size)
}

func BenchmarkHashMapSet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	m := NewHashMap[int, any]()
	for n := 0; n < size; n++ {
		m.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkHashMapSet(b, m, size)
}

func BenchmarkHashMapRemove100(b *testing.B) {
	b.StopTimer()
	size := 100
	m := NewHashMap[int, any]()
	for n := 0; n < size; n++ {
		m.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkHashMapRemove(b, m, size)
}

func BenchmarkHashMapRemove1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	m := NewHashMap[int, any]()
	for n := 0; n < size; n++ {
		m.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkHashMapRemove(b, m, size)
}

func BenchmarkHashMapRemove10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	m := NewHashMap[int, any]()
	for n := 0; n < size; n++ {
		m.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkHashMapRemove(b, m, size)
}
