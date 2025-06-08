package xin

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/askasoft/pango/test/assert"
	"github.com/askasoft/pango/tpl"
)

// testNewHttpServer create a http.Server instance
func testNewHttpServer(engine *Engine) *http.Server {
	server := &http.Server{
		Handler:           engine,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
	}
	return server
}

// testRun attaches the router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func testRun(engine *Engine, address string) (err error) {
	engine.Logger.Infof("Listening and serving HTTP on %s", address)
	server := testNewHttpServer(engine)
	server.Addr = address
	err = server.ListenAndServe()
	if err != nil {
		engine.Logger.Errorf("Listening and serving HTTP on %s failed: %v", err)
	}
	return
}

// testRunTLS attaches the router to a http.Server and starts listening and serving HTTPS (secure) requests.
// It is a shortcut for http.ListenAndServeTLS(addr, certFile, keyFile, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func testRunTLS(engine *Engine, address, certFile, keyFile string) (err error) {
	engine.Logger.Infof("Listening and serving HTTPS on %s\n", address)

	server := testNewHttpServer(engine)
	server.Addr = address
	err = server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		engine.Logger.Errorf("Listen and serve HTTPs on %s failed: %v", address, err)
	}
	return
}

// testRunListener attaches the router to a http.Server and starts listening and serving HTTP requests
// through the specified net.Listener
func testRunListener(engine *Engine, listener net.Listener) (err error) {
	engine.Logger.Infof("Listening and serving HTTP on listener what's bind with address@%s", listener.Addr())

	err = http.Serve(listener, engine)
	if err != nil {
		engine.Logger.Errorf("Listen and serve HTTP on listener what's bind with address@%s failed: %v", listener.Addr(), err)
	}
	return err
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}

func testLoadHTML(t *testing.T, router *Engine) {
	router.HTMLTemplates.Delims("{[{", "}]}")

	router.HTMLTemplates.Funcs(tpl.FuncMap{
		"formatAsDate": formatAsDate,
	})
	err := router.HTMLTemplates.Load("./testdata/template/")
	assert.NoError(t, err)
}

func setupHTMLFiles(t *testing.T, tls bool, loadMethod func(*Engine)) *httptest.Server {
	router := New()
	testLoadHTML(t, router)
	loadMethod(router)
	router.GET("/test", func(c *Context) {
		c.HTML(http.StatusOK, "hello", map[string]string{"name": "world"})
	})
	router.GET("/raw", func(c *Context) {
		c.HTML(http.StatusOK, "raw", map[string]any{
			"now": time.Date(2017, 07, 01, 0, 0, 0, 0, time.UTC),
		})
	})

	var ts *httptest.Server

	if tls {
		ts = httptest.NewTLSServer(router)
	} else {
		ts = httptest.NewServer(router)
	}

	return ts
}

func TestLoadHTMLGlobDebugMode(t *testing.T) {
	ts := setupHTMLFiles(
		t,
		false,
		func(router *Engine) {
			router.HTMLTemplates.Load("./testdata/template/")
		},
	)
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/test", ts.URL))
	if err != nil {
		fmt.Println(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp))
}

func TestLoadHTMLGlobTestMode(t *testing.T) {
	ts := setupHTMLFiles(
		t,
		false,
		func(router *Engine) {
			router.HTMLTemplates.Load("./testdata/template/")
		},
	)
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/test", ts.URL))
	if err != nil {
		fmt.Println(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp))
}

func TestLoadHTMLGlobReleaseMode(t *testing.T) {
	ts := setupHTMLFiles(
		t,
		false,
		func(router *Engine) {
			router.HTMLTemplates.Load("./testdata/template/")
		},
	)
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/test", ts.URL))
	if err != nil {
		fmt.Println(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp))
}

func TestLoadHTMLGlobUsingTLS(t *testing.T) {
	ts := setupHTMLFiles(
		t,
		true,
		func(router *Engine) {
			router.HTMLTemplates.Load("./testdata/template/")
		},
	)
	defer ts.Close()

	// Use InsecureSkipVerify for avoiding `x509: certificate signed by unknown authority` error
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Get(fmt.Sprintf("%s/test", ts.URL))
	if err != nil {
		fmt.Println(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp))
}

func TestLoadHTMLGlobFromFuncMap(t *testing.T) {
	ts := setupHTMLFiles(
		t,
		false,
		func(router *Engine) {
			router.HTMLTemplates.Load("./testdata/template/")
		},
	)
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/raw", ts.URL))
	if err != nil {
		fmt.Println(err)
	}

	resp, _ := io.ReadAll(res.Body)
	assert.Equal(t, "Date: 2017/07/01\n", string(resp))
}

func TestCreateEngine(t *testing.T) {
	router := New()
	assert.Equal(t, "/", router.basePath)
	assert.Equal(t, router.engine, router)
	assert.Empty(t, router.Handlers)
}

func TestAddRoute(t *testing.T) {
	router := New()
	router.addRoute("GET", "/", HandlersChain{func(_ *Context) {}})

	assert.Len(t, router.trees, 1)
	assert.NotNil(t, router.trees.get("GET"))
	assert.Nil(t, router.trees.get("POST"))

	router.addRoute("POST", "/", HandlersChain{func(_ *Context) {}})

	assert.Len(t, router.trees, 2)
	assert.NotNil(t, router.trees.get("GET"))
	assert.NotNil(t, router.trees.get("POST"))

	router.addRoute("POST", "/post", HandlersChain{func(_ *Context) {}})
	assert.Len(t, router.trees, 2)
}

func TestAddRouteFails(t *testing.T) {
	router := New()
	assert.Panics(t, func() { router.addRoute("", "/", HandlersChain{func(_ *Context) {}}) })
	assert.Panics(t, func() { router.addRoute("GET", "a", HandlersChain{func(_ *Context) {}}) })
	assert.Panics(t, func() { router.addRoute("GET", "/", HandlersChain{}) })

	router.addRoute("POST", "/post", HandlersChain{func(_ *Context) {}})
	assert.Panics(t, func() {
		router.addRoute("POST", "/post", HandlersChain{func(_ *Context) {}})
	})
}

func TestNoRouteWithoutGlobalHandlers(t *testing.T) {
	var middleware0 HandlerFunc = func(c *Context) {}
	var middleware1 HandlerFunc = func(c *Context) {}

	router := New()

	router.NoRoute(middleware0)
	assert.Nil(t, router.Handlers)
	assert.Len(t, router.noRoute, 1)
	assert.Len(t, router.allNoRoute, 1)
	compareFunc(t, router.noRoute[0], middleware0)
	compareFunc(t, router.allNoRoute[0], middleware0)

	router.NoRoute(middleware1, middleware0)
	assert.Len(t, router.noRoute, 2)
	assert.Len(t, router.allNoRoute, 2)
	compareFunc(t, router.noRoute[0], middleware1)
	compareFunc(t, router.allNoRoute[0], middleware1)
	compareFunc(t, router.noRoute[1], middleware0)
	compareFunc(t, router.allNoRoute[1], middleware0)
}

func TestNoRouteWithGlobalHandlers(t *testing.T) {
	var middleware0 HandlerFunc = func(c *Context) {}
	var middleware1 HandlerFunc = func(c *Context) {}
	var middleware2 HandlerFunc = func(c *Context) {}

	router := New()
	router.Use(middleware2)

	router.NoRoute(middleware0)
	assert.Len(t, router.allNoRoute, 2)
	assert.Len(t, router.Handlers, 1)
	assert.Len(t, router.noRoute, 1)

	compareFunc(t, router.Handlers[0], middleware2)
	compareFunc(t, router.noRoute[0], middleware0)
	compareFunc(t, router.allNoRoute[0], middleware2)
	compareFunc(t, router.allNoRoute[1], middleware0)

	router.Use(middleware1)
	assert.Len(t, router.allNoRoute, 3)
	assert.Len(t, router.Handlers, 2)
	assert.Len(t, router.noRoute, 1)

	compareFunc(t, router.Handlers[0], middleware2)
	compareFunc(t, router.Handlers[1], middleware1)
	compareFunc(t, router.noRoute[0], middleware0)
	compareFunc(t, router.allNoRoute[0], middleware2)
	compareFunc(t, router.allNoRoute[1], middleware1)
	compareFunc(t, router.allNoRoute[2], middleware0)
}

func TestNoMethodWithoutGlobalHandlers(t *testing.T) {
	var middleware0 HandlerFunc = func(c *Context) {}
	var middleware1 HandlerFunc = func(c *Context) {}

	router := New()

	router.NoMethod(middleware0)
	assert.Empty(t, router.Handlers)
	assert.Len(t, router.noMethod, 1)
	assert.Len(t, router.allNoMethod, 1)
	compareFunc(t, router.noMethod[0], middleware0)
	compareFunc(t, router.allNoMethod[0], middleware0)

	router.NoMethod(middleware1, middleware0)
	assert.Len(t, router.noMethod, 2)
	assert.Len(t, router.allNoMethod, 2)
	compareFunc(t, router.noMethod[0], middleware1)
	compareFunc(t, router.allNoMethod[0], middleware1)
	compareFunc(t, router.noMethod[1], middleware0)
	compareFunc(t, router.allNoMethod[1], middleware0)
}

func TestRebuild404Handlers(t *testing.T) {
}

func TestNoMethodWithGlobalHandlers(t *testing.T) {
	var middleware0 HandlerFunc = func(c *Context) {}
	var middleware1 HandlerFunc = func(c *Context) {}
	var middleware2 HandlerFunc = func(c *Context) {}

	router := New()
	router.Use(middleware2)

	router.NoMethod(middleware0)
	assert.Len(t, router.allNoMethod, 2)
	assert.Len(t, router.Handlers, 1)
	assert.Len(t, router.noMethod, 1)

	compareFunc(t, router.Handlers[0], middleware2)
	compareFunc(t, router.noMethod[0], middleware0)
	compareFunc(t, router.allNoMethod[0], middleware2)
	compareFunc(t, router.allNoMethod[1], middleware0)

	router.Use(middleware1)
	assert.Len(t, router.allNoMethod, 3)
	assert.Len(t, router.Handlers, 2)
	assert.Len(t, router.noMethod, 1)

	compareFunc(t, router.Handlers[0], middleware2)
	compareFunc(t, router.Handlers[1], middleware1)
	compareFunc(t, router.noMethod[0], middleware0)
	compareFunc(t, router.allNoMethod[0], middleware2)
	compareFunc(t, router.allNoMethod[1], middleware1)
	compareFunc(t, router.allNoMethod[2], middleware0)
}

func compareFunc(t *testing.T, a, b any) {
	sf1 := reflect.ValueOf(a)
	sf2 := reflect.ValueOf(b)
	if sf1.Pointer() != sf2.Pointer() {
		t.Error("different functions")
	}
}

func TestListOfRoutes(t *testing.T) {
	router := New()
	router.GET("/favicon.ico", handlerTest1)
	router.GET("/", handlerTest1)
	group := router.Group("/users")
	{
		group.GET("/", handlerTest2)
		group.GET("/:id", handlerTest1)
		group.POST("/:id", handlerTest2)
	}
	Static(router, "/static", ".")

	list := router.Routes()

	assert.Len(t, list, 7)
	assertRoutePresent(t, list, RouteInfo{
		Method:  "GET",
		Path:    "/favicon.ico",
		Handler: "^(.*/vendor/)?github.com/askasoft/pango/xin.handlerTest1$",
	})
	assertRoutePresent(t, list, RouteInfo{
		Method:  "GET",
		Path:    "/",
		Handler: "^(.*/vendor/)?github.com/askasoft/pango/xin.handlerTest1$",
	})
	assertRoutePresent(t, list, RouteInfo{
		Method:  "GET",
		Path:    "/users/",
		Handler: "^(.*/vendor/)?github.com/askasoft/pango/xin.handlerTest2$",
	})
	assertRoutePresent(t, list, RouteInfo{
		Method:  "GET",
		Path:    "/users/:id",
		Handler: "^(.*/vendor/)?github.com/askasoft/pango/xin.handlerTest1$",
	})
	assertRoutePresent(t, list, RouteInfo{
		Method:  "POST",
		Path:    "/users/:id",
		Handler: "^(.*/vendor/)?github.com/askasoft/pango/xin.handlerTest2$",
	})
}

func TestEngineHandleContext(t *testing.T) {
	r := New()
	r.GET("/", func(c *Context) {
		c.Request.URL.Path = "/v2"
		r.HandleContext(c)
	})
	v2 := r.Group("/v2")
	{
		v2.GET("/", func(c *Context) {})
	}

	assert.NotPanics(t, func() {
		w := performRequest(r, "GET", "/")
		assert.Equal(t, 301, w.Code)
	})
}

func TestEngineHandleContextManyReEntries(t *testing.T) {
	expectValue := 10000

	var handlerCounter, middlewareCounter int64

	r := New()
	r.Use(func(c *Context) {
		atomic.AddInt64(&middlewareCounter, 1)
	})
	r.GET("/:count", func(c *Context) {
		countStr := c.Param("count")
		count, err := strconv.Atoi(countStr)
		assert.NoError(t, err)

		n, err := c.Writer.Write([]byte("."))
		assert.NoError(t, err)
		assert.Equal(t, 1, n)

		switch {
		case count > 0:
			c.Request.URL.Path = "/" + strconv.Itoa(count-1)
			r.HandleContext(c)
		}
	}, func(c *Context) {
		atomic.AddInt64(&handlerCounter, 1)
	})

	assert.NotPanics(t, func() {
		w := performRequest(r, "GET", "/"+strconv.Itoa(expectValue-1)) // include 0 value
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, expectValue, w.Body.Len())
	})

	assert.Equal(t, int64(expectValue), handlerCounter)
	assert.Equal(t, int64(expectValue), middlewareCounter)
}

func TestEngineHandleContextPreventsMiddlewareReEntry(t *testing.T) {
	// given
	var handlerCounterV1, handlerCounterV2, middlewareCounterV1 int64

	r := New()
	v1 := r.Group("/v1")
	{
		v1.Use(func(c *Context) {
			atomic.AddInt64(&middlewareCounterV1, 1)
		})
		v1.GET("/test", func(c *Context) {
			atomic.AddInt64(&handlerCounterV1, 1)
			c.Status(http.StatusOK)
		})
	}

	v2 := r.Group("/v2")
	{
		v2.GET("/test", func(c *Context) {
			c.Request.URL.Path = "/v1/test"
			r.HandleContext(c)
		}, func(c *Context) {
			atomic.AddInt64(&handlerCounterV2, 1)
		})
	}

	// when
	responseV1 := performRequest(r, "GET", "/v1/test")
	responseV2 := performRequest(r, "GET", "/v2/test")

	// then
	assert.Equal(t, 200, responseV1.Code)
	assert.Equal(t, 200, responseV2.Code)
	assert.Equal(t, int64(2), handlerCounterV1)
	assert.Equal(t, int64(2), middlewareCounterV1)
	assert.Equal(t, int64(1), handlerCounterV2)
}

func TestPrepareTrustedCIRDsWith(t *testing.T) {
	r := New()

	// valid ipv4 cidr
	{
		expectedTrustedCIDRs := []*net.IPNet{parseCIDR("0.0.0.0/0")}
		err := r.SetTrustedProxies([]string{"0.0.0.0/0"})

		assert.NoError(t, err)
		assert.Equal(t, expectedTrustedCIDRs, r.trustedProxies)
	}

	// invalid ipv4 cidr
	{
		err := r.SetTrustedProxies([]string{"192.168.1.33/33"})

		assert.Error(t, err)
	}

	// valid ipv4 address
	{
		expectedTrustedCIDRs := []*net.IPNet{parseCIDR("192.168.1.33/32")}

		err := r.SetTrustedProxies([]string{"192.168.1.33"})

		assert.NoError(t, err)
		assert.Equal(t, expectedTrustedCIDRs, r.trustedProxies)
	}

	// invalid ipv4 address
	{
		err := r.SetTrustedProxies([]string{"192.168.1.256"})

		assert.Error(t, err)
	}

	// valid ipv6 address
	{
		expectedTrustedCIDRs := []*net.IPNet{parseCIDR("2002:0000:0000:1234:abcd:ffff:c0a8:0101/128")}
		err := r.SetTrustedProxies([]string{"2002:0000:0000:1234:abcd:ffff:c0a8:0101"})

		assert.NoError(t, err)
		assert.Equal(t, expectedTrustedCIDRs, r.trustedProxies)
	}

	// invalid ipv6 address
	{
		err := r.SetTrustedProxies([]string{"gggg:0000:0000:1234:abcd:ffff:c0a8:0101"})

		assert.Error(t, err)
	}

	// valid ipv6 cidr
	{
		expectedTrustedCIDRs := []*net.IPNet{parseCIDR("::/0")}
		err := r.SetTrustedProxies([]string{"::/0"})

		assert.NoError(t, err)
		assert.Equal(t, expectedTrustedCIDRs, r.trustedProxies)
	}

	// invalid ipv6 cidr
	{
		err := r.SetTrustedProxies([]string{"gggg:0000:0000:1234:abcd:ffff:c0a8:0101/129"})

		assert.Error(t, err)
	}

	// valid combination
	{
		expectedTrustedCIDRs := []*net.IPNet{
			parseCIDR("::/0"),
			parseCIDR("192.168.0.0/16"),
			parseCIDR("172.16.0.1/32"),
		}
		err := r.SetTrustedProxies([]string{
			"::/0",
			"192.168.0.0/16",
			"172.16.0.1",
		})

		assert.NoError(t, err)
		assert.Equal(t, expectedTrustedCIDRs, r.trustedProxies)
	}

	// invalid combination
	{
		err := r.SetTrustedProxies([]string{
			"::/0",
			"192.168.0.0/16",
			"172.16.0.256",
		})

		assert.Error(t, err)
	}

	// nil value
	{
		err := r.SetTrustedProxies(nil)

		assert.Empty(t, r.trustedProxies)
		assert.Nil(t, err)
	}
}

func parseCIDR(cidr string) *net.IPNet {
	_, parsedCIDR, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Println(err)
	}
	return parsedCIDR
}

func assertRoutePresent(t *testing.T, gotRoutes RoutesInfo, wantRoute RouteInfo) {
	for _, gotRoute := range gotRoutes {
		if gotRoute.Path == wantRoute.Path && gotRoute.Method == wantRoute.Method {
			assert.Regexp(t, wantRoute.Handler, gotRoute.Handler)
			return
		}
	}
	t.Errorf("route not found: %v", wantRoute)
}

func handlerTest1(c *Context) {}
func handlerTest2(c *Context) {}

func TestMethodNotAllowedNoRoute(t *testing.T) {
	g := New()
	g.HandleMethodNotAllowed = true

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	assert.NotPanics(t, func() { g.ServeHTTP(resp, req) })
	assert.Equal(t, http.StatusNotFound, resp.Code)
}
