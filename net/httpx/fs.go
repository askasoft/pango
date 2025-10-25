package httpx

import (
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/askasoft/pango/asg"
)

// StripPrefix returns a handler that serves HTTP requests by removing the
// given prefix from the request URL's Path (and RawPath if set) and invoking
// the handler h. StripPrefix handles a request for a path that doesn't begin
// with prefix by replying with an HTTP 404 not found error. The prefix must
// match exactly: if the prefix in the request contains escaped characters
// the reply is also an HTTP 404 not found error.
func StripPrefix(prefix string, h http.Handler) http.Handler {
	return http.StripPrefix(prefix, h)
}

// StringReplace returns a handler that serves HTTP requests by replacing the
// request URL's Path (and RawPath if set) (use strings.Replace(path, src, rep) and invoking
// the handler hh.
func StringReplace(src, rep string, hh http.Handler) http.Handler {
	if src == "" || src == rep {
		return hh
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.Replace(r.URL.Path, src, rep, 1)
		rp := strings.Replace(r.URL.RawPath, src, rep, 1)
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = p
		r2.URL.RawPath = rp
		hh.ServeHTTP(w, r2)
	})
}

// RegexpReplace returns a handler that serves HTTP requests by replacing the
// request URL's Path (and RawPath if set) (use Regexp.MustCompile(expr).ReplaceAllString(path, rep) and invoking
// the handler hh.
func RegexpReplace(expr, repl string, hh http.Handler) http.Handler {
	re := regexp.MustCompile(expr)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := re.ReplaceAllString(r.URL.Path, repl)
		rp := re.ReplaceAllString(r.URL.RawPath, repl)
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = p
		r2.URL.RawPath = rp
		hh.ServeHTTP(w, r2)
	})
}

// ----------------------------------------------------------------
// Dir returns a http.FileSystem that can be used by http.FileServer().
// if browsable == true, then it works the same as http.Dir() otherwise it returns
// a filesystem that prevents http.FileServer() to list the directory files.
func Dir(root string, browsable ...bool) http.FileSystem {
	fs := http.Dir(root)

	if asg.First(browsable) {
		return fs
	}
	return onlyFilesFS{fs}
}

// FS returns a http.FileSystem that can be used by http.FileServer().
// if browsable == true, then it works the same as http.FS() otherwise it returns
// a filesystem that prevents http.FileServer() to list the directory files.
func FS(fsys fs.FS, browsable ...bool) http.FileSystem {
	fs := http.FS(fsys)

	if asg.First(browsable) {
		return fs
	}
	return onlyFilesFS{fs}
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
