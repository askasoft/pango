package xin

import (
	"bytes"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/askasoft/pango/net/httpx"
)

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

func ServeFSHandler(prefix string, hfs http.FileSystem, filePath string) HandlerFunc {
	fileServer := httpx.FileServer(prefix, hfs, filePath)
	return func(c *Context) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func ServeFSCHandler(prefix string, hfsc func(c *Context) http.FileSystem, filePath string) HandlerFunc {
	return func(c *Context) {
		hfs := hfsc(c)
		fileServer := httpx.FileServer(prefix, hfs, filePath)
		fileServer.ServeHTTP(c.Writer, c.Request)
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

	r.GET(relativePath, handlers...)
	r.HEAD(relativePath, handlers...)
}

// Static serves files from the given file system root.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// example:
//
//	xin.Static(r, "/static", "/var/www", xin.NewCacheControlSetter("public, max-age=31536000").Handler())
func Static(r IRoutes, relativePath, root string, handlers ...HandlerFunc) {
	StaticFS(r, relativePath, httpx.Dir(root), "", handlers...)
}

// StaticFS works just like `Static()` but a custom `http.FileSystem` can be used instead.
func StaticFS(r IRoutes, relativePath string, hfs http.FileSystem, filePath string, handlers ...HandlerFunc) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}

	prefix := path.Join(r.BasePath(), relativePath)

	handler := ServeFSHandler(prefix, hfs, filePath)
	handlers = append(handlers, handler)

	// Register GET and HEAD handlers
	urlPattern := path.Join(relativePath, "/*path")
	r.GET(urlPattern, handlers...)
	r.HEAD(urlPattern, handlers...)
}

// StaticFSC works just like `StaticFS()` but a dynamic `http.FileSystem` can be used instead.
func StaticFSC(r IRoutes, relativePath string, hfsc func(c *Context) http.FileSystem, filePath string, handlers ...HandlerFunc) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}

	prefix := path.Join(r.BasePath(), relativePath)

	handler := ServeFSCHandler(prefix, hfsc, filePath)
	handlers = append(handlers, handler)

	// Register GET and HEAD handlers
	urlPattern := path.Join(relativePath, "/*path")
	r.GET(urlPattern, handlers...)
	r.HEAD(urlPattern, handlers...)
}

// StaticFSFile registers a single route in order to serve a single file of the filesystem.
// xin.StaticFSFile(r, "favicon.ico", hfs, "./resources/favicon.ico", xin.NewCacheControlSetter("public, max-age=31536000").Handler())
func StaticFSFile(r IRoutes, relativePath string, hfs http.FileSystem, filePath string, handlers ...HandlerFunc) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}

	handler := ServeFSFileHandler(hfs, filePath)
	handlers = append(handlers, handler)

	r.GET(relativePath, handlers...)
	r.HEAD(relativePath, handlers...)
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

	r.GET(relativePath, handlers...)
	r.HEAD(relativePath, handlers...)
}

// -----------------------------------------------

// CacheControlSetter set Cache-Control header when statusCode == 200
type CacheControlSetter struct {
	CacheControl string
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
	return NewHeaderWriter(w, "Cache-Control", ccs.CacheControl)
}

func (ccs *CacheControlSetter) Handle(c *Context) {
	c.Writer = ccs.WrapWriter(c.Writer)
}

func (ccs *CacheControlSetter) Handler() HandlerFunc {
	return ccs.Handle
}
