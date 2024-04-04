package log

import (
	"strings"
	"sync"
	"testing"
	"time"
)

const benchmarkTestEventCount = 10000

func BenchmarkLogEventNewWithStackTrace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < benchmarkTestEventCount; n++ {
			le := &Event{}
			le.Msg = ""
			le.Level = LevelError
			le.CallerDepth(5, true)
		}
	}
}

func BenchmarkLogEventNewWithoutStackTrace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < benchmarkTestEventCount; n++ {
			le := &Event{}
			le.Msg = ""
			le.Level = LevelError
			le.CallerDepth(5, false)
		}
	}
}

func BenchmarkLogEventNewWithoutPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < benchmarkTestEventCount; n++ {
			le := &Event{}
			le.Msg = ""
			le.CallerDepth(5, false)
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
			le.Msg = ""
			le.CallerDepth(5, false)
			testEventPool.Put(le)
		}
	}
}

func BenchmarkLogEventNewWithPoolParallel(b *testing.B) {
	// eventPool log event pool
	var eventPool = &sync.Pool{
		New: func() any {
			return &Event{}
		},
	}

	b.RunParallel(func(pb *testing.PB) {
		sb := &strings.Builder{}
		for pb.Next() {
			le := eventPool.Get().(*Event)
			le.Logger = &logger{}
			le.Level = LevelInfo
			le.Msg = "simple"
			le.Time = time.Now()
			TextFmtSimple.Write(sb, le)
			eventPool.Put(le)
		}
	})
}

func BenchmarkLogEventNewWithoutPoolParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		sb := &strings.Builder{}
		for pb.Next() {
			le := &Event{}
			le.Logger = &logger{}
			le.Level = LevelInfo
			le.Msg = "simple"
			le.Time = time.Now()
			TextFmtSimple.Write(sb, le)
		}
	})
}
