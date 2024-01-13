package wcu

import (
	"io"
	"os"
	"path"

	"github.com/askasoft/pango/str"
)

type WFile struct {
	file *os.File
	read io.Reader
}

func (wf *WFile) Read(p []byte) (n int, err error) {
	return wf.read.Read(p)
}

func (wf *WFile) Close() error {
	return wf.file.Close()
}

func OpenFile(filename string, charset string) (*WFile, error) {
	fr, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	dr, err := Transform(fr, charset)
	if err != nil {
		fr.Close()
		return nil, err
	}

	return &WFile{file: fr, read: dr}, nil
}

func IsHTMLFile(filename string) bool {
	ext := path.Ext(filename)
	html := (str.EqualFold(ext, ".html") || str.EqualFold(ext, ".htm"))
	return html
}

func DetectAndOpenFile(filename string) (*WFile, error) {
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

	return &WFile{file: fr, read: dr}, nil
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
