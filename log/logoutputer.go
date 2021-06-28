package log

import "github.com/pandafw/pango/bye"

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
