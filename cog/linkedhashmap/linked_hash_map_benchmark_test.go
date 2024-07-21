//go:build go1.18
// +build go1.18

package linkedhashmap

import "testing"

func benchmarkLinkedHashMapGet(b *testing.B, m *LinkedHashMap[int, any], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			m.Get(n)
		}
	}
}

func benchmarkLinkedHashMapSet(b *testing.B, m *LinkedHashMap[int, any], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			m.Set(n, struct{}{})
		}
	}
}

func benchmarkLinkedHashMapDel(b *testing.B, m *LinkedHashMap[int, any], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			m.Remove(n)
		}
	}
}

func BenchmarkLinkedHashMapGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	m := NewLinkedHashMap[int, any]()
	for n := 0; n < size; n++ {
		m.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkLinkedHashMapGet(b, m, size)
}

func BenchmarkLinkedHashMapGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	m := NewLinkedHashMap[int, any]()
	for n := 0; n < size; n++ {
		m.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkLinkedHashMapGet(b, m, size)
}

func BenchmarkLinkedHashMapGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	m := NewLinkedHashMap[int, any]()
	for n := 0; n < size; n++ {
		m.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkLinkedHashMapGet(b, m, size)
}

func BenchmarkLinkedHashMapPut100(b *testing.B) {
	b.StopTimer()
	size := 100
	m := NewLinkedHashMap[int, any]()
	b.StartTimer()
	benchmarkLinkedHashMapSet(b, m, size)
}

func BenchmarkLinkedHashMapPut1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	m := NewLinkedHashMap[int, any]()
	for n := 0; n < size; n++ {
		m.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkLinkedHashMapSet(b, m, size)
}

func BenchmarkLinkedHashMapPut10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	m := NewLinkedHashMap[int, any]()
	for n := 0; n < size; n++ {
		m.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkLinkedHashMapSet(b, m, size)
}

func BenchmarkLinkedHashMapRemove100(b *testing.B) {
	b.StopTimer()
	size := 100
	m := NewLinkedHashMap[int, any]()
	for n := 0; n < size; n++ {
		m.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkLinkedHashMapDel(b, m, size)
}

func BenchmarkLinkedHashMapRemove1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	m := NewLinkedHashMap[int, any]()
	for n := 0; n < size; n++ {
		m.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkLinkedHashMapDel(b, m, size)
}

func BenchmarkLinkedHashMapRemove10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	m := NewLinkedHashMap[int, any]()
	for n := 0; n < size; n++ {
		m.Set(n, struct{}{})
	}
	b.StartTimer()
	benchmarkLinkedHashMapDel(b, m, size)
}
