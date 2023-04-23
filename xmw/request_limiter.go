package xmw

import (
	"io"

	"github.com/askasoft/pango/xin"
)

// MaxBytesError is returned by MaxBytesReader when its read limit is exceeded.
type MaxBytesError struct {
	Limit int64
}

func (e *MaxBytesError) Error() string {
	// Due to Hyrum's law, this text cannot be changed.
	return "http: request body too large"
}

// RequestLimiter http request limit middleware
type RequestLimiter struct {
	MaxBodySize int64
}

// NewRequestLimiter create a default RequestLimiter
func NewRequestLimiter(maxBodySize int64) *RequestLimiter {
	return &RequestLimiter{MaxBodySize: maxBodySize}
}

// Handler returns the xin.HandlerFunc
func (rl *RequestLimiter) Handler() xin.HandlerFunc {
	return func(c *xin.Context) {
		rl.handle(c)
	}
}

func MaxBytesReader(c *xin.Context, r io.ReadCloser, n int64) io.ReadCloser {
	if n <= 0 {
		return r
	}
	return &maxBytesReader{c: c, r: r, i: n, n: n}
}

// handle process xin request
func (rl *RequestLimiter) handle(c *xin.Context) {
	if rl.MaxBodySize > 0 {
		c.Request.Body = MaxBytesReader(c, c.Request.Body, rl.MaxBodySize)
	}
	c.Next()
}

type maxBytesReader struct {
	c   *xin.Context
	r   io.ReadCloser // underlying reader
	i   int64         // max bytes initially, for MaxBytesError
	n   int64         // max bytes remaining
	err error         // sticky error
}

func (l *maxBytesReader) Read(p []byte) (n int, err error) {
	if l.err != nil {
		return 0, l.err
	}
	if len(p) == 0 {
		return 0, nil
	}
	// If they asked for a 32KB byte read but only 5 bytes are
	// remaining, no need to read 32KB. 6 bytes will answer the
	// question of the whether we hit the limit or go past it.
	if int64(len(p)) > l.n+1 {
		p = p[:l.n+1]
	}
	n, err = l.r.Read(p)

	if int64(n) <= l.n {
		l.n -= int64(n)
		l.err = err
		return n, err
	}

	n = int(l.n)
	l.n = 0

	l.err = &MaxBytesError{l.i}
	l.c.AddError(l.err)
	return n, l.err
}

func (l *maxBytesReader) Close() error {
	return l.r.Close()
}
