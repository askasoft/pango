package xin

import (
	"bytes"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/askasoft/pango/net/httpx"
)

// Dir returns a http.FileSystem that can be used by http.FileServer().
// if browsable == true, then it works the same as http.Dir() otherwise it returns
// a filesystem that prevents http.FileServer() to list the directory files.
func Dir(root string, browsable ...bool) http.FileSystem {
	return httpx.Dir(root, browsable...)
}

// FS returns a http.FileSystem that can be used by http.FileServer().
// if browsable == true, then it works the same as http.FS() otherwise it returns
// a filesystem that prevents http.FileServer() to list the directory files.
func FS(fsys fs.FS, browsable ...bool) http.FileSystem {
	return httpx.FS(fsys, browsable...)
}

func ServeFileHandler(filePath string) HandlerFunc {
	return func(c *Context) {
		http.ServeFile(c.Writer, c.Request, filePath)
	}
}

func ServeFSFileHandler(hfs http.FileSystem, filePath string) HandlerFunc {
	return func(c *Context) {
		org := c.Request.URL.Path

		defer func(url string) {
			c.Request.URL.Path = url
		}(org)

		c.Request.URL.Path = filePath
		http.FileServer(hfs).ServeHTTP(c.Writer, c.Request)
	}
}

func ServeFSFuncFileHandler(hfsc func(c *Context) http.FileSystem, filePath string) HandlerFunc {
	return func(c *Context) {
		org := c.Request.URL.Path

		defer func(url string) {
			c.Request.URL.Path = url
		}(org)

		c.Request.URL.Path = filePath
		hfs := hfsc(c)
		http.FileServer(hfs).ServeHTTP(c.Writer, c.Request)
	}
}

func ServeFSHandler(prefix string, hfs http.FileSystem) HandlerFunc {
	fsv := httpx.StripPrefix(prefix, http.FileServer(hfs))
	return func(c *Context) {
		fsv.ServeHTTP(c.Writer, c.Request)
	}
}

func ServeFSFuncHandler(prefix string, hfsc func(c *Context) http.FileSystem) HandlerFunc {
	return func(c *Context) {
		hfs := hfsc(c)
		fsv := httpx.StripPrefix(prefix, http.FileServer(hfs))
		fsv.ServeHTTP(c.Writer, c.Request)
	}
}

func ServeContentHandler(data []byte, modtime time.Time) HandlerFunc {
	if modtime.IsZero() {
		modtime = time.Now()
	}

	return func(c *Context) {
		name := path.Base(c.Request.URL.Path)
		http.ServeContent(c.Writer, c.Request, name, modtime, bytes.NewReader(data))
	}
}

// StaticFile registers a single route in order to serve a single file of the local filesystem.
// example:
//
//	xin.StaticFile(r, "favicon.ico", "./resources/favicon.ico", xin.NewCacheControlSetter("public, max-age=31536000").Handler())
func StaticFile(r IRoutes, relativePath, filePath string, handlers ...HandlerFunc) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}

	handler := ServeFileHandler(filePath)
	handlers = append(handlers, handler)

	r.HEAD(relativePath, handlers...)
	r.GET(relativePath, handlers...)
}

// Static serves files from the local file system directory.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// example:
//
//	xin.Static(r, "/static", "/var/www", xin.NewCacheControlSetter("public, max-age=31536000").Handler())
func Static(r IRoutes, relativePath, localPath string, handlers ...HandlerFunc) {
	StaticFS(r, relativePath, httpx.Dir(localPath), handlers...)
}

// StaticFS works just like `Static()` but a custom `http.FileSystem` can be used instead.
func StaticFS(r IRoutes, relativePath string, hfs http.FileSystem, handlers ...HandlerFunc) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}

	prefix := path.Join(r.BasePath(), relativePath)

	handler := ServeFSHandler(prefix, hfs)
	handlers = append(handlers, handler)

	// Register GET and HEAD handlers
	urlPattern := path.Join(relativePath, "/*path")
	r.HEAD(urlPattern, handlers...)
	r.GET(urlPattern, handlers...)
}

// StaticFSFunc works just like `StaticFS()` but a dynamic `http.FileSystem` can be used instead.
func StaticFSFunc(r IRoutes, relativePath string, hfsc func(c *Context) http.FileSystem, handlers ...HandlerFunc) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}

	prefix := path.Join(r.BasePath(), relativePath)

	handler := ServeFSFuncHandler(prefix, hfsc)
	handlers = append(handlers, handler)

	// Register GET and HEAD handlers
	urlPattern := path.Join(relativePath, "/*path")
	r.HEAD(urlPattern, handlers...)
	r.GET(urlPattern, handlers...)
}

// StaticFSFile registers a single route in order to serve a single file of the filesystem.
// xin.StaticFSFile(r, "favicon.ico", hfs, "./resources/favicon.ico", xin.NewCacheControlSetter("public, max-age=31536000").Handler())
func StaticFSFile(r IRoutes, relativePath string, hfs http.FileSystem, filePath string, handlers ...HandlerFunc) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}

	handler := ServeFSFileHandler(hfs, filePath)
	handlers = append(handlers, handler)

	r.HEAD(relativePath, handlers...)
	r.GET(relativePath, handlers...)
}

// StaticFSFuncFile works just like `StaticFSFile()` but a dynamic `http.FileSystem` can be used instead.
func StaticFSFuncFile(r IRoutes, relativePath string, hfsc func(c *Context) http.FileSystem, filePath string, handlers ...HandlerFunc) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}

	handler := ServeFSFuncFileHandler(hfsc, filePath)
	handlers = append(handlers, handler)

	r.HEAD(relativePath, handlers...)
	r.GET(relativePath, handlers...)
}

// StaticContent registers a single route in order to serve a single file of the data.
// example:
//
//	//go:embed favicon.ico
//	var favicon []byte
//	xin.StaticContent(r, "favicon.ico", favicon, time.Now(), xin.NewCacheControlSetter("public, max-age=31536000").Handler())
func StaticContent(r IRoutes, relativePath string, data []byte, modtime time.Time, handlers ...HandlerFunc) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static content")
	}

	handler := ServeContentHandler(data, modtime)
	handlers = append(handlers, handler)

	r.HEAD(relativePath, handlers...)
	r.GET(relativePath, handlers...)
}

// -----------------------------------------------

// CacheControlSetter set Cache-Control header when statusCode == 200
type CacheControlSetter struct {
	CacheControl string
	Overwrite    bool
}

func NewCacheControlSetter(cacheControls ...string) *CacheControlSetter {
	ccs := &CacheControlSetter{}
	ccs.SetCacheControl(cacheControls...)
	return ccs
}

func (ccs *CacheControlSetter) SetCacheControl(cacheControls ...string) {
	ccs.CacheControl = strings.Join(cacheControls, ", ")
}

func (ccs *CacheControlSetter) WrapWriter(w ResponseWriter) ResponseWriter {
	if ccs.CacheControl == "" {
		return w
	}
	return NewHeaderWriter(w, "Cache-Control", ccs.CacheControl, ccs.Overwrite)
}

func (ccs *CacheControlSetter) Handle(c *Context) {
	c.Writer = ccs.WrapWriter(c.Writer)
}

// -----------------------------------------------

// DisableAcceptRanges overwrite header "Accept-Ranges: none" when statusCode == 200
// http.FileServer always set "Accept-Ranges: bytes" header,
// overwrite this header to prevent 206 Partial Content Download.
func DisableAcceptRanges(c *Context) {
	c.Writer = NewHeaderWriter(c.Writer, "Accept-Ranges", "none", true)
}
