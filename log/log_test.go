package log

import (
	"fmt"
	"testing"
)

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

func TestLogFuncs(t *testing.T) {
	fmt.Println("\n\n--------------- Funcs ---------------------")
	log.AddWriter(&ConsoleWriter{Level: LevelTrace})
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
	fmt.Println("\n\n-------------- GetLogger() ------------------")
	log0 := GetLogger("some")
	testLoggerCalls(log0)
}

func TestLogNewLog(t *testing.T) {
	fmt.Println("\n\n-------------- NewLog() ---------------------")
	log1 := NewLog()
	log1.SetLevel(LevelTrace)
	log1.AddWriter(&ConsoleWriter{Level: LevelTrace})
	testLoggerCalls(log1)
}

func TestLogNewLogGetLogger(t *testing.T) {
	fmt.Println("\n\n-------------- NewLog().GetLogger() ---------")
	log1 := NewLog()
	log1.SetLevel(LevelTrace)
	log1.AddWriter(&ConsoleWriter{Level: LevelTrace})
	log2 := log1.GetLogger("hello")
	testLoggerCalls(log2)
}

func TestLogNewLogParam(t *testing.T) {
	fmt.Println("\n\n-------------- NewLog().Param ----------------")
	log1 := NewLog()
	log1.SetFormatter(NewFormatter("%X{key} %X{nil} - %m%T%n"))
	log1.SetLevel(LevelTrace)
	log1.AddWriter(&ConsoleWriter{Level: LevelTrace})
	log1.SetParam("key", "val")
	testLoggerCalls(log1)
}
