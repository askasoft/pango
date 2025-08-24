package log

import (
	"os"

	"github.com/askasoft/pango/str"
)

// Logger logger interface
type Logger interface {
	// GetLogger create a new logger with name
	GetLogger(name string) Logger

	// Outputer return a io.Writer for go log.SetOutput
	// callerSkip: default is 1 (means +1)
	// if the outputer is used by go std log, set callerSkip to 2
	// example:
	//
	//	import (
	//	  golog "log"
	//	  "github.com/askasoft/pango/log"
	//	)
	//	golog.SetOutput(log.Outputer("GO", log.LevelInfo, 3))
	GetOutputer(name string, lvl Level, callerSkip ...int) Outputer

	// GetName return the logger's name
	GetName() string

	// GetLevel return the logger's level
	GetLevel() Level

	// GetTraceLevel return the logger's trace level
	GetTraceLevel() Level

	// GetCallerSkip return the logger's caller skip
	GetCallerSkip() int

	// SetCallerSkip set the logger's caller skip (!!SLOW!!), 0: disable runtime.Caller()
	SetCallerSkip(n int)

	// GetProp get logger property
	GetProp(k string) any

	// SetProp set logger property
	SetProp(k string, v any)

	// GetProps get logger properties
	GetProps() map[string]any

	// SetProps set logger properties
	SetProps(map[string]any)

	// IsLevelEnabled is specified level enabled
	IsLevelEnabled(lvl Level) bool

	// IsFatalEnabled is FATAL level enabled
	IsFatalEnabled() bool

	// Fatal print a fatal message, close the logs and call [os.Exit](code)
	Fatal(code int, v ...any)

	// Fatalf format and print a fatal message, close the logs and call [os.Exit](code)
	Fatalf(code int, f string, v ...any)

	// IsErrorEnabled is ERROR level enabled
	IsErrorEnabled() bool

	// Error print a message at error level.
	Error(v ...any)

	// Errorf format and print a message at error level.
	Errorf(f string, v ...any)

	// IsWarnEnabled is WARN level enabled
	IsWarnEnabled() bool

	// Warn print a message at warning level.
	Warn(v ...any)

	// Warnf format and print a message at warning level.
	Warnf(f string, v ...any)

	// IsInfoEnabled is INFO level enabled
	IsInfoEnabled() bool

	// Info print a message at info level.
	Info(v ...any)

	// Infof format and print a message at info level.
	Infof(f string, v ...any)

	// IsDebugEnabled is DEBUG level enabled
	IsDebugEnabled() bool

	// Debug print a message at debug level.
	Debug(v ...any)

	// Debugf format and print a message at debug level.
	Debugf(f string, v ...any)

	// IsTraceEnabled is TRACE level enabled
	IsTraceEnabled() bool

	// Trace print a message at trace level.
	Trace(v ...any)

	// Tracef format and print a message at trace level.
	Tracef(f string, v ...any)

	// Log print a message at specified level.
	Log(lvl Level, v ...any)

	// Logf format and print a message at specified level.
	Logf(lvl Level, f string, v ...any)

	// Write write a log event
	Write(le *Event)
}

// logger logger interface implement
type logger struct {
	log   *Log
	name  string
	skip  int
	props map[string]any
}

// GetLogger create a new logger with name
func (lg *logger) GetLogger(name string) Logger {
	return &logger{
		log:   lg.log,
		name:  str.IfEmpty(name, "_"),
		skip:  lg.skip,
		props: lg.props,
	}
}

// Outputer return a io.Writer for go log.SetOutput
// callerSkip: default is 1 (means +1)
// if the outputer is used by go std log, set callerSkip to 2
// example:
//
//	import (
//	  golog "log"
//	  "github.com/askasoft/pango/log"
//	)
//	golog.SetOutput(log.Outputer("GO", log.LevelInfo, 3))
func (lg *logger) GetOutputer(name string, lvl Level, callerSkip ...int) Outputer {
	return lg.log.GetOutputer(name, lvl, callerSkip...)
}

// GetName return the logger's name
func (lg *logger) GetName() string {
	return lg.name
}

// GetCallerSkip return the logger's caller skip
func (lg *logger) GetCallerSkip() int {
	return lg.skip
}

// SetCallerSkip set the logger's caller skip (!!SLOW!!), 0: disable runtime.Caller()
func (lg *logger) SetCallerSkip(n int) {
	lg.skip = n
}

// GetLevel return the logger's level
func (lg *logger) GetLevel() Level {
	lvl := lg.log.GetLoggerLevel(lg.name)
	if lvl == LevelNone {
		return lg.log.GetLevel()
	}
	return lvl
}

// GetTraceLevel return the logger's trace level
func (lg *logger) GetTraceLevel() Level {
	return lg.log.GetTraceLevel()
}

// GetProp get logger property
func (lg *logger) GetProp(k string) any {
	return lg.props[k]
}

// SetProp set logger property
func (lg *logger) SetProp(k string, v any) {
	// copy on write for async
	lg.props = cloneAndSetProp(lg.props, k, v)
}

// GetProps get logger properties
func (lg *logger) GetProps() map[string]any {
	return lg.props
}

// SetProps set logger properties
func (lg *logger) SetProps(props map[string]any) {
	if props == nil {
		props = lg.log.props
	}
	lg.props = props
}

// IsLevelEnabled is specified level enabled
func (lg *logger) IsLevelEnabled(lvl Level) bool {
	return lg.GetLevel() >= lvl
}

// Log print a message at specified level.
func (lg *logger) Log(lvl Level, v ...any) {
	lg._log(lvl, v...)
}

// Logf format and print a message at specified level.
func (lg *logger) Logf(lvl Level, f string, v ...any) {
	lg._logf(lvl, f, v...)
}

// IsFatalEnabled is FATAL level enabled
func (lg *logger) IsFatalEnabled() bool {
	return lg.IsLevelEnabled(LevelFatal)
}

// Fatal print a message at fatal level, close the logs and call [os.Exit](code).
func (lg *logger) Fatal(code int, v ...any) {
	lg._log(LevelFatal, v...)
	lg.log.Close()
	os.Exit(code)
}

// Fatalf format and print a message at fatal level, close the logs and call [os.Exit](code).
func (lg *logger) Fatalf(code int, f string, v ...any) {
	lg._logf(LevelFatal, f, v...)
	lg.log.Close()
	os.Exit(code)
}

// IsErrorEnabled is ERROR level enabled
func (lg *logger) IsErrorEnabled() bool {
	return lg.IsLevelEnabled(LevelError)
}

// Error print a message at error level.
func (lg *logger) Error(v ...any) {
	lg._log(LevelError, v...)
}

// Errorf format and print a message at error level.
func (lg *logger) Errorf(f string, v ...any) {
	lg._logf(LevelError, f, v...)
}

// IsWarnEnabled is WARN level enabled
func (lg *logger) IsWarnEnabled() bool {
	return lg.IsLevelEnabled(LevelWarn)
}

// Warn print a message at warning level.
func (lg *logger) Warn(v ...any) {
	lg._log(LevelWarn, v...)
}

// Warnf format and print a message at warning level.
func (lg *logger) Warnf(f string, v ...any) {
	lg._logf(LevelWarn, f, v...)
}

// IsInfoEnabled is INFO level enabled
func (lg *logger) IsInfoEnabled() bool {
	return lg.IsLevelEnabled(LevelInfo)
}

// Info print a message at info level.
func (lg *logger) Info(v ...any) {
	lg._log(LevelInfo, v...)
}

// Infof format and print a message at info level.
func (lg *logger) Infof(f string, v ...any) {
	lg._logf(LevelInfo, f, v...)
}

// IsDebugEnabled is DEBUG level enabled
func (lg *logger) IsDebugEnabled() bool {
	return lg.IsLevelEnabled(LevelDebug)
}

// Debug print a message at debug level.
func (lg *logger) Debug(v ...any) {
	lg._log(LevelDebug, v...)
}

// Debugf format print a message at debug level.
func (lg *logger) Debugf(f string, v ...any) {
	lg._logf(LevelDebug, f, v...)
}

// IsTraceEnabled is TRACE level enabled
func (lg *logger) IsTraceEnabled() bool {
	return lg.IsLevelEnabled(LevelTrace)
}

// Trace print a message at trace level.
func (lg *logger) Trace(v ...any) {
	lg._log(LevelTrace, v...)
}

// Tracef format and print a message at trace level.
func (lg *logger) Tracef(f string, v ...any) {
	lg._logf(LevelTrace, f, v...)
}

// Write write a log event
func (lg *logger) Write(le *Event) {
	if lg.IsLevelEnabled(le.Level) {
		lg.log._write(le)
	}
}

func (lg *logger) _log(lvl Level, v ...any) {
	if lg.IsLevelEnabled(lvl) {
		s := _printv(v...)
		le := newLogEvent(lg, lvl, s)
		lg.log._write(le)
	}
}

func (lg *logger) _logf(lvl Level, f string, v ...any) {
	if lg.IsLevelEnabled(lvl) {
		s := _printf(f, v...)
		le := newLogEvent(lg, lvl, s)
		lg.log._write(le)
	}
}
