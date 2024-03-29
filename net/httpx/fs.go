package httpx

import (
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"strings"
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

// AppendPrefix returns a handler that serves HTTP requests by appending the
// given prefix from the request URL's Path (and RawPath if set) and invoking
// the handler hh.
func AppendPrefix(prefix string, hh http.Handler) http.Handler {
	if prefix == "" {
		return hh
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := prefix + r.URL.Path
		rp := prefix + r.URL.RawPath
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = p
		r2.URL.RawPath = rp
		hh.ServeHTTP(w, r2)
	})
}

// URLReplace returns a handler that serves HTTP requests by replacing the
// request URL's Path (and RawPath if set) (use strings.Replace(path, src, des) and invoking
// the handler hh.
func URLReplace(src, des string, hh http.Handler) http.Handler {
	if src == "" || src == des {
		return hh
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.Replace(r.URL.Path, src, des, 1)
		rp := strings.Replace(r.URL.RawPath, src, des, 1)
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = p
		r2.URL.RawPath = rp
		hh.ServeHTTP(w, r2)
	})
}

func FileServer(prefix string, hfs http.FileSystem, filePath string) http.Handler {
	fileServer := http.FileServer(hfs)
	if prefix == "" || prefix == "/" {
		fileServer = AppendPrefix(filePath, fileServer)
	} else if filePath == "" || filePath == "." {
		fileServer = http.StripPrefix(prefix, fileServer)
	} else {
		fileServer = URLReplace(prefix, filePath, fileServer)
	}
	return fileServer
}

//----------------------------------------------------------------

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
