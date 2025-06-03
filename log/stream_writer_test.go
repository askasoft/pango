package log

import (
	"testing"
)

func testConsoleCalls(log Logger, loop int) {
	for i := 0; i < loop; i++ {
		log.Fatal("fatal<", i, ">")
		log.Error("error<", i, ">")
		log.Warn("warn <", i, ">")
		log.Info("info <", i, ">")
		log.Debug("debug<", i, ">")
		log.Trace("trace<", i, ">")

		log.Fatalf("fatal(%d)", i)
		log.Errorf("error(%d)", i)
		log.Warnf("warn (%d)", i)
		log.Infof("info (%d)", i)
		log.Debugf("debug(%d)", i)
		log.Tracef("trace(%d)", i)
	}
}

// Test console info level filter
func TestConsoleFilterInfo(t *testing.T) {
	log := NewLog()

	sw := &StreamWriter{Color: true}
	sw.Filter = NewLevelFilter(LevelInfo)
	sw.SetFormat("[%c] %l - %m%n")

	log.SetWriter(sw)
	testConsoleCalls(log, 1)
	log.Close()
}

// Test console name filter
func TestConsoleFilterName(t *testing.T) {
	log := NewLog()
	sw := &StreamWriter{Color: true}
	sw.Filter = NewNameFilter("out")
	sw.SetFormat("[%c] %l - %m%n")
	log.SetWriter(sw)

	log1 := log.GetLogger("OUT")
	testConsoleCalls(log1, 1)

	log2 := log.GetLogger("out")
	testConsoleCalls(log2, 1)
	log.Close()
}

// Test console multi filter
func TestConsoleFilterMulti(t *testing.T) {
	log := NewLog()
	sw := &StreamWriter{Color: true}
	sw.SetFilter("level:warn name:out")
	sw.SetFormat("[%c] %l - %m%n")
	log.SetWriter(sw)

	log1 := log.GetLogger("OUT")
	testConsoleCalls(log1, 1)

	log2 := log.GetLogger("out")
	testConsoleCalls(log2, 1)
	log.Close()
}

// Test console without color
func TestConsoleNoColor(t *testing.T) {
	log := NewLog()

	sw := &StreamWriter{Color: false}
	sw.SetFormat("[%c] %l - %m%n")
	log.SetWriter(sw)
	testConsoleCalls(log, 1)
	log.Close()
}

// Test console strace trace
func TestConsoleStackTrace(t *testing.T) {
	log := NewLog()
	log.SetWriter(&StreamWriter{Color: true})
	testConsoleCalls(log, 1)
	log.Close()
}
