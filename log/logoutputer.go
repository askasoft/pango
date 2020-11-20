package log

import (
	"unsafe"
)

// outputer a io.Writer implement for go log.SetOutput
type outputer struct {
	level  int
	logger Logger
}

// Write io.Writer implement
func (o *outputer) Write(p []byte) (int, error) {
	o.logger.Log(o.level, o.string(p))
	return len(p), nil
}

// string cast []byte to string
func (o *outputer) string(p []byte) string {
	return *(*string)(unsafe.Pointer(&p))
}
