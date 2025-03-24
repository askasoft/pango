package xmw

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/askasoft/pango/xin"
)

func TestRequestSizeLimiter(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", strings.NewReader(`{"a": "1345678901"}`))

	w := httptest.NewRecorder()
	router := xin.New()
	router.Use(NewRequestSizeLimiter(10).Handler())
	router.POST("/", func(c *xin.Context) {
		m := map[string]string{}
		if err := c.MustBindJSON(&m); err == nil {
			c.String(200, "OK")
		}
	})

	router.ServeHTTP(w, req)

	val := w.Result().StatusCode
	exp := http.StatusRequestEntityTooLarge
	if val != exp {
		t.Errorf("%v = %d, want %d", req.URL.String(), val, exp)
	}
}
