package httpx

import (
	"io/fs"
	"net/http"
	"os"
	"time"
)

// Dir returns a http.FileSystem that can be used by http.FileServer().
// if browsable == true, then it works the same as http.Dir() otherwise it returns
// a filesystem that prevents http.FileServer() to list the directory files.
func Dir(root string, browsable ...bool) http.FileSystem {
	fs := http.Dir(root)
	if len(browsable) > 0 && browsable[0] {
		return fs
	}
	return &onlyFilesFS{fs}
}

// FS returns a http.FileSystem that can be used by http.FileServer().
// if browsable == true, then it works the same as http.FS() otherwise it returns
// a filesystem that prevents http.FileServer() to list the directory files.
func FS(fsys fs.FS, browsable ...bool) http.FileSystem {
	fs := http.FS(fsys)
	if len(browsable) > 0 && browsable[0] {
		return fs
	}
	return &onlyFilesFS{fs}
}

type onlyFilesFS struct {
	fs http.FileSystem
}

// Open conforms to http.Filesystem.
func (fs onlyFilesFS) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return onlyFile{f}, nil
}

type onlyFile struct {
	http.File
}

// Readdir overrides the http.File default implementation.
func (f onlyFile) Readdir(count int) ([]os.FileInfo, error) {
	// this disables directory listing
	return nil, nil
}

//----------------------------------------------------------------

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
