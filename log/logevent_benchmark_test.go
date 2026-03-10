package log

import (
	"strings"
	"sync"
	"testing"
	"time"
)

const benchmarkTestEventCount = 10000

func BenchmarkLogEventNewWithStackTrace(b *testing.B) {
	for b.Loop() {
		for range benchmarkTestEventCount {
			le := &Event{}
			le.Level = LevelInfo
			le.Message = "simple"
			le.Time = time.Now()
			le.CallerSkip(5, true)
			nopWriter.Write(le)
		}
	}
}

func BenchmarkLogEventNewWithoutStackTrace(b *testing.B) {
	for b.Loop() {
		for range benchmarkTestEventCount {
			le := &Event{}
			le.Level = LevelInfo
			le.Message = "simple"
			le.Time = time.Now()
			le.CallerSkip(5, false)
			nopWriter.Write(le)
		}
	}
}

func BenchmarkLogEventNewWithoutPool(b *testing.B) {
	sb := &strings.Builder{}
	for b.Loop() {
		for range benchmarkTestEventCount {
			le := &Event{}
			le.Level = LevelInfo
			le.Message = "simple"
			le.Time = time.Now()
			le.CallerSkip(5, false)
			sb.Reset()
			TextFmtSimple.Write(sb, le)
			nopWriter.Write(le)
		}
	}
}

func BenchmarkLogEventNewWithPool(b *testing.B) {
	// eventPool log event pool
	var eventPool = &sync.Pool{
		New: func() any {
			return &Event{}
		},
	}

	sb := &strings.Builder{}
	for b.Loop() {
		for range benchmarkTestEventCount {
			le := eventPool.Get().(*Event)
			le.Level = LevelInfo
			le.Message = "simple"
			le.Time = time.Now()
			le.CallerSkip(5, false)
			sb.Reset()
			TextFmtSimple.Write(sb, le)
			nopWriter.Write(le)
			eventPool.Put(le)
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
			le.Level = LevelInfo
			le.Message = "simple"
			le.Time = time.Now()
			sb.Reset()
			TextFmtSimple.Write(sb, le)
			nopWriter.Write(le)
			eventPool.Put(le)
		}
	})
}

func BenchmarkLogEventNewWithoutPoolParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		sb := &strings.Builder{}
		for pb.Next() {
			le := &Event{}
			le.Level = LevelInfo
			le.Message = "simple"
			le.Time = time.Now()
			sb.Reset()
			TextFmtSimple.Write(sb, le)
			nopWriter.Write(le)
		}
	})
}
