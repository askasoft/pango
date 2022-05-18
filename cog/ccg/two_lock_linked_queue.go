//go:build go1.18
// +build go1.18

package ccg

import (
	"sync"
)

// LQueue is a concurrent unbounded queue which uses two-Lock concurrent queue qlgorithm.
type LQueue[T any] struct {
	head  *lnode[T]
	tail  *lnode[T]
	hlock sync.Mutex
	tlock sync.Mutex
}

type lnode[T any] struct {
	data T
	next *lnode[T]
}

// NewLQueue returns an empty LQueue.
func NewLQueue[T any]() *LQueue[T] {
	n := &lnode[T]{}
	return &LQueue[T]{head: n, tail: n}
}

// IsEmpty returns true if the container length == 0
func (q *LQueue[T]) IsEmpty() bool {
	_, ok := q.Peek()
	return !ok
}

// Clear clears the container
func (q *LQueue[T]) Clear() {
	q.tlock.Lock()
	defer q.tlock.Unlock()

	q.hlock.Lock()
	defer q.hlock.Unlock()

	//q.head.data = nil
	q.head.next = nil
	q.tail = q.head
}

// Push adds items of vs to the tail of queue
func (q *LQueue[T]) Push(vs ...T) {
	q.tlock.Lock()
	defer q.tlock.Unlock()

	for _, v := range vs {
		n := &lnode[T]{data: v}
		q.tail.next = n // Link node at the end of the linked list
		q.tail = n      // Swing Tail to node
	}
}

// Peek Retrieves, but does not remove, the head of this queue, or returns (nil, false) if this queue is empty.
func (q *LQueue[T]) Peek() (T, bool) {
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
func (q *LQueue[T]) Poll() (T, bool) {
	q.hlock.Lock()
	defer q.hlock.Unlock()

	head := q.head.next
	if head == nil {
		var v T
		return v, false
	}

	v := head.data
	//head.data = nil
	q.head = head
	return v, true
}
