package xin

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRouterInvalidStatic(t *testing.T) {
	router := New()
	assert.Panics(t, func() {
		Static(router, "/path/:param", "/")
	})

	assert.Panics(t, func() {
		Static(router, "/path/*param", "/")
	})
}

func TestRouterInvalidStaticFile(t *testing.T) {
	router := New()
	assert.Panics(t, func() {
		StaticFile(router, "/path/:param", "favicon.ico")
	})

	assert.Panics(t, func() {
		StaticFile(router, "/path/*param", "favicon.ico")
	})
}

func TestRouterStaticFSNotFound(t *testing.T) {
	router := New()
	StaticFS(router, "/", http.Dir("/thisreallydoesntexist/"), "")
	router.NoRoute(func(c *Context) {
		c.String(404, "non existent")
	})

	w := performRequest(router, http.MethodGet, "/nonexistent")
	assert.Equal(t, 404, w.Result().StatusCode)
	// assert.Equal(t, "non existent", w.Body.String())

	w = performRequest(router, http.MethodHead, "/nonexistent")
	assert.Equal(t, 404, w.Result().StatusCode)
	// assert.Equal(t, "non existent", w.Body.String())
}

func TestRouterStaticFSFileNotFound(t *testing.T) {
	router := New()

	StaticFS(router, "/", http.Dir("."), "")

	assert.NotPanics(t, func() {
		performRequest(router, http.MethodGet, "/nonexistent")
	})
}

// Reproduction test for the bug of issue #1805
func TestMiddlewareCalledOnceByRouterStaticFSNotFound(t *testing.T) {
	router := New()

	// Middleware must be called just only once by per request.
	middlewareCalledNum := 0
	router.Use(func(c *Context) {
		middlewareCalledNum++
	})

	StaticFS(router, "/", http.Dir("/thisreallydoesntexist/"), "")

	// First access
	performRequest(router, http.MethodGet, "/nonexistent")
	assert.Equal(t, 1, middlewareCalledNum)

	// Second access
	performRequest(router, http.MethodHead, "/nonexistent")
	assert.Equal(t, 2, middlewareCalledNum)
}

// TestHandleStaticFile - ensure the static file handles properly
func TestRouteStaticFile(t *testing.T) {
	// SETUP file
	testRoot, _ := os.Getwd()
	f, err := os.CreateTemp(testRoot, "")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())
	_, err = f.WriteString("Xin Web Framework")
	assert.NoError(t, err)
	f.Close()

	dir, filename := filepath.Split(f.Name())

	// SETUP xin
	router := New()
	Static(router, "/using_static", dir)
	StaticFile(router, "/result", f.Name())

	w := performRequest(router, http.MethodGet, "/using_static/"+filename)
	w2 := performRequest(router, http.MethodGet, "/result")

	assert.Equal(t, w, w2)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Xin Web Framework", w.Body.String())
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))

	w3 := performRequest(router, http.MethodHead, "/using_static/"+filename)
	w4 := performRequest(router, http.MethodHead, "/result")

	assert.Equal(t, w3, w4)
	assert.Equal(t, http.StatusOK, w3.Code)
}

// TestHandleStaticDir - ensure the root/sub dir handles properly
func TestRouteStaticListingDir(t *testing.T) {
	router := New()

	StaticFS(router, "/", http.Dir("./"), "")

	w := performRequest(router, http.MethodGet, "/")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "xin.go")
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
}

// TestHandleHeadToDir - ensure the root/sub dir handles properly
// func TestRouteStaticNoListing(t *testing.T) {
// 	router := New()
// 	router.Static("/", "./")

// 	w := performRequest(router, http.MethodGet, "/")

// 	assert.Equal(t, http.StatusNotFound, w.Code)
// 	assert.NotContains(t, w.Body.String(), "xin.go")
// }

func TestRouterMiddlewareAndStatic(t *testing.T) {
	router := New()
	static := router.Group("/", func(c *Context) {
		c.Writer.Header().Add("Last-Modified", "Mon, 02 Jan 2006 15:04:05 MST")
		c.Writer.Header().Add("Expires", "Mon, 02 Jan 2006 15:04:05 MST")
		c.Writer.Header().Add("X-XIN", "Xin Framework")
	})

	Static(static, "/", "./")

	w := performRequest(router, http.MethodGet, "/xin.go")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "package xin")
	// Content-Type='text/plain; charset=utf-8' when go version <= 1.16,
	// else, Content-Type='text/x-go; charset=utf-8'
	assert.NotEqual(t, "", w.Header().Get("Content-Type"))
	assert.NotEqual(t, w.Header().Get("Last-Modified"), "Mon, 02 Jan 2006 15:04:05 MST")
	assert.Equal(t, "Mon, 02 Jan 2006 15:04:05 MST", w.Header().Get("Expires"))
	assert.Equal(t, "Xin Framework", w.Header().Get("x-XIN"))
}

//go:embed testdata
var testdata embed.FS

//go:embed testdata/files/file1.txt
var file1 []byte

func testGetFile(t *testing.T, r *Engine, path string, cache string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, nil)
	r.ServeHTTP(w, req)

	if 200 != w.Code {
		t.Errorf("w.Code = %v, want %v", w.Code, 200)
	}
	if cache != w.Header().Get("Cache-Control") {
		t.Errorf("Header[Cache-Control] = %v, want %v", w.Header().Get("Cache-Control"), cache)
	}
	if filepath.Base(path) != w.Body.String() {
		t.Errorf(`Body = %v, want %v`, w.Body.String(), filepath.Base(path))
	}
}

func TestRouterStatic(t *testing.T) {
	r := Default()
	Static(r, "/", "testdata", NewCacheControlSetter("private").Handler())
	testGetFile(t, r, "/root1.txt", "private")
	testGetFile(t, r, "/files/file1.txt", "private")
}

func TestRouterStaticFile(t *testing.T) {
	r := Default()
	StaticFile(r, "/root1.txt", "testdata/root1.txt", NewCacheControlSetter("public").Handler())
	testGetFile(t, r, "/root1.txt", "public")
}

func TestRouterStaticFS_AppendPrefix(t *testing.T) {
	r := Default()
	StaticFS(r, "", FS(testdata), "/testdata", NewCacheControlSetter("private").Handler())
	testGetFile(t, r, "/root1.txt", "private")
	testGetFile(t, r, "/files/file1.txt", "private")
}

func TestRouterStaticFS_AppendPrefix2(t *testing.T) {
	r := Default()
	StaticFS(r, "/", FS(testdata), "/testdata", NewCacheControlSetter("private").Handler())
	testGetFile(t, r, "/root1.txt", "private")
	testGetFile(t, r, "/files/file1.txt", "private")
}

func TestRouterStaticFS_StripPrefix(t *testing.T) {
	r := Default()
	g := r.Group("/web")

	StaticFS(g, "/", FS(testdata), "", NewCacheControlSetter("private").Handler())
	testGetFile(t, r, "/web/testdata/root1.txt", "private")
	testGetFile(t, r, "/web/testdata/files/file1.txt", "private")
}

func TestRouterStaticFS_URLReplace(t *testing.T) {
	r := New()
	StaticFS(r, "/data", FS(testdata), "/testdata", NewCacheControlSetter("private").Handler())
	testGetFile(t, r, "/data/root1.txt", "private")
	testGetFile(t, r, "/data/files/file1.txt", "private")
}

func TestRouterStaticFS_URLReplace2(t *testing.T) {
	r := Default()
	g := r.Group("web")
	StaticFS(g, "/", FS(testdata), "/testdata", NewCacheControlSetter("private").Handler())
	testGetFile(t, r, "/web/root1.txt", "private")
	testGetFile(t, r, "/web/files/file1.txt", "private")
}

func TestRouterStaticFSFile(t *testing.T) {
	r := Default()
	StaticFSFile(r, "/root1.txt", FS(testdata), "testdata/root1.txt", NewCacheControlSetter("public").Handler())
	testGetFile(t, r, "/root1.txt", "public")
}

func TestRouterStaticContent(t *testing.T) {
	r := Default()
	StaticContent(r, "/files/file1.txt", file1, time.Now(), NewCacheControlSetter("no-store").Handler())
	testGetFile(t, r, "/files/file1.txt", "no-store")
}
