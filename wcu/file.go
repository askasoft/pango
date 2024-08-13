package wcu

import (
	"io"
	"os"
	"path/filepath"

	"github.com/askasoft/pango/str"
)

// rc read closer
type rc struct {
	r io.Reader
	c io.Closer
}

func (rc *rc) Read(p []byte) (n int, err error) {
	return rc.r.Read(p)
}

func (rc *rc) Close() error {
	return rc.c.Close()
}

// OpenFile use specified charset to open file and return a transformed ReadCloser.
// return (nil, error) if failed to open file or charset is unsupported
func OpenFile(filename string, charset string) (io.ReadCloser, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	r, _, err := Transform(f, charset)
	if err != nil {
		f.Close()
		return nil, err
	}

	return &rc{r: r, c: f}, nil
}

// ReadFile use specified charset to read file.
// return error if failed to open/read file or charset is unsupported
func ReadFile(filename string, charset string) ([]byte, error) {
	wf, err := OpenFile(filename, charset)
	if err != nil {
		return nil, err
	}
	defer wf.Close()

	return io.ReadAll(wf)
}

func IsHTMLFile(filename string) bool {
	ext := filepath.Ext(filename)
	html := (str.EqualFold(ext, ".html") || str.EqualFold(ext, ".htm"))
	return html
}

// DetectCharsetFile detect file's charset
func DetectCharsetFile(filename string, limit int) (string, error) {
	fr, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer fr.Close()

	html := IsHTMLFile(filename)
	_, cs, err := DetectCharsetReader(fr, limit, html)
	return cs, err
}

// DetectAndOpenFile use specified charsets to open file, return (transformed reader, encoding name, error).
// if charsets is not specified, or specified charset is unsupported,
// detect the charset of the file, and open a transformed ReadCloser.
// return (nil, "", err) if failed to open file, charset is unsupported or read error occurred.
func DetectAndOpenFile(filename string, detect int, charsets ...string) (io.ReadCloser, string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}

	for _, cs := range charsets {
		if cs != "" {
			r, enc, err := Transform(f, cs)
			if err == nil {
				return &rc{r: r, c: f}, enc, nil
			}
		}
	}

	html := IsHTMLFile(filename)
	r, enc, err := DetectAndTransform(f, detect, html)
	if err != nil {
		f.Close()
		return nil, enc, err
	}

	return &rc{r: r, c: f}, enc, nil
}

// DetectAndReadFile detect the charset of the file, and open a transformed ReadCloser
func DetectAndReadFile(filename string, detect int, charsets ...string) ([]byte, string, error) {
	r, enc, err := DetectAndOpenFile(filename, detect, charsets...)
	if err != nil {
		return nil, enc, err
	}
	defer r.Close()

	bs, err := io.ReadAll(r)
	return bs, enc, err
}
