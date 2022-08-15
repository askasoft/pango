package log

import (
	"fmt"
)

// Logger logger interface
type Logger interface {
	GetLogger(name string) Logger
	GetOutputer(name string, lvl Level, callerDepth ...int) Outputer
	GetName() string
	SetName(name string)
	GetLevel() Level
	GetTraceLevel() Level
	GetFormatter() Formatter
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
}

// logger logger interface implement
type logger struct {
	log   *Log
	name  string
	depth int
	props map[string]any
}

// GetLogger create a new logger with name
func (l *logger) GetLogger(name string) Logger {
	return l.log.GetLogger(name)
}

// Outputer return a io.Writer for go log.SetOutput
// callerDepth: default is 1 (means +1)
// if the outputer is used by go std log, set callerDepth to 2
// example:
//
//	import (
//	  golog "log"
//	  "github.com/pandafw/pango/log"
//	)
//	golog.SetOutput(log.Outputer("GO", log.LevelInfo, 3))
func (l *logger) GetOutputer(name string, lvl Level, callerDepth ...int) Outputer {
	return l.log.GetOutputer(name, lvl, callerDepth...)
}

// GetName return the logger's name
func (l *logger) GetName() string {
	return l.name
}

// SetName set the logger's name
func (l *logger) SetName(name string) {
	l.name = name
}

// GetCallerDepth return the logger's depth
func (l *logger) GetCallerDepth() int {
	return l.depth
}

// SetCallerDepth set the logger's caller depth (!!SLOW!!), 0: disable runtime.Caller()
func (l *logger) SetCallerDepth(d int) {
	l.depth = d
}

// GetLevel return the logger's level
func (l *logger) GetLevel() Level {
	lvl := l.log.getLoggerLevel(l.name)
	if lvl == LevelNone {
		return l.log.GetLevel()
	}
	return lvl
}

// GetTraceLevel return the logger's trace level
func (l *logger) GetTraceLevel() Level {
	return l.log.GetTraceLevel()
}

// GetProp get logger property
func (l *logger) GetProp(k string) any {
	ps := l.props
	if ps != nil {
		if v, ok := ps[k]; ok {
			return v
		}
	}
	return l.log.GetProp(k)
}

// SetProp set logger property
func (l *logger) SetProp(k string, v any) {
	// copy on write for async
	om := l.props

	nm := make(map[string]any)
	for k, v := range om {
		nm[k] = v
	}
	nm[k] = v

	l.props = nm
}

// GetProps get logger properties
func (l *logger) GetProps() map[string]any {
	tm := l.props
	if tm == nil {
		return l.log.GetProps()
	}

	// props
	var nm map[string]any

	// parent props
	pm := l.log.logger.props
	if pm != nil {
		// new return props
		nm = make(map[string]any, len(tm)+len(pm))
		for k, v := range pm {
			nm[k] = v
		}
	}

	// self props
	if nm == nil {
		nm = make(map[string]any, len(tm))
	}
	for k, v := range tm {
		nm[k] = v
	}
	return nm
}

// SetProps set logger properties
func (l *logger) SetProps(props map[string]any) {
	l.props = props
}

// GetFormatter get logger formatter
func (l *logger) GetFormatter() Formatter {
	return l.log.GetFormatter()
}

// IsLevelEnabled is specified level enabled
func (l *logger) IsLevelEnabled(lvl Level) bool {
	return l.GetLevel() >= lvl
}

// Log log a message at specified level.
func (l *logger) Log(lvl Level, v ...any) {
	l._log(lvl, v...)
}

// Logf format and log a message at specified level.
func (l *logger) Logf(lvl Level, f string, v ...any) {
	l._logf(lvl, f, v...)
}

// IsFatalEnabled is FATAL level enabled
func (l *logger) IsFatalEnabled() bool {
	return l.IsLevelEnabled(LevelFatal)
}

// Fatal log a message at fatal level.
func (l *logger) Fatal(v ...any) {
	l._log(LevelFatal, v...)
}

// Fatalf format and log a message at fatal level.
func (l *logger) Fatalf(f string, v ...any) {
	l._logf(LevelFatal, f, v...)
}

// IsErrorEnabled is ERROR level enabled
func (l *logger) IsErrorEnabled() bool {
	return l.IsLevelEnabled(LevelError)
}

// Error log a message at error level.
func (l *logger) Error(v ...any) {
	l._log(LevelError, v...)
}

// Errorf format and log a message at error level.
func (l *logger) Errorf(f string, v ...any) {
	l._logf(LevelError, f, v...)
}

// IsWarnEnabled is WARN level enabled
func (l *logger) IsWarnEnabled() bool {
	return l.IsLevelEnabled(LevelWarn)
}

// Warn log a message at warning level.
func (l *logger) Warn(v ...any) {
	l._log(LevelWarn, v...)
}

// Warnf format and log a message at warning level.
func (l *logger) Warnf(f string, v ...any) {
	l._logf(LevelWarn, f, v...)
}

// IsInfoEnabled is INFO level enabled
func (l *logger) IsInfoEnabled() bool {
	return l.IsLevelEnabled(LevelInfo)
}

// Info log a message at info level.
func (l *logger) Info(v ...any) {
	l._log(LevelInfo, v...)
}

// Infof format and log a message at info level.
func (l *logger) Infof(f string, v ...any) {
	l._logf(LevelInfo, f, v...)
}

// IsDebugEnabled is DEBUG level enabled
func (l *logger) IsDebugEnabled() bool {
	return l.IsLevelEnabled(LevelDebug)
}

// Debug log a message at debug level.
func (l *logger) Debug(v ...any) {
	l._log(LevelDebug, v...)
}

// Debugf format log a message at debug level.
func (l *logger) Debugf(f string, v ...any) {
	l._logf(LevelDebug, f, v...)
}

// IsTraceEnabled is TRACE level enabled
func (l *logger) IsTraceEnabled() bool {
	return l.IsLevelEnabled(LevelTrace)
}

// Trace log a message at trace level.
func (l *logger) Trace(v ...any) {
	l._log(LevelTrace, v...)
}

// Tracef format and log a message at trace level.
func (l *logger) Tracef(f string, v ...any) {
	l._logf(LevelTrace, f, v...)
}

func (l *logger) _log(lvl Level, v ...any) {
	if l.IsLevelEnabled(lvl) {
		s := l._printv(v...)
		le := newEvent(l, lvl, s)
		l.log.write(le)
	}
}

func (l *logger) _logf(lvl Level, f string, v ...any) {
	if l.IsLevelEnabled(lvl) {
		s := l._printf(f, v...)
		le := newEvent(l, lvl, s)
		l.log.write(le)
	}
}

func (l *logger) _printv(v ...any) string {
	if len(v) == 0 {
		return ""
	}
	return fmt.Sprint(v...)
}

func (l *logger) _printf(f string, v ...any) string {
	return fmt.Sprintf(f, v...)
}
