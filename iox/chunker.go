package iox

import "io"

// MimeChunkSize As required by RFC 2045, 6.8. (page 25) for base64.
const MimeChunkSize = 76

// PemChunkSize PEM chunk size per RFC 1421 section 4.3.2.4.
const PemChunkSize = 64

// ChunkLineWriter limits text to n characters per line
type ChunkLineWriter struct {
	Writer io.Writer
	EOL    string
	Length int // chunk line size

	size int // internal writed line size
}

// NewMimeChunkWriter create a writer for split base64 to 76 characters per line
func NewMimeChunkWriter(w io.Writer) *ChunkLineWriter {
	return &ChunkLineWriter{Writer: w, EOL: CRLF, Length: MimeChunkSize}
}

// NewPemChunkWriter create a writer for split base64 to 76 characters per line
func NewPemChunkWriter(w io.Writer) *ChunkLineWriter {
	return &ChunkLineWriter{Writer: w, EOL: CRLF, Length: PemChunkSize}
}

// Write implements io.Writer
func (w *ChunkLineWriter) Write(p []byte) (n int, err error) {
	n = 0
	for len(p)+w.size > w.Length {
		_, err = w.Writer.Write(p[:w.Length-w.size])
		_, err = w.Writer.Write([]byte(w.EOL))
		p = p[w.Length-w.size:]
		n += w.Length - w.size
		w.size = 0
	}

	_, err = w.Writer.Write(p)
	w.size += len(p)

	return n + len(p), err
}
