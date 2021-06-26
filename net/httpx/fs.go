package httpx

import (
	"net/http"
	"net/url"
	"path/filepath"
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
