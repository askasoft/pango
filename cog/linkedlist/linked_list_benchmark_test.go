package linkedlist

import "testing"

func benchmarkLinkedListGet(b *testing.B, list *LinkedList[int], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			list.Get(n)
		}
	}
}

func benchmarkLinkedListAdd(b *testing.B, list *LinkedList[int], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			list.Add(n)
		}
	}
}

func benchmarkLinkedListDelete(b *testing.B, list *LinkedList[int], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			list.Remove(n)
		}
	}
}

func BenchmarkLinkedListGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	list := NewLinkedList[int]()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkLinkedListGet(b, list, size)
}

func BenchmarkLinkedListGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	list := NewLinkedList[int]()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkLinkedListGet(b, list, size)
}

func BenchmarkLinkedListGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	list := NewLinkedList[int]()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkLinkedListGet(b, list, size)
}

func BenchmarkLinkedListAdd100(b *testing.B) {
	b.StopTimer()
	size := 100
	list := NewLinkedList[int]()
	b.StartTimer()
	benchmarkLinkedListAdd(b, list, size)
}

func BenchmarkLinkedListAdd1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	list := NewLinkedList[int]()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkLinkedListAdd(b, list, size)
}

func BenchmarkLinkedListAdd10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	list := NewLinkedList[int]()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkLinkedListAdd(b, list, size)
}

func BenchmarkLinkedListDelete100(b *testing.B) {
	b.StopTimer()
	size := 100
	list := NewLinkedList[int]()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkLinkedListDelete(b, list, size)
}

func BenchmarkLinkedListDelete1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	list := NewLinkedList[int]()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkLinkedListDelete(b, list, size)
}

func BenchmarkLinkedListDelete10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	list := NewLinkedList[int]()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkLinkedListDelete(b, list, size)
}
