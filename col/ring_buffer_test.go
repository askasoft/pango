package col

import "testing"

func TestRingBufferSimple(t *testing.T) {
	rb := NewRingBuffer()

	for i := 0; i < minQueueCap; i++ {
		rb.Push(i)
	}

	for i := 0; i < minQueueCap; i++ {
		v, _ := rb.Peek()
		if v.(int) != i {
			t.Error("peek", i, "had value", v)
		}

		x, _ := rb.Poll()
		if x != i {
			t.Error("remove", i, "had value", x)
		}
	}
}

func TestRingBufferWrapping(t *testing.T) {
	rb := NewRingBuffer()

	for i := 0; i < minQueueCap; i++ {
		rb.Push(i)
	}
	for i := 0; i < 3; i++ {
		rb.Poll()
		rb.Push(minQueueCap + i)
	}

	for i := 0; i < minQueueCap; i++ {
		v, _ := rb.Peek()
		if v.(int) != i+3 {
			t.Error("peek", i, "had value", v)
		}
		rb.Poll()
	}
}

func TestRingBufferLength(t *testing.T) {
	rb := NewRingBuffer()

	if rb.Len() != 0 {
		t.Error("empty queue length not 0")
	}

	for i := 0; i < 1000; i++ {
		rb.Push(i)
		if rb.Len() != i+1 {
			t.Error("adding: queue with", i, "elements has length", rb.Len())
		}
	}
	for i := 0; i < 1000; i++ {
		rb.Poll()
		if rb.Len() != 1000-i-1 {
			t.Error("removing: queue with", 1000-i-i, "elements has length", rb.Len())
		}
	}
}

func TestRingBufferGet(t *testing.T) {
	rb := NewRingBuffer()

	for i := 0; i < 1000; i++ {
		rb.Push(i)
		for j := 0; j < rb.Len(); j++ {
			v := rb.Get(j).(int)
			if v != j {
				t.Errorf("[%d] index %d = %d, want %d", i, j, v, j)
			}
		}
	}
}

func TestRingBufferGetNegative(t *testing.T) {
	rb := NewRingBuffer()

	for i := 0; i < 1000; i++ {
		rb.Push(i)
		for j := 1; j <= rb.Len(); j++ {
			if rb.Get(-j).(int) != rb.Len()-j {
				t.Errorf("index %d doesn't contain %d", -j, rb.Len()-j)
			}
		}
	}
}

func TestRingBufferGetOutOfRangePanics(t *testing.T) {
	rb := NewRingBuffer()

	rb.Push(1, 2, 3)

	assertPanics(t, "should panic when negative index", func() {
		rb.Get(-4)
	})

	assertPanics(t, "should panic when index greater than length", func() {
		rb.Get(4)
	})
}

func TestRingBufferPeekEmpty(t *testing.T) {
	rb := NewRingBuffer()

	if _, ok := rb.Peek(); ok {
		t.Error("should return false when peeking empty queue")
	}

	rb.Push(1)
	rb.Poll()

	if _, ok := rb.Peek(); ok {
		t.Error("should return false when peeking empty queue")
	}
}

func TestRingBufferPollEmpty(t *testing.T) {
	rb := NewRingBuffer()

	if _, ok := rb.Poll(); ok {
		t.Error("should return false when removing empty queue")
	}

	rb.Push(1)
	rb.Poll()

	if _, ok := rb.Poll(); ok {
		t.Error("should return false when removing empty queue")
	}
}
