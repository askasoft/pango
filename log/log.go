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

// package logger
var _logger Logger = newPkgLogger()

func newPkgLogger() Logger {
	l := _log.GetLogger("")
	l.SetCallerDepth(l.GetCallerDepth() + 1)
	return l
}

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
func Outputer(name string, lvl Level) io.Writer {
	return _log.Outputer(name, lvl)
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

// SetFormatter set the formatter.
func SetFormatter(lf Formatter) {
	_log.SetFormatter(lf)
	_logger.SetFormatter(lf)
}

// SetWriter set the writer.
func SetWriter(lw Writer) {
	_log.SetWriter(lw)
}

// Close will remove all writers and stop async goroutine
func Close() {
	_log.Close()
}

// GetLevel return the logger's level
func GetLevel() Level {
	return _log.GetLevel()
}

// SetLevel set the logger's level
func SetLevel(lvl Level) {
	_log.SetLevel(lvl)
	_logger.SetLevel(lvl)
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
	_logger.SetCallerDepth(d + 1)
}

// IsFatalEnabled is FATAL level enabled
func IsFatalEnabled() bool {
	return _logger.IsFatalEnabled()
}

// Fatal log a message at fatal level.
func Fatal(v ...interface{}) {
	_logger.Fatal(v...)
}

// Fatalf format and log a message at fatal level.
func Fatalf(f string, v ...interface{}) {
	_logger.Fatalf(f, v...)
}

// IsErrorEnabled is ERROR level enabled
func IsErrorEnabled() bool {
	return _logger.IsErrorEnabled()
}

// Error log a message at error level.
func Error(v ...interface{}) {
	_logger.Error(v...)
}

// Errorf format and log a message at error level.
func Errorf(f string, v ...interface{}) {
	_logger.Errorf(f, v...)
}

// IsWarnEnabled is WARN level enabled
func IsWarnEnabled() bool {
	return _logger.IsWarnEnabled()
}

// Warn log a message at warning level.
func Warn(v ...interface{}) {
	_logger.Warn(v...)
}

// Warnf format and log a message at warning level.
func Warnf(f string, v ...interface{}) {
	_logger.Warnf(f, v...)
}

// IsInfoEnabled is INFO level enabled
func IsInfoEnabled() bool {
	return _logger.IsInfoEnabled()
}

// Info log a message at info level.
func Info(v ...interface{}) {
	_logger.Info(v...)
}

// Infof format and log a message at info level.
func Infof(f string, v ...interface{}) {
	_logger.Infof(f, v...)
}

// IsDebugEnabled is DEBUG level enabled
func IsDebugEnabled() bool {
	return _logger.IsDebugEnabled()
}

// Debug log a message at debug level.
func Debug(v ...interface{}) {
	_logger.Debug(v...)
}

// Debugf format log a message at debug level.
func Debugf(f string, v ...interface{}) {
	_logger.Debugf(f, v...)
}

// IsTraceEnabled is TRACE level enabled
func IsTraceEnabled() bool {
	return _logger.IsTraceEnabled()
}

// Trace log a message at trace level.
func Trace(v ...interface{}) {
	_logger.Trace(v...)
}

// Tracef format and log a message at trace level.
func Tracef(f string, v ...interface{}) {
	_logger.Tracef(f, v...)
}
