package iox

import (
	"bufio"
	"errors"
	"io"

	"github.com/askasoft/pango/str"
)

// Discard is an io.Writer on which all Write calls succeed
// without doing anything.
var Discard = io.Discard

// Copy copies from src to dst until either EOF is reached
// on src or an error occurs. It returns the number of bytes
// copied and the first error encountered while copying, if any.
//
// A successful Copy returns err == nil, not err == EOF.
// Because Copy is defined to read from src until EOF, it does
// not treat an EOF from Read as an error to be reported.
//
// If src implements the WriterTo interface,
// the copy is implemented by calling src.WriteTo(dst).
// Otherwise, if dst implements the ReaderFrom interface,
// the copy is implemented by calling dst.ReadFrom(src).
func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}

// CopyBuffer is identical to Copy except that it stages through the
// provided buffer (if one is required) rather than allocating a
// temporary one. If buf is nil, one is allocated; otherwise if it has
// zero length, CopyBuffer panics.
//
// If either src implements WriterTo or dst implements ReaderFrom,
// buf will not be used to perform the copy.
func CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	return io.CopyBuffer(dst, src, buf)
}

// CopyN copies n bytes (or until an error) from src to dst.
// It returns the number of bytes copied and the earliest
// error encountered while copying.
// On return, written == n if and only if err == nil.
//
// If dst implements the ReaderFrom interface,
// the copy is implemented using it.
func CopyN(dst io.Writer, src io.Reader, n int64) (written int64, err error) {
	return io.CopyN(dst, src, n)
}

// Drain drain the reader
func Drain(r io.Reader) {
	io.Copy(Discard, r) //nolint: errcheck
}

// DrainAndClose drain and close the reader
func DrainAndClose(r io.ReadCloser) {
	Drain(r)
	r.Close()
}

// NewSectionReader returns a SectionReader that reads from r
// starting at offset off and stops with EOF after n bytes.
func NewSectionReader(r io.ReaderAt, off int64, n int64) *io.SectionReader {
	return io.NewSectionReader(r, off, n)
}

// NopCloser returns a ReadCloser with a no-op Close method wrapping
// the provided Reader r.
func NopCloser(r io.Reader) io.ReadCloser {
	return io.NopCloser(r)
}

// ReadAll reads from r until an error or EOF and returns the data it read.
// A successful call returns err == nil, not err == EOF. Because ReadAll is
// defined to read from src until EOF, it does not treat an EOF from Read
// as an error to be reported.
func ReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

// ReadAtLeast reads from r into buf until it has read at least n bytes.
// It returns the number of bytes copied and an error if fewer bytes were read.
// The error is EOF only if no bytes were read.
// If an EOF happens after reading fewer than n bytes,
// ReadAtLeast returns ErrUnexpectedEOF.
// If n is greater than the length of buf, ReadAtLeast returns ErrShortBuffer.
// The return size >= n if and only if err == nil.
// If r returns an error having read at least n bytes, the error is dropped.
func ReadAtLeast(r io.Reader, buf []byte, n int) (int, error) {
	return io.ReadAtLeast(r, buf, n)
}

// ReadFull reads exactly len(buf) bytes from r into buf.
// It returns the number of bytes copied and an error if fewer bytes were read.
// The error is EOF only if no bytes were read.
// If an EOF happens after reading some but not all the bytes,
// ReadFull returns ErrUnexpectedEOF.
// The return size == len(buf) if and only if err == nil.
// If r returns an error having read at least len(buf) bytes, the error is dropped.
func ReadFull(r io.Reader, buf []byte) (int, error) {
	return io.ReadFull(r, buf)
}

// SkipBOM skip bom and return a reader
func SkipBOM(r io.Reader) (io.Reader, error) {
	br := bufio.NewReader(r)
	c, _, err := br.ReadRune()

	if errors.Is(err, io.EOF) {
		return br, nil
	}
	if err != nil {
		return br, err
	}

	if c != BOM {
		// Not a BOM -- put the rune back
		err = br.UnreadRune()
	}
	return br, err
}

// TeeReader returns a Reader that writes to w what it reads from r.
// All reads from r performed through it are matched with
// corresponding writes to w. There is no internal buffering -
// the write must complete before the read completes.
// Any error encountered while writing is reported as a read error.
func TeeReader(r io.Reader, w io.Writer) io.Reader {
	return io.TeeReader(r, w)
}

// WriteString writes the contents of the string s to w, which accepts a slice of bytes.
// If w implements StringWriter, its WriteString method is invoked directly.
// Otherwise, w.Write is called exactly once.
func WriteString(w io.Writer, s string) (n int, err error) {
	if sw, ok := w.(io.StringWriter); ok {
		return sw.WriteString(s)
	}
	return w.Write(str.UnsafeBytes(s))
}
