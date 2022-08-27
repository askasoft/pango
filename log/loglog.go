package log

import (
	"sync"
	"time"
)

// Log is default logger in application.
// it can contain several writers and log message into all writers.
type Log struct {
	logger *logger
	level  Level
	trace  Level

	writer Writer
	levels map[string]Level
	logfmt Formatter

	mutex sync.Mutex
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
		writer: NewConsoleWriter(),
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

// GetWriter get the log writer
func (log *Log) GetWriter() Writer {
	return log.writer
}

// SetWriter set the log writer
func (log *Log) SetWriter(lw Writer) {
	log.writer = lw
}

// SwitchWriter use lw to replace the log writer
func (log *Log) SwitchWriter(lw Writer) {
	log.mutex.Lock()
	defer log.mutex.Unlock()

	ow := log.writer

	if osw, ok := ow.(*SyncWriter); ok {
		osw.SetWriter(lw)
		log.writer = lw
		return
	}

	if oaw, ok := ow.(*AsyncWriter); ok {
		oaw.SetWriter(lw)
		oaw.StopAfter(time.Second)
		log.writer = lw
		return
	}

	ow.Close()
	log.writer = lw
}

// Flush flush all chan data.
func (log *Log) Flush() {
	log.writer.Flush()
}

// Close close logger, flush all data and close the writer.
func (log *Log) Close() {
	log.writer.Close()
}

// write write a log event
func (log *Log) write(le *Event) {
	log.writer.Write(le) //nolint: errcheck
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
	return log.logger.name
}

// SetName set the logger's name
func (log *Log) SetName(name string) {
	log.logger.name = name
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
func (log *Log) GetProp(k string) any {
	ps := log.logger.props
	if ps == nil {
		return nil
	}
	return ps[k]
}

// SetProp set logger property
func (log *Log) SetProp(k string, v any) {
	log.logger.SetProp(k, v)
}

// GetProps get logger properties
func (log *Log) GetProps() map[string]any {
	tm := log.logger.props
	if tm == nil {
		return nil
	}

	// new return props
	nm := make(map[string]any, len(tm))
	for k, v := range tm {
		nm[k] = v
	}
	return nm
}

// SetProps set logger properties
func (log *Log) SetProps(props map[string]any) {
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
func (log *Log) Log(lvl Level, v ...any) {
	log.logger._log(lvl, v...)
}

// Logf format and log a message at specified level.
func (log *Log) Logf(lvl Level, f string, v ...any) {
	log.logger._logf(lvl, f, v...)
}

// IsFatalEnabled is FATAL level enabled
func (log *Log) IsFatalEnabled() bool {
	return log.IsLevelEnabled(LevelFatal)
}

// Fatal log a message at fatal level.
func (log *Log) Fatal(v ...any) {
	log.logger._log(LevelFatal, v...)
}

// Fatalf format and log a message at fatal level.
func (log *Log) Fatalf(f string, v ...any) {
	log.logger._logf(LevelFatal, f, v...)
}

// IsErrorEnabled is ERROR level enabled
func (log *Log) IsErrorEnabled() bool {
	return log.IsLevelEnabled(LevelError)
}

// Error log a message at error level.
func (log *Log) Error(v ...any) {
	log.logger._log(LevelError, v...)
}

// Errorf format and log a message at error level.
func (log *Log) Errorf(f string, v ...any) {
	log.logger._logf(LevelError, f, v...)
}

// IsWarnEnabled is WARN level enabled
func (log *Log) IsWarnEnabled() bool {
	return log.IsLevelEnabled(LevelWarn)
}

// Warn log a message at warning level.
func (log *Log) Warn(v ...any) {
	log.logger._log(LevelWarn, v...)
}

// Warnf format and log a message at warning level.
func (log *Log) Warnf(f string, v ...any) {
	log.logger._logf(LevelWarn, f, v...)
}

// IsInfoEnabled is INFO level enabled
func (log *Log) IsInfoEnabled() bool {
	return log.IsLevelEnabled(LevelInfo)
}

// Info log a message at info level.
func (log *Log) Info(v ...any) {
	log.logger._log(LevelInfo, v...)
}

// Infof format and log a message at info level.
func (log *Log) Infof(f string, v ...any) {
	log.logger._logf(LevelInfo, f, v...)
}

// IsDebugEnabled is DEBUG level enabled
func (log *Log) IsDebugEnabled() bool {
	return log.IsLevelEnabled(LevelDebug)
}

// Debug log a message at debug level.
func (log *Log) Debug(v ...any) {
	log.logger._log(LevelDebug, v...)
}

// Debugf format log a message at debug level.
func (log *Log) Debugf(f string, v ...any) {
	log.logger._logf(LevelDebug, f, v...)
}

// IsTraceEnabled is TRACE level enabled
func (log *Log) IsTraceEnabled() bool {
	return log.IsLevelEnabled(LevelTrace)
}

// Trace log a message at trace level.
func (log *Log) Trace(v ...any) {
	log.logger._log(LevelTrace, v...)
}

// Tracef format and log a message at trace level.
func (log *Log) Tracef(f string, v ...any) {
	log.logger._logf(LevelTrace, f, v...)
}
