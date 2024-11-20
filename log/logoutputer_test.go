package log

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"testing"

	golog "log"
)

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

func TestGoLogOutputGlobal(t *testing.T) {
	fmt.Println("\n\n--------------- TestGoLogOutputGlobal ---------------------")
	SetWriter(NewConsoleWriter())
	golog.SetOutput(GetOutputer("golog", LevelInfo))
	golog.Print("hello", "golog")
}

func TestGoLogOutputNewLog(t *testing.T) {
	fmt.Println("\n\n--------------- TestGoLogOutputNewLog ---------------------")

	lg := NewLog()
	lg.SetWriter(NewConsoleWriter())
	golog.SetOutput(lg.GetOutputer("std", LevelInfo))
	golog.Print("hello", "golog")
}

func TestGoLogCallerGlobal(t *testing.T) {
	sb := &strings.Builder{}

	sw := &StreamWriter{Output: sb}
	sw.SetFormat("%l %S:%L %F() - %m")
	SetWriter(sw)

	golog.SetFlags(0)
	golog.SetOutput(GetOutputer("golog", LevelInfo, 3))
	file, line, ffun := testGetCaller(1)
	golog.Print("hello", "golog")
	Close()

	a := sb.String()
	w := fmt.Sprintf("INFO %s:%d %s() - hellogolog\n", file, line, ffun)
	if a != w {
		t.Errorf("output = %v\n, want = %v", a, w)
	}
}

func TestGoLogCallerNewLog(t *testing.T) {
	sb := &strings.Builder{}

	lg := NewLog()
	sw := &StreamWriter{Output: sb}
	sw.SetFormat("%l %S:%L %F() - %m")
	lg.SetWriter(sw)

	golog.SetFlags(0)
	golog.SetOutput(lg.GetOutputer("std", LevelInfo, 3))
	file, line, ffun := testGetCaller(1)
	golog.Print("hello", "golog")
	lg.Close()

	a := sb.String()
	w := fmt.Sprintf("INFO %s:%d %s() - hellogolog\n", file, line, ffun)
	if a != w {
		t.Errorf("output = %v\n, want = %v", a, w)
	}
}

func TestIoWriterCallerGlobal(t *testing.T) {
	sb := &strings.Builder{}

	sw := &StreamWriter{Output: sb}
	sw.SetFormat("%l %S:%L %F() - %m%n")
	SetWriter(sw)

	iow := GetOutputer("iow", LevelInfo)
	file, line, ffun := testGetCaller(1)
	iow.Write(([]byte)("hello writer"))
	Close()

	a := sb.String()
	w := fmt.Sprintf("INFO %s:%d %s() - hello writer"+EOL, file, line, ffun)
	if a != w {
		t.Errorf("output = %v\n, want = %v", a, w)
	}
}

func TestIoWriterFileCallerNewLog(t *testing.T) {
	sb := &strings.Builder{}

	lg := NewLog()
	sw := &StreamWriter{Output: sb}
	sw.SetFormat("%l %S:%L %F() - %m%n")
	lg.SetWriter(sw)

	iow := lg.GetOutputer("iow", LevelInfo)
	file, line, ffun := testGetCaller(1)
	iow.Write(([]byte)("hello writer"))
	lg.Close()

	a := sb.String()
	w := fmt.Sprintf("INFO %s:%d %s() - hello writer"+EOL, file, line, ffun)
	if a != w {
		t.Errorf("output = %v\n, want = %v", a, w)
	}
}
