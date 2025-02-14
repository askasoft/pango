package twolockqueue

import "testing"

func TestTwoLockQueueSimple(t *testing.T) {
	rb := NewTwoLockQueue[int]()

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

func TestTwoLockQueuePeekEmpty(t *testing.T) {
	rb := NewTwoLockQueue[int]()

	if _, ok := rb.Peek(); ok {
		t.Error("should return false when peeking empty queue")
	}

	rb.Push(1)
	rb.Poll()

	if _, ok := rb.Peek(); ok {
		t.Error("should return false when peeking empty queue")
	}
}

func TestTwoLockQueuePollEmpty(t *testing.T) {
	rb := NewTwoLockQueue[int]()

	if _, ok := rb.Poll(); ok {
		t.Error("should return false when removing empty queue")
	}

	rb.Push(1)
	rb.Poll()

	if _, ok := rb.Poll(); ok {
		t.Error("should return false when removing empty queue")
	}
}
