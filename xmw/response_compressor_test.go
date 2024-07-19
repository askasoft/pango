package xmw

import (
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/askasoft/pango/xin"
)

func TestResponseCompressorProxiedFlag(t *testing.T) {
	fs := []ProxiedFlag{
		ProxiedOff,
		ProxiedAny,
		ProxiedAuth,
		ProxiedExpired,
		ProxiedNoCache,
		ProxiedNoStore,
		ProxiedPrivate,
		ProxiedNoLastModified,
		ProxiedNoETag,
	}

	for i, f := range fs {
		e := 0
		if i > 0 {
			e = 1 << (i - 1)
		}
		a := int(f)
		if e != a {
			t.Errorf("%v = %v, want %v", f, a, e)
		}
	}
}

func assertResponseCompressorHeader(t *testing.T, rr *httptest.ResponseRecorder, sc int, hce, hvary string) {
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

func assertResponseCompressorIgnore(t *testing.T, rr *httptest.ResponseRecorder, body string) {
	assertResponseCompressorHeader(t, rr, http.StatusOK, "", "")
	if rr.Body.String() != body {
		t.Errorf(`Body = %v, want %v`, rr.Body.String(), body)
	}
}

func assertResponseCompressorEnable(t *testing.T, rr *httptest.ResponseRecorder, enc, body string) {
	assertResponseCompressorHeader(t, rr, http.StatusOK, enc, "Accept-Encoding")

	if len(body) == rr.Body.Len() {
		t.Errorf("len(body) = rr.Body.Len() = %v", len(body))
	}

	var ur io.ReadCloser
	var err error

	if enc == "deflate" {
		ur, err = zlib.NewReader(rr.Body)
		if err != nil {
			t.Fatalf("zlib.NewReader(rr.Body) = %v", err)
			return
		}
	} else {
		ur, err = gzip.NewReader(rr.Body)
		if err != nil {
			t.Fatalf("gzip.NewReader(rr.Body) = %v", err)
			return
		}
	}
	defer ur.Close()

	bdec, _ := io.ReadAll(ur)
	if body != string(bdec) {
		t.Errorf("BODY = %v, want %v", string(bdec), body)
	}
}

func TestResponseCompressorGzip(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Encoding", "gzip, deflate")

	w := httptest.NewRecorder()
	body := strings.Repeat("This is a Test!\n", 1000)
	router := xin.New()
	router.Use(DefaultResponseCompressor().Handler())
	router.GET("/", func(c *xin.Context) {
		c.String(200, body)
	})

	router.ServeHTTP(w, req)

	assertResponseCompressorEnable(t, w, "gzip", body)
}

func TestResponseCompressorZlib(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Encoding", "deflate")

	w := httptest.NewRecorder()
	body := strings.Repeat("This is a Test!\n", 1000)
	router := xin.New()
	router.Use(DefaultResponseCompressor().Handler())
	router.GET("/", func(c *xin.Context) {
		c.String(200, body)
	})

	router.ServeHTTP(w, req)

	assertResponseCompressorEnable(t, w, "deflate", body)
}

func TestResponseCompressorIgnore_HTTP_1_0(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.ProtoMajor = 1
	req.ProtoMinor = 0
	req.Header.Add("Accept-Encoding", "gzip")

	router := xin.New()
	rc := DefaultResponseCompressor()
	router.Use(rc.Handler())

	body := strings.Repeat("This is http 1.0!\n", 1000)
	router.GET("/", func(c *xin.Context) {
		c.String(200, body)
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assertResponseCompressorIgnore(t, rr, body)
}

func TestResponseCompressorIgnoreSmallSize(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Encoding", "gzip")

	router := xin.New()
	rc := DefaultResponseCompressor()
	router.Use(rc.Handler())

	body := "this is a TEXT!"
	router.GET("/", func(c *xin.Context) {
		c.String(200, body)
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assertResponseCompressorIgnore(t, rr, body)
}

func TestResponseCompressorIngoreContentType(t *testing.T) {
	req, _ := http.NewRequest("GET", "/image.png", nil)
	req.Header.Add("Accept-Encoding", "gzip")

	body := strings.Repeat("This is a PNG!\n", 1000)
	router := xin.New()
	router.Use(DefaultResponseCompressor().Handler())
	router.GET("/image.png", func(c *xin.Context) {
		c.Header("Content-Type", "image/png")
		c.String(200, body)
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assertResponseCompressorIgnore(t, rr, body)
}

func TestResponseCompressorIgnorePathPrefix(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/books", nil)
	req.Header.Add("Accept-Encoding", "gzip")

	router := xin.New()
	rc := DefaultResponseCompressor()
	rc.IgnorePathPrefix("/api/")
	router.Use(rc.Handler())

	body := strings.Repeat("This is books!\n", 1000)
	router.GET("/api/books", func(c *xin.Context) {
		c.String(200, body)
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assertResponseCompressorIgnore(t, rr, body)
}

func TestResponseCompressorIgnorePathRegexp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/books", nil)
	req.Header.Add("Accept-Encoding", "gzip")

	router := xin.New()
	rc := DefaultResponseCompressor()
	rc.IgnorePathRegexp("/.*/books")
	router.Use(rc.Handler())

	body := strings.Repeat("This is books!\n", 1000)
	router.GET("/api/books", func(c *xin.Context) {
		c.String(200, body)
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assertResponseCompressorIgnore(t, rr, body)
}

func testResponseCompressorIgnoreProxied(t *testing.T, proxied string, hand xin.HandlerFunc) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Via", "test")

	router := xin.New()
	rc := DefaultResponseCompressor()
	rc.SetProxied(proxied)
	router.Use(rc.Handler())

	body := strings.Repeat("This is proxy!\n", 1000)
	router.GET("/", func(c *xin.Context) {
		if hand != nil {
			hand(c)
		}
		c.String(200, body)
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assertResponseCompressorIgnore(t, rr, body)
}

func TestResponseCompressorIgnoreProxiedOff(t *testing.T) {
	testResponseCompressorIgnoreProxied(t, "off", nil)
}

func TestResponseCompressorIgnoreProxiedExpired(t *testing.T) {
	testResponseCompressorIgnoreProxied(t, "expired", nil)
}

func TestResponseCompressorIgnoreProxiedNoCache(t *testing.T) {
	testResponseCompressorIgnoreProxied(t, "no-Cache", nil)
}

func TestResponseCompressorIgnoreProxiedNoStore(t *testing.T) {
	testResponseCompressorIgnoreProxied(t, "No-Store", nil)
}

func TestResponseCompressorIgnoreProxiedPrivate(t *testing.T) {
	testResponseCompressorIgnoreProxied(t, "Private", nil)
}

func TestResponseCompressorIgnoreProxiedNoLastModified(t *testing.T) {
	testResponseCompressorIgnoreProxied(t, "no_Last_Modified", func(c *xin.Context) {
		c.Header("Last-Modified", "Wed, 21 Oct 2015 07:28:00 GMT")
	})
}

func TestResponseCompressorIgnoreProxiedNoETag(t *testing.T) {
	testResponseCompressorIgnoreProxied(t, "No_ETag", func(c *xin.Context) {
		c.Header("ETag", "13932423424")
	})
}

func testResponseCompressorEnableProxied(t *testing.T, proxied string, hand xin.HandlerFunc) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Via", "test")

	router := xin.New()
	rc := DefaultResponseCompressor()
	rc.SetProxied(proxied)
	router.Use(rc.Handler())

	body := strings.Repeat("This is proxy!\n", 1000)
	router.GET("/", func(c *xin.Context) {
		if hand != nil {
			hand(c)
		}
		c.String(200, body)
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assertResponseCompressorEnable(t, rr, "gzip", body)
}

func TestResponseCompressorEnableProxiedExpired(t *testing.T) {
	testResponseCompressorEnableProxied(t, "Expired", func(c *xin.Context) {
		c.Header("Expires", "Wed, 21 Oct 2015 07:28:00 GMT")
	})
}

func TestResponseCompressorEnableProxiedNoCache(t *testing.T) {
	testResponseCompressorEnableProxied(t, "No-Cache", func(c *xin.Context) {
		c.Header("Cache-Control", "no-cache")
	})
}

func TestResponseCompressorEnableProxiedNoStore(t *testing.T) {
	testResponseCompressorEnableProxied(t, "No-Store", func(c *xin.Context) {
		c.Header("Cache-Control", "no-store")
	})
}

func TestResponseCompressorEnableProxiedPrivate(t *testing.T) {
	testResponseCompressorEnableProxied(t, "Private", func(c *xin.Context) {
		c.Header("Cache-Control", "Private")
	})
}

func TestResponseCompressorEnableProxiedNoLastModified(t *testing.T) {
	testResponseCompressorEnableProxied(t, "No_Last_Modified", nil)
}

func TestResponseCompressorEnableProxiedNoETag(t *testing.T) {
	testResponseCompressorEnableProxied(t, "No_ETag", nil)
}

func TestResponseCompressorEnableProxiedAny(t *testing.T) {
	testResponseCompressorEnableProxied(t, "Any", nil)
	fmt.Printf("%d\n", ProxiedOff)
}
