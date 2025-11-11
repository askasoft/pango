package xin

import (
	"fmt"
	"net"
	"net/http"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/net/netx"
	"github.com/askasoft/pango/ref"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin/render"
	"github.com/askasoft/pango/xin/validate"
)

const (
	defaultMultipartMemory = 32 << 20 // 32 MB
)

var (
	DefaultRemoteIPHeaders = []string{"X-Forwarded-For", "X-Real-IP"}
	DefaultSSLProxyHeaders = map[string]string{"X-Forwarded-Proto": "https"}
	DefaultTrustedProxies  = netx.IntranetCIDRs
)

// regex used by redirectTrailingSlash()
var (
	regUnsafePrefix  = regexp.MustCompile("[^a-zA-Z0-9/-]+")
	regRepeatedSlash = regexp.MustCompile("/{2,}")
)

// HandlerFunc defines the handler used by xin middleware as return value.
type HandlerFunc func(*Context)

// HandlersChain defines a HandlerFunc slice.
type HandlersChain []HandlerFunc

// Last returns the last handler in the chain. i.e. the last handler is the main one.
func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}

// RouteInfo represents a request route's specification which contains method and path and its handler.
type RouteInfo struct {
	Method      string
	Path        string
	Handler     string
	HandlerFunc HandlerFunc
}

// RoutesInfo defines a RouteInfo slice.
type RoutesInfo []RouteInfo

// Trusted platforms
const (
	// PlatformGoogleAppEngine when running on Google App Engine. Trust X-Appengine-Remote-Addr
	// for determining the client's IP
	PlatformGoogleAppEngine = "X-Appengine-Remote-Addr"

	// PlatformCloudflare when using Cloudflare's CDN. Trust CF-Connecting-IP for determining
	// the client's IP
	PlatformCloudflare = "CF-Connecting-IP"
)

// Engine is the framework's instance, it contains the muxer, middleware and configuration settings.
// Create an instance of Engine, by using New() or Default()
type Engine struct {
	RouterGroup

	// RedirectTrailingSlash enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 307 for all other request methods.
	RedirectTrailingSlash bool

	// RedirectFixedPath if enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 307 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// HandleMethodNotAllowed if enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// UseRawPath if enabled, the url.RawPath will be used to find parameters.
	UseRawPath bool

	// UnescapePathValues if true, the path value will be unescaped.
	// If UseRawPath is false (by default), the UnescapePathValues effectively is true,
	// as url.Path gonna be used, which is already unescaped.
	UnescapePathValues bool

	// RemoveExtraSlash a parameter can be parsed from the URL even with extra slashes.
	RemoveExtraSlash bool

	// RemoteIPHeaders list of headers used to obtain the client IP when
	// `(*xin.Context).Request.RemoteAddr` is matched by at least one of the
	// network origins of list defined by `(*xin.Engine).SetTrustedProxies()`.
	RemoteIPHeaders []string

	// TrustedIPHeader if set to a constant of value xin.Platform*, trusts the headers set by
	// that platform to determine the client IP
	TrustedIPHeader string

	// SSLProxyHeaders is set of header keys with associated values that would indicate a valid https request.
	// Useful when behind a Proxy Server (Apache, Nginx).
	// Default is `map[string]string{"X-Forwarded-Proto": "https"}`.
	SSLProxyHeaders map[string]string

	// MaxMultipartMemory value of 'maxMemory' param that is given to http.Request's ParseMultipartForm
	// method call.
	MaxMultipartMemory int64

	// HTMLRenderer html templates renderer
	HTMLRenderer render.HTMLRenderer

	// Validator struct validator
	Validator validate.StructValidator

	// Logger
	Logger log.Logger

	secureJSONPrefix string
	allNoRoute       HandlersChain
	allNoMethod      HandlersChain
	noRoute          HandlersChain
	noMethod         HandlersChain
	pool             sync.Pool
	trees            methodTrees
	maxParams        uint16
	maxSections      uint16
	trustedProxies   []*net.IPNet
}

// New returns a new blank Engine instance without any middleware attached.
// By default, the configuration is:
// - RedirectTrailingSlash:  true
// - RedirectFixedPath:      false
// - HandleMethodNotAllowed: false
// - ForwardedByClientIP:    true
// - UseRawPath:             false
// - UnescapePathValues:     true
func New() *Engine {
	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      false,
		HandleMethodNotAllowed: false,
		RemoteIPHeaders:        DefaultRemoteIPHeaders,
		TrustedIPHeader:        "",
		SSLProxyHeaders:        DefaultSSLProxyHeaders,
		UseRawPath:             false,
		RemoveExtraSlash:       false,
		UnescapePathValues:     true,
		MaxMultipartMemory:     defaultMultipartMemory,
		Validator:              validate.NewStructValidator(),
		Logger:                 log.GetLogger("XIN"),
		trees:                  make(methodTrees, 0, 9),
		secureJSONPrefix:       ")]}',\n",
	}
	engine.engine = engine
	engine.pool.New = engine.allocateContext
	engine.SetTrustedProxies(DefaultTrustedProxies) //nolint: errcheck
	return engine
}

// Default returns an Engine instance with the Recovery middleware already attached.
func Default() *Engine {
	engine := New()
	engine.Use(Recovery())
	engine.Logger.Info("Creating an Engine instance with the Recovery middleware already attached.")
	return engine
}

func (engine *Engine) allocateContext() any {
	v := make(Params, 0, engine.maxParams)
	skippedNodes := make([]skippedNode, 0, engine.maxSections)

	c := &Context{
		engine:       engine,
		params:       &v,
		skippedNodes: &skippedNodes,
		Logger:       engine.Logger.GetLogger("XINC"),
	}
	return c
}

// SecureJSONPrefix sets the secureJSONPrefix used in Context.SecureJSON.
// Prefixing the JSON string in this manner is used to help prevent JSON Hijacking.
// The prefix renders the string syntactically invalid as a script so that it cannot be hijacked.
// This prefix should be stripped before parsing the string as JSON.
// The default prefix is ")]}',\n".
func (engine *Engine) SecureJSONPrefix(prefix string) *Engine {
	engine.secureJSONPrefix = prefix
	return engine
}

// NoRoute adds handlers for NoRoute. It returns a 404 code by default.
func (engine *Engine) NoRoute(handlers ...HandlerFunc) {
	engine.noRoute = handlers
	engine.rebuild404Handlers()
}

// NoMethod sets the handlers called when Engine.HandleMethodNotAllowed = true.
func (engine *Engine) NoMethod(handlers ...HandlerFunc) {
	engine.noMethod = handlers
	engine.rebuild405Handlers()
}

// Use attaches a global middleware to the router. i.e. the middleware attached through Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func (engine *Engine) Use(middleware ...HandlerFunc) IRoutes {
	engine.RouterGroup.Use(middleware...)
	engine.rebuild404Handlers()
	engine.rebuild405Handlers()
	return engine
}

func (engine *Engine) rebuild404Handlers() {
	engine.allNoRoute = engine.combineHandlers(engine.noRoute)
}

func (engine *Engine) rebuild405Handlers() {
	engine.allNoMethod = engine.combineHandlers(engine.noMethod)
}

func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
	if path[0] != '/' {
		panic("path must begin with '/'")
	}
	if method == "" {
		panic("HTTP method can not be empty")
	}
	if len(handlers) == 0 {
		panic("there must be at least one handler")
	}

	if engine.Logger.IsInfoEnabled() {
		nuHandlers := len(handlers)
		handlerName := ref.NameOfFunc(handlers.Last())
		engine.Logger.Infof("%-6s %-25s --> %s (%d handlers)", method, path, handlerName, nuHandlers)
	}

	root := engine.trees.get(method)
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		engine.trees = append(engine.trees, methodTree{method: method, root: root})
	}
	root.addRoute(path, handlers)

	if paramsCount := countParams(path); paramsCount > engine.maxParams {
		engine.maxParams = paramsCount
	}

	if sectionsCount := countSections(path); sectionsCount > engine.maxSections {
		engine.maxSections = sectionsCount
	}
}

// Routes returns a slice of registered routes, including some useful information, such as:
// the http method, path and the handler name.
func (engine *Engine) Routes() (routes RoutesInfo) {
	for _, tree := range engine.trees {
		routes = iterate("", tree.method, routes, tree.root)
	}
	return routes
}

func iterate(path, method string, routes RoutesInfo, root *node) RoutesInfo {
	path += root.path
	if len(root.handlers) > 0 {
		handlerFunc := root.handlers.Last()
		routes = append(routes, RouteInfo{
			Method:      method,
			Path:        path,
			Handler:     ref.NameOfFunc(handlerFunc),
			HandlerFunc: handlerFunc,
		})
	}
	for _, child := range root.children {
		routes = iterate(path, method, routes, child)
	}
	return routes
}

// SetTrustedProxies set a list of network origins (IPv4 addresses,
// IPv4 CIDRs, IPv6 addresses or IPv6 CIDRs) from which to trust
// request's headers that contain alternative client IP when
// `(*xin.Engine).ForwardedByClientIP` is `true`.
// `TrustedProxies` feature is enabled by default, and it also trusts all intranet proxies by default.
// If you want to disable this feature, use Engine.SetTrustedProxies(nil),
// then Context.ClientIP() will return the remote address directly.
func (engine *Engine) SetTrustedProxies(cidrs []string) error {
	ipnets, err := netx.ParseCIDRs(cidrs)
	if err != nil {
		return err
	}

	engine.trustedProxies = ipnets
	return nil
}

// isTrustedProxy will check whether the IP address is included in the trusted list according to Engine.trustedProxies
func (engine *Engine) isTrustedProxy(ip net.IP) bool {
	for _, cidr := range engine.trustedProxies {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

// validateClientIP will parse X-Forwarded-For header and return the trusted client IP address
func (engine *Engine) validateClientIP(header string) (clientIP string, valid bool) {
	if header == "" {
		return "", false
	}

	items := str.FieldsRune(header, ',')
	for i := len(items) - 1; i >= 0; i-- {
		ipStr := str.TrimSpace(items[i])
		ip := net.ParseIP(ipStr)
		if ip == nil {
			break
		}

		// X-Forwarded-For is appended by proxy
		// Check IPs in reverse order and stop when find untrusted proxy
		if (i == 0) || (!engine.isTrustedProxy(ip)) {
			return ipStr, true
		}
	}
	return "", false
}

// ServeHTTP conforms to the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := engine.pool.Get().(*Context)
	c.writermem.reset(w, c.Logger)
	c.Request = req
	c.Context = req.Context()

	engine.handleHTTPRequest(c)

	c.reset()
	engine.pool.Put(c)
}

// HandleContext re-enters a context that has been rewritten.
// This can be done by setting c.Request.URL.Path to your new target.
// Disclaimer: You can loop yourself to deal with this, use wisely.
func (engine *Engine) HandleContext(c *Context) {
	oi, oh := c.index, c.handlers

	engine.handleHTTPRequest(c)

	c.index, c.handlers = oi, oh
}

func (engine *Engine) handleHTTPRequest(c *Context) {
	c.reset()

	httpMethod := c.Request.Method
	rPath := c.Request.URL.Path
	unescape := false
	if engine.UseRawPath && len(c.Request.URL.RawPath) > 0 {
		rPath = c.Request.URL.RawPath
		unescape = engine.UnescapePathValues
	}

	if engine.RemoveExtraSlash {
		rPath = cleanPath(rPath)
	}

	// Find root of the tree for the given HTTP method
	t := engine.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method != httpMethod {
			continue
		}
		root := t[i].root
		// Find route in tree
		value := root.getValue(rPath, c.params, c.skippedNodes, unescape)
		if value.params != nil {
			c.Params = *value.params
		}
		if value.handlers != nil {
			c.handlers = value.handlers
			c.fullPath = value.fullPath
			c.Next()
			c.writermem.WriteHeaderNow()
			return
		}
		if httpMethod != http.MethodConnect && rPath != "/" {
			if value.tsr && engine.RedirectTrailingSlash {
				redirectTrailingSlash(c)
				return
			}
			if engine.RedirectFixedPath && redirectFixedPath(c, root) {
				return
			}
		}
		break
	}

	if engine.HandleMethodNotAllowed && len(t) > 0 {
		// According to RFC 7231 section 6.5.5, MUST generate an Allow header field in response
		// containing a list of the target resource's currently supported methods.
		allowed := make([]string, 0, len(t)-1)
		for _, tree := range engine.trees {
			if tree.method == httpMethod {
				continue
			}
			if value := tree.root.getValue(rPath, nil, c.skippedNodes, unescape); value.handlers != nil {
				allowed = append(allowed, tree.method)
			}
		}
		if len(allowed) > 0 {
			c.handlers = engine.allNoMethod
			c.writermem.Header().Set("Allow", strings.Join(allowed, ", "))
			serveError(c, http.StatusMethodNotAllowed)
			return
		}
	}

	c.handlers = engine.allNoRoute
	serveError(c, http.StatusNotFound)
}

func serveError(c *Context, code int) {
	c.writermem.status = code
	c.Next()
	if c.writermem.Written() {
		return
	}
	if c.writermem.Status() == code {
		c.writermem.Header().Set("Content-Type", MIMEPlain)
		body := fmt.Sprintf("%d %s", code, http.StatusText(code))
		_, err := c.Writer.Write(str.UnsafeBytes(body))
		if err != nil {
			c.Logger.Warnf("cannot write message to writer during serve error: %v", err)
		}
		return
	}
	c.writermem.WriteHeaderNow()
}

func redirectTrailingSlash(c *Context) {
	req := c.Request
	p := req.URL.Path
	if prefix := path.Clean(c.Request.Header.Get("X-Forwarded-Prefix")); prefix != "." {
		prefix = regUnsafePrefix.ReplaceAllString(prefix, "")
		prefix = regRepeatedSlash.ReplaceAllString(prefix, "/")

		p = prefix + "/" + req.URL.Path
	}
	req.URL.Path = p + "/"
	if length := len(p); length > 1 && p[length-1] == '/' {
		req.URL.Path = p[:length-1]
	}
	redirectRequest(c)
}

func redirectFixedPath(c *Context, root *node) bool {
	req := c.Request
	rPath := req.URL.Path

	if fixedPath, ok := root.findCaseInsensitivePath(cleanPath(rPath), true); ok {
		req.URL.Path = str.UnsafeString(fixedPath)
		redirectRequest(c)
		return true
	}
	return false
}

func redirectRequest(c *Context) {
	req := c.Request
	rPath := req.URL.Path
	rURL := req.URL.String()

	code := http.StatusMovedPermanently // Permanent redirect, request with GET method
	if req.Method != http.MethodGet {
		code = http.StatusTemporaryRedirect
	}

	c.Logger.Debugf("redirect request %d: %s --> %s", code, rPath, rURL)
	http.Redirect(c.Writer, req, rURL, code)
	c.writermem.WriteHeaderNow()
}
