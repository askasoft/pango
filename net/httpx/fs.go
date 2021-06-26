package httpx

import (
	"net/http"
	"os"
	"time"
)

// FixedModTimeFS returns a FileSystem with fixed ModTime
func FixedModTimeFS(hfs http.FileSystem, mt time.Time) http.FileSystem {
	return &StaticFS{hfs, mt}
}

// StaticFS a FileSystem with fixed ModTime
type StaticFS struct {
	http.FileSystem
	FixedModTime time.Time
}

// Open implements http.FileSystem.Open()
func (sfs *StaticFS) Open(name string) (http.File, error) {
	file, err := sfs.FileSystem.Open(name)
	return &StaticFile{File: file, modTime: sfs.FixedModTime}, err
}

// StaticFile a File with fixed ModTime
type StaticFile struct {
	http.File
	modTime time.Time
}

// Stat implements File.Stat()
func (sf *StaticFile) Stat() (os.FileInfo, error) {
	fi, err := sf.File.Stat()
	return &StaticFileInfo{FileInfo: fi, modTime: sf.modTime}, err
}

// StaticFileInfo a FileInfo with fixed ModTime
type StaticFileInfo struct {
	os.FileInfo
	modTime time.Time
}

// ModTime implements FileInfo.ModTime()
func (sfi *StaticFileInfo) ModTime() time.Time {
	mt := sfi.FileInfo.ModTime()
	if mt.IsZero() {
		return sfi.modTime
	}
	return mt
}
