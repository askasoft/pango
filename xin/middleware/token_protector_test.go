package middleware

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/askasoft/pango/xin"
)

func TestTokenProtectorFail(t *testing.T) {
	w := httptest.NewRecorder()
	router := xin.New()
	tp := NewTokenProtector("1234567890123456")
	router.Use(tp.Handle)
	router.POST("/", func(c *xin.Context) {
		c.String(200, "OK")
	})

	req, _ := http.NewRequest("POST", "/", nil)
	req.AddCookie(&http.Cookie{
		Name: TokenCookieName,
	})
	router.ServeHTTP(w, req)

	val := w.Result().StatusCode
	exp := http.StatusForbidden
	if val != exp {
		t.Errorf("%v = %v, want %v", req.URL.String(), val, exp)
	}
}

func TestTokenProtectorOK(t *testing.T) {
	w := httptest.NewRecorder()
	router := xin.New()
	tp := NewTokenProtector("1234567890123456")
	router.Use(tp.Handle)
	router.GET("/", func(c *xin.Context) {
		c.String(200, tp.RefreshToken(c))
	})
	router.POST("/", func(c *xin.Context) {
		c.String(200, "OK")
	})

	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)
	ts := w.Body.String()

	req2, _ := http.NewRequest("POST", "/", nil)
	for _, c := range w.Result().Cookies() {
		req2.AddCookie(c)
	}
	req2.PostForm = url.Values{}
	req2.PostForm.Add("_token_", ts)

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Result().StatusCode != http.StatusOK {
		t.Errorf("%v = %v, want %v", req2.URL.String(), w2.Result().StatusCode, http.StatusOK)
	}

	val := w2.Body.String()
	exp := "OK"
	if val != exp {
		t.Errorf("%v = %q, want %q", req2.URL.String(), val, exp)
	}
}
