package twolockqueue

import (
	"sync"
)

// TwoLockQueue is a concurrent unbounded queue which uses two-Lock concurrent queue algorithm.
type TwoLockQueue[T any] struct {
	head  *lnode[T]
	tail  *lnode[T]
	hlock sync.Mutex
	tlock sync.Mutex
}

type lnode[T any] struct {
	data T
	next *lnode[T]
}

// NewTwoLockQueue returns an empty TwoLockQueue.
func NewTwoLockQueue[T any]() *TwoLockQueue[T] {
	n := &lnode[T]{}
	return &TwoLockQueue[T]{head: n, tail: n}
}

// IsEmpty returns true if the container length == 0
func (q *TwoLockQueue[T]) IsEmpty() bool {
	_, ok := q.Peek()
	return !ok
}

// Clear clears the container
func (q *TwoLockQueue[T]) Clear() {
	q.tlock.Lock()
	defer q.tlock.Unlock()

	q.hlock.Lock()
	defer q.hlock.Unlock()

	var v T
	q.head.data = v
	q.head.next = nil
	q.tail = q.head
}

// Push adds items of vs to the tail of queue
func (q *TwoLockQueue[T]) Push(vs ...T) {
	q.tlock.Lock()
	defer q.tlock.Unlock()

	for _, v := range vs {
		n := &lnode[T]{data: v}
		q.tail.next = n // Link node at the end of the linked list
		q.tail = n      // Swing Tail to node
	}
}

// Peek Retrieves, but does not remove, the head of this queue, or returns (nil, false) if this queue is empty.
func (q *TwoLockQueue[T]) Peek() (T, bool) {
	q.hlock.Lock()
	defer q.hlock.Unlock()

	head := q.head.next
	if head == nil {
		var v T
		return v, false
	}

	return head.data, true
}

// Poll Retrieves and removes the head of this queue, or returns (nil, false) if this queue is empty.
func (q *TwoLockQueue[T]) Poll() (T, bool) {
	q.hlock.Lock()
	defer q.hlock.Unlock()

	var v T
	head := q.head.next
	if head == nil {
		return v, false
	}

	o := head.data
	head.data = v
	q.head = head
	return o, true
}
