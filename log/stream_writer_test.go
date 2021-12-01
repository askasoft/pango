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
	log.SetWriter(&StreamWriter{Color: true, Logfil: NewLevelFilter(LevelInfo)})
	log.SetFormatter(NewTextFormatter("[%c] %l - %m%n"))
	testConsoleCalls(log, 1)
	log.Close()
}

// Test console name filter
func TestConsoleFilterName(t *testing.T) {
	log := NewLog()
	log.SetWriter(&StreamWriter{Color: true, Logfil: NewNameFilter("out")})
	log.SetFormatter(NewTextFormatter("[%c] %l - %m%n"))

	log1 := log.GetLogger("OUT")
	testConsoleCalls(log1, 1)

	log2 := log.GetLogger("out")
	testConsoleCalls(log2, 1)
	log.Close()
}

// Test console multi filter
func TestConsoleFilterMulti(t *testing.T) {
	log := NewLog()
	log.SetWriter(&StreamWriter{Color: true, Logfil: NewMultiFilter(NewLevelFilter(LevelWarn), NewNameFilter("out"))})
	log.SetFormatter(NewTextFormatter("[%c] %l - %m%n"))

	log1 := log.GetLogger("OUT")
	testConsoleCalls(log1, 1)

	log2 := log.GetLogger("out")
	testConsoleCalls(log2, 1)
	log.Close()
}

// Test console without color
func TestConsoleNoColor(t *testing.T) {
	log := NewLog()
	log.SetWriter(&StreamWriter{Color: false})
	log.SetFormatter(NewTextFormatter("[%c] %l - %m%n"))
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
