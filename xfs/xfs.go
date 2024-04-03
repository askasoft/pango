package xfs

import (
	"bytes"
	"io/fs"
	"time"
)

type XFS interface {
	fs.FS

	// FindFile find a file
	FindFile(id string) (*File, error)

	// SaveFile save a file
	SaveFile(id string, filename string, modTime time.Time, data []byte) (*File, error)

	// ReadFile read file data
	ReadFile(fid string) ([]byte, error)

	// DeleteFile delete file by id
	DeleteFile(id string) error

	// DeleteFiles delete file by ids
	DeleteFiles(ids ...string) (int64, error)

	// DeletePrefix delete files by prefix
	DeletePrefix(prefix string) (int64, error)

	// DeleteBefore delete files by time
	DeleteBefore(before time.Time) (int64, error)

	// DeletePrefixBefore delete files by prefix and time
	DeletePrefixBefore(prefix string, before time.Time) (int64, error)

	// DeleteWhere delete files by customized where filter
	DeleteWhere(where string, args ...any) (int64, error)

	// DeleteAll delete all files
	DeleteAll() (int64, error)

	// Truncate truncate files
	Truncate() error
}

//----------------------------------------------------

// FSFile implements fs.File interface
type FSFile struct {
	XFS  XFS
	File *File
	r    *bytes.Reader
}

func (f *FSFile) open() error {
	if f.r == nil {
		data, err := f.XFS.ReadFile(f.File.ID)
		if err != nil {
			return err
		}

		f.File.Data = data
		f.r = bytes.NewReader(f.File.Data)
	}
	return nil
}

func (f *FSFile) Close() error {
	return nil
}

func (f *FSFile) Read(p []byte) (int, error) {
	err := f.open()
	if err != nil {
		return 0, err
	}
	return f.r.Read(p)
}

func (f *FSFile) Seek(offset int64, whence int) (int64, error) {
	err := f.open()
	if err != nil {
		return 0, err
	}
	return f.r.Seek(offset, whence)
}

func (f *FSFile) Readdir(count int) ([]fs.FileInfo, error) {
	return nil, fs.ErrInvalid
}

func (f *FSFile) Stat() (fs.FileInfo, error) {
	return &FSFileInfo{f.File}, nil
}

//----------------------------------------------------

// FSFileInfo implements fs.FileInfo interface
type FSFileInfo struct {
	f *File
}

// base name of the file
func (fi *FSFileInfo) Name() string {
	return fi.f.Name
}

// length in bytes for regular files; system-dependent for others
func (fi *FSFileInfo) Size() int64 {
	return fi.f.Size
}

// file mode bits
func (fi *FSFileInfo) Mode() fs.FileMode {
	return fs.FileMode(0400)
}

// modification time
func (fi *FSFileInfo) ModTime() time.Time {
	return fi.f.Time
}

// abbreviation for Mode().IsDir()
func (fi *FSFileInfo) IsDir() bool {
	return false
}

// underlying data source (can return nil)
func (fi *FSFileInfo) Sys() any {
	return nil
}
