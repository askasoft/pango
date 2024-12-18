package log

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/askasoft/pango/mag"
	"github.com/askasoft/pango/str"
)

// Log is default logger in application.
// it can contain several writers and log message into all writers.
type Log struct {
	name  string
	depth int
	props map[string]any

	level Level
	trace Level

	writer Writer
	levels map[string]Level

	mutex sync.Mutex
}

var emptyProps = make(map[string]any)

// NewLog returns a new Log.
func NewLog() *Log {
	log := &Log{
		name:   "_",
		depth:  5,
		props:  emptyProps,
		level:  LevelTrace,
		trace:  LevelError,
		levels: make(map[string]Level),
		writer: NewStdoutWriter(),
	}

	runtime.SetFinalizer(log, finalClose)
	return log
}

func finalClose(log *Log) {
	log.Close()
}

// SetLevels set the logger levels
func (log *Log) SetLevels(lvls map[string]Level) {
	log.levels = lvls
}

// GetLoggerLevel get the named logger level
func (log *Log) GetLoggerLevel(name string) Level {
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
		name:  str.IfEmpty(name, "_"),
		depth: log.depth,
		props: log.props,
	}
}

// GetWriter get the log writer
func (log *Log) GetWriter() Writer {
	return log.writer
}

// SetWriter set the log writer
func (log *Log) SetWriter(lw Writer) {
	log.mutex.Lock()
	defer log.mutex.Unlock()

	log.writer = lw
}

// SwitchWriter use lw to replace the log writer
func (log *Log) SwitchWriter(lw Writer) {
	log.mutex.Lock()
	defer log.mutex.Unlock()

	ow := log.writer
	log.writer = lw

	if osw, ok := ow.(*SyncWriter); ok {
		osw.SetWriter(lw)
		return
	}

	if oaw, ok := ow.(*AsyncWriter); ok {
		oaw.SetWriter(lw)
		oaw.Stop()
		return
	}

	ow.Close()
}

// Flush flush all chan data.
func (log *Log) Flush() {
	safeFlush(log.writer)
}

// Close close logger, flush all data and close the writer.
func (log *Log) Close() {
	safeClose(log.writer)
}

// Outputer return a io.Writer for go log.SetOutput
// callerDepth: default is 1 (means +1)
// if the outputer is used by go std log, set callerDepth to 2
// example:
//
//	import (
//	  golog "log"
//	  "github.com/askasoft/pango/log"
//	)
//	golog.SetOutput(log.Outputer("GO", log.LevelInfo, 3))
func (log *Log) GetOutputer(name string, lvl Level, callerDepth ...int) Outputer {
	lg := log.GetLogger(name)
	cd := 1
	if len(callerDepth) > 0 {
		cd = callerDepth[0]
	}
	lg.SetCallerDepth(lg.GetCallerDepth() + cd)
	return &outputer{logger: lg, level: lvl}
}

/*----------------------------------------------------
 logger interface implements
----------------------------------------------------*/

// GetName return the logger's name
func (log *Log) GetName() string {
	return log.name
}

// GetCallerDepth return the logger's depth
func (log *Log) GetCallerDepth() int {
	return log.depth
}

// SetCallerDepth set the logger's caller depth (!!SLOW!!), 0: disable runtime.Caller()
func (log *Log) SetCallerDepth(d int) {
	log.depth = d
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
func (log *Log) GetProp(k string) any {
	return log.props
}

// SetProp set logger property
func (log *Log) SetProp(k string, v any) {
	// copy on write for async
	log.props = cloneAndSetProp(log.props, k, v)
}

// GetProps get logger properties
func (log *Log) GetProps() map[string]any {
	return log.props
}

// SetProps set logger properties
func (log *Log) SetProps(props map[string]any) {
	if props == nil {
		props = emptyProps
	}
	log.props = props
}

// IsLevelEnabled is specified level enabled
func (log *Log) IsLevelEnabled(lvl Level) bool {
	return log.level >= lvl
}

// Log log a message at specified level.
func (log *Log) Log(lvl Level, v ...any) {
	log._log(lvl, v...)
}

// Logf format and log a message at specified level.
func (log *Log) Logf(lvl Level, f string, v ...any) {
	log._logf(lvl, f, v...)
}

// IsFatalEnabled is FATAL level enabled
func (log *Log) IsFatalEnabled() bool {
	return log.IsLevelEnabled(LevelFatal)
}

// Fatal log a message at fatal level.
func (log *Log) Fatal(v ...any) {
	log._log(LevelFatal, v...)
}

// Fatalf format and log a message at fatal level.
func (log *Log) Fatalf(f string, v ...any) {
	log._logf(LevelFatal, f, v...)
}

// IsErrorEnabled is ERROR level enabled
func (log *Log) IsErrorEnabled() bool {
	return log.IsLevelEnabled(LevelError)
}

// Error log a message at error level.
func (log *Log) Error(v ...any) {
	log._log(LevelError, v...)
}

// Errorf format and log a message at error level.
func (log *Log) Errorf(f string, v ...any) {
	log._logf(LevelError, f, v...)
}

// IsWarnEnabled is WARN level enabled
func (log *Log) IsWarnEnabled() bool {
	return log.IsLevelEnabled(LevelWarn)
}

// Warn log a message at warning level.
func (log *Log) Warn(v ...any) {
	log._log(LevelWarn, v...)
}

// Warnf format and log a message at warning level.
func (log *Log) Warnf(f string, v ...any) {
	log._logf(LevelWarn, f, v...)
}

// IsInfoEnabled is INFO level enabled
func (log *Log) IsInfoEnabled() bool {
	return log.IsLevelEnabled(LevelInfo)
}

// Info log a message at info level.
func (log *Log) Info(v ...any) {
	log._log(LevelInfo, v...)
}

// Infof format and log a message at info level.
func (log *Log) Infof(f string, v ...any) {
	log._logf(LevelInfo, f, v...)
}

// IsDebugEnabled is DEBUG level enabled
func (log *Log) IsDebugEnabled() bool {
	return log.IsLevelEnabled(LevelDebug)
}

// Debug log a message at debug level.
func (log *Log) Debug(v ...any) {
	log._log(LevelDebug, v...)
}

// Debugf format log a message at debug level.
func (log *Log) Debugf(f string, v ...any) {
	log._logf(LevelDebug, f, v...)
}

// IsTraceEnabled is TRACE level enabled
func (log *Log) IsTraceEnabled() bool {
	return log.IsLevelEnabled(LevelTrace)
}

// Trace log a message at trace level.
func (log *Log) Trace(v ...any) {
	log._log(LevelTrace, v...)
}

// Tracef format and log a message at trace level.
func (log *Log) Tracef(f string, v ...any) {
	log._logf(LevelTrace, f, v...)
}

// Write write a log event
func (log *Log) Write(le *Event) {
	if log.IsLevelEnabled(le.Level) {
		log._write(le)
	}
}

func (log *Log) _log(lvl Level, v ...any) {
	if log.IsLevelEnabled(lvl) {
		s := _printv(v...)
		le := NewEvent(log, lvl, s)
		log._write(le)
	}
}

func (log *Log) _logf(lvl Level, f string, v ...any) {
	if log.IsLevelEnabled(lvl) {
		s := _printf(f, v...)
		le := NewEvent(log, lvl, s)
		log._write(le)
	}
}

func (log *Log) _write(le *Event) {
	safeWrite(log.writer, le)
}

func _printv(v ...any) string {
	if len(v) == 0 {
		return ""
	}
	return fmt.Sprint(v...)
}

func _printf(f string, v ...any) string {
	return fmt.Sprintf(f, v...)
}

func cloneAndSetProp(om map[string]any, k string, v any) map[string]any {
	// copy on write for async
	nm := make(map[string]any, len(om)+1)
	mag.Copy(nm, om)
	nm[k] = v

	return nm
}
