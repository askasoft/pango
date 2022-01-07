package ccu

import (
	"runtime"
	"strconv"
	"testing"
)

func benchmarkQueue(b *testing.B, name string, q Queue) {
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
	benchmarkQueue(b, "lock-free queue", NewCQueue())
}

func BenchmarkTwoLockQueue(b *testing.B) {
	benchmarkQueue(b, "two-lock queue", NewCQueue())
}

func BenchmarkUnboundedChan(b *testing.B) {
	ubc := NewUnboundedChan(1000)
	benchmarkQueue(b, "unbounded chan", ubc)
	ubc.Close()
}
