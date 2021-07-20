package log

import (
	"fmt"
)

// Logger logger interface
type Logger interface {
	GetName() string
	GetLevel() Level
	GetTraceLevel() Level
	GetFormatter() Formatter
	GetCallerDepth() int
	SetCallerDepth(d int)
	GetProp(k string) interface{}
	SetProp(k string, v interface{})
	GetProps() map[string]interface{}
	SetProps(map[string]interface{})
	IsLevelEnabled(lvl Level) bool
	Log(lvl Level, v ...interface{})
	Logf(lvl Level, f string, v ...interface{})
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
	log   *Log
	name  string
	depth int
	props map[string]interface{}
}

// GetName return the logger's name
func (l *logger) GetName() string {
	return l.name
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
func (l *logger) GetProp(k string) interface{} {
	ps := l.props
	if ps != nil {
		if v, ok := ps[k]; ok {
			return v
		}
	}
	return l.log.GetProp(k)
}

// SetProp set logger property
func (l *logger) SetProp(k string, v interface{}) {
	// copy on write for async
	om := l.props
	nm := make(map[string]interface{})

	if om != nil {
		for k, v := range om {
			nm[k] = v
		}
	}
	nm[k] = v
	l.props = nm
}

// GetProps get logger properties
func (l *logger) GetProps() map[string]interface{} {
	tm := l.props
	if tm == nil {
		return l.log.GetProps()
	}

	// props
	var nm map[string]interface{}

	// parent props
	pm := l.log.logger.props
	if pm != nil {
		// new return props
		nm = make(map[string]interface{}, len(tm)+len(pm))
		if pm != nil {
			for k, v := range pm {
				nm[k] = v
			}
		}
	}

	// self props
	if nm == nil {
		nm = make(map[string]interface{}, len(tm))
	}
	for k, v := range tm {
		nm[k] = v
	}
	return nm
}

// SetProps set logger properties
func (l *logger) SetProps(props map[string]interface{}) {
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
func (l *logger) Log(lvl Level, v ...interface{}) {
	l._log(lvl, v...)
}

// Logf format and log a message at specified level.
func (l *logger) Logf(lvl Level, f string, v ...interface{}) {
	l._logf(lvl, f, v...)
}

// IsFatalEnabled is FATAL level enabled
func (l *logger) IsFatalEnabled() bool {
	return l.IsLevelEnabled(LevelFatal)
}

// Fatal log a message at fatal level.
func (l *logger) Fatal(v ...interface{}) {
	l._log(LevelFatal, v...)
}

// Fatalf format and log a message at fatal level.
func (l *logger) Fatalf(f string, v ...interface{}) {
	l._logf(LevelFatal, f, v...)
}

// IsErrorEnabled is ERROR level enabled
func (l *logger) IsErrorEnabled() bool {
	return l.IsLevelEnabled(LevelError)
}

// Error log a message at error level.
func (l *logger) Error(v ...interface{}) {
	l._log(LevelError, v...)
}

// Errorf format and log a message at error level.
func (l *logger) Errorf(f string, v ...interface{}) {
	l._logf(LevelError, f, v...)
}

// IsWarnEnabled is WARN level enabled
func (l *logger) IsWarnEnabled() bool {
	return l.IsLevelEnabled(LevelWarn)
}

// Warn log a message at warning level.
func (l *logger) Warn(v ...interface{}) {
	l._log(LevelWarn, v...)
}

// Warnf format and log a message at warning level.
func (l *logger) Warnf(f string, v ...interface{}) {
	l._logf(LevelWarn, f, v...)
}

// IsInfoEnabled is INFO level enabled
func (l *logger) IsInfoEnabled() bool {
	return l.IsLevelEnabled(LevelInfo)
}

// Info log a message at info level.
func (l *logger) Info(v ...interface{}) {
	l._log(LevelInfo, v...)
}

// Infof format and log a message at info level.
func (l *logger) Infof(f string, v ...interface{}) {
	l._logf(LevelInfo, f, v...)
}

// IsDebugEnabled is DEBUG level enabled
func (l *logger) IsDebugEnabled() bool {
	return l.IsLevelEnabled(LevelDebug)
}

// Debug log a message at debug level.
func (l *logger) Debug(v ...interface{}) {
	l._log(LevelDebug, v...)
}

// Debugf format log a message at debug level.
func (l *logger) Debugf(f string, v ...interface{}) {
	l._logf(LevelDebug, f, v...)
}

// IsTraceEnabled is TRACE level enabled
func (l *logger) IsTraceEnabled() bool {
	return l.IsLevelEnabled(LevelTrace)
}

// Trace log a message at trace level.
func (l *logger) Trace(v ...interface{}) {
	l._log(LevelTrace, v...)
}

// Tracef format and log a message at trace level.
func (l *logger) Tracef(f string, v ...interface{}) {
	l._logf(LevelTrace, f, v...)
}

func (l *logger) _log(lvl Level, v ...interface{}) {
	if l.IsLevelEnabled(lvl) {
		s := l._printv(v...)
		le := newEvent(l, lvl, s)
		l.log.write(le)
	}
}

func (l *logger) _logf(lvl Level, f string, v ...interface{}) {
	if l.IsLevelEnabled(lvl) {
		s := l._printf(f, v...)
		le := newEvent(l, lvl, s)
		l.log.write(le)
	}
}

func (l *logger) _printv(v ...interface{}) string {
	if len(v) == 0 {
		return ""
	}
	return fmt.Sprint(v...)
}

func (l *logger) _printf(f string, v ...interface{}) string {
	return fmt.Sprintf(f, v...)
}
