// Package log provide a general log interface
// Usage:
//
// import "github.com/askasoft/pango/log"
//
//	log := log.NewLog()
//	log.SetWriter(log.NewAsyncWriter(log.NewConsoleWriter()))
//
// Use it like this:
//
//	log.Fatal("fatal")
//	log.Error("error")
//	log.Warn("warning")
//	log.Info("info")
//	log.Debug("debug")
//	log.Trace("trace")
//
// A Logger with name:
//
//	log := log.GetLogger("foo")
//	log.Debug("hello")
package log

import (
	"fmt"
	"os"
)

//--------------------------------------------------------------------
// package functions
//

// default Log instance
var _log = NewLog()

// Default returns the default Log instance used by the package-level functions.
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

// GetOutputer return a io.Writer for go log.SetOutput
// callerSkip: default is 1 (means +1)
// if the outputer is used by go std log, set callerSkip to 2
// example:
//
//	import (
//	  golog "log"
//	  "github.com/askasoft/pango/log"
//	)
//	golog.SetOutput(log.Outputer("GO", log.LevelInfo, 2))
func GetOutputer(name string, lvl Level, callerSkip ...int) Outputer {
	return _log.GetOutputer(name, lvl, callerSkip...)
}

// GetWriter get the writer
func GetWriter() Writer {
	return _log.GetWriter()
}

// SetWriter set the writer.
func SetWriter(lw Writer) {
	_log.SetWriter(lw)
}

// SwitchWriter use lw to replace the log writer
func SwitchWriter(lw Writer) {
	_log.SwitchWriter(lw)
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

// GetCallerSkip return the logger's caller skip
func GetCallerSkip() int {
	return _log.GetCallerSkip()
}

// SetCallerSkip set the logger's caller skip (!!SLOW!!), 0: disable runtime.Caller()
func SetCallerSkip(n int) {
	_log.SetCallerSkip(n)
}

// GetProp get logger property
func GetProp(k string) any {
	return _log.GetProp(k)
}

// SetProp set logger property
func SetProp(k string, v any) {
	_log.SetProp(k, v)
}

// GetProps get logger properties
func GetProps() map[string]any {
	return _log.GetProps()
}

// SetProps set logger properties
func SetProps(props map[string]any) {
	_log.SetProps(props)
}

// Flush flush all chan data.
func Flush() {
	_log.Flush()
}

// Close will remove all writers and stop async goroutine
func Close() {
	_log.Close()
}

// IsFatalEnabled is FATAL level enabled
func IsFatalEnabled() bool {
	return _log.IsFatalEnabled()
}

// Fatal log a message at fatal level.
func Fatal(v ...any) {
	_log._log(LevelFatal, v...)
}

// Fatalf format and log a message at fatal level.
func Fatalf(f string, v ...any) {
	_log._logf(LevelFatal, f, v...)
}

// IsErrorEnabled is ERROR level enabled
func IsErrorEnabled() bool {
	return _log.IsErrorEnabled()
}

// Error log a message at error level.
func Error(v ...any) {
	_log._log(LevelError, v...)
}

// Errorf format and log a message at error level.
func Errorf(f string, v ...any) {
	_log._logf(LevelError, f, v...)
}

// IsWarnEnabled is WARN level enabled
func IsWarnEnabled() bool {
	return _log.IsWarnEnabled()
}

// Warn log a message at warning level.
func Warn(v ...any) {
	_log._log(LevelWarn, v...)
}

// Warnf format and log a message at warning level.
func Warnf(f string, v ...any) {
	_log._logf(LevelWarn, f, v...)
}

// IsInfoEnabled is INFO level enabled
func IsInfoEnabled() bool {
	return _log.IsInfoEnabled()
}

// Info log a message at info level.
func Info(v ...any) {
	_log._log(LevelInfo, v...)
}

// Infof format and log a message at info level.
func Infof(f string, v ...any) {
	_log._logf(LevelInfo, f, v...)
}

// IsDebugEnabled is DEBUG level enabled
func IsDebugEnabled() bool {
	return _log.IsDebugEnabled()
}

// Debug log a message at debug level.
func Debug(v ...any) {
	_log._log(LevelDebug, v...)
}

// Debugf format log a message at debug level.
func Debugf(f string, v ...any) {
	_log._logf(LevelDebug, f, v...)
}

// IsTraceEnabled is TRACE level enabled
func IsTraceEnabled() bool {
	return _log.IsTraceEnabled()
}

// Trace log a message at trace level.
func Trace(v ...any) {
	_log._log(LevelTrace, v...)
}

// Tracef format and log a message at trace level.
func Tracef(f string, v ...any) {
	_log._logf(LevelTrace, f, v...)
}

// Perror print error message to stderr
func Perror(a any) {
	fmt.Fprintln(os.Stderr, a)
}

// Perrorf print formatted error message to stderr
func Perrorf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintln(os.Stderr)
}
