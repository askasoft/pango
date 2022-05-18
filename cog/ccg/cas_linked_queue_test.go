//go:build go1.18
// +build go1.18

package ccg

import "testing"

func TestCQueueSimple(t *testing.T) {
	rb := NewCQueue[int]()

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

func TestCQueuePeekEmpty(t *testing.T) {
	rb := NewCQueue[int]()

	if _, ok := rb.Peek(); ok {
		t.Error("should return false when peeking empty queue")
	}

	rb.Push(1)
	rb.Poll()

	if _, ok := rb.Peek(); ok {
		t.Error("should return false when peeking empty queue")
	}
}

func TestCQueuePollEmpty(t *testing.T) {
	rb := NewCQueue[int]()

	if _, ok := rb.Poll(); ok {
		t.Error("should return false when removing empty queue")
	}

	rb.Push(1)
	rb.Poll()

	if _, ok := rb.Poll(); ok {
		t.Error("should return false when removing empty queue")
	}
}
