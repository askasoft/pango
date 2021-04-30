package iox

import "io"

// MimeChunkSize As required by RFC 2045, 6.8. (page 25) for base64.
const MimeChunkSize = 76

// PemChunkSize PEM chunk size per RFC 1421 section 4.3.2.4.
const PemChunkSize = 64

// ChunkLineWriter limits text to n characters per line
type ChunkLineWriter struct {
	Writer   io.Writer
	EOL      string
	LineSize int // chunk line size

	written int // internal written line size
}

// NewMimeChunkWriter create a writer for split base64 to 76 characters per line
func NewMimeChunkWriter(w io.Writer) *ChunkLineWriter {
	return &ChunkLineWriter{Writer: w, EOL: CRLF, LineSize: MimeChunkSize}
}

// NewPemChunkWriter create a writer for split base64 to 76 characters per line
func NewPemChunkWriter(w io.Writer) *ChunkLineWriter {
	return &ChunkLineWriter{Writer: w, EOL: CRLF, LineSize: PemChunkSize}
}

// Write implements io.Writer
func (cw *ChunkLineWriter) Write(p []byte) (n int, err error) {
	n = 0
	for len(p)+cw.written > cw.LineSize {
		_, err = cw.Writer.Write(p[:cw.LineSize-cw.written])
		_, err = cw.Writer.Write([]byte(cw.EOL))
		p = p[cw.LineSize-cw.written:]
		n += cw.LineSize - cw.written
		cw.written = 0
	}

	_, err = cw.Writer.Write(p)
	cw.written += len(p)
	n += len(p)
	return
}
