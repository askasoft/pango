package log

import (
	"fmt"
	"testing"
)

func newLogTestWriter() Writer {
	return &ConsoleWriter{Color: true}
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
	log.SetWriter(newLogTestWriter())
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
	log1.SetWriter(newLogTestWriter())
	testLoggerCalls(log1)
}

func TestLogNewLogGetLogger(t *testing.T) {
	fmt.Println("\n\n-------------- TestLogNewLogGetLogger ---------")
	log1 := NewLog()
	log1.SetLevel(LevelTrace)
	log1.SetWriter(newLogTestWriter())
	log2 := log1.GetLogger("hello")
	testLoggerCalls(log2)
}

func TestLogNewLogParam(t *testing.T) {
	fmt.Println("\n\n-------------- TestLogNewLogParam ----------------")
	log1 := NewLog()
	log1.SetFormatter(NewTextFormatter("%X{key} %X{nil} - %m%T%n"))
	log1.SetLevel(LevelTrace)
	log1.SetWriter(newLogTestWriter())
	log1.SetParam("key", "val")
	testLoggerCalls(log1)
}
