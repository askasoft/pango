package col

import (
	"testing"

	"github.com/pandafw/pango/cmp"
)

func benchmarkSortedListGet(b *testing.B, list *SortedList, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			list.Get(n)
		}
	}
}

func benchmarkSortedListAdd(b *testing.B, list *SortedList, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			list.Add(n)
		}
	}
}

func benchmarkSortedListDel(b *testing.B, list *SortedList, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			list.Delete(n)
		}
	}
}

func BenchmarkSortedListGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	list := NewSortedList(cmp.LessInt)
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkSortedListGet(b, list, size)
}

func BenchmarkSortedListGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	list := NewSortedList(cmp.LessInt)
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkSortedListGet(b, list, size)
}

func BenchmarkSortedListGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	list := NewSortedList(cmp.LessInt)
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkSortedListGet(b, list, size)
}

func BenchmarkSortedListGet100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	list := NewSortedList(cmp.LessInt)
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkSortedListGet(b, list, size)
}

func BenchmarkSortedListAdd100(b *testing.B) {
	b.StopTimer()
	size := 100
	list := NewSortedList(cmp.LessInt)
	b.StartTimer()
	benchmarkSortedListAdd(b, list, size)
}

func BenchmarkSortedListAdd1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	list := NewSortedList(cmp.LessInt)
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkSortedListAdd(b, list, size)
}

func BenchmarkSortedListAdd10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	list := NewSortedList(cmp.LessInt)
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkSortedListAdd(b, list, size)
}

func BenchmarkSortedListAdd100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	list := NewSortedList(cmp.LessInt)
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkSortedListAdd(b, list, size)
}

func BenchmarkSortedListRemove100(b *testing.B) {
	b.StopTimer()
	size := 100
	list := NewSortedList(cmp.LessInt)
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkSortedListDel(b, list, size)
}

func BenchmarkSortedListRemove1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	list := NewSortedList(cmp.LessInt)
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkSortedListDel(b, list, size)
}

func BenchmarkSortedListRemove10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	list := NewSortedList(cmp.LessInt)
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkSortedListDel(b, list, size)
}

func BenchmarkSortedListRemove100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	list := NewSortedList(cmp.LessInt)
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkSortedListDel(b, list, size)
}
