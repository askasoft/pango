package log

import (
	"strings"
	"testing"
	"time"
)

func TestEventCaller(t *testing.T) {
	le := newEvent(&logger{}, LevelInfo, "caller")
	le.when = time.Time{}
	le.Caller(2, false)

	if le.file != "logevent_test.go" {
		t.Errorf("le.file = %v, want %v", le.file, "logevent_test.go")
	}
	if le._func != "log.TestEventCaller" {
		t.Errorf("le._func = %v, want %v", le._func, "log.TestEventCaller")
	}
	if le.line == 0 {
		t.Errorf("le.line = %v, want != %v", le.line, 0)
	}
}

func BenchmarkEventPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		sb := &strings.Builder{}
		for pb.Next() {
			le := eventPool.Get().(*Event)
			le.logger = &logger{}
			le.level = LevelInfo
			le.msg = "simple"
			le.when = time.Now()
			TextFmtSimple.Write(sb, le)
			le.clear()
			eventPool.Put(le)
		}
	})
}

func BenchmarkEventNew(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		sb := &strings.Builder{}
		for pb.Next() {
			le := &Event{}
			le.logger = &logger{}
			le.level = LevelInfo
			le.msg = "simple"
			le.when = time.Now()
			TextFmtSimple.Write(sb, le)
		}
	})
}
