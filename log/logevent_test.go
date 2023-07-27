package log

import (
	"strings"
	"sync"
	"testing"
	"time"
)

func TestEventCaller(t *testing.T) {
	le := newEvent(&logger{}, LevelInfo, "caller")
	le.When = time.Time{}
	le.CallerDepth(2, false)

	if le.File != "logevent_test.go" {
		t.Errorf("le.file = %v, want %v", le.File, "logevent_test.go")
	}
	if le.Func != "log.TestEventCaller" {
		t.Errorf("le._func = %v, want %v", le.Func, "log.TestEventCaller")
	}
	if le.Line == 0 {
		t.Errorf("le.line = %v, want != %v", le.Line, 0)
	}
}

func BenchmarkEventPool(b *testing.B) {
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
			le.When = time.Now()
			TextFmtSimple.Write(sb, le)
			eventPool.Put(le)
		}
	})
}

func BenchmarkEventNew(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		sb := &strings.Builder{}
		for pb.Next() {
			le := &Event{}
			le.Logger = &logger{}
			le.Level = LevelInfo
			le.Msg = "simple"
			le.When = time.Now()
			TextFmtSimple.Write(sb, le)
		}
	})
}
