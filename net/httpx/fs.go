package httpx

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// PrefixFS append a prefix to a http.FileSystem.
func PrefixFS(hfs http.FileSystem, prefix string) http.FileSystem {
	return prefixFS{hfs, prefix}
}

type prefixFS struct {
	hfs    http.FileSystem
	prefix string
}

func (pfs prefixFS) Open(name string) (http.File, error) {
	p := filepath.ToSlash(filepath.Join(pfs.prefix, name))
	return pfs.hfs.Open(p)
}

// PrefixHandler returns a handler that serves HTTP requests by appending the
// given prefix from the request URL's Path (and RawPath if set) and invoking
// the handler h.
func PrefixHandler(prefix string, h http.Handler) http.Handler {
	if prefix == "" {
		return h
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
		h.ServeHTTP(w, r2)
	})
}

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
