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

// LimitWriter returns a Writer that writes limited bytes to the underlying writer.
// The underlying implementation is a *LimitedWriter.
func LimitWriter(w io.Writer, n int64) io.Writer { return &LimitedWriter{w, n} }

// LimitedWriter implements io.Writer and writes the data to an io.Writer, but
// limits the total bytes written to it, discards the remaining bytes.
type LimitedWriter struct {
	W io.Writer // underlying writer
	N int64     // max bytes remaining
}

func (lw *LimitedWriter) Write(data []byte) (int, error) {
	if lw.N <= 0 {
		return len(data), nil
	}

	n := lw.N - int64(len(data))
	if n >= 0 {
		lw.N = n
		return lw.W.Write(data)
	}

	_, err := lw.W.Write(data[:int(lw.N)])
	lw.N = 0
	return len(data), err
}
