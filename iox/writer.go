package iox

import (
	"bytes"
	"io"
	"sync"
	"unicode"

	"github.com/askasoft/pango/str"
)

// RepeatWrite repeat write bytes s.
func RepeatWrite(w io.Writer, s []byte, count int) (int, error) {
	if count <= 0 {
		return 0, nil
	}

	total := 0
	for range count {
		n, err := w.Write(s)
		total += n
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

// RepeatWriteString repeat write string s.
func RepeatWriteString(w io.Writer, s string, count int) (int, error) {
	if count <= 0 {
		return 0, nil
	}

	total := 0
	for range count {
		n, err := WriteString(w, s)
		total += n
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

// ProxyWriter proxy writer
type ProxyWriter struct {
	W io.Writer
}

func (pw *ProxyWriter) Write(p []byte) (int, error) {
	return pw.W.Write(p)
}

func (pw *ProxyWriter) WriteString(s string) (int, error) {
	return WriteString(pw.W, s)
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
	return
}

// linePrefixWriter a prefix line wrap writer
type linePrefixWriter struct {
	w      io.Writer
	prefix string
	lastlf bool
}

// LinePrefixWriter return a prefix line wrap writer
func LinePrefixWriter(w io.Writer, prefix string) io.Writer {
	return &linePrefixWriter{w, prefix, true}
}

// Write io.Writer implement
func (lpw *linePrefixWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}

	n = len(p)
	ps := str.UnsafeBytes(lpw.prefix)

	for len(p) > 0 {
		if lpw.lastlf {
			if _, err = lpw.w.Write(ps); err != nil {
				return
			}
		}

		i := bytes.IndexByte(p, '\n')
		if i < 0 {
			lpw.lastlf = false
			if _, err = lpw.w.Write(p); err != nil {
				return
			}
			break
		}

		lpw.lastlf = true
		if _, err = lpw.w.Write(p[:i+1]); err != nil {
			return
		}
		p = p[i+1:]
	}
	return
}

// WriteString io.StringWriter implement
func (lpw *linePrefixWriter) WriteString(s string) (n int, err error) {
	if s == "" {
		return
	}

	if sw, ok := lpw.w.(io.StringWriter); ok {
		n = len(s)
		for s != "" {
			if lpw.lastlf {
				if _, err = sw.WriteString(lpw.prefix); err != nil {
					return
				}
			}

			i := str.IndexByte(s, '\n')
			if i < 0 {
				lpw.lastlf = false
				if _, err = sw.WriteString(s); err != nil {
					return
				}
				break
			}

			lpw.lastlf = true
			if _, err = sw.WriteString(s[:i+1]); err != nil {
				return
			}
			s = s[i+1:]
		}
	} else {
		n, err = lpw.Write(str.UnsafeBytes(s))
	}
	return
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

	s := str.Strip(str.UnsafeString(p))
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

type CompactWriter struct {
	w io.Writer       // underlying writer
	f func(rune) bool // is compact rune function
	r rune            // replace rune
	c bool            // last rune is compact or not
}

// SpaceCompactor return a space compact writer
func SpaceCompactWriter(w io.Writer) *CompactWriter {
	return NewCompactWriter(w, unicode.IsSpace, ' ')
}

// NewCompactWriter return a compact writer
func NewCompactWriter(w io.Writer, f func(rune) bool, r rune) *CompactWriter {
	return &CompactWriter{w: w, f: f, r: r, c: true}
}

func (cw *CompactWriter) Reset(c bool) {
	cw.c = c
}

func (cw *CompactWriter) Write(p []byte) (int, error) {
	return cw.WriteString(str.UnsafeString(p))
}

func (cw *CompactWriter) WriteString(s string) (n int, err error) {
	n = len(s)
	if n == 0 {
		return
	}

	var sb str.Builder
	for _, r := range s {
		c := cw.f(r)
		if cw.c {
			if !c {
				sb.WriteRune(r)
			}
		} else {
			if c {
				r = cw.r
			}
			sb.WriteRune(r)
		}
		cw.c = c
	}
	if sb.Len() > 0 {
		_, err = WriteString(cw.w, sb.String())
	}
	return
}
