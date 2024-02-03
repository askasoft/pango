package xwf

import (
	"io/fs"
	"time"
)

// fileinfo implements fs.FileInfo interface
type fileinfo struct {
	f *File
}

// base name of the file
func (fi *fileinfo) Name() string {
	return fi.f.Name
}

// length in bytes for regular files; system-dependent for others
func (fi *fileinfo) Size() int64 {
	return fi.f.Size
}

// file mode bits
func (fi *fileinfo) Mode() fs.FileMode {
	return fs.FileMode(0400)
}

// modification time
func (fi *fileinfo) ModTime() time.Time {
	return fi.f.UpdatedAt
}

// abbreviation for Mode().IsDir()
func (fi *fileinfo) IsDir() bool {
	return false
}

// underlying data source (can return nil)
func (fi *fileinfo) Sys() any {
	return nil
}
