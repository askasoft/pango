package xmw

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/askasoft/pango/xin"
)

func assertGzipHeader(t *testing.T, rr *httptest.ResponseRecorder, sc int, hce, hvary string) {
	if sc != rr.Code {
		t.Errorf("rr.Code = %v, want %v", rr.Code, sc)
	}
	if rr.Header().Get("Content-Encoding") != hce {
		t.Errorf(`Header[Content-Encoding] = %v, want %v`, rr.Header().Get("Content-Encoding"), hce)
	}
	if rr.Header().Get("Vary") != hvary {
		t.Errorf(`Header[Vary] = %v, want %v`, rr.Header().Get("Vary"), hce)
	}
}

func assertGzipIgnore(t *testing.T, rr *httptest.ResponseRecorder, body string) {
	assertGzipHeader(t, rr, http.StatusOK, "", "")
	if rr.Body.String() != body {
		t.Errorf(`Body = %v, want %v`, rr.Body.String(), body)
	}
}

func assertGzipEnable(t *testing.T, rr *httptest.ResponseRecorder, body string) {
	assertGzipHeader(t, rr, http.StatusOK, "gzip", "Accept-Encoding")

	if len(body) == rr.Body.Len() {
		t.Errorf("len(body) = rr.Body.Len() = %v", len(body))
	}

	gr, err := gzip.NewReader(rr.Body)
	if err != nil {
		t.Errorf("gzip.NewReader(rr.Body) = %v", err)
		return
	}
	defer gr.Close()

	bdec, _ := io.ReadAll(gr)
	if body != string(bdec) {
		t.Errorf("BODY = %v, want %v", string(bdec), body)
	}
}

func TestGzip(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Encoding", "gzip, deflate")

	w := httptest.NewRecorder()
	body := strings.Repeat("This is a Test!\n", 1000)
	router := xin.New()
	router.Use(DefaultHTTPGziper().Handler())
	router.GET("/", func(c *xin.Context) {
		c.String(200, body)
	})

	router.ServeHTTP(w, req)

	assertGzipEnable(t, w, body)
}

func TestGzipIgnore_HTTP_1_0(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.ProtoMajor = 1
	req.ProtoMinor = 0
	req.Header.Add("Accept-Encoding", "gzip")

	router := xin.New()
	zp := DefaultHTTPGziper()
	router.Use(zp.Handler())

	body := strings.Repeat("This is http 1.0!\n", 1000)
	router.GET("/", func(c *xin.Context) {
		c.String(200, body)
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assertGzipIgnore(t, rr, body)
}

func TestGzipIgnoreSmallSize(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Encoding", "gzip")

	router := xin.New()
	zp := DefaultHTTPGziper()
	router.Use(zp.Handler())

	body := "this is a TEXT!"
	router.GET("/", func(c *xin.Context) {
		c.String(200, body)
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assertGzipIgnore(t, rr, body)
}

func TestGzipIngoreContentType(t *testing.T) {
	req, _ := http.NewRequest("GET", "/image.png", nil)
	req.Header.Add("Accept-Encoding", "gzip")

	body := strings.Repeat("This is a PNG!\n", 1000)
	router := xin.New()
	router.Use(DefaultHTTPGziper().Handler())
	router.GET("/image.png", func(c *xin.Context) {
		c.Header("Content-Type", "image/png")
		c.String(200, body)
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assertGzipIgnore(t, rr, body)
}

func TestGzipIgnorePathPrefix(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/books", nil)
	req.Header.Add("Accept-Encoding", "gzip")

	router := xin.New()
	zp := DefaultHTTPGziper()
	zp.IgnorePathPrefix("/api/")
	router.Use(zp.Handler())

	body := strings.Repeat("This is books!\n", 1000)
	router.GET("/api/books", func(c *xin.Context) {
		c.String(200, body)
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assertGzipIgnore(t, rr, body)
}

func TestGzipIgnorePathRegexp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/books", nil)
	req.Header.Add("Accept-Encoding", "gzip")

	router := xin.New()
	zp := DefaultHTTPGziper()
	zp.IgnorePathRegexp("/.*/books")
	router.Use(zp.Handler())

	body := strings.Repeat("This is books!\n", 1000)
	router.GET("/api/books", func(c *xin.Context) {
		c.String(200, body)
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assertGzipIgnore(t, rr, body)
}

func testGzipIgnoreProxied(t *testing.T, proxied string, hand xin.HandlerFunc) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Via", "test")

	router := xin.New()
	zp := DefaultHTTPGziper()
	zp.SetProxied(proxied)
	router.Use(zp.Handler())

	body := strings.Repeat("This is proxy!\n", 1000)
	router.GET("/", func(c *xin.Context) {
		if hand != nil {
			hand(c)
		}
		c.String(200, body)
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assertGzipIgnore(t, rr, body)
}

func TestGzipIgnoreProxiedOff(t *testing.T) {
	testGzipIgnoreProxied(t, "off", nil)
}

func TestGzipIgnoreProxiedExpired(t *testing.T) {
	testGzipIgnoreProxied(t, "expired", nil)
}

func TestGzipIgnoreProxiedNoCache(t *testing.T) {
	testGzipIgnoreProxied(t, "no-Cache", nil)
}

func TestGzipIgnoreProxiedNoStore(t *testing.T) {
	testGzipIgnoreProxied(t, "No-Store", nil)
}

func TestGzipIgnoreProxiedPrivate(t *testing.T) {
	testGzipIgnoreProxied(t, "Private", nil)
}

func TestGzipIgnoreProxiedNoLastModified(t *testing.T) {
	testGzipIgnoreProxied(t, "no_Last_Modified", func(c *xin.Context) {
		c.Header("Last-Modified", "Wed, 21 Oct 2015 07:28:00 GMT")
	})
}

func TestGzipIgnoreProxiedNoETag(t *testing.T) {
	testGzipIgnoreProxied(t, "No_ETag", func(c *xin.Context) {
		c.Header("ETag", "13932423424")
	})
}

func testGzipEnableProxied(t *testing.T, proxied string, hand xin.HandlerFunc) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Via", "test")

	router := xin.New()
	zp := DefaultHTTPGziper()
	zp.SetProxied(proxied)
	router.Use(zp.Handler())

	body := strings.Repeat("This is proxy!\n", 1000)
	router.GET("/", func(c *xin.Context) {
		if hand != nil {
			hand(c)
		}
		c.String(200, body)
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assertGzipEnable(t, rr, body)
}

func TestGzipEnableProxiedExpired(t *testing.T) {
	testGzipEnableProxied(t, "Expired", func(c *xin.Context) {
		c.Header("Expires", "Wed, 21 Oct 2015 07:28:00 GMT")
	})
}

func TestGzipEnableProxiedNoCache(t *testing.T) {
	testGzipEnableProxied(t, "No-Cache", func(c *xin.Context) {
		c.Header("Cache-Control", "no-cache")
	})
}

func TestGzipEnableProxiedNoStore(t *testing.T) {
	testGzipEnableProxied(t, "No-Store", func(c *xin.Context) {
		c.Header("Cache-Control", "no-store")
	})
}

func TestGzipEnableProxiedPrivate(t *testing.T) {
	testGzipEnableProxied(t, "Private", func(c *xin.Context) {
		c.Header("Cache-Control", "Private")
	})
}

func TestGzipEnableProxiedNoLastModified(t *testing.T) {
	testGzipEnableProxied(t, "No_Last_Modified", nil)
}

func TestGzipEnableProxiedNoETag(t *testing.T) {
	testGzipEnableProxied(t, "No_ETag", nil)
}

func TestGzipEnableProxiedAny(t *testing.T) {
	testGzipEnableProxied(t, "Any", nil)
	fmt.Printf("%d\n", ProxiedOff)
}
