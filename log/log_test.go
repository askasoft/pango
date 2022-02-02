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
