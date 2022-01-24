package gin

import (
	"bytes"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/pandafw/pango/net/httpx"
)

// Public1Year Cache-Control: public, max-age=31536000
const Public1Year = "public, max-age=31536000"

var (
	// reg match english letters for http method name
	regEnLetter = regexp.MustCompile("^[A-Z]+$")

	// anyMethods for RouterGroup Any method
	anyMethods = []string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodConnect,
		http.MethodTrace,
	}
)

// IRouter defines all router handle interface includes single and group router.
type IRouter interface {
	IRoutes
	Group(string, ...HandlerFunc) *RouterGroup
}

// IRoutes defines all router handle interface.
type IRoutes interface {
	Use(...HandlerFunc) IRoutes

	Handle(string, string, ...HandlerFunc) IRoutes
	Any(string, ...HandlerFunc) IRoutes
	GET(string, ...HandlerFunc) IRoutes
	POST(string, ...HandlerFunc) IRoutes
	DELETE(string, ...HandlerFunc) IRoutes
	PATCH(string, ...HandlerFunc) IRoutes
	PUT(string, ...HandlerFunc) IRoutes
	OPTIONS(string, ...HandlerFunc) IRoutes
	HEAD(string, ...HandlerFunc) IRoutes

	StaticFile(string, string, ...string) IRoutes
	Static(string, string, ...string) IRoutes
	StaticFS(string, string, http.FileSystem, ...string) IRoutes
	StaticFSFile(relativePath, filePath string, hfs http.FileSystem, cacheControls ...string) IRoutes
	StaticContent(relativePath string, data []byte, modtime time.Time, cacheControls ...string) IRoutes
}

// RouterGroup is used internally to configure router, a RouterGroup is associated with
// a prefix and an array of handlers (middleware).
type RouterGroup struct {
	Handlers HandlersChain
	basePath string
	engine   *Engine
	root     bool
}

// Use adds middleware to the group, see example code in GitHub.
func (group *RouterGroup) Use(middleware ...HandlerFunc) IRoutes {
	group.Handlers = append(group.Handlers, middleware...)
	return group.returnObj()
}

// Group creates a new router group. You should add all the routes that have common middlewares or the same path prefix.
// For example, all the routes that use a common middleware for authorization could be grouped.
func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		Handlers: group.combineHandlers(handlers),
		basePath: group.calculateAbsolutePath(relativePath),
		engine:   group.engine,
	}
}

// BasePath returns the base path of router group.
// For example, if v := router.Group("/rest/n/v1/api"), v.BasePath() is "/rest/n/v1/api".
func (group *RouterGroup) BasePath() string {
	return group.basePath
}

func (group *RouterGroup) handle(httpMethod, relativePath string, handlers HandlersChain) IRoutes {
	absolutePath := group.calculateAbsolutePath(relativePath)
	handlers = group.combineHandlers(handlers)
	group.engine.addRoute(httpMethod, absolutePath, handlers)
	return group.returnObj()
}

// Handle registers a new request handle and middleware with the given path and method.
// The last handler should be the real handler, the other ones should be middleware that can and should be shared among different routes.
// See the example code in GitHub.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (group *RouterGroup) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes {
	if matched := regEnLetter.MatchString(httpMethod); !matched {
		panic("http method " + httpMethod + " is not valid")
	}
	return group.handle(httpMethod, relativePath, handlers)
}

// POST is a shortcut for router.Handle("POST", path, handle).
func (group *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodPost, relativePath, handlers)
}

// GET is a shortcut for router.Handle("GET", path, handle).
func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodGet, relativePath, handlers)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func (group *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodDelete, relativePath, handlers)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func (group *RouterGroup) PATCH(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodPatch, relativePath, handlers)
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func (group *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodPut, relativePath, handlers)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func (group *RouterGroup) OPTIONS(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodOptions, relativePath, handlers)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func (group *RouterGroup) HEAD(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodHead, relativePath, handlers)
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (group *RouterGroup) Any(relativePath string, handlers ...HandlerFunc) IRoutes {
	for _, method := range anyMethods {
		group.handle(method, relativePath, handlers)
	}

	return group.returnObj()
}

func getCacheControlWriter(c *Context, cacheControl string) http.ResponseWriter {
	if cacheControl == "" {
		return c.Writer
	}
	h := map[string]string{"Cache-Control": cacheControl}
	return httpx.NewHeaderAppender(c.Writer, h)
}

// StaticFile registers a single route in order to serve a single file of the local filesystem.
// router.StaticFile("favicon.ico", "./resources/favicon.ico", "public, max-age=31536000")
func (group *RouterGroup) StaticFile(relativePath, localPath string, cacheControls ...string) IRoutes {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}

	var handler func(c *Context)

	cacheControl := strings.Join(cacheControls, ", ")
	if cacheControl == "" {
		handler = func(c *Context) {
			c.File(localPath)
		}
	} else {
		handler = func(c *Context) {
			ccw := getCacheControlWriter(c, cacheControl)
			http.ServeFile(ccw, c.Request, localPath)
		}
	}
	group.GET(relativePath, handler)
	group.HEAD(relativePath, handler)
	return group.returnObj()
}

// Static serves files from the given file system root.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use :
//     router.Static("/static", "/var/www")
func (group *RouterGroup) Static(relativePath, root string, cacheControls ...string) IRoutes {
	return group.StaticFS(relativePath, "", http.Dir(root), cacheControls...)
}

// StaticFS works just like `Static()` but a custom `http.FileSystem` can be used instead.
func (group *RouterGroup) StaticFS(relativePath string, localPath string, hfs http.FileSystem, cacheControls ...string) IRoutes {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}

	prefix := path.Join(group.BasePath(), relativePath)
	fileServer := http.FileServer(hfs)
	if prefix == "" || prefix == "/" {
		fileServer = httpx.AppendPrefix(localPath, fileServer)
	} else if localPath == "" || localPath == "." {
		fileServer = http.StripPrefix(prefix, fileServer)
	} else {
		fileServer = httpx.URLReplace(prefix, localPath, fileServer)
	}

	var handler func(c *Context)

	cacheControl := strings.Join(cacheControls, ", ")
	if cacheControl == "" {
		handler = func(c *Context) {
			fileServer.ServeHTTP(c.Writer, c.Request)
		}
	} else {
		handler = func(c *Context) {
			ccw := getCacheControlWriter(c, cacheControl)
			fileServer.ServeHTTP(ccw, c.Request)
		}
	}

	urlPattern := path.Join(relativePath, "/*path")

	// Register GET and HEAD handlers
	group.GET(urlPattern, handler)
	group.HEAD(urlPattern, handler)
	return group.returnObj()
}

// StaticFSFile registers a single route in order to serve a single file of the filesystem.
// ginfile.StaticFSFile(gin, "favicon.ico", "./resources/favicon.ico", hfs, gin.Public1Year)
func (group *RouterGroup) StaticFSFile(relativePath, filePath string, hfs http.FileSystem, cacheControls ...string) IRoutes {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}

	cacheControl := strings.Join(cacheControls, ", ")
	handler := func(c *Context) {
		defer func(old string) {
			c.Request.URL.Path = old
		}(c.Request.URL.Path)

		c.Request.URL.Path = filePath

		var hrw http.ResponseWriter
		if cacheControl == "" {
			hrw = c.Writer
		} else {
			hrw = getCacheControlWriter(c, cacheControl)
		}

		http.FileServer(hfs).ServeHTTP(hrw, c.Request)
	}

	group.GET(relativePath, handler)
	group.HEAD(relativePath, handler)
	return group.returnObj()
}

// StaticContent registers a single route in order to serve a single file of the data.
// //go:embed favicon.ico
// var favicon []byte
// group.StaticContent("favicon.ico", favicon, time.Now(), "public")
func (group *RouterGroup) StaticContent(relativePath string, data []byte, modtime time.Time, cacheControls ...string) IRoutes {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}

	if modtime.IsZero() {
		modtime = time.Now()
	}

	cacheControl := strings.Join(cacheControls, ", ")
	handler := func(c *Context) {
		if cacheControl != "" {
			c.Header("Cache-Control", cacheControl)
		}
		name := filepath.Base(c.Request.URL.Path)
		http.ServeContent(c.Writer, c.Request, name, modtime, bytes.NewReader(data))
	}
	group.GET(relativePath, handler)
	group.HEAD(relativePath, handler)
	return group.returnObj()
}

func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(group.Handlers) + len(handlers)
	assert1(finalSize < int(abortIndex), "too many handlers")
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
}

func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return joinPaths(group.basePath, relativePath)
}

func (group *RouterGroup) returnObj() IRoutes {
	if group.root {
		return group.engine
	}
	return group
}
