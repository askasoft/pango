package log

import (
	"fmt"
)

// Logger logger interface
type Logger interface {
	GetLevel() int
	SetLevel(l int)
	GetCallerDepth() int
	SetCallerDepth(d int)
	GetTraceLevel() int
	SetTraceLevel(t int)
	GetName() string
	SetName(n string)
	GetParam(k string) interface{}
	SetParam(k string, v interface{})
	GetFormatter() Formatter
	SetFormatter(lf Formatter)
	IsAsync() bool
	IsFatalEnabled() bool
	Fatal(v ...interface{})
	Fatalf(f string, v ...interface{})
	IsErrorEnabled() bool
	Error(v ...interface{})
	Errorf(f string, v ...interface{})
	IsWarnEnabled() bool
	Warn(v ...interface{})
	Warnf(f string, v ...interface{})
	IsInfoEnabled() bool
	Info(v ...interface{})
	Infof(f string, v ...interface{})
	IsDebugEnabled() bool
	Debug(v ...interface{})
	Debugf(f string, v ...interface{})
	IsTraceEnabled() bool
	Trace(v ...interface{})
	Tracef(f string, v ...interface{})
}

// logger logger interface implement
type logger struct {
	name   string
	level  int
	depth  int
	trace  int
	log    *Log
	logfmt Formatter
	vars   map[string]interface{}
}

// GetName return the logger's name
func (l *logger) GetName() string {
	return l.name
}

// SetName set the logger's name
func (l *logger) SetName(n string) {
	l.name = n
}

// GetLevel return the logger's level
func (l *logger) GetLevel() int {
	return l.level
}

// SetLevel set the logger's level
func (l *logger) SetLevel(lvl int) {
	l.level = lvl
}

// GetCallerDepth return the logger's depth
func (l *logger) GetCallerDepth() int {
	return l.depth
}

// SetCallerDepth set the logger's depth, 0: disable runtime.Caller()
func (l *logger) SetCallerDepth(d int) {
	l.depth = d
}

// GetTraceLevel return the logger's trace level
func (l *logger) GetTraceLevel() int {
	return l.trace
}

// SetTraceLevel set the logger's trace level
func (l *logger) SetTraceLevel(lvl int) {
	l.trace = lvl
}

// GetParam get logger k's value
func (l *logger) GetParam(k string) interface{} {
	if l.vars == nil {
		return nil
	}
	return l.vars[k]
}

// SetParam set logger k's value
func (l *logger) SetParam(k string, v interface{}) {
	if l.vars == nil {
		l.vars = make(map[string]interface{})
	}
	l.vars[k] = v
}

// GetFormatter get logger formatter
func (l *logger) GetFormatter() Formatter {
	return l.logfmt
}

// SetFormatter set logger formatter
func (l *logger) SetFormatter(lf Formatter) {
	l.logfmt = lf
}

// IsFatalEnabled is FATAL level enabled
func (l *logger) IsFatalEnabled() bool {
	return l.level >= LevelFatal
}

// IsAsync return the logger's async
func (l *logger) IsAsync() bool {
	return l.log.async
}

// Fatal log a message at fatal level.
func (l *logger) Fatal(v ...interface{}) {
	if l.IsFatalEnabled() {
		s := printv(v...)
		le := newEvent(l, LevelFatal, s)
		l.log.Log(le)
	}
}

// Fatalf format and log a message at fatal level.
func (l *logger) Fatalf(f string, v ...interface{}) {
	if l.IsFatalEnabled() {
		s := printf(f, v...)
		le := newEvent(l, LevelFatal, s)
		l.log.Log(le)
	}
}

// IsErrorEnabled is ERROR level enabled
func (l *logger) IsErrorEnabled() bool {
	return l.level >= LevelError
}

// Error log a message at error level.
func (l *logger) Error(v ...interface{}) {
	if l.IsErrorEnabled() {
		s := printv(v...)
		le := newEvent(l, LevelError, s)
		l.log.Log(le)
	}
}

// Errorf format and log a message at error level.
func (l *logger) Errorf(f string, v ...interface{}) {
	if l.IsErrorEnabled() {
		s := printf(f, v...)
		le := newEvent(l, LevelError, s)
		l.log.Log(le)
	}
}

// IsWarnEnabled is WARN level enabled
func (l *logger) IsWarnEnabled() bool {
	return l.level >= LevelWarn
}

// Warn log a message at warning level.
func (l *logger) Warn(v ...interface{}) {
	if l.IsWarnEnabled() {
		s := printv(v...)
		le := newEvent(l, LevelWarn, s)
		l.log.Log(le)
	}
}

// Warnf format and log a message at warning level.
func (l *logger) Warnf(f string, v ...interface{}) {
	if l.IsWarnEnabled() {
		s := printf(f, v...)
		le := newEvent(l, LevelWarn, s)
		l.log.Log(le)
	}
}

// IsInfoEnabled is INFO level enabled
func (l *logger) IsInfoEnabled() bool {
	return l.level >= LevelInfo
}

// Info log a message at info level.
func (l *logger) Info(v ...interface{}) {
	if l.IsInfoEnabled() {
		s := printv(v...)
		le := newEvent(l, LevelInfo, s)
		l.log.Log(le)
	}
}

// Infof format and log a message at info level.
func (l *logger) Infof(f string, v ...interface{}) {
	if l.IsInfoEnabled() {
		s := printf(f, v...)
		le := newEvent(l, LevelInfo, s)
		l.log.Log(le)
	}
}

// IsDebugEnabled is DEBUG level enabled
func (l *logger) IsDebugEnabled() bool {
	return l.level >= LevelDebug
}

// Debug log a message at debug level.
func (l *logger) Debug(v ...interface{}) {
	if l.IsDebugEnabled() {
		s := printv(v...)
		le := newEvent(l, LevelDebug, s)
		l.log.Log(le)
	}
}

// Debugf format log a message at debug level.
func (l *logger) Debugf(f string, v ...interface{}) {
	if l.IsDebugEnabled() {
		s := printf(f, v...)
		le := newEvent(l, LevelDebug, s)
		l.log.Log(le)
	}
}

// IsTraceEnabled is TRACE level enabled
func (l *logger) IsTraceEnabled() bool {
	return l.level >= LevelTrace
}

// Trace log a message at trace level.
func (l *logger) Trace(v ...interface{}) {
	if l.IsTraceEnabled() {
		s := printv(v...)
		le := newEvent(l, LevelTrace, s)
		l.log.Log(le)
	}
}

// Tracef format and log a message at trace level.
func (l *logger) Tracef(f string, v ...interface{}) {
	if l.IsTraceEnabled() {
		s := printf(f, v...)
		le := newEvent(l, LevelTrace, s)
		l.log.Log(le)
	}
}

func printv(v ...interface{}) string {
	if len(v) == 0 {
		return ""
	}
	return fmt.Sprint(v...)
}

func printf(f string, v ...interface{}) string {
	return fmt.Sprintf(f, v...)
}
