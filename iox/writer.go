package iox

import (
	"io"
	"sync"

	"github.com/askasoft/pango/bye"
	"github.com/askasoft/pango/str"
)

// wrapWriter a prefix/suffix wrap writer
type wrapWriter struct {
	w      io.Writer
	prefix string
	suffix string
}

// WrapWriter return a prefix/suffix wrap writer
func WrapWriter(w io.Writer, prefix, suffix string) io.Writer {
	return &wrapWriter{w, prefix, suffix}
}

// Write io.Writer implement
func (ww *wrapWriter) Write(p []byte) (n int, err error) {
	if ww.prefix != "" {
		if _, err = ww.w.Write(str.UnsafeBytes(ww.prefix)); err != nil {
			return
		}
	}

	if n, err = ww.w.Write(p); err != nil {
		return
	}

	if ww.suffix != "" {
		_, err = ww.w.Write(str.UnsafeBytes(ww.suffix))
	}
	return
}

// WriteString io.StringWriter implement
func (ww *wrapWriter) WriteString(s string) (n int, err error) {
	if sw, ok := ww.w.(io.StringWriter); ok {
		if ww.prefix != "" {
			if _, err = sw.WriteString(ww.prefix); err != nil {
				return
			}
		}

		if n, err = sw.WriteString(s); err != nil {
			return
		}

		if ww.suffix != "" {
			_, err = sw.WriteString(ww.suffix)
		}
	} else {
		n, err = ww.Write(str.UnsafeBytes(s))
	}
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

// LimitWriter returns a Writer that writes limited bytes to the underlying writer.
// The underlying implementation is a *limitWriter.
func LimitWriter(w io.Writer, n int64) io.Writer { return &limitWriter{w, n} }

// limitWriter implements io.Writer and writes the data to an io.Writer, but
// limits the total bytes written to it, discards the remaining bytes.
type limitWriter struct {
	w io.Writer // underlying writer
	n int64     // max bytes remaining
}

func (lw *limitWriter) Write(data []byte) (int, error) {
	if lw.n <= 0 {
		return len(data), nil
	}

	n := lw.n - int64(len(data))
	if n >= 0 {
		lw.n = n
		return lw.w.Write(data)
	}

	_, err := lw.w.Write(data[:int(lw.n)])
	lw.n = 0
	return len(data), err
}

type lineWriter struct {
	w   io.Writer
	eol string
}

// LineWriter return a eol append writer
func LineWriter(w io.Writer, eol ...string) io.Writer {
	lw := &lineWriter{w: w}
	if len(eol) > 0 {
		lw.eol = eol[0]
	} else {
		lw.eol = "\n"
	}
	return lw
}

func (lw *lineWriter) Write(p []byte) (n int, err error) {
	n, err = lw.w.Write(p)
	if err != nil {
		return
	}

	_, err = lw.w.Write(str.UnsafeBytes(lw.eol))
	return
}

func (lw *lineWriter) WriteString(s string) (n int, err error) {
	if sw, ok := lw.w.(io.StringWriter); ok {
		n, err = sw.WriteString(s)
		if err != nil {
			return
		}
		_, err = sw.WriteString(lw.eol)
	} else {
		n, err = lw.Write(str.UnsafeBytes(s))
	}
	return n, err
}

type stripWriter struct {
	w io.Writer
}

// StripWriter return a string strip writer
func StripWriter(w io.Writer) io.Writer {
	return &stripWriter{w}
}

func (sw *stripWriter) Write(p []byte) (n int, err error) {
	n = len(p)

	s := str.Strip(bye.UnsafeString(p))
	if s == "" {
		return
	}

	_, err = sw.w.Write(str.UnsafeBytes(s))
	return
}

func (sw *stripWriter) WriteString(s string) (n int, err error) {
	n = len(s)

	s = str.Strip(s)
	if s == "" {
		return
	}

	if w, ok := sw.w.(io.StringWriter); ok {
		_, err = w.WriteString(s)
	} else {
		_, err = sw.Write(str.UnsafeBytes(s))
	}
	return
}
