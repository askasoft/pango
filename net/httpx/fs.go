package httpx

import (
	"net/http"
	"os"
	"time"
)

// FixedModTimeFS returns a FileSystem with fixed ModTime
func FixedModTimeFS(hfs http.FileSystem, mt time.Time) http.FileSystem {
	return &staticFS{hfs, mt}
}

// staticFS a FileSystem with fixed ModTime
type staticFS struct {
	http.FileSystem
	FixedModTime time.Time
}

// Open implements http.FileSystem.Open()
func (sfs *staticFS) Open(name string) (http.File, error) {
	file, err := sfs.FileSystem.Open(name)
	return &staticFile{File: file, modTime: sfs.FixedModTime}, err
}

// staticFile a File with fixed ModTime
type staticFile struct {
	http.File
	modTime time.Time
}

// Stat implements http.File.Stat()
func (sf *staticFile) Stat() (os.FileInfo, error) {
	fi, err := sf.File.Stat()
	return &staticFileInfo{FileInfo: fi, modTime: sf.modTime}, err
}

// staticFileInfo a FileInfo with fixed ModTime
type staticFileInfo struct {
	os.FileInfo
	modTime time.Time
}

// ModTime implements FileInfo.ModTime()
func (sfi *staticFileInfo) ModTime() time.Time {
	mt := sfi.FileInfo.ModTime()
	if mt.IsZero() {
		return sfi.modTime
	}
	return mt
}
