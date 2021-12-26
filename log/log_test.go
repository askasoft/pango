package log

import (
	"fmt"
	"path"
	"runtime"
	"testing"
)

func skipTest(t *testing.T, msg string) {
	fmt.Println(msg)
	t.Skip(msg)
}

func testGetCaller(offset int) (string, int, string) {
	rpc := make([]uintptr, 1)
	n := runtime.Callers(2, rpc)
	if n > 0 {
		frames := runtime.CallersFrames(rpc)
		frame, _ := frames.Next()
		_, ffun := path.Split(frame.Function)
		_, file := path.Split(frame.File)
		line := frame.Line + offset
		return file, line, ffun
	}
	return "???", 0, "???"
}

// Try each log level in decreasing order of priority.
func testLoggerCalls(l Logger) {
	for i := 0; i < 1; i++ {
		l.Fatal("hello", "fatal")
		l.Error("hello", "error")
		l.Warn("hello", "warning")
		l.Info("hello", "info")
		l.Debug("hello", "debug")
		l.Trace("hello", "trace")
	}
}
