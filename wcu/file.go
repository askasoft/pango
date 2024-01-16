package wcu

import (
	"io"
	"os"
	"path"

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

func OpenFile(filename string, charset string) (io.ReadCloser, error) {
	fr, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	dr, err := Transform(fr, charset)
	if err != nil {
		fr.Close()
		return nil, err
	}

	return &rc{r: dr, c: fr}, nil
}

func IsHTMLFile(filename string) bool {
	ext := path.Ext(filename)
	html := (str.EqualFold(ext, ".html") || str.EqualFold(ext, ".htm"))
	return html
}

func DetectAndOpenFile(filename string) (io.ReadCloser, error) {
	fr, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	html := IsHTMLFile(filename)
	dr, err := DetectAndTransform(fr, html)
	if err != nil {
		fr.Close()
		return nil, err
	}

	return &rc{r: dr, c: fr}, nil
}

func ReadFile(filename string, charset string) ([]byte, error) {
	wf, err := OpenFile(filename, charset)
	if err != nil {
		return nil, err
	}
	defer wf.Close()

	return io.ReadAll(wf)
}

func DetectAndReadFile(filename string) ([]byte, error) {
	wf, err := DetectAndOpenFile(filename)
	if err != nil {
		return nil, err
	}
	defer wf.Close()

	return io.ReadAll(wf)
}
