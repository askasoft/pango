package enc

import "io"

// MimeChunkSize As required by RFC 2045, 6.8. (page 25) for base64.
const MimeChunkSize = 76

// PemChunkSize PEM chunk size per RFC 1421 section 4.3.2.4.
const PemChunkSize = 64

// Base64LineWriter limits text encoded in base64 to 76 characters per line
type Base64LineWriter struct {
	Writer io.Writer
	Length int
	size   int
}

// NewBase64LineWriter create a writer for split base64 to 76 characters per line
func NewBase64LineWriter(w io.Writer) *Base64LineWriter {
	return &Base64LineWriter{Writer: w, Length: MimeChunkSize}
}

// Write implements io.Writer
func (w *Base64LineWriter) Write(p []byte) (int, error) {
	n := 0
	for len(p)+w.size > w.Length {
		w.Writer.Write(p[:w.Length-w.size])
		w.Writer.Write([]byte("\r\n"))
		p = p[w.Length-w.size:]
		n += w.Length - w.size
		w.size = 0
	}

	w.Writer.Write(p)
	w.size += len(p)

	return n + len(p), nil
}
