package ginfile

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Public1Year Cache-Control: public, max-age=31536000
const Public1Year = "public, max-age=31536000"

// Static serves files from the given file system root.
func Static(g *gin.RouterGroup, relativePath, localPath, cacheControl string) {
	StaticFS(g, relativePath, http.Dir(localPath), cacheControl)
}

// StaticFile registers a single route in order to serve a single file of the local filesystem.
// ginfile.StaticFSFile(gin, "favicon.ico", "./resources/favicon.ico", "public")
func StaticFile(g *gin.RouterGroup, relativePath, localPath, cacheControl string) {
	dir := filepath.Dir(localPath)
	StaticFSFile(g, relativePath, localPath, http.Dir(dir), cacheControl)
}

// StaticFS works just like `Static()` but a custom `http.FileSystem` can be used instead.
// Gin by default user: gin.Dir()
func StaticFS(g *gin.RouterGroup, relativePath string, hfs http.FileSystem, cacheControl string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}

	prefix := path.Join(g.BasePath(), relativePath)
	fileServer := http.StripPrefix(prefix, http.FileServer(hfs))

	handler := func(c *gin.Context) {
		file := c.Param("path")
		// Check if file exists and/or if we have permission to access it
		f, err := hfs.Open(file)
		if err != nil {
			msg, code := toHTTPError(err)
			http.Error(c.Writer, msg, code)
			return
		}
		f.Close()

		if cacheControl != "" {
			c.Header("Cache-Control", cacheControl)
		}
		fileServer.ServeHTTP(c.Writer, c.Request)
	}

	urlPattern := path.Join(relativePath, "/*path")

	// Register GET and HEAD handlers
	g.GET(urlPattern, handler)
	g.HEAD(urlPattern, handler)
}

// StaticFSFile registers a single route in order to serve a single file of the filesystem.
// ginfile.StaticFSFile(gin, "favicon.ico", "./resources/favicon.ico", hfs, ginfile.Public1Year)
func StaticFSFile(g *gin.RouterGroup, relativePath, filePath string, hfs http.FileSystem, cacheControl string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}
	handler := func(c *gin.Context) {
		name := filepath.Base(relativePath)
		f, err := hfs.Open(filePath)
		if err != nil {
			msg, code := toHTTPError(err)
			http.Error(c.Writer, msg, code)
			return
		}
		f.Close()

		if cacheControl != "" {
			c.Header("Cache-Control", cacheControl)
		}

		t := time.Time{}
		if fi, err := f.Stat(); err == nil {
			t = fi.ModTime()
		}
		http.ServeContent(c.Writer, c.Request, name, t, f)
	}
	g.GET(relativePath, handler)
	g.HEAD(relativePath, handler)
}

// toHTTPError returns a non-specific HTTP error message and status code
// for a given non-nil error value. It's important that toHTTPError does not
// actually return err.Error(), since msg and httpStatus are returned to users,
// and historically Go's ServeContent always returned just "404 Not Found" for
// all errors. We don't want to start leaking information in error messages.
func toHTTPError(err error) (msg string, httpStatus int) {
	if os.IsNotExist(err) {
		return "404 page not found", http.StatusNotFound
	}
	if os.IsPermission(err) {
		return "403 Forbidden", http.StatusForbidden
	}
	// Default:
	return "500 Internal Server Error", http.StatusInternalServerError
}
