package iox

import (
	"io"
	"sync"
)

// WrapWriter a prefix/suffix wrap writer
type wrapWriter struct {
	writer io.Writer
	prefix string
	suffix string
}

// WrapWriter return a prefix/suffix wrap writer
func WrapWriter(writer io.Writer, prefix, suffix string) io.Writer {
	return &wrapWriter{writer, prefix, suffix}
}

// Write io.Writer implement
func (ww *wrapWriter) Write(p []byte) (n int, err error) {
	if _, err = ww.writer.Write([]byte(ww.prefix)); err != nil {
		return
	}

	if n, err = ww.writer.Write(p); err != nil {
		return
	}

	_, err = ww.writer.Write([]byte(ww.suffix))
	return
}

// syncWriter synchronize writer
type syncWriter struct {
	w io.Writer
	m sync.Mutex
}

// SyncWriter return a synchronized writer
func SyncWriter(w io.Writer) io.Writer {
	return &syncWriter{w: w}
}

func (sw *syncWriter) Write(p []byte) (int, error) {
	sw.m.Lock()
	defer sw.m.Unlock()
	return sw.w.Write(p)
}
