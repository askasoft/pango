package xmw

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/askasoft/pango/xin"
)

func newHTTPSRedirectorServer(sh *HTTPSRedirector) *xin.Engine {
	r := xin.New()
	r.Use(sh.Handle)
	r.GET("/foo", func(c *xin.Context) {
		c.String(200, "bar")
	})
	return r
}

func TestHTTPSRedirectorDisabled(t *testing.T) {
	s := newHTTPSRedirectorServer(&HTTPSRedirector{
		disabled: true,
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.RemoteAddr = "127.0.0.1:1000"
	req.Host = "www.example.com"
	req.URL.Scheme = "http"

	s.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
}

func TestHTTPSRedirectorNoConfig(t *testing.T) {
	s := newHTTPSRedirectorServer(&HTTPSRedirector{})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.RemoteAddr = "127.0.0.1:1000"
	req.Host = "www.example.com"
	req.URL.Scheme = "http"

	s.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusMovedPermanently)
	expect(t, res.Header().Get("Location"), "https://www.example.com/foo")
}

func TestHTTPSRedirectorWithHost(t *testing.T) {
	s := newHTTPSRedirectorServer(&HTTPSRedirector{
		SSLHost: "secure.example.com",
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.RemoteAddr = "127.0.0.1:1000"
	req.Host = "www.example.com"
	req.URL.Scheme = "http"

	s.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusMovedPermanently)
	expect(t, res.Header().Get("Location"), "https://secure.example.com/foo")
}

func TestHTTPSRedirectorNoProxyHeaders(t *testing.T) {
	s := newHTTPSRedirectorServer(&HTTPSRedirector{})
	s.SSLProxyHeaders = map[string]string{}

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.RemoteAddr = "127.0.0.1:1000"
	req.Host = "www.example.com"
	req.URL.Scheme = "http"
	req.Header.Add("X-Forwarded-Proto", "https")

	s.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusMovedPermanently)
	expect(t, res.Header().Get("Location"), "https://www.example.com/foo")
}

func TestHTTPSRedirectorWithProxyHeaders(t *testing.T) {
	s := newHTTPSRedirectorServer(NewHTTPSRedirector())

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.RemoteAddr = "127.0.0.1:1000"
	req.Host = "www.example.com"
	req.URL.Scheme = "http"
	req.Header.Add("X-Forwarded-Proto", "https")

	s.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
}

func TestHTTPSRedirectorWithProxyHeadersDisabled(t *testing.T) {
	s := newHTTPSRedirectorServer(&HTTPSRedirector{
		disabled: true,
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.RemoteAddr = "127.0.0.1:1000"
	req.Host = "www.example.com"
	req.URL.Scheme = "http"
	req.Header.Add("X-Forwarded-Proto", "http")

	s.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
}

func TestHTTPSRedirectorWithProxyAndHost(t *testing.T) {
	s := newHTTPSRedirectorServer(&HTTPSRedirector{
		SSLHost: "secure.example.com",
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.RemoteAddr = "127.0.0.1:1000"
	req.Host = "www.example.com"
	req.URL.Scheme = "http"
	req.Header.Add("X-Forwarded-Proto", "https")

	s.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
}

func TestHTTPSRedirectorCustomBadProxyAndHost(t *testing.T) {
	s := newHTTPSRedirectorServer(&HTTPSRedirector{
		SSLHost: "secure.example.com",
	})
	s.SSLProxyHeaders = map[string]string{"X-Forwarded-Proto": "bad"}

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.RemoteAddr = "127.0.0.1:1000"
	req.Host = "www.example.com"
	req.URL.Scheme = "http"
	req.Header.Add("X-Forwarded-Proto", "https")

	s.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusMovedPermanently)
	expect(t, res.Header().Get("Location"), "https://secure.example.com/foo")
}

func TestHTTPSRedirectorCustomBadProxyAndHostWithTempRedirect(t *testing.T) {
	s := newHTTPSRedirectorServer(&HTTPSRedirector{
		SSLHost:           "secure.example.com",
		TemporaryRedirect: true,
	})
	s.SSLProxyHeaders = map[string]string{"X-Forwarded-Proto": "bad"}

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.RemoteAddr = "127.0.0.1:1000"
	req.Host = "www.example.com"
	req.URL.Scheme = "http"
	req.Header.Add("X-Forwarded-Proto", "https")

	s.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusTemporaryRedirect)
	expect(t, res.Header().Get("Location"), "https://secure.example.com/foo")
}

/* Test Helpers */
func expect(t *testing.T, a any, b any) {
	if a != b {
		t.Errorf("Expected [%v] (type %v) - Got [%v] (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
