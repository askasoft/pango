package ginfile

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata
var testdata embed.FS

//go:embed testdata/d1/d1f1.txt
var d1f1 []byte

func testGetFile(t *testing.T, r *gin.Engine, path string, cache string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, cache, w.Header().Get("Cache-Control"))
	assert.Equal(t, filepath.Base(path), w.Body.String())
}

func TestStatic(t *testing.T) {
	r := gin.Default()
	Static(&r.RouterGroup, "/", "testdata", "private")
	testGetFile(t, r, "/r1.txt", "private")
	testGetFile(t, r, "/d1/d1f1.txt", "private")
}

func TestStaticFile(t *testing.T) {
	r := gin.Default()
	StaticFile(&r.RouterGroup, "/r1.txt", "testdata/r1.txt", "public")
	testGetFile(t, r, "/r1.txt", "public")
}

func TestStaticFS_AppendPrefix(t *testing.T) {
	r := gin.Default()
	StaticFS(&r.RouterGroup, "", "/testdata", http.FS(testdata), "private")
	testGetFile(t, r, "/r1.txt", "private")
	testGetFile(t, r, "/d1/d1f1.txt", "private")
}

func TestStaticFS_AppendPrefix2(t *testing.T) {
	r := gin.Default()
	StaticFS(&r.RouterGroup, "/", "/testdata", http.FS(testdata), "private")
	testGetFile(t, r, "/r1.txt", "private")
	testGetFile(t, r, "/d1/d1f1.txt", "private")
}

func TestStaticFS_StripPrefix(t *testing.T) {
	r := gin.Default()
	g := r.Group("/web")

	StaticFS(g, "/", "", http.FS(testdata), "private")
	testGetFile(t, r, "/web/testdata/r1.txt", "private")
	testGetFile(t, r, "/web/testdata/d1/d1f1.txt", "private")
}

func TestStaticFS_URLReplace(t *testing.T) {
	r := gin.Default()
	StaticFS(&r.RouterGroup, "/data", "/testdata", http.FS(testdata), "private")
	testGetFile(t, r, "/data/r1.txt", "private")
	testGetFile(t, r, "/data/d1/d1f1.txt", "private")
}

func TestStaticFS_URLReplace2(t *testing.T) {
	r := gin.Default()
	g := r.Group("web")
	StaticFS(g, "/", "/testdata", http.FS(testdata), "private")
	testGetFile(t, r, "/web/r1.txt", "private")
	testGetFile(t, r, "/web/d1/d1f1.txt", "private")
}

func TestStaticFS_URLReplace4(t *testing.T) {
	r := gin.Default()
	g := r.Group("/web")

	StaticFS(g, "/", "/testdata", http.FS(testdata), "private")
	testGetFile(t, r, "/web/testdata/r1.txt", "private")
	testGetFile(t, r, "/web/testdata/d1/d1f1.txt", "private")
}

func TestStaticFSFile(t *testing.T) {
	r := gin.Default()
	StaticFSFile(&r.RouterGroup, "/r1.txt", "testdata/r1.txt", http.FS(testdata), "public")
	testGetFile(t, r, "/r1.txt", "public")
}

func TestStaticContent(t *testing.T) {
	r := gin.Default()
	StaticContent(&r.RouterGroup, "/d1/d1f1.txt", d1f1, "no-store")
	testGetFile(t, r, "/d1/d1f1.txt", "no-store")
}
