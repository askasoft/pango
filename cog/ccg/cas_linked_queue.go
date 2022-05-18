//go:build go1.18
// +build go1.18

package ccg

import (
	"sync/atomic"
	"unsafe"
)

// CQueue is a lock-free (CAS) unbounded queue.
// https://www.cs.rochester.edu/u/scott/papers/1996_PODC_queues.pdf
type CQueue[T any] struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

type cnode[T any] struct {
	data T
	next unsafe.Pointer
}

// NewCQueue returns an empty queue.
func NewCQueue[T any]() *CQueue[T] {
	n := unsafe.Pointer(&cnode[T]{})
	return &CQueue[T]{head: n, tail: n}
}

// IsEmpty returns true if the container length == 0
func (q *CQueue[T]) IsEmpty() bool {
	_, ok := q.Peek()
	return !ok
}

// Clear clears the container
func (q *CQueue[T]) Clear() {
	for {
		if _, ok := q.Poll(); !ok {
			return
		}
	}
}

// Push adds items of vs to the tail of queue
func (q *CQueue[T]) Push(vs ...T) {
	for _, v := range vs {
		n := &cnode[T]{data: v}
		for {
			tail := load[T](&q.tail)
			next := load[T](&tail.next)
			if tail == load[T](&q.tail) { // are tail and next consistent?
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
func (q *CQueue[T]) Peek() (T, bool) {
	for {
		head := load[T](&q.head)
		tail := load[T](&q.tail)
		next := load[T](&head.next)
		if head == load[T](&q.head) { // are head, tail, and next consistent?
			if head != tail {
				return next.data, true
			}

			if next == nil { // is queue empty?
				var v T
				return v, false
			}

			// tail is falling behind. try to advance it
			cas(&q.tail, tail, next)
		}
	}
}

// Poll Retrieves and removes the head of this queue, or returns (nil, false) if this queue is empty.
func (q *CQueue[T]) Poll() (T, bool) {
	for {
		head := load[T](&q.head)
		tail := load[T](&q.tail)
		next := load[T](&head.next)
		if head == load[T](&q.head) { // are head, tail, and next consistent?
			if head != tail {
				// read value before CAS otherwise another dequeue might free the next node
				v := next.data
				if cas(&q.head, head, next) {
					return v, true
				}
			} else {
				// is queue empty?
				if next == nil {
					var v T
					return v, false
				}

				// tail is falling behind. try to advance it
				cas(&q.tail, tail, next)
			}
		}
	}
}

func load[T any](p *unsafe.Pointer) (n *cnode[T]) {
	return (*cnode[T])(atomic.LoadPointer(p))
}

func cas[T any](p *unsafe.Pointer, old, new *cnode[T]) (ok bool) {
	return atomic.CompareAndSwapPointer(p, unsafe.Pointer(old), unsafe.Pointer(new))
}
