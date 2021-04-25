package log

import (
	"io"
	"sync"
)

// Log is default logger in application.
// it can contain several writers and log message into all writers.
type Log struct {
	logger
	async   bool
	evtChan chan *Event
	sigChan chan string
	waitg   sync.WaitGroup
	writer  Writer
	mutex   sync.Mutex
	levels  map[string]Level
}

// NewLog returns a new Log.
func NewLog() *Log {
	return newLog(5)
}

func newLog(depth int) *Log {
	log := &Log{}
	log.log = log
	log.level = LevelTrace
	log.depth = depth
	log.trace = LevelError
	log.levels = make(map[string]Level)
	return log
}

// SetLevels set the logger levels
func (log *Log) SetLevels(lvls map[string]Level) {
	log.levels = lvls
}

// getLoggerLevel get the named logger level
func (log *Log) getLoggerLevel(name string) Level {
	level := log.levels[name]
	if level == LevelNone {
		level = log.GetLevel()
	}
	return level
}

// GetLogger returns a new Logger with name
func (log *Log) GetLogger(name string) Logger {
	level := log.getLoggerLevel(name)
	return &logger{
		name:   name,
		log:    log,
		logfmt: log.logfmt,
		depth:  log.depth,
		level:  level,
		trace:  log.trace,
	}
}

// Async set the log to asynchronous and start the goroutine
// if size < 1 then stop async goroutine
func (log *Log) Async(size int) *Log {
	log.mutex.Lock()
	defer log.mutex.Unlock()

	if size < 1 {
		if log.async {
			log.stopAsync()
		}
		return log
	}

	if log.async {
		if size == len(log.evtChan) {
			return log
		}
		log.stopAsync()
	}

	log.async = true
	log.evtChan = make(chan *Event, size)
	log.sigChan = make(chan string, 1)
	go log.startAsync()
	return log
}

// GetWriter get the log writer
func (log *Log) GetWriter() Writer {
	return log.writer
}

// SetWriter set the log writer
func (log *Log) SetWriter(lw Writer) {
	log.mutex.Lock()
	defer log.mutex.Unlock()

	log.close()
	log.writer = lw
}

// Flush flush all chan data.
func (log *Log) Flush() {
	log.mutex.Lock()
	defer log.mutex.Unlock()

	if log.async {
		log.execSignal("flush")
		return
	}

	log.flush()
}

// startAsync start async log goroutine
func (log *Log) startAsync() {
	done := false
	for {
		select {
		case le := <-log.evtChan:
			log.write(le)
		case sg := <-log.sigChan:
			// Now should only send "flush" or "close" to bl.sigChan
			log.flush()
			switch sg {
			case "close":
				log.close()
				done = true
			case "done":
				done = true
			}
			log.waitg.Done()
		}
		if done {
			break
		}
	}
}

// stopAsync flush and stop async goroutine
func (log *Log) stopAsync() {
	log.execSignal("done")

	log.async = false
	log.drain()
	close(log.evtChan)
	close(log.sigChan)
}

// execSignal send a signal and wait for done
func (log *Log) execSignal(sig string) {
	log.waitg.Add(1)
	log.sigChan <- sig
	log.waitg.Wait()
}

func (log *Log) write(le *Event) {
	if log.writer != nil {
		log.writer.Write(le)
	}
	putEvent(le)
}

// submit submit a log event
func (log *Log) submit(le *Event) {
	if log.async {
		log.evtChan <- le
		return
	}

	log.mutex.Lock()
	log.write(le)
	log.mutex.Unlock()
}

func (log *Log) drain() {
	for {
		if len(log.evtChan) > 0 {
			le := <-log.evtChan
			log.write(le)
			continue
		}
		break
	}
}

func (log *Log) flush() {
	if log.async {
		log.drain()
	}

	if log.writer != nil {
		log.writer.Flush()
	}
}

func (log *Log) close() {
	if log.writer != nil {
		log.writer.Close()
		log.writer = nil
	}
}

// Close close logger, flush all chan data and close the writer.
func (log *Log) Close() {
	log.mutex.Lock()
	defer log.mutex.Unlock()

	if log.async {
		log.execSignal("close")
		close(log.evtChan)
		close(log.sigChan)
		log.async = false
		return
	}

	log.flush()
	log.close()
}

// Outputer return a io.Writer for go log.SetOutput
func (log *Log) Outputer(name string, lvl Level) io.Writer {
	lg := log.GetLogger(name)
	lg.SetCallerDepth(lg.GetCallerDepth() + 2)
	return &outputer{logger: lg, level: lvl}
}
