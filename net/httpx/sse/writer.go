package sse

import (
	"io"

	"github.com/askasoft/pango/str"
)

type stringWriter interface {
	io.Writer
	WriteString(string) (int, error)
}

type stringWrapper struct {
	io.Writer
}

func (w stringWrapper) WriteString(s string) (int, error) {
	return w.Write(str.UnsafeBytes(s))
}

func wrapWriter(writer io.Writer) stringWriter {
	if w, ok := writer.(stringWriter); ok {
		return w
	}

	return stringWrapper{writer}
}
