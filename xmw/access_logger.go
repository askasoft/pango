package xmw

import (
	"time"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/xin"
)

const (
	AccessLogStartKey = "X_ACCESS_START"
	AccessLogEndKey   = "X_ACCESS_END"
)

// AccessLogger access logger middleware for XIN
type AccessLogger struct {
	writer   AccessLogWriter
	disabled bool
}

// DefaultAccessLogger create a access logger middleware for XIN
// Equals: NewAccessLogger(NewAccessLogTextWriter(xin.Logger.GetOutputer("XAL", log.LevelTrace), AccessLogTextFormat))
func DefaultAccessLogger(xin *xin.Engine) *AccessLogger {
	return NewAccessLogger(NewAccessLogTextWriter(xin.Logger.GetOutputer("XAL", log.LevelTrace), AccessLogTextFormat))
}

// NewAccessLogger create a log middleware for xin access logger
func NewAccessLogger(writer AccessLogWriter) *AccessLogger {
	return &AccessLogger{writer: writer}
}

// Disable disable the logger or not
func (al *AccessLogger) Disable(disabled bool) {
	al.disabled = disabled
}

// Handle process xin request
func (al *AccessLogger) Handle(c *xin.Context) {
	w := al.writer
	if w == nil || al.disabled {
		c.Next()
		return
	}

	c.Set(AccessLogStartKey, time.Now())
	c.Next()
	c.Set(AccessLogEndKey, time.Now())

	w.Write(c)
}

// SetWriter set the access logger writer
func (al *AccessLogger) SetWriter(alw AccessLogWriter) {
	al.writer = alw
}
