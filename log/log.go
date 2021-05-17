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
)

//--------------------------------------------------------------------
// package functions
//

// default Log instance
var _log = NewLog()

// Default get default Log
func Default() *Log {
	return _log
}

// Config config log by configuration file
func Config(filename string) error {
	return _log.Config(filename)
}

// GetLogger returns a new logger
func GetLogger(name string) Logger {
	return _log.GetLogger(name)
}

// Outputer return a io.Writer for go log.SetOutput
// callerDepth: default is 1 (means +1)
// if the outputer is used by go std log, set callerDepth to 2
// example:
//   import (
//     golog "log"
//     "github.com/pandafw/pango/log"
//   )
//   golog.SetOutput(log.Outputer("GO", log.LevelInfo, 2))
//
func Outputer(name string, lvl Level, callerDepth ...int) io.Writer {
	return _log.Outputer(name, lvl, callerDepth...)
}

// Async set the log to asynchronous and start the goroutine
// if size < 1 then stop async goroutine
func Async(size int) {
	_log.Async(size)
}

// IsAsync return the logger's async
func IsAsync() bool {
	return _log.async
}

// GetFormatter get the formatter
func GetFormatter() Formatter {
	return _log.GetFormatter()
}

// SetFormatter set the formatter.
func SetFormatter(lf Formatter) {
	_log.SetFormatter(lf)
}

// GetWriter get the writer
func GetWriter() Writer {
	return _log.GetWriter()
}

// SetWriter set the writer.
func SetWriter(lw Writer) {
	_log.SetWriter(lw)
}

// GetLevel return the logger's level
func GetLevel() Level {
	return _log.GetLevel()
}

// SetLevel set the logger's level
func SetLevel(lvl Level) {
	_log.SetLevel(lvl)
}

// SetLevels set the logger levels
func SetLevels(lvls map[string]Level) {
	_log.SetLevels(lvls)
}

// GetCallerDepth return the logger's caller depth
func GetCallerDepth() int {
	return _log.GetCallerDepth()
}

// SetCallerDepth set the logger's caller depth (!!SLOW!!), 0: disable runtime.Caller()
func SetCallerDepth(d int) {
	_log.SetCallerDepth(d)
}

// GetProp get logger property
func GetProp(k string) interface{} {
	return _log.GetProp(k)
}

// SetProp set logger property
func SetProp(k string, v interface{}) {
	_log.SetProp(k, v)
}

// GetProps get logger properties
func GetProps() map[string]interface{} {
	return _log.GetProps()
}

// SetProps set logger properties
func SetProps(props map[string]interface{}) {
	_log.SetProps(props)
}

// Close will remove all writers and stop async goroutine
func Close() {
	_log.Close()
}

// IsFatalEnabled is FATAL level enabled
func IsFatalEnabled() bool {
	return _log.logger.IsFatalEnabled()
}

// Fatal log a message at fatal level.
func Fatal(v ...interface{}) {
	_log.logger._log(LevelFatal, v...)
}

// Fatalf format and log a message at fatal level.
func Fatalf(f string, v ...interface{}) {
	_log.logger._logf(LevelFatal, f, v...)
}

// IsErrorEnabled is ERROR level enabled
func IsErrorEnabled() bool {
	return _log.logger.IsErrorEnabled()
}

// Error log a message at error level.
func Error(v ...interface{}) {
	_log.logger._log(LevelError, v...)
}

// Errorf format and log a message at error level.
func Errorf(f string, v ...interface{}) {
	_log.logger._logf(LevelError, f, v...)
}

// IsWarnEnabled is WARN level enabled
func IsWarnEnabled() bool {
	return _log.logger.IsWarnEnabled()
}

// Warn log a message at warning level.
func Warn(v ...interface{}) {
	_log.logger._log(LevelWarn, v...)
}

// Warnf format and log a message at warning level.
func Warnf(f string, v ...interface{}) {
	_log.logger._logf(LevelWarn, f, v...)
}

// IsInfoEnabled is INFO level enabled
func IsInfoEnabled() bool {
	return _log.logger.IsInfoEnabled()
}

// Info log a message at info level.
func Info(v ...interface{}) {
	_log.logger._log(LevelInfo, v...)
}

// Infof format and log a message at info level.
func Infof(f string, v ...interface{}) {
	_log.logger._logf(LevelInfo, f, v...)
}

// IsDebugEnabled is DEBUG level enabled
func IsDebugEnabled() bool {
	return _log.logger.IsDebugEnabled()
}

// Debug log a message at debug level.
func Debug(v ...interface{}) {
	_log.logger._log(LevelDebug, v...)
}

// Debugf format log a message at debug level.
func Debugf(f string, v ...interface{}) {
	_log.logger._logf(LevelDebug, f, v...)
}

// IsTraceEnabled is TRACE level enabled
func IsTraceEnabled() bool {
	return _log.logger.IsTraceEnabled()
}

// Trace log a message at trace level.
func Trace(v ...interface{}) {
	_log.logger._log(LevelTrace, v...)
}

// Tracef format and log a message at trace level.
func Tracef(f string, v ...interface{}) {
	_log.logger._logf(LevelTrace, f, v...)
}
