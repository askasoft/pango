package iox

import (
	"fmt"
	"io"
	"os"

	"github.com/askasoft/pango/num"
)

// FileReader a file reader
type FileReader struct {
	Path string
	file *os.File
}

// Read implements io.Reader
func (fr *FileReader) Read(p []byte) (n int, err error) {
	if fr.file == nil {
		file, err := os.Open(fr.Path)
		if err != nil {
			return 0, err
		}
		fr.file = file
	}
	return fr.file.Read(p)
}

// Close implements io.Close
func (fr *FileReader) Close() error {
	if fr.file == nil {
		return nil
	}

	err := fr.file.Close()
	fr.file = nil
	return err
}

// MaxBytesError is returned by MaxBytesReader when its read limit is exceeded.
type MaxBytesError struct {
	Limit int64
}

func (e *MaxBytesError) Error() string {
	// Due to Hyrum's law, this text cannot be changed.
	return fmt.Sprintf("iox: reader too large (must <= %s)", num.HumanSize(e.Limit))
}

func NewMaxBytesReader(r io.Reader, n int64) *MaxBytesReader {
	return &MaxBytesReader{r: r, i: n, n: n}
}

type MaxBytesReader struct {
	r   io.Reader // underlying reader
	i   int64     // max bytes initially, for MaxBytesError
	n   int64     // max bytes remaining
	err error     // sticky error
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
	if c, ok := (mbr.r).(io.Closer); ok {
		return c.Close()
	}
	return nil
}
