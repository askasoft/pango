package log

import (
	"encoding/json"
	"fmt"
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

func TestEventJsonMarshall(t *testing.T) {
	le := newEvent(&logger{}, LevelInfo, "caller")
	le.When = time.Now()
	le.CallerDepth(2, false)

	bs, _ := json.Marshal(le)
	fmt.Println(string(bs))
}
