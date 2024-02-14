package xfs

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
)

// HFS converts xfs to a http.FileSystem implementation.
func HFS(xfs XFS) http.FileSystem {
	return &hfs{xfs}
}

// hfs implements http.FileSystem interface
type hfs struct {
	xfs XFS
}

func (hfs *hfs) Open(name string) (http.File, error) {
	f, err := hfs.xfs.Open(name)
	if err != nil {
		return nil, err
	}
	return hfile{f}, nil
}

type hfile struct {
	file fs.File
}

var errMissingSeek = errors.New("io.File missing Seek method")

func (hf hfile) Close() error               { return hf.file.Close() }
func (hf hfile) Read(b []byte) (int, error) { return hf.file.Read(b) }
func (hf hfile) Stat() (fs.FileInfo, error) { return hf.file.Stat() }

func (hf hfile) Seek(offset int64, whence int) (int64, error) {
	s, ok := hf.file.(io.Seeker)
	if !ok {
		return 0, errMissingSeek
	}
	return s.Seek(offset, whence)
}

// Readdir overrides the http.File default implementation.
func (hf hfile) Readdir(count int) ([]os.FileInfo, error) {
	// this disables directory listing
	return nil, nil
}
