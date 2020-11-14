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
		log.Info("info", i)
		log.Debug("debug", i)
		log.Trace("trace", i)

		log.Fatalf("fatal(%d)", i)
		log.Errorf("error(%d)", i)
		log.Warnf("warning(%d)", i)
		log.Infof("info(%d)", i)
		log.Debugf("debug(%d)", i)
		log.Tracef("trace(%d)", i)
	}
}

// Test console info level
func TestConsoleWarn(t *testing.T) {
	log1 := NewLog()
	log1.AddWriter(&ConsoleWriter{Level: LevelWarn, Color: true})
	testConsoleCalls(log1, 1)
}

// Test console info level
func TestConsoleInfo(t *testing.T) {
	log1 := NewLog()
	log1.AddWriter(&ConsoleWriter{Level: LevelInfo, Color: true})
	testConsoleCalls(log1, 1)
}

// Test console trace level
func TestConsoleTrace(t *testing.T) {
	log1 := NewLog()
	log1.AddWriter(&ConsoleWriter{Level: LevelTrace, Color: true})
	testConsoleCalls(log1, 1)
}

// Test console without color
func TestConsoleNoColor(t *testing.T) {
	log := NewLog()
	log.AddWriter(&ConsoleWriter{Level: LevelTrace, Color: false})
	testConsoleCalls(log, 1)
}

// Test console async
func TestConsoleAsync(t *testing.T) {
	log := NewLog()
	log.AddWriter(&ConsoleWriter{Level: LevelTrace, Color: true})
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
	log1.AddWriter(&ConsoleWriter{Level: LevelTrace, Color: true})
	testConsoleCalls(log1, 1)
}
