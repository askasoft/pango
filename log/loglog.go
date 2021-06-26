package log

import (
	"io"
	"sync"
)

// Log is default logger in application.
// it can contain several writers and log message into all writers.
type Log struct {
	logger *logger
	level  Level
	trace  Level

	async   bool
	evtChan chan *Event
	sigChan chan string
	waitg   sync.WaitGroup
	writer  Writer
	mutex   sync.Mutex
	levels  map[string]Level
	logfmt  Formatter
}

// NewLog returns a new Log.
func NewLog() *Log {
	return newLog(5)
}

func newLog(depth int) *Log {
	log := &Log{
		logger: &logger{
			depth: depth,
		},
		level:  LevelTrace,
		trace:  LevelError,
		levels: make(map[string]Level),
	}
	log.logger.log = log
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
	return &logger{
		log:   log,
		name:  name,
		depth: log.logger.depth,
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
// callerDepth: default is 1 (means +1)
// if the outputer is used by go std log, set callerDepth to 2
// example:
//   import (
//     golog "log"
//     "github.com/pandafw/pango/log"
//   )
//   golog.SetOutput(log.Outputer("GO", log.LevelInfo, 3))
//
func (log *Log) Outputer(name string, lvl Level, callerDepth ...int) io.Writer {
	lg := log.GetLogger(name)
	cd := 1
	if len(callerDepth) > 0 {
		cd = callerDepth[0]
	}
	lg.SetCallerDepth(lg.GetCallerDepth() + cd)
	return &outputer{logger: lg, level: lvl}
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
	lw := log.writer
	if lw != nil {
		lw.Write(le)
	}

	// put event back to pool
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
	for len(log.evtChan) > 0 {
		le := <-log.evtChan
		log.write(le)
	}
}

func (log *Log) flush() {
	if log.async {
		log.drain()
	}

	lw := log.writer
	if lw != nil {
		lw.Flush()
	}
}

func (log *Log) close() {
	lw := log.writer
	if lw != nil {
		lw.Close()
	}
}

/*----------------------------------------------------
 logger interface implements
----------------------------------------------------*/

// GetName return the logger's name
func (log *Log) GetName() string {
	return log.logger.name
}

// GetCallerDepth return the logger's depth
func (log *Log) GetCallerDepth() int {
	return log.logger.depth
}

// SetCallerDepth set the logger's caller depth (!!SLOW!!), 0: disable runtime.Caller()
func (log *Log) SetCallerDepth(d int) {
	log.logger.depth = d
}

// GetLevel return the logger's level
func (log *Log) GetLevel() Level {
	return log.level
}

// SetLevel set the logger's level
func (log *Log) SetLevel(lvl Level) {
	log.level = lvl
}

// GetTraceLevel return the logger's trace level
func (log *Log) GetTraceLevel() Level {
	return log.trace
}

// SetTraceLevel set the logger's trace level
func (log *Log) SetTraceLevel(lvl Level) {
	log.trace = lvl
}

// GetProp get logger property
func (log *Log) GetProp(k string) interface{} {
	ps := log.logger.props
	if ps == nil {
		return nil
	}
	return ps[k]
}

// SetProp set logger property
func (log *Log) SetProp(k string, v interface{}) {
	log.logger.SetProp(k, v)
}

// GetProps get logger properties
func (log *Log) GetProps() map[string]interface{} {
	tm := log.logger.props
	if tm == nil {
		return nil
	}

	// new return props
	nm := make(map[string]interface{}, len(tm))
	for k, v := range tm {
		nm[k] = v
	}
	return nm
}

// SetProps set logger properties
func (log *Log) SetProps(props map[string]interface{}) {
	log.logger.SetProps(props)
}

// GetFormatter get logger formatter
func (log *Log) GetFormatter() Formatter {
	return log.logfmt
}

// SetFormatter set logger formatter
func (log *Log) SetFormatter(lf Formatter) {
	log.logfmt = lf
}

// IsLevelEnabled is specified level enabled
func (log *Log) IsLevelEnabled(lvl Level) bool {
	return log.level > lvl
}

// Log log a message at specified level.
func (log *Log) Log(lvl Level, v ...interface{}) {
	log.logger._log(lvl, v...)
}

// Logf format and log a message at specified level.
func (log *Log) Logf(lvl Level, f string, v ...interface{}) {
	log.logger._logf(lvl, f, v...)
}

// IsFatalEnabled is FATAL level enabled
func (log *Log) IsFatalEnabled() bool {
	return log.IsLevelEnabled(LevelFatal)
}

// Fatal log a message at fatal level.
func (log *Log) Fatal(v ...interface{}) {
	log.logger._log(LevelFatal, v...)
}

// Fatalf format and log a message at fatal level.
func (log *Log) Fatalf(f string, v ...interface{}) {
	log.logger._logf(LevelFatal, f, v...)
}

// IsErrorEnabled is ERROR level enabled
func (log *Log) IsErrorEnabled() bool {
	return log.IsLevelEnabled(LevelError)
}

// Error log a message at error level.
func (log *Log) Error(v ...interface{}) {
	log.logger._log(LevelError, v...)
}

// Errorf format and log a message at error level.
func (log *Log) Errorf(f string, v ...interface{}) {
	log.logger._logf(LevelError, f, v...)
}

// IsWarnEnabled is WARN level enabled
func (log *Log) IsWarnEnabled() bool {
	return log.IsLevelEnabled(LevelWarn)
}

// Warn log a message at warning level.
func (log *Log) Warn(v ...interface{}) {
	log.logger._log(LevelWarn, v...)
}

// Warnf format and log a message at warning level.
func (log *Log) Warnf(f string, v ...interface{}) {
	log.logger._logf(LevelWarn, f, v...)
}

// IsInfoEnabled is INFO level enabled
func (log *Log) IsInfoEnabled() bool {
	return log.IsLevelEnabled(LevelInfo)
}

// Info log a message at info level.
func (log *Log) Info(v ...interface{}) {
	log.logger._log(LevelInfo, v...)
}

// Infof format and log a message at info level.
func (log *Log) Infof(f string, v ...interface{}) {
	log.logger._logf(LevelInfo, f, v...)
}

// IsDebugEnabled is DEBUG level enabled
func (log *Log) IsDebugEnabled() bool {
	return log.IsLevelEnabled(LevelDebug)
}

// Debug log a message at debug level.
func (log *Log) Debug(v ...interface{}) {
	log.logger._log(LevelDebug, v...)
}

// Debugf format log a message at debug level.
func (log *Log) Debugf(f string, v ...interface{}) {
	log.logger._logf(LevelDebug, f, v...)
}

// IsTraceEnabled is TRACE level enabled
func (log *Log) IsTraceEnabled() bool {
	return log.IsLevelEnabled(LevelTrace)
}

// Trace log a message at trace level.
func (log *Log) Trace(v ...interface{}) {
	log.logger._log(LevelTrace, v...)
}

// Tracef format and log a message at trace level.
func (log *Log) Tracef(f string, v ...interface{}) {
	log.logger._logf(LevelTrace, f, v...)
}
