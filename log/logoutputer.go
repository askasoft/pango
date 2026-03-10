package log

import (
	"fmt"
	"io"
)

// Outputer interface for io.Writer
type Outputer interface {
	io.Writer
	Printf(format string, args ...any)
}

// outputer a io.Writer implement for go log.SetOutput
type outputer struct {
	logger Logger
	level  Level
	skip   int
}

// Write io.Writer implement
func (o *outputer) Write(p []byte) (int, error) {
	if o.logger.IsLevelEnabled(o.level) {
		le := NewEvent(o.logger, o.level, string(p))
		if o.skip > 0 {
			le.CallerSkip(o.skip, o.logger.GetTraceLevel() >= o.level)
		}
		o.logger.Write(le)
	}
	return len(p), nil
}

// Printf printf implement
func (o *outputer) Printf(format string, args ...any) {
	if o.logger.IsLevelEnabled(o.level) {
		s := fmt.Sprintf(format, args...)
		le := NewEvent(o.logger, o.level, s)
		if o.skip > 0 {
			le.CallerSkip(o.skip, o.logger.GetTraceLevel() >= o.level)
		}
		o.logger.Write(le)
	}
}
