package xmw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/askasoft/pango/xin"
)

func TestResponseHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	w := httptest.NewRecorder()
	router := xin.New()
	router.Use(NewResponseHeader(map[string]string{"Access-Control-Allow-Origin": "*"}).Handler())
	router.GET("/", func(c *xin.Context) {
		c.String(200, "OK")
	})

	router.ServeHTTP(w, req)

	val := w.Header().Get("Access-Control-Allow-Origin")
	exp := "*"
	if val != exp {
		t.Errorf("%v = %q, want %q", req.URL.String(), val, exp)
	}
}
