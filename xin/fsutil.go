package xin

import (
	"bytes"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/askasoft/pango/net/httpx"
)

type WriterWrapper func(http.ResponseWriter) http.ResponseWriter

func ServeFileHandler(filePath string, wws ...WriterWrapper) HandlerFunc {
	if len(wws) == 0 {
		return func(c *Context) {
			http.ServeFile(c.Writer, c.Request, filePath)
		}
	}

	return func(c *Context) {
		ww := wws[0](c.Writer)
		http.ServeFile(ww, c.Request, filePath)
	}
}

func ServeFSFileHandler(hfs http.FileSystem, filePath string, wws ...WriterWrapper) HandlerFunc {
	if len(wws) == 0 {
		return func(c *Context) {
			org := c.Request.URL.Path

			defer func(url string) {
				c.Request.URL.Path = url
			}(org)

			c.Request.URL.Path = filePath
			http.FileServer(hfs).ServeHTTP(c.Writer, c.Request)
		}
	}

	return func(c *Context) {
		org := c.Request.URL.Path

		defer func(url string) {
			c.Request.URL.Path = url
		}(org)

		c.Request.URL.Path = filePath

		ww := wws[0](c.Writer)

		http.FileServer(hfs).ServeHTTP(ww, c.Request)
	}
}

func ServerFSHandler(prefix string, hfs http.FileSystem, filePath string, wws ...WriterWrapper) HandlerFunc {
	fileServer := http.FileServer(hfs)
	if prefix == "" || prefix == "/" {
		fileServer = httpx.AppendPrefix(filePath, fileServer)
	} else if filePath == "" || filePath == "." {
		fileServer = http.StripPrefix(prefix, fileServer)
	} else {
		fileServer = httpx.URLReplace(prefix, filePath, fileServer)
	}

	if len(wws) == 0 {
		return func(c *Context) {
			fileServer.ServeHTTP(c.Writer, c.Request)
		}
	}

	return func(c *Context) {
		ww := wws[0](c.Writer)
		fileServer.ServeHTTP(ww, c.Request)
	}
}

func ServeContent(data []byte, modtime time.Time, wws ...WriterWrapper) HandlerFunc {
	if modtime.IsZero() {
		modtime = time.Now()
	}

	if len(wws) == 0 {
		return func(c *Context) {
			name := filepath.Base(c.Request.URL.Path)
			http.ServeContent(c.Writer, c.Request, name, modtime, bytes.NewReader(data))
		}
	}

	return func(c *Context) {
		ww := wws[0](c.Writer)
		name := filepath.Base(c.Request.URL.Path)
		http.ServeContent(ww, c.Request, name, modtime, bytes.NewReader(data))
	}
}

// StaticFile registers a single route in order to serve a single file of the local filesystem.
// router.StaticFile("favicon.ico", "./resources/favicon.ico", "public, max-age=31536000")
func StaticFile(r IRoutes, relativePath, filePath string, wws ...WriterWrapper) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}

	handler := ServeFileHandler(filePath, wws...)

	r.GET(relativePath, handler)
	r.HEAD(relativePath, handler)
}

// Static serves files from the given file system root.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use :
//
//	router.Static("/static", "/var/www")
func Static(r IRoutes, relativePath, root string, wws ...WriterWrapper) {
	StaticFS(r, relativePath, httpx.Dir(root), "", wws...)
}

// StaticFS works just like `Static()` but a custom `http.FileSystem` can be used instead.
func StaticFS(r IRoutes, relativePath string, hfs http.FileSystem, filePath string, wws ...WriterWrapper) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}

	prefix := path.Join(r.BasePath(), relativePath)

	handler := ServerFSHandler(prefix, hfs, filePath, wws...)

	urlPattern := path.Join(relativePath, "/*path")

	// Register GET and HEAD handlers
	r.GET(urlPattern, handler)
	r.HEAD(urlPattern, handler)
}

// StaticFSFile registers a single route in order to serve a single file of the filesystem.
// router.StaticFSFile("favicon.ico", "./resources/favicon.ico", hfs, "public, max-age=31536000")
func StaticFSFile(r IRoutes, relativePath string, hfs http.FileSystem, filePath string, wws ...WriterWrapper) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}

	handler := ServeFSFileHandler(hfs, filePath, wws...)

	r.GET(relativePath, handler)
	r.HEAD(relativePath, handler)
}

// StaticContent registers a single route in order to serve a single file of the data.
// //go:embed favicon.ico
// var favicon []byte
// router.StaticContent("favicon.ico", favicon, time.Now(), "public, max-age=31536000")
func StaticContent(r IRoutes, relativePath string, data []byte, modtime time.Time, wws ...WriterWrapper) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static content")
	}

	handler := ServeContent(data, modtime, wws...)

	r.GET(relativePath, handler)
	r.HEAD(relativePath, handler)
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

func (ccs *CacheControlSetter) WrapWriter(hrw http.ResponseWriter) http.ResponseWriter {
	if ccs.CacheControl == "" {
		return hrw
	}
	return httpx.NewHeaderWriter(hrw, "Cache-Control", ccs.CacheControl)
}

func (ccs *CacheControlSetter) WriterWrapper() WriterWrapper {
	return ccs.WrapWriter
}
