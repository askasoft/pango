package log

import (
	"fmt"
)

// Logger logger interface
type Logger interface {
	GetName() string
	GetLevel() int
	SetLevel(l int)
	GetCallerDepth() int
	SetCallerDepth(d int)
	GetTraceLevel() int
	SetTraceLevel(t int)
	GetProp(k string) interface{}
	SetProp(k string, v interface{})
	GetProps() map[string]interface{}
	SetProps(map[string]interface{})
	GetFormatter() Formatter
	SetFormatter(lf Formatter)
	IsAsync() bool
	Lock()
	Unlock()
	IsLevelEnabled(lvl int) bool
	Log(lvl int, v ...interface{})
	Logf(lvl int, f string, v ...interface{})
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
	props  map[string]interface{}
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

// SetCallerDepth set the logger's caller depth (!!SLOW!!), 0: disable runtime.Caller()
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

// GetProp get logger property
func (l *logger) GetProp(k string) interface{} {
	if l.props == nil {
		return nil
	}
	return l.props[k]
}

// SetProp set logger property
func (l *logger) SetProp(k string, v interface{}) {
	if l.props == nil {
		l.props = make(map[string]interface{})
	}
	l.props[k] = v
}

// GetProps get logger properties
func (l *logger) GetProps() map[string]interface{} {
	return l.props
}

// SetProps set logger properties
func (l *logger) SetProps(props map[string]interface{}) {
	l.props = props
}

// GetFormatter get logger formatter
func (l *logger) GetFormatter() Formatter {
	return l.logfmt
}

// SetFormatter set logger formatter
func (l *logger) SetFormatter(lf Formatter) {
	l.logfmt = lf
}

// IsAsync return the logger's async
func (l *logger) IsAsync() bool {
	return l.log.async
}

// Lock lock the log
func (l *logger) Lock() {
	l.log.Lock()
}

// Unlock unlock the log
func (l *logger) Unlock() {
	l.log.Unlock()
}

// IsLevelEnabled is specified level enabled
func (l *logger) IsLevelEnabled(lvl int) bool {
	return l.level >= lvl
}

// Log log a message at specified level.
func (l *logger) Log(lvl int, v ...interface{}) {
	if l.IsLevelEnabled(LevelFatal) {
		s := printv(v...)
		le := newEvent(l, lvl, s)
		l.log.submit(le)
	}
}

// Logf format and log a message at specified level.
func (l *logger) Logf(lvl int, f string, v ...interface{}) {
	if l.IsLevelEnabled(lvl) {
		s := printf(f, v...)
		le := newEvent(l, lvl, s)
		l.log.submit(le)
	}
}

// IsFatalEnabled is FATAL level enabled
func (l *logger) IsFatalEnabled() bool {
	return l.IsLevelEnabled(LevelFatal)
}

// Fatal log a message at fatal level.
func (l *logger) Fatal(v ...interface{}) {
	l.Log(LevelFatal, v...)
}

// Fatalf format and log a message at fatal level.
func (l *logger) Fatalf(f string, v ...interface{}) {
	l.Logf(LevelFatal, f, v...)
}

// IsErrorEnabled is ERROR level enabled
func (l *logger) IsErrorEnabled() bool {
	return l.IsLevelEnabled(LevelError)
}

// Error log a message at error level.
func (l *logger) Error(v ...interface{}) {
	l.Log(LevelError, v...)
}

// Errorf format and log a message at error level.
func (l *logger) Errorf(f string, v ...interface{}) {
	l.Logf(LevelError, f, v...)
}

// IsWarnEnabled is WARN level enabled
func (l *logger) IsWarnEnabled() bool {
	return l.IsLevelEnabled(LevelWarn)
}

// Warn log a message at warning level.
func (l *logger) Warn(v ...interface{}) {
	l.Log(LevelWarn, v...)
}

// Warnf format and log a message at warning level.
func (l *logger) Warnf(f string, v ...interface{}) {
	l.Logf(LevelWarn, f, v...)
}

// IsInfoEnabled is INFO level enabled
func (l *logger) IsInfoEnabled() bool {
	return l.IsLevelEnabled(LevelInfo)
}

// Info log a message at info level.
func (l *logger) Info(v ...interface{}) {
	l.Log(LevelInfo, v...)
}

// Infof format and log a message at info level.
func (l *logger) Infof(f string, v ...interface{}) {
	l.Logf(LevelInfo, f, v...)
}

// IsDebugEnabled is DEBUG level enabled
func (l *logger) IsDebugEnabled() bool {
	return l.IsLevelEnabled(LevelDebug)
}

// Debug log a message at debug level.
func (l *logger) Debug(v ...interface{}) {
	l.Log(LevelDebug, v...)
}

// Debugf format log a message at debug level.
func (l *logger) Debugf(f string, v ...interface{}) {
	l.Logf(LevelDebug, f, v...)
}

// IsTraceEnabled is TRACE level enabled
func (l *logger) IsTraceEnabled() bool {
	return l.IsLevelEnabled(LevelTrace)
}

// Trace log a message at trace level.
func (l *logger) Trace(v ...interface{}) {
	l.Log(LevelTrace, v...)
}

// Tracef format and log a message at trace level.
func (l *logger) Tracef(f string, v ...interface{}) {
	l.Logf(LevelTrace, f, v...)
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
