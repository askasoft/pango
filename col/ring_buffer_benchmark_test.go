package col

import "testing"

func benchmarkRingBufferGet(b *testing.B, rb *RingBuffer, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			rb.Get(n)
		}
	}
}

func benchmarkRingBufferPush(b *testing.B, rb *RingBuffer, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			rb.Push(n)
		}
	}
}

func benchmarkRingBufferPoll(b *testing.B, rb *RingBuffer, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			rb.Poll()
		}
	}
}

func BenchmarkRingBufferGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	rb := NewRingBuffer()
	for n := 0; n < size; n++ {
		rb.Push(n)
	}
	b.StartTimer()
	benchmarkRingBufferGet(b, rb, size)
}

func BenchmarkRingBufferGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	rb := NewRingBuffer()
	for n := 0; n < size; n++ {
		rb.Push(n)
	}
	b.StartTimer()
	benchmarkRingBufferGet(b, rb, size)
}

func BenchmarkRingBufferGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	rb := NewRingBuffer()
	for n := 0; n < size; n++ {
		rb.Push(n)
	}
	b.StartTimer()
	benchmarkRingBufferGet(b, rb, size)
}

func BenchmarkRingBufferPush100(b *testing.B) {
	b.StopTimer()
	size := 100
	rb := NewRingBuffer()
	b.StartTimer()
	benchmarkRingBufferPush(b, rb, size)
}

func BenchmarkRingBufferPush1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	rb := NewRingBuffer()
	for n := 0; n < size; n++ {
		rb.Push(n)
	}
	b.StartTimer()
	benchmarkRingBufferPush(b, rb, size)
}

func BenchmarkRingBufferPush10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	rb := NewRingBuffer()
	for n := 0; n < size; n++ {
		rb.Push(n)
	}
	b.StartTimer()
	benchmarkRingBufferPush(b, rb, size)
}

func BenchmarkRingBufferPoll100(b *testing.B) {
	b.StopTimer()
	size := 100
	rb := NewRingBuffer()
	for n := 0; n < size; n++ {
		rb.Push(n)
	}
	b.StartTimer()
	benchmarkRingBufferPoll(b, rb, size)
}

func BenchmarkRingBufferPoll1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	rb := NewRingBuffer()
	for n := 0; n < size; n++ {
		rb.Push(n)
	}
	b.StartTimer()
	benchmarkRingBufferPoll(b, rb, size)
}

func BenchmarkRingBufferPoll10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	rb := NewRingBuffer()
	for n := 0; n < size; n++ {
		rb.Push(n)
	}
	b.StartTimer()
	benchmarkRingBufferPoll(b, rb, size)
}
