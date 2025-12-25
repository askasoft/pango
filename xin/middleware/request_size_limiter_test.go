package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin"
)

func TestRequestSizeLimiterJSON(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", strings.NewReader(`{"a": "1345678901"}`))

	w := httptest.NewRecorder()
	router := xin.New()
	router.Use(NewRequestSizeLimiter(10).Handle)
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

func TestRequestSizeLimiterFormOK(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", strings.NewReader(`a=abc`+str.Repeat("0123456789", 1024*1024)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	router := xin.New()
	router.Use(NewRequestSizeLimiter(100 * 1024 * 1024).Handle)
	router.POST("/", func(c *xin.Context) {
		a := c.PostForm("a")
		if a == "" {
			c.String(500, "NG")
		} else {
			c.String(200, "OK")
		}
	})

	router.ServeHTTP(w, req)

	val := w.Result().StatusCode
	exp := http.StatusOK
	if val != exp {
		t.Errorf("%v = %d, want %d", req.URL.String(), val, exp)
	}
}

func TestRequestSizeLimiterFormNG1(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", strings.NewReader(`a=abc`+str.Repeat("0123456789", 1024*1024)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	router := xin.New()
	router.Use(NewRequestSizeLimiter(1024 * 1024).Handle)
	router.POST("/", func(c *xin.Context) {
		a := c.PostForm("a")
		if a == "" {
			c.String(500, "NG")
		} else {
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

type strReader struct {
	sr *strings.Reader
}

func (sr *strReader) Read(b []byte) (n int, err error) {
	return sr.sr.Read(b)
}

func TestRequestSizeLimiterFormNG2(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", &strReader{strings.NewReader(`a=abc` + str.Repeat("0123456789", 1024*1024))})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	router := xin.New()
	router.Use(NewRequestSizeLimiter(1024 * 1024).Handle)
	router.POST("/", func(c *xin.Context) {
		a := c.PostForm("a")
		if a != "" {
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
