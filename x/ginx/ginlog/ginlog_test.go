package ginlog

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func TestTextLog(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := gin.New()

	writer := io.MultiWriter(buffer, os.Stdout)
	router.Use(New(writer, DefaultTextLogFormat).Handler())

	router.Any("/example", func(c *gin.Context) {})

	buffer.Reset()
	performRequest(router, "GET", "/example?a=100")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/example")
	assert.Contains(t, buffer.String(), "a=100")

	buffer.Reset()
	performRequest(router, "POST", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "POST")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "PUT", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "PUT")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "DELETE", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "DELETE")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "PATCH", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "PATCH")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "HEAD", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "HEAD")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "OPTIONS", "/example")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "OPTIONS")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	performRequest(router, "GET", "/notfound")
	assert.Contains(t, buffer.String(), "404")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/notfound")
}

func TestJSONLog(t *testing.T) {
	result := make(map[string]interface{})
	buffer := new(bytes.Buffer)
	router := gin.New()

	writer := io.MultiWriter(buffer, os.Stdout)
	router.Use(New(writer, DefaultJSONLogFormat).Handler())

	router.Any("/example", func(c *gin.Context) {})

	buffer.Reset()
	performRequest(router, "GET", "/example?a=100")
	json.Unmarshal(buffer.Bytes(), &result)
	assert.Equal(t, result["status"], float64(200))
	assert.Equal(t, result["method"], "GET")
	assert.Equal(t, result["url"], "/example?a=100")

	buffer.Reset()
	performRequest(router, "POST", "/example")
	json.Unmarshal(buffer.Bytes(), &result)
	assert.Equal(t, result["status"], float64(200))
	assert.Equal(t, result["method"], "POST")
	assert.Equal(t, result["url"], "/example")

	buffer.Reset()
	performRequest(router, "PUT", "/example")
	json.Unmarshal(buffer.Bytes(), &result)
	assert.Equal(t, result["status"], float64(200))
	assert.Equal(t, result["method"], "PUT")
	assert.Equal(t, result["url"], "/example")

	buffer.Reset()
	performRequest(router, "DELETE", "/example")
	json.Unmarshal(buffer.Bytes(), &result)
	assert.Equal(t, result["status"], float64(200))
	assert.Equal(t, result["method"], "DELETE")
	assert.Equal(t, result["url"], "/example")

	buffer.Reset()
	performRequest(router, "PATCH", "/example")
	json.Unmarshal(buffer.Bytes(), &result)
	assert.Equal(t, result["status"], float64(200))
	assert.Equal(t, result["method"], "PATCH")
	assert.Equal(t, result["url"], "/example")

	buffer.Reset()
	performRequest(router, "HEAD", "/example")
	json.Unmarshal(buffer.Bytes(), &result)
	assert.Equal(t, result["status"], float64(200))
	assert.Equal(t, result["method"], "HEAD")
	assert.Equal(t, result["url"], "/example")

	buffer.Reset()
	performRequest(router, "OPTIONS", "/example")
	json.Unmarshal(buffer.Bytes(), &result)
	assert.Equal(t, result["status"], float64(200))
	assert.Equal(t, result["method"], "OPTIONS")
	assert.Equal(t, result["url"], "/example")

	buffer.Reset()
	performRequest(router, "GET", "/notfound")
	json.Unmarshal(buffer.Bytes(), &result)
	assert.Equal(t, result["status"], float64(404))
	assert.Equal(t, result["method"], "GET")
	assert.Equal(t, result["url"], "/notfound")
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
