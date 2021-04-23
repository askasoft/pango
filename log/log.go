// Package log provide a general log interface
// Usage:
//
// import "github.com/pandafw/pango/log"
//
//	log := log.NewLog()
//	log.SetWriter(&log.StreamWriter{Color: true})
//	log.Async(1000)
//
// Use it like this:
//	log.Fatal("fatal")
//	log.Error("error")
//	log.Warn("warning")
//	log.Info("info")
//	log.Debug("debug")
//	log.Trace("trace")
//
// A Logger with name:
//	log := log.GetLogger("foo")
//	log.Debug("hello")
//
package log

import (
	"io"
	"sync"

	"github.com/fsnotify/fsnotify"
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
	levels  map[string]int
	watcher *fsnotify.Watcher
}

// GetLoggerLevel get the named logger level
func (log *Log) GetLoggerLevel(name string) int {
	level := log.levels[name]
	if level == LevelNone {
		level = log.GetLevel()
	}
	return level
}

// SetLoggerLevel set the named logger level
func (log *Log) SetLoggerLevel(name string, level int) {
	log.levels[name] = level
}

// GetLogger returns a new Logger with name
func (log *Log) GetLogger(name string) Logger {
	if name == "" || name == ROOT {
		return log
	}

	level := log.GetLoggerLevel(name)
	return &logger{name: name, log: log, logfmt: log.logfmt, depth: log.depth, level: level, trace: log.trace}
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

	if log.watcher != nil {
		log.watcher.Close()
		log.watcher = nil
	}

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
func (log *Log) Outputer(lvl int) io.Writer {
	lg := log.GetLogger("golog")
	lg.SetCallerDepth(lg.GetCallerDepth() + 2)
	return &outputer{logger: lg, level: lvl}
}

//--------------------------------------------------------------------
// package functions
//

// ROOT default logger name
const ROOT = "_"

// NewLog returns a new Log.
func NewLog() *Log {
	return newLog(5)
}

// default package Log instance
var log = newLog(6)

func newLog(depth int) *Log {
	log := &Log{}
	log.log = log
	log.name = ROOT
	log.level = LevelTrace
	log.depth = depth
	log.trace = LevelError
	log.levels = make(map[string]int)
	return log
}

// Default get default Log
func Default() *Log {
	return log
}

// GetLogger returns a new logger
func GetLogger(name string) Logger {
	if name == "" || name == ROOT {
		return log
	}

	l := log.GetLogger(name)
	l.SetCallerDepth(GetCallerDepth() - 1)
	return l
}

// Outputer return a io.Writer for go log.SetOutput
func Outputer(lvl int) io.Writer {
	lg := GetLogger("golog")
	lg.SetCallerDepth(lg.GetCallerDepth() + 2)
	return &outputer{logger: lg, level: lvl}
}

// Async set the Log with Async mode and hold msglen messages
func Async(msgLen int) {
	log.Async(msgLen)
}

// IsAsync return the logger's async
func IsAsync() bool {
	return log.async
}

// SetFormatter set the formatter.
func SetFormatter(lf Formatter) {
	log.SetFormatter(lf)
}

// SetWriter set the writer.
func SetWriter(lw Writer) {
	log.SetWriter(lw)
}

// Close will remove all writers and stop async goroutine
func Close() {
	log.Close()
}

// GetLevel return the logger's level
func GetLevel() int {
	return log.GetLevel()
}

// SetLevel set the logger's level
func SetLevel(lvl int) {
	log.SetLevel(lvl)
}

// GetCallerDepth return the logger's caller depth
func GetCallerDepth() int {
	return log.GetCallerDepth()
}

// SetCallerDepth set the logger's caller depth (!!SLOW!!), 0: disable runtime.Caller()
func SetCallerDepth(d int) {
	log.SetCallerDepth(d)
}

// IsFatalEnabled is FATAL level enabled
func IsFatalEnabled() bool {
	return log.IsFatalEnabled()
}

// Fatal log a message at fatal level.
func Fatal(v ...interface{}) {
	log.Fatal(v...)
}

// Fatalf format and log a message at fatal level.
func Fatalf(f string, v ...interface{}) {
	log.Fatalf(f, v...)
}

// IsErrorEnabled is ERROR level enabled
func IsErrorEnabled() bool {
	return log.IsErrorEnabled()
}

// Error log a message at error level.
func Error(v ...interface{}) {
	log.Error(v...)
}

// Errorf format and log a message at error level.
func Errorf(f string, v ...interface{}) {
	log.Errorf(f, v...)
}

// IsWarnEnabled is WARN level enabled
func IsWarnEnabled() bool {
	return log.IsWarnEnabled()

}

// Warn log a message at warning level.
func Warn(v ...interface{}) {
	log.Warn(v...)
}

// Warnf format and log a message at warning level.
func Warnf(f string, v ...interface{}) {
	log.Warnf(f, v...)
}

// IsInfoEnabled is INFO level enabled
func IsInfoEnabled() bool {
	return log.IsInfoEnabled()
}

// Info log a message at info level.
func Info(v ...interface{}) {
	log.Info(v...)
}

// Infof format and log a message at info level.
func Infof(f string, v ...interface{}) {
	log.Infof(f, v...)
}

// IsDebugEnabled is DEBUG level enabled
func IsDebugEnabled() bool {
	return log.IsDebugEnabled()
}

// Debug log a message at debug level.
func Debug(v ...interface{}) {
	log.Debug(v...)
}

// Debugf format log a message at debug level.
func Debugf(f string, v ...interface{}) {
	log.Debugf(f, v...)
}

// IsTraceEnabled is TRACE level enabled
func IsTraceEnabled() bool {
	return log.IsTraceEnabled()
}

// Trace log a message at trace level.
func Trace(v ...interface{}) {
	log.Trace(v...)
}

// Tracef format and log a message at trace level.
func Tracef(f string, v ...interface{}) {
	log.Tracef(f, v...)
}
