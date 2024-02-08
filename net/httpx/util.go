package httpx

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

// SaveMultipartFile save multipart file to the specific local file 'dst'.
func SaveMultipartFile(file *multipart.FileHeader, dst string) error {
	dir := path.Dir(dst)
	if err := os.MkdirAll(dir, os.FileMode(0770)); err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
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
