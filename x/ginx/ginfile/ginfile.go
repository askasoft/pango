package ginfile

import (
	"bytes"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pandafw/pango/net/httpx"
)

// Public1Year Cache-Control: public, max-age=31536000
const Public1Year = "public, max-age=31536000"

func getCacheControlWriter(c *gin.Context, cacheControl string) http.ResponseWriter {
	if cacheControl == "" {
		return c.Writer
	}
	h := map[string]string{"Cache-Control": cacheControl}
	return httpx.NewHeaderAppender(c.Writer, h)
}

// Static serves files from the given file system root.
func Static(g *gin.RouterGroup, relativePath, localPath, cacheControl string) {
	StaticFS(g, relativePath, "", http.Dir(localPath), cacheControl)
}

// StaticFile registers a single route in order to serve a single file of the local filesystem.
// ginfile.StaticFSFile(gin, "favicon.ico", "./resources/favicon.ico", "public")
func StaticFile(g *gin.RouterGroup, relativePath, localPath, cacheControl string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}

	handler := func(c *gin.Context) {
		ccw := getCacheControlWriter(c, cacheControl)
		http.ServeFile(ccw, c.Request, localPath)
	}
	g.GET(relativePath, handler)
	g.HEAD(relativePath, handler)
}

// StaticFS works just like `Static()` but a custom `http.FileSystem` can be used instead.
func StaticFS(g *gin.RouterGroup, relativePath string, localPath string, hfs http.FileSystem, cacheControl string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}

	prefix := path.Join(g.BasePath(), relativePath)
	fileServer := http.FileServer(hfs)
	if prefix == "" || prefix == "/" {
		fileServer = httpx.AppendPrefix(localPath, fileServer)
	} else if localPath == "" || localPath == "." {
		fileServer = http.StripPrefix(prefix, fileServer)
	} else {
		fileServer = httpx.URLReplace(prefix, localPath, fileServer)
	}

	handler := func(c *gin.Context) {
		ccw := getCacheControlWriter(c, cacheControl)
		fileServer.ServeHTTP(ccw, c.Request)
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
		defer func(old string) {
			c.Request.URL.Path = old
		}(c.Request.URL.Path)

		c.Request.URL.Path = filePath
		ccw := getCacheControlWriter(c, cacheControl)

		http.FileServer(hfs).ServeHTTP(ccw, c.Request)
	}

	g.GET(relativePath, handler)
	g.HEAD(relativePath, handler)
}

// StaticContent registers a single route in order to serve a single file of the data.
// //go:embed favicon.ico
// var favicon []byte
// ginfile.StaticContent(gin, "favicon.ico", favicon, "public")
func StaticContent(g *gin.RouterGroup, relativePath string, data []byte, cacheControl string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}

	t := time.Time{}
	handler := func(c *gin.Context) {
		if cacheControl != "" {
			c.Header("Cache-Control", cacheControl)
		}
		name := filepath.Base(c.Request.URL.Path)
		http.ServeContent(c.Writer, c.Request, name, t, bytes.NewReader(data))
	}
	g.GET(relativePath, handler)
	g.HEAD(relativePath, handler)
}
