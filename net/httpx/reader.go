package httpx

import (
	"fmt"
	"io"

	"github.com/askasoft/pango/num"
)

// MaxBytesError is returned by MaxBytesReader when its read limit is exceeded.
type MaxBytesError struct {
	Limit int64
}

func (e *MaxBytesError) Error() string {
	// Due to Hyrum's law, this text cannot be changed.
	return fmt.Sprintf("http: request body too large (must <= %s)", num.HumanSize(float64(e.Limit)))
}

func NewMaxBytesReader(r io.ReadCloser, n int64) *MaxBytesReader {
	return &MaxBytesReader{r: r, i: n, n: n}
}

type MaxBytesReader struct {
	r   io.ReadCloser // underlying reader
	i   int64         // max bytes initially, for MaxBytesError
	n   int64         // max bytes remaining
	err error         // sticky error
}

func (mbr *MaxBytesReader) Error() error {
	return mbr.err
}

func (mbr *MaxBytesReader) Read(p []byte) (n int, err error) {
	if mbr.err != nil {
		return 0, mbr.err
	}
	if len(p) == 0 {
		return 0, nil
	}

	// If they asked for a 32KB byte read but only 5 bytes are
	// remaining, no need to read 32KB. 6 bytes will answer the
	// question of the whether we hit the limit or go past it.
	if int64(len(p)) > mbr.n+1 {
		p = p[:mbr.n+1]
	}
	n, err = mbr.r.Read(p)

	if int64(n) <= mbr.n {
		mbr.n -= int64(n)
		mbr.err = err
		return n, err
	}

	n = int(mbr.n)
	mbr.n = 0

	mbr.err = &MaxBytesError{Limit: mbr.i}
	return n, mbr.err
}

func (mbr *MaxBytesReader) Close() error {
	return mbr.r.Close()
}
