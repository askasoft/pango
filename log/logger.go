package log

import (
	"github.com/askasoft/pango/str"
)

// Logger logger interface
type Logger interface {
	GetLogger(name string) Logger
	GetOutputer(name string, lvl Level, callerDepth ...int) Outputer
	GetName() string
	GetLevel() Level
	GetTraceLevel() Level
	GetCallerDepth() int
	SetCallerDepth(d int)
	GetProp(k string) any
	SetProp(k string, v any)
	GetProps() map[string]any
	SetProps(map[string]any)
	IsLevelEnabled(lvl Level) bool
	Log(lvl Level, v ...any)
	Logf(lvl Level, f string, v ...any)
	IsFatalEnabled() bool
	Fatal(v ...any)
	Fatalf(f string, v ...any)
	IsErrorEnabled() bool
	Error(v ...any)
	Errorf(f string, v ...any)
	IsWarnEnabled() bool
	Warn(v ...any)
	Warnf(f string, v ...any)
	IsInfoEnabled() bool
	Info(v ...any)
	Infof(f string, v ...any)
	IsDebugEnabled() bool
	Debug(v ...any)
	Debugf(f string, v ...any)
	IsTraceEnabled() bool
	Trace(v ...any)
	Tracef(f string, v ...any)
	Write(le *Event)
}

// logger logger interface implement
type logger struct {
	log   *Log
	name  string
	depth int
	props map[string]any
}

// GetLogger create a new logger with name
func (lg *logger) GetLogger(name string) Logger {
	return &logger{
		log:   lg.log,
		name:  str.IfEmpty(name, "_"),
		depth: lg.depth,
		props: lg.props,
	}
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
func (lg *logger) GetOutputer(name string, lvl Level, callerDepth ...int) Outputer {
	return lg.log.GetOutputer(name, lvl, callerDepth...)
}

// GetName return the logger's name
func (lg *logger) GetName() string {
	return lg.name
}

// GetCallerDepth return the logger's depth
func (lg *logger) GetCallerDepth() int {
	return lg.depth
}

// SetCallerDepth set the logger's caller depth (!!SLOW!!), 0: disable runtime.Caller()
func (lg *logger) SetCallerDepth(d int) {
	lg.depth = d
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

// Log log a message at specified level.
func (lg *logger) Log(lvl Level, v ...any) {
	lg._log(lvl, v...)
}

// Logf format and log a message at specified level.
func (lg *logger) Logf(lvl Level, f string, v ...any) {
	lg._logf(lvl, f, v...)
}

// IsFatalEnabled is FATAL level enabled
func (lg *logger) IsFatalEnabled() bool {
	return lg.IsLevelEnabled(LevelFatal)
}

// Fatal log a message at fatal level.
func (lg *logger) Fatal(v ...any) {
	lg._log(LevelFatal, v...)
}

// Fatalf format and log a message at fatal level.
func (lg *logger) Fatalf(f string, v ...any) {
	lg._logf(LevelFatal, f, v...)
}

// IsErrorEnabled is ERROR level enabled
func (lg *logger) IsErrorEnabled() bool {
	return lg.IsLevelEnabled(LevelError)
}

// Error log a message at error level.
func (lg *logger) Error(v ...any) {
	lg._log(LevelError, v...)
}

// Errorf format and log a message at error level.
func (lg *logger) Errorf(f string, v ...any) {
	lg._logf(LevelError, f, v...)
}

// IsWarnEnabled is WARN level enabled
func (lg *logger) IsWarnEnabled() bool {
	return lg.IsLevelEnabled(LevelWarn)
}

// Warn log a message at warning level.
func (lg *logger) Warn(v ...any) {
	lg._log(LevelWarn, v...)
}

// Warnf format and log a message at warning level.
func (lg *logger) Warnf(f string, v ...any) {
	lg._logf(LevelWarn, f, v...)
}

// IsInfoEnabled is INFO level enabled
func (lg *logger) IsInfoEnabled() bool {
	return lg.IsLevelEnabled(LevelInfo)
}

// Info log a message at info level.
func (lg *logger) Info(v ...any) {
	lg._log(LevelInfo, v...)
}

// Infof format and log a message at info level.
func (lg *logger) Infof(f string, v ...any) {
	lg._logf(LevelInfo, f, v...)
}

// IsDebugEnabled is DEBUG level enabled
func (lg *logger) IsDebugEnabled() bool {
	return lg.IsLevelEnabled(LevelDebug)
}

// Debug log a message at debug level.
func (lg *logger) Debug(v ...any) {
	lg._log(LevelDebug, v...)
}

// Debugf format log a message at debug level.
func (lg *logger) Debugf(f string, v ...any) {
	lg._logf(LevelDebug, f, v...)
}

// IsTraceEnabled is TRACE level enabled
func (lg *logger) IsTraceEnabled() bool {
	return lg.IsLevelEnabled(LevelTrace)
}

// Trace log a message at trace level.
func (lg *logger) Trace(v ...any) {
	lg._log(LevelTrace, v...)
}

// Tracef format and log a message at trace level.
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
		le := NewEvent(lg, lvl, s)
		lg.log._write(le)
	}
}

func (lg *logger) _logf(lvl Level, f string, v ...any) {
	if lg.IsLevelEnabled(lvl) {
		s := _printf(f, v...)
		le := NewEvent(lg, lvl, s)
		lg.log._write(le)
	}
}
