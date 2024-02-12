package gormfile

import (
	"bytes"
	"io/fs"

	"github.com/askasoft/pango/xwa/xfu"
)

type File = xfu.File
type FileResult = xfu.FileResult
type FilesResult = xfu.FilesResult

type file struct {
	hfs *hfs
	f   *File
	r   *bytes.Reader
}

func (f *file) open() error {
	if f.r == nil {
		data, err := ReadFile(f.hfs.db, f.hfs.tn, f.f.ID)
		if err != nil {
			return err
		}

		f.f.Data = data
		f.r = bytes.NewReader(f.f.Data)
	}
	return nil
}

func (f *file) Close() error {
	return nil
}

func (f *file) Read(p []byte) (int, error) {
	err := f.open()
	if err != nil {
		return 0, err
	}
	return f.r.Read(p)
}

func (f *file) Seek(offset int64, whence int) (int64, error) {
	err := f.open()
	if err != nil {
		return 0, err
	}
	return f.r.Seek(offset, whence)
}

func (f *file) Readdir(count int) ([]fs.FileInfo, error) {
	return nil, fs.ErrInvalid
}

func (f *file) Stat() (fs.FileInfo, error) {
	return &fileinfo{f.f}, nil
}
