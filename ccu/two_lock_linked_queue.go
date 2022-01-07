package ccu

import (
	"sync"
)

// LQueue is a concurrent unbounded queue which uses two-Lock concurrent queue qlgorithm.
type LQueue struct {
	head  *lnode
	tail  *lnode
	hlock sync.Mutex
	tlock sync.Mutex
}

type lnode struct {
	data T
	next *lnode
}

// NewLQueue returns an empty LQueue.
func NewLQueue() *LQueue {
	n := &lnode{}
	return &LQueue{head: n, tail: n}
}

// IsEmpty returns true if the container length == 0
func (q *LQueue) IsEmpty() bool {
	_, ok := q.Peek()
	return !ok
}

// Clear clears the container
func (q *LQueue) Clear() {
	q.tlock.Lock()
	defer q.tlock.Unlock()

	q.hlock.Lock()
	defer q.hlock.Unlock()

	q.head.data = nil
	q.head.next = nil
	q.tail = q.head
}

// Push adds items of vs to the tail of queue
func (q *LQueue) Push(vs ...T) {
	q.tlock.Lock()
	defer q.tlock.Unlock()

	for _, v := range vs {
		n := &lnode{data: v}
		q.tail.next = n // Link node at the end of the linked list
		q.tail = n      // Swing Tail to node
	}
}

// Peek Retrieves, but does not remove, the head of this queue, or returns (nil, false) if this queue is empty.
func (q *LQueue) Peek() (T, bool) {
	q.hlock.Lock()
	defer q.hlock.Unlock()

	head := q.head.next
	if head == nil {
		return nil, false
	}

	return head.data, true
}

// Poll Retrieves and removes the head of this queue, or returns (nil, false) if this queue is empty.
func (q *LQueue) Poll() (T, bool) {
	q.hlock.Lock()
	defer q.hlock.Unlock()

	head := q.head.next
	if head == nil {
		return nil, false
	}

	v := head.data
	head.data = nil
	q.head = head
	return v, true
}
