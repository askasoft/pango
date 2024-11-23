package queuetest

import (
	"runtime"
	"strconv"
	"testing"

	"github.com/askasoft/pango/cog/ccg/casqueue"
	"github.com/askasoft/pango/cog/ccg/twolockqueue"
	"github.com/askasoft/pango/cog/ccg/unboundedchan"
)

// queue queue interface
type queue[T any] interface {
	// Push adds items of vs to the tail of queue
	Push(vs ...T)

	// Poll Retrieves and removes the head of this queue, or returns (nil, false) if this queue is empty.
	Poll() (T, bool)
}

func benchmarkQueue(b *testing.B, name string, q queue[int]) {
	for _, cpus := range []int{4, 8, 32, 1024} {
		runtime.GOMAXPROCS(cpus)
		b.Run(name+"#"+strconv.Itoa(cpus), func(b *testing.B) {
			b.ResetTimer()

			for i := 0; i < cpus; i++ {
				q.Push(i)
			}

			b.RunParallel(func(pb *testing.PB) {
				for c := 0; pb.Next(); c++ {
					if c&1 == 0 {
						q.Push(c)
					} else {
						q.Poll()
					}
				}
			})
		})
	}
}

func BenchmarkLockFreeQueue(b *testing.B) {
	benchmarkQueue(b, "lock-free queue", casqueue.NewCasQueue[int]())
}

func BenchmarkTwoLockQueue(b *testing.B) {
	benchmarkQueue(b, "two-lock queue", twolockqueue.NewTwoLockQueue[int]())
}

func BenchmarkUnboundedChan(b *testing.B) {
	ubc := unboundedchan.NewUnboundedChan[int](1000)
	benchmarkQueue(b, "unbounded chan", ubc)
	ubc.Close()
}
