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

func testNewConsoleWriter() Writer {
	return &StreamWriter{Color: true}
}

func testNewFileWriter(path, format string) Writer {
	fw := &FileWriter{Path: path}
	fw.SetFormat(format)
	return fw
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

func TestLogNilWriter(t *testing.T) {
	fmt.Println("\n\n--------------- TestLogNilWriter ---------------------")
	log0 := GetLogger("some")
	testLoggerCalls(log0)
}

func TestLogFuncs(t *testing.T) {
	fmt.Println("\n\n--------------- TestLogFuncs ---------------------")
	SetWriter(testNewConsoleWriter())
	SetLevel(LevelTrace)
	for i := 0; i < 1; i++ {
		Fatal("fatal")
		Error("error")
		Warn("warning")
		Info("info")
		Debug("debug")
		Trace("trace")
	}
}

func TestLogGetLogger(t *testing.T) {
	fmt.Println("\n\n-------------- TestLogGetLogger ------------------")
	log0 := GetLogger("some")
	testLoggerCalls(log0)
}

func TestLogNewLog(t *testing.T) {
	fmt.Println("\n\n-------------- TestLogNewLog ---------------------")
	log1 := NewLog()
	log1.SetLevel(LevelTrace)
	log1.SetWriter(testNewConsoleWriter())
	testLoggerCalls(log1)
}

func TestLogNewLogGetLogger(t *testing.T) {
	fmt.Println("\n\n-------------- TestLogNewLogGetLogger ---------")
	log1 := NewLog()
	log1.SetLevel(LevelTrace)
	log1.SetWriter(testNewConsoleWriter())
	log2 := log1.GetLogger("hello")
	testLoggerCalls(log2)
}

func TestLogNewLogProp(t *testing.T) {
	fmt.Println("\n\n-------------- TestLogNewLogProp ----------------")
	log1 := NewLog()
	log1.SetFormatter(NewTextFormatter("%x{key} %x{nil} - %m%T%n"))
	log1.SetLevel(LevelTrace)
	log1.SetWriter(testNewConsoleWriter())
	log1.SetProp("key", "val")
	testLoggerCalls(log1)
}

func TestLogDefault(t *testing.T) {
	fmt.Println("\n\n-------------- TestLogDefault ----------------")
	log1 := Default()
	log1.SetLevel(LevelTrace)
	log1.SetWriter(testNewConsoleWriter())
	log1.SetProp("key", "val")
	testLoggerCalls(log1)

	for i := 0; i < 1; i++ {
		Fatal("hello", "fatal")
		Error("hello", "error")
		Warn("hello", "warning")
		Info("hello", "info")
		Debug("hello", "debug")
		Trace("hello", "trace")
	}
}
