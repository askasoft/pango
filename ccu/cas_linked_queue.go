package ccu

import (
	"sync/atomic"
	"unsafe"
)

// CQueue is a lock-free (CAS) unbounded queue.
// https://www.cs.rochester.edu/u/scott/papers/1996_PODC_queues.pdf
type CQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

type cnode struct {
	data T
	next unsafe.Pointer
}

// NewCQueue returns an empty queue.
func NewCQueue() *CQueue {
	n := unsafe.Pointer(&cnode{})
	return &CQueue{head: n, tail: n}
}

// IsEmpty returns true if the container length == 0
func (q *CQueue) IsEmpty() bool {
	_, ok := q.Peek()
	return !ok
}

// Clear clears the container
func (q *CQueue) Clear() {
	for {
		if _, ok := q.Poll(); !ok {
			return
		}
	}
}

// Push adds items of vs to the tail of queue
func (q *CQueue) Push(vs ...T) {
	for _, v := range vs {
		n := &cnode{data: v}
		for {
			tail := load(&q.tail)
			next := load(&tail.next)
			if tail == load(&q.tail) { // are tail and next consistent?
				if next == nil {
					if cas(&tail.next, next, n) {
						cas(&q.tail, tail, n) // Push is done.  try to swing tail to the inserted node
						return
					}
				} else { // tail was not pointing to the last node
					// try to swing Tail to the next node
					cas(&q.tail, tail, next)
				}
			}
		}
	}
}

// Peek Retrieves, but does not remove, the head of this queue, or returns (nil, false) if this queue is empty.
func (q *CQueue) Peek() (T, bool) {
	for {
		head := load(&q.head)
		tail := load(&q.tail)
		next := load(&head.next)
		if head == load(&q.head) { // are head, tail, and next consistent?
			if head != tail {
				return next.data, true
			}

			if next == nil { // is queue empty?
				return nil, false
			}

			// tail is falling behind. try to advance it
			cas(&q.tail, tail, next)
		}
	}
}

// Poll Retrieves and removes the head of this queue, or returns (nil, false) if this queue is empty.
func (q *CQueue) Poll() (T, bool) {
	for {
		head := load(&q.head)
		tail := load(&q.tail)
		next := load(&head.next)
		if head == load(&q.head) { // are head, tail, and next consistent?
			if head != tail {
				// read value before CAS otherwise another dequeue might free the next node
				v := next.data
				if cas(&q.head, head, next) {
					return v, true
				}
			} else {
				// is queue empty?
				if next == nil {
					return nil, false
				}

				// tail is falling behind. try to advance it
				cas(&q.tail, tail, next)
			}
		}
	}
}

func load(p *unsafe.Pointer) (n *cnode) {
	return (*cnode)(atomic.LoadPointer(p))
}

func cas(p *unsafe.Pointer, old, new *cnode) (ok bool) {
	return atomic.CompareAndSwapPointer(p, unsafe.Pointer(old), unsafe.Pointer(new))
}
