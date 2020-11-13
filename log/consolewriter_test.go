package log

import (
	"testing"
	"time"
)

func testConsoleCalls(log Logger, loop int) {
	for i := 0; i < loop; i++ {
		log.Fatal("fatal", i)
		log.Error("error", i)
		log.Warn("warning", i)
		log.Debug("debug", i)
		log.Trace("trace", i)
		log.Fatalf("fatal(%d)", i)
		log.Errorf("error(%d)", i)
		log.Warnf("warning(%d)", i)
		log.Debugf("debug(%d)", i)
		log.Tracef("trace(%d)", i)
	}
}

// Test console info level
func TestConsoleInfo(t *testing.T) {
	log1 := NewLog()
	log1.AddWriter(&ConsoleWriter{Level: LevelInfo})
	testConsoleCalls(log1, 10)
}

// Test console trace level
func TestConsoleTrace(t *testing.T) {
	log1 := NewLog()
	log1.AddWriter(&ConsoleWriter{Level: LevelTrace})
	testConsoleCalls(log1, 10)
}

// Test console with color
func TestConsoleWithColor(t *testing.T) {
	log := NewLog()
	log.AddWriter(&ConsoleWriter{Level: LevelTrace, Color: true})
	testConsoleCalls(log, 10)
}

// Test console async
func TestConsoleAsync(t *testing.T) {
	log := NewLog()
	log.AddWriter(&ConsoleWriter{Level: LevelTrace})
	log.Async(100)
	go testConsoleCalls(log, 100)
	go testConsoleCalls(log, 100)
	time.Sleep(1 * time.Second)
	for len(log.evtChan) != 0 {
		time.Sleep(1 * time.Millisecond)
	}
	log.Close()
}

// Test console strace trace
func TestConsoleStackTrace(t *testing.T) {
	log1 := NewLog()
	log1.AddWriter(&ConsoleWriter{Level: LevelTrace})
	testConsoleCalls(log1, 1)
}
