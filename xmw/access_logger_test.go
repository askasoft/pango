package xmw

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/pandafw/pango/str"
	"github.com/pandafw/pango/xin"
)

func logPerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func logAssertContains(t *testing.T, msg string, body string, ss ...string) {
	for _, s := range ss {
		if !str.Contains(body, s) {
			t.Errorf(`%s access log does not contains %q`, msg, s)
		}
	}
}

func TestTextLog(t *testing.T) {
	router := xin.New()

	buffer := new(bytes.Buffer)
	writer := io.MultiWriter(buffer, os.Stdout)
	router.Use(NewAccessLogger(writer, DefaultTextLogFormat).Handler())

	router.Any("/example", func(c *xin.Context) {})

	buffer.Reset()
	logPerformRequest(router, "GET", "/example?a=100")
	logAssertContains(t, "GET /example?a=100", buffer.String(), "200", "GET", "/example", "a=100")

	buffer.Reset()
	logPerformRequest(router, "POST", "/example")
	logAssertContains(t, "POST /example", buffer.String(), "200", "POST", "/example")

	buffer.Reset()
	logPerformRequest(router, "PUT", "/example")
	logAssertContains(t, "PUT /example", buffer.String(), "200", "PUT", "/example")

	buffer.Reset()
	logPerformRequest(router, "DELETE", "/example")
	logAssertContains(t, "DELETE /example", buffer.String(), "200", "DELETE", "/example")

	buffer.Reset()
	logPerformRequest(router, "PATCH", "/example")
	logAssertContains(t, "PATCH /example", buffer.String(), "200", "PATCH", "/example")

	buffer.Reset()
	logPerformRequest(router, "HEAD", "/example")
	logAssertContains(t, "HEAD /example", buffer.String(), "200", "HEAD", "/example")

	buffer.Reset()
	logPerformRequest(router, "OPTIONS", "/example")
	logAssertContains(t, "OPTIONS /example", buffer.String(), "200", "OPTIONS", "/example")

	buffer.Reset()
	logPerformRequest(router, "GET", "/notfound")
	logAssertContains(t, "GET /notfound", buffer.String(), "404", "GET", "/notfound")
}

func assertJSONResult(t *testing.T, result map[string]interface{}, sc int, method string, url string) {
	if result["status"] != float64(sc) {
		t.Errorf("status = %v, want %v", result["status"], sc)
	}
	if result["method"] != method {
		t.Errorf("method = %v, want %v", result["method"], sc)
	}
	if result["url"] != url {
		t.Errorf("url = %v, want %v", result["url"], sc)
	}
}

func TestJSONLog(t *testing.T) {
	result := make(map[string]interface{})
	buffer := new(bytes.Buffer)
	router := xin.New()

	writer := io.MultiWriter(buffer, os.Stdout)
	router.Use(NewAccessLogger(writer, DefaultJSONLogFormat).Handler())

	router.Any("/example", func(c *xin.Context) {})

	buffer.Reset()
	logPerformRequest(router, "GET", "/example?a=100")
	json.Unmarshal(buffer.Bytes(), &result)
	assertJSONResult(t, result, 200, "GET", "/example?a=100")

	buffer.Reset()
	logPerformRequest(router, "POST", "/example")
	json.Unmarshal(buffer.Bytes(), &result)
	assertJSONResult(t, result, 200, "POST", "/example")

	buffer.Reset()
	logPerformRequest(router, "PUT", "/example")
	json.Unmarshal(buffer.Bytes(), &result)
	assertJSONResult(t, result, 200, "PUT", "/example")

	buffer.Reset()
	logPerformRequest(router, "DELETE", "/example")
	json.Unmarshal(buffer.Bytes(), &result)
	assertJSONResult(t, result, 200, "DELETE", "/example")

	buffer.Reset()
	logPerformRequest(router, "PATCH", "/example")
	json.Unmarshal(buffer.Bytes(), &result)
	assertJSONResult(t, result, 200, "PATCH", "/example")

	buffer.Reset()
	logPerformRequest(router, "HEAD", "/example")
	json.Unmarshal(buffer.Bytes(), &result)
	assertJSONResult(t, result, 200, "HEAD", "/example")

	buffer.Reset()
	logPerformRequest(router, "OPTIONS", "/example")
	json.Unmarshal(buffer.Bytes(), &result)
	assertJSONResult(t, result, 200, "OPTIONS", "/example")

	buffer.Reset()
	logPerformRequest(router, "GET", "/notfound")
	json.Unmarshal(buffer.Bytes(), &result)
	assertJSONResult(t, result, 404, "GET", "/notfound")
}
