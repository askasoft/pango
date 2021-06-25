package httpx

import (
	"net/http"
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
