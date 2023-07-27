package log

import (
	"sync"
	"testing"
)

const benchmarkTestEventCount = 10000

func BenchmarkLogEventNewWithStackTrace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < benchmarkTestEventCount; n++ {
			le := &Event{}
			le.msg = ""
			le.level = LevelError
			le.caller(5, true)
		}
	}
}

func BenchmarkLogEventNewWithoutStackTrace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < benchmarkTestEventCount; n++ {
			le := &Event{}
			le.msg = ""
			le.level = LevelError
			le.caller(5, false)
		}
	}
}

func BenchmarkLogEventNewWithoutPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < benchmarkTestEventCount; n++ {
			le := &Event{}
			le.msg = ""
			le.caller(5, false)
		}
	}
}

var testEventPool = sync.Pool{
	New: func() any {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		return &Event{}
	},
}

func BenchmarkLogEventNewWithPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < benchmarkTestEventCount; n++ {
			le := testEventPool.Get().(*Event)
			le.msg = ""
			le.caller(5, false)
		}
	}
}
