package iox

import (
	"io"
	"sync"
)

// WriterWrapFunc a writer wrapper function
type WriterWrapFunc func(w io.Writer) io.Writer

// WriteCloserWrapFunc a write closer wrapper function
type WriteCloserWrapFunc func(w io.WriteCloser) io.WriteCloser

// WrapWriter a prefix/suffix append writer
type WrapWriter struct {
	Writer io.Writer
	Prefix string
	Suffix string
}

// Write io.Writer implement
func (ww *WrapWriter) Write(p []byte) (int, error) {
	ww.Writer.Write([]byte(ww.Prefix))
	n, err := ww.Writer.Write(p)
	ww.Writer.Write([]byte(ww.Suffix))
	return n, err
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
