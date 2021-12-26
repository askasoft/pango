package col

import "fmt"

// minQueueCap is smallest capacity that queue may have.
const minQueueCap = 32

// RingBuffer A fast Golang queue using a ring-buffer, based on the version suggested by Dariusz Górecki.
// Using this instead of other, simpler, queue implementations (slice+append or linked list) provides substantial memory and time benefits, and fewer GC pauses.
// The queue implemented here is as fast as it is in part because it is not thread-safe.
type RingBuffer struct {
	data                  []T
	head, tail, len, grow int
}

// NewRingBuffer constructs and returns a new RingBuffer.
// Example: NewRingBuffer(initialSize, growSize)
func NewRingBuffer(ns ...int) *RingBuffer {
	initSize := minQueueCap
	growSize := minQueueCap

	if len(ns) > 0 {
		if ns[0] > initSize {
			initSize = roundup(ns[0], minQueueCap)
		}
	}
	if len(ns) > 1 {
		if ns[1] > growSize {
			growSize = roundup(ns[1], minQueueCap)
		}
	}

	return &RingBuffer{
		data: make([]T, initSize),
		grow: growSize,
	}
}

// Len returns the number of elements currently stored in the queue.
func (rb *RingBuffer) Len() int {
	return rb.len
}

// IsEmpty returns true if the container length == 0
func (rb *RingBuffer) IsEmpty() bool {
	return rb.len == 0
}

// GetGrowSize returns the grow size of queue buffer
func (rb *RingBuffer) GetGrowSize() int {
	return rb.grow
}

// SetGrowSize sets the buffer's grow size
func (rb *RingBuffer) SetGrowSize(n int) {
	rb.grow = roundup(n, minQueueCap)
}

// Push adds items of vs to the tail of queue
func (rb *RingBuffer) Push(vs ...T) {
	if rb.len+len(vs) > len(rb.data) {
		rb.resize(rb.len + len(vs))
	}

	for _, v := range vs {
		rb.data[rb.tail] = v
		rb.len++

		rb.tail++
		if rb.tail >= len(rb.data) {
			rb.tail = 0
		}
	}
}

// Get returns the value at the specified index in this queue. If the index is
// invalid, the call will panic. This method accepts both positive and
// negative index values. Index 0 refers to the first element, and
// index -1 refers to the last.
func (rb *RingBuffer) Get(index int) T {
	index = rb.checkItemIndex(index)

	index += rb.head
	if index >= len(rb.data) {
		index -= len(rb.data)
	}

	return rb.data[index]
}

// Peek Retrieves, but does not remove, the head of this queue, or returns (nil, false) if this queue is empty.
func (rb *RingBuffer) Peek() (T, bool) {
	if rb.len <= 0 {
		return nil, false
	}
	return rb.data[rb.head], true
}

// MustPeek Retrieves, but does not remove, the head of this queue, panic if this queue is empty.
func (rb *RingBuffer) MustPeek() T {
	if v, ok := rb.Peek(); ok {
		return v
	}

	panic("RingBuffer: MustPeek() called on empty queue")
}

// Poll Retrieves and removes the head of this queue, or returns (nil, false) if this queue is empty.
func (rb *RingBuffer) Poll() (T, bool) {
	if rb.len <= 0 {
		return nil, false
	}

	v := rb.data[rb.head]
	rb.data[rb.head] = nil

	rb.head++
	if rb.head >= len(rb.data) {
		rb.head = 0
	}
	rb.len--

	// Resize down if data is less than 1/4 full.
	if len(rb.data) > minQueueCap && (rb.len<<2) <= len(rb.data) {
		rb.resize(len(rb.data))
	}
	return v, true
}

// MustPoll Retrieves and removes the head of this queue, panic if this queue is empty.
func (rb *RingBuffer) MustPoll() T {
	if v, ok := rb.Poll(); ok {
		return v
	}

	panic("RingBuffer: MustPoll() called on empty queue")
}

//-----------------------------------------------------------
// resizes the queue to fit exactly twice its current contents
// this can result in shrinking if the queue is less than 1/4 full
func (rb *RingBuffer) resize(n int) {
	data := make([]T, roundup(n, rb.grow))

	if rb.tail > rb.head {
		copy(data, rb.data[rb.head:rb.tail])
	} else {
		n := copy(data, rb.data[rb.head:])
		copy(data[n:], rb.data[:rb.tail])
	}

	rb.head = 0
	rb.tail = rb.len
	rb.data = data
}

func (rb *RingBuffer) checkItemIndex(index int) int {
	len := rb.Len()
	if index >= len || index < -len {
		panic(fmt.Sprintf("RingBuffer out of bounds: index=%d, len=%d", index, len))
	}

	if index < 0 {
		index += len
	}
	return index
}
