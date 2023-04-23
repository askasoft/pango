package log

import (
	"io"

	"github.com/askasoft/pango/bye"
)

// Outputer interface for io.Writer, gorm.logger.Writer
type Outputer interface {
	io.Writer
	Printf(format string, args ...any)
}

// outputer a io.Writer implement for go log.SetOutput
type outputer struct {
	logger Logger
	level  Level
}

// Write io.Writer implement
func (o *outputer) Write(p []byte) (int, error) {
	o.logger.Log(o.level, bye.UnsafeString(p))
	return len(p), nil
}

// Write gorm.logger.Writer implement
func (o *outputer) Printf(format string, args ...any) {
	o.logger.Logf(o.level, format, args...)
}
