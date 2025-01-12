package casqueue

import "testing"

func TestCasQueueSimple(t *testing.T) {
	rb := NewCasQueue[int]()

	for i := 0; i < 100; i++ {
		rb.Push(i)
	}

	for i := 0; i < 100; i++ {
		v, _ := rb.Peek()
		if v != i {
			t.Error("peek", i, "had value", v)
		}

		x, _ := rb.Poll()
		if x != i {
			t.Error("remove", i, "had value", x)
		}
	}
}

func TestCasQueuePeekEmpty(t *testing.T) {
	rb := NewCasQueue[int]()

	if _, ok := rb.Peek(); ok {
		t.Error("should return false when peeking empty queue")
	}

	rb.Push(1)
	rb.Poll()

	if _, ok := rb.Peek(); ok {
		t.Error("should return false when peeking empty queue")
	}
}

func TestCasQueuePollEmpty(t *testing.T) {
	rb := NewCasQueue[int]()

	if _, ok := rb.Poll(); ok {
		t.Error("should return false when removing empty queue")
	}

	rb.Push(1)
	rb.Poll()

	if _, ok := rb.Poll(); ok {
		t.Error("should return false when removing empty queue")
	}
}
