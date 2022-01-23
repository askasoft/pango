package gmw

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/pandafw/pango/gin"
	"github.com/pandafw/pango/str"
)

func dumpPerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	fmt.Println(strings.Repeat("-", 78))
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func dumpAssertContains(t *testing.T, msg string, body string, ss ...string) {
	for _, s := range ss {
		if !str.Contains(body, s) {
			t.Errorf(`%s http dump does not contains %q`, msg, s)
		}
	}
}

func TestHttpDump(t *testing.T) {
	router := gin.New()

	buffer := new(bytes.Buffer)
	writer := io.MultiWriter(buffer, os.Stdout)
	router.Use(NewDumper(writer).Handler())

	router.Any("/example", func(c *gin.Context) {
		c.String(http.StatusOK, c.Request.URL.String())
	})

	buffer.Reset()
	dumpPerformRequest(router, "GET", "/example?a=100")
	dumpAssertContains(t, "GET /example?a=100", buffer.String(), "GET /example?a=100 HTTP/1.1", "HTTP/1.1 200 OK")

	buffer.Reset()
	dumpPerformRequest(router, "POST", "/example")
	dumpAssertContains(t, "POST /example", buffer.String(), "POST /example HTTP/1.1", "HTTP/1.1 200 OK")

	buffer.Reset()
	dumpPerformRequest(router, "PUT", "/example")
	dumpAssertContains(t, "PUT /example", buffer.String(), "PUT /example HTTP/1.1", "HTTP/1.1 200 OK")

	buffer.Reset()
	dumpPerformRequest(router, "DELETE", "/example")
	dumpAssertContains(t, "DELETE /example", buffer.String(), "DELETE /example HTTP/1.1", "HTTP/1.1 200 OK")

	buffer.Reset()
	dumpPerformRequest(router, "PATCH", "/example")
	dumpAssertContains(t, "PATCH /example", buffer.String(), "PATCH /example HTTP/1.1", "HTTP/1.1 200 OK")

	buffer.Reset()
	dumpPerformRequest(router, "HEAD", "/example")
	dumpAssertContains(t, "HEAD /example", buffer.String(), "HEAD /example HTTP/1.1", "HTTP/1.1 200 OK")

	buffer.Reset()
	dumpPerformRequest(router, "OPTIONS", "/example")
	dumpAssertContains(t, "OPTIONS /example", buffer.String(), "OPTIONS /example HTTP/1.1", "HTTP/1.1 200 OK")

	buffer.Reset()
	dumpPerformRequest(router, "GET", "/notfound")
	dumpAssertContains(t, "GET /notfound", buffer.String(), "GET /notfound HTTP/1.1", "HTTP/1.1 404 Not Found")
}
