package xin

import (
	"errors"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/net/httpx"
	"github.com/askasoft/pango/net/httpx/sse"
	"github.com/askasoft/pango/ref"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin/binding"
	"github.com/askasoft/pango/xin/render"
)

// Content-Type MIME of the most common data formats.
const (
	MIMEJSON              = binding.MIMEJSON
	MIMEHTML              = binding.MIMEHTML
	MIMEXML               = binding.MIMEXML
	MIMEXML2              = binding.MIMEXML2
	MIMEPlain             = binding.MIMEPlain
	MIMEPOSTForm          = binding.MIMEPOSTForm
	MIMEMultipartPOSTForm = binding.MIMEMultipartPOSTForm
)

// BodyBytesKey indicates a default body bytes key.
const BodyBytesKey = "XIN_BODY_BYTES"

// ContextKey is the key that a Context returns itself for.
const ContextKey = "XIN_CONTEXT"

// RequestKey is the key that a Context returns it's Request.
const RequestKey = "XIN_REQUEST"

// abortIndex represents a typical value used in abort functions.
const abortIndex = 128

// Context is the most important part of xin. It allows us to pass variables between middleware,
// manage the flow, validate the JSON of a request and render a JSON response for example.
type Context struct {
	writermem responseWriter
	Request   *http.Request
	Writer    ResponseWriter

	Params   Params
	handlers HandlersChain
	index    int
	fullPath string

	engine       *Engine
	params       *Params
	skippedNodes *[]skippedNode

	// locale string for the context of each request.
	Locale string

	// Errors is a list of errors attached to all the handlers/middlewares who used this context.
	Errors []error

	// Accepted defines a list of manually accepted formats for content negotiation.
	Accepted []string

	// attrs is a key/value pair exclusively for the context of each request.
	attrs map[string]any

	// queryCache caches the query result from c.Request.URL.Query().
	queryCache url.Values

	// formCache caches c.Request.PostForm, which contains the parsed form data from POST, PATCH,
	// or PUT body parameters.
	formCache url.Values

	// Logger
	Logger log.Logger
}

/************************************/
/********** CONTEXT CREATION ********/
/************************************/

func (c *Context) reset() {
	c.Writer = &c.writermem
	c.Params = c.Params[:0]
	c.handlers = nil
	c.index = -1
	c.fullPath = ""

	c.Locale = ""
	c.Errors = c.Errors[:0]
	c.Accepted = nil

	clear(c.attrs)
	c.queryCache = nil
	c.formCache = nil
	*c.params = (*c.params)[:0]
	*c.skippedNodes = (*c.skippedNodes)[:0]

	c.Logger.SetProps(nil)
}

// Copy returns a copy of the current context that can be safely used outside the request's scope.
// This has to be used when the context has to be passed to a goroutine.
func (c *Context) Copy() *Context {
	cp := Context{
		writermem: c.writermem,
		Request:   c.Request,
		engine:    c.engine,
		Locale:    c.Locale,
		Logger:    c.Logger,
	}

	cp.writermem.ResponseWriter = nil
	cp.Writer = &cp.writermem
	cp.index = abortIndex
	cp.handlers = nil
	cp.fullPath = c.fullPath

	cp.attrs = make(map[string]any, len(c.attrs))
	for k, v := range c.attrs {
		cp.attrs[k] = v
	}
	cp.Params = make([]Param, len(c.Params))
	copy(cp.Params, c.Params)
	return &cp
}

// HandlerName returns the main handler's name. For example if the handler is "handleGetUsers()",
// this function will return "main.handleGetUsers".
func (c *Context) HandlerName() string {
	return ref.NameOfFunc(c.handlers.Last())
}

// HandlerNames returns a list of all registered handlers for this context in descending order,
// following the semantics of HandlerName()
func (c *Context) HandlerNames() []string {
	hn := make([]string, 0, len(c.handlers))
	for _, val := range c.handlers {
		hn = append(hn, ref.NameOfFunc(val))
	}
	return hn
}

// Handler returns the main handler.
func (c *Context) Handler() HandlerFunc {
	return c.handlers.Last()
}

// FullPath returns a matched route full path. For not found routes
// returns an empty string.
//
//	router.GET("/user/:id", func(c *xin.Context) {
//	    c.FullPath() == "/user/:id" // true
//	})
func (c *Context) FullPath() string {
	return c.fullPath
}

/************************************/
/*********** FLOW CONTROL ***********/
/************************************/

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in GitHub.
func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.index++
	}
}

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized.
// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
// for this request are not called.
func (c *Context) Abort() {
	c.index = abortIndex
}

// AbortWithStatus calls `Abort()` and c.String(code, "%d %s", code, http.StatusText(code)).
// For example, a failed attempt to authenticate a request could use: context.AbortWithStatus(401).
func (c *Context) AbortWithStatus(code int) {
	c.Abort()
	c.String(code, "%d %s", code, http.StatusText(code))
}

// AbortWithStatusText calls `Abort()` and then `String` internally.
// This method stops the chain, writes the status code and return a text body.
// It also sets the Content-Type as "text/plain".
func (c *Context) AbortWithStatusText(code int, format string, values ...any) {
	c.Abort()
	c.String(code, format, values...)
}

// AbortWithStatusJSON calls `Abort()` and then `JSON` internally.
// This method stops the chain, writes the status code and return a JSON body.
// It also sets the Content-Type as "application/json".
func (c *Context) AbortWithStatusJSON(code int, jsonObj any) {
	c.Abort()
	c.JSON(code, jsonObj)
}

// AbortWithError calls `AbortWithStatus()` and `Error()` internally.
// This method stops the chain, writes the status code and pushes the specified error to `c.Errors`.
// See Context.Error() for more details.
func (c *Context) AbortWithError(code int, err error) {
	c.AbortWithStatus(code)
	c.AddError(err)
}

/************************************/
/********* ERROR MANAGEMENT *********/
/************************************/

// AddError attaches errors to the current context. The error is pushed to a list of errors.
// It's a good idea to call Error for each error that occurred during the resolution of a request.
// A middleware can be used to collect all the errors and push them to a database together,
// print a log, or append it in the HTTP response.
// Error will panic if err is nil.
func (c *Context) AddError(errs ...error) {
	c.Errors = append(c.Errors, errs...)
}

// SetError clear and attaches new errors to the current context.
func (c *Context) SetError(errs ...error) {
	c.Errors = errs
}

/************************************/
/******** METADATA MANAGEMENT********/
/************************************/

// Attrs get key/value map exclusively for this context.
func (c *Context) Attrs() map[string]any {
	if c.attrs == nil {
		c.attrs = make(map[string]any)
	}
	return c.attrs
}

// Del is used to delete a stored key/value pair for this context.
func (c *Context) Del(key string) {
	if c.attrs == nil {
		return
	}

	delete(c.attrs, key)
}

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.attrs if it was not used previously.
func (c *Context) Set(key string, value any) {
	if c.attrs == nil {
		c.attrs = make(map[string]any)
	}

	c.attrs[key] = value
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exist it returns (nil, false)
func (c *Context) Get(key string) (value any, exists bool) {
	value, exists = c.attrs[key]
	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (c *Context) MustGet(key string) any {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// GetString returns the value associated with the key as a string.
func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

// GetBool returns the value associated with the key as a boolean.
func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt returns the value associated with the key as an integer.
func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetInt64 returns the value associated with the key as an integer.
func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

// GetUint returns the value associated with the key as an unsigned integer.
func (c *Context) GetUint(key string) (ui uint) {
	if val, ok := c.Get(key); ok && val != nil {
		ui, _ = val.(uint)
	}
	return
}

// GetUint64 returns the value associated with the key as an unsigned integer.
func (c *Context) GetUint64(key string) (ui64 uint64) {
	if val, ok := c.Get(key); ok && val != nil {
		ui64, _ = val.(uint64)
	}
	return
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// GetTime returns the value associated with the key as time.
func (c *Context) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

// GetDuration returns the value associated with the key as a duration.
func (c *Context) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

// GetStrings returns the value associated with the key as a slice of strings.
func (c *Context) GetStrings(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *Context) GetStringMap(key string) (sm map[string]any) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]any)
	}
	return
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *Context) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

// GetStringMapStrings returns the value associated with the key as a map to a slice of strings.
func (c *Context) GetStringMapStrings(key string) (smss map[string][]string) {
	if val, ok := c.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}

/************************************/
/************ INPUT DATA ************/
/************************************/
func (c *Context) initQueryCache() {
	if c.queryCache == nil {
		if c.Request != nil && c.Request.URL != nil {
			c.queryCache = c.Request.URL.Query()
		} else {
			c.queryCache = url.Values{}
		}
	}
}

func (c *Context) initFormCache() {
	if c.formCache == nil {
		req := c.Request
		if err := req.ParseMultipartForm(c.engine.MaxMultipartMemory); err != nil {
			if !errors.Is(err, http.ErrNotMultipart) {
				c.Logger.Warnf("parse multipart form error: %v", err)
			}
		}

		c.formCache = req.PostForm
		if c.formCache == nil {
			c.formCache = make(url.Values)
		}
	}
}

// getStringMapFromCache is an internal method and returns a map which satisfy conditions.
func getStringMapFromCache(cache map[string][]string, key string) (map[string]string, bool) {
	dict := make(map[string]string)
	for n, v := range cache {
		if i := strings.IndexByte(n, '.'); i >= 1 && n[0:i] == key {
			dict[n[i+1:]] = v[0]
		} else if i := strings.IndexByte(n, '['); i >= 1 && n[0:i] == key {
			if j := strings.IndexByte(n[i+1:], ']'); j >= 1 {
				dict[n[i+1:][:j]] = v[0]
			}
		}
	}
	return dict, len(dict) > 0
}

// getStringsMapFromCache is an internal method and returns a map which satisfy conditions.
func getStringsMapFromCache(cache map[string][]string, key string) (map[string][]string, bool) {
	dict := make(map[string][]string)
	for n, v := range cache {
		k := ""
		if i := strings.IndexByte(n, '.'); i >= 1 && n[0:i] == key {
			k = n[i+1:]
		} else if i := strings.IndexByte(n, '['); i >= 1 && n[0:i] == key {
			if j := strings.IndexByte(n[i+1:], ']'); j >= 1 {
				k = n[i+1:][:j]
			}
		}

		if k != "" {
			dict[k] = append(dict[k], v...)
		}
	}
	return dict, len(dict) > 0
}

// Param returns the value of the URL param.
// It is a shortcut for c.Params.ByName(key)
//
//	router.GET("/user/:id", func(c *xin.Context) {
//	    // a GET request to /user/john
//	    id := c.Param("id") // id == "john"
//	    // a GET request to /user/john/
//	    id := c.Param("id") // id == "/john/"
//	})
func (c *Context) Param(key string) string {
	return c.Params.ByName(key)
}

// AddParam adds param to context and
// replaces path param key with given value for e2e testing purposes
// Example Route: "/user/:id"
// AddParam("id", 1)
// Result: "/user/1"
func (c *Context) AddParam(key, value string) {
	c.Params = append(c.Params, Param{Key: key, Value: value})
}

// Querys returns the url query map
func (c *Context) Querys() map[string][]string {
	c.initQueryCache()
	return c.queryCache
}

// Query returns the keyed url query value if it exists,
// otherwise it returns first non-empty value of defs or an empty string `("")`.
// It is shortcut for `c.Request.URL.Query().Get(key)`
//
//	    GET /path?id=1234&name=Manu&value=
//		   c.Query("id") == "1234"
//		   c.Query("name") == "Manu"
//		   c.Query("value") == ""
//		   c.Query("wtf", "none") == "none"
func (c *Context) Query(key string, defs ...string) string {
	if value, ok := c.GetQuery(key); ok {
		return value
	}
	return str.NonEmpty(defs...)
}

// GetQuery is like Query(), it returns the keyed url query value
// if it exists `(value, true)` (even when the value is an empty string),
// otherwise it returns `("", false)`.
// It is shortcut for `c.Request.URL.Query().Get(key)`
//
//	GET /?name=Manu&lastname=
//	("Manu", true) == c.GetQuery("name")
//	("", false) == c.GetQuery("id")
//	("", true) == c.GetQuery("lastname")
func (c *Context) GetQuery(key string) (string, bool) {
	if values, ok := c.GetQueryArray(key); ok {
		return values[0], ok
	}
	return "", false
}

// QueryArray returns a slice of strings for a given query key.
// The length of the slice depends on the number of params with the given key.
func (c *Context) QueryArray(key string) (values []string) {
	values, _ = c.GetQueryArray(key)
	return
}

// GetQueryArray returns a slice of strings for a given query key, plus
// a boolean value whether at least one value exists for the given key.
func (c *Context) GetQueryArray(key string) (values []string, ok bool) {
	c.initQueryCache()
	values, ok = c.queryCache[key]
	return
}

// QueryMap returns a map for a given query key.
func (c *Context) QueryMap(key string) (dict map[string]string) {
	dict, _ = c.GetQueryMap(key)
	return
}

// GetQueryMap returns a map for a given query key, plus a boolean value
// whether at least one value exists for the given key.
func (c *Context) GetQueryMap(key string) (map[string]string, bool) {
	c.initQueryCache()
	return getStringMapFromCache(c.queryCache, key)
}

// QueryMapArray returns a map for a given query key.
func (c *Context) QueryMapArray(key string) (dict map[string][]string) {
	dict, _ = c.GetQueryMapArray(key)
	return
}

// GetQueryMapArray returns a map for a given query key, plus a boolean value
// whether at least one value exists for the given key.
func (c *Context) GetQueryMapArray(key string) (map[string][]string, bool) {
	c.initQueryCache()
	return getStringsMapFromCache(c.queryCache, key)
}

// PostForms returns the POST urlencoded form or multipart form value map
func (c *Context) PostForms() map[string][]string {
	c.initFormCache()
	return c.formCache
}

// PostForm returns the specified key from a POST urlencoded form or multipart form
// when it exists, otherwise it returns the first non-empty value of defs or an empty string `("")`.
func (c *Context) PostForm(key string, defs ...string) string {
	if value, ok := c.GetPostForm(key); ok {
		return value
	}
	return str.NonEmpty(defs...)
}

// GetPostForm is like PostForm(key). It returns the specified key from a POST urlencoded
// form or multipart form when it exists `(value, true)` (even when the value is an empty string),
// otherwise it returns ("", false).
// For example, during a PATCH request to update the user's email:
//
// email=mail@example.com  -->  ("mail@example.com", true) := GetPostForm("email")
// email=                  -->  ("", true) := GetPostForm("email")
// =                       -->  ("", false) := GetPostForm("email")
func (c *Context) GetPostForm(key string) (string, bool) {
	if values, ok := c.GetPostFormArray(key); ok {
		return values[0], ok
	}
	return "", false
}

// PostFormArray returns a slice of strings for a given form key.
// The length of the slice depends on the number of params with the given key.
func (c *Context) PostFormArray(key string) (values []string) {
	values, _ = c.GetPostFormArray(key)
	return
}

// GetPostFormArray returns a slice of strings for a given form key, plus
// a boolean value whether at least one value exists for the given key.
func (c *Context) GetPostFormArray(key string) (values []string, ok bool) {
	c.initFormCache()
	values, ok = c.formCache[key]
	return
}

// PostFormMap returns a map for a given form key.
func (c *Context) PostFormMap(key string) (dict map[string]string) {
	dict, _ = c.GetPostFormMap(key)
	return
}

// GetPostFormMap returns a map for a given form key, plus a boolean value
// whether at least one value exists for the given key.
func (c *Context) GetPostFormMap(key string) (map[string]string, bool) {
	c.initFormCache()
	return getStringMapFromCache(c.formCache, key)
}

// PostFormMapArray returns a map for a given form key.
func (c *Context) PostFormMapArray(key string) (dict map[string][]string) {
	dict, _ = c.GetPostFormMapArray(key)
	return
}

// GetPostFormMapArray returns a map for a given form key, plus a boolean value
// whether at least one value exists for the given key.
func (c *Context) GetPostFormMapArray(key string) (map[string][]string, bool) {
	c.initFormCache()
	return getStringsMapFromCache(c.formCache, key)
}

// FormFile returns the first file for the provided form key.
func (c *Context) FormFile(key string) (*multipart.FileHeader, error) {
	mf, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}

	if fs := mf.File[key]; len(fs) > 0 {
		return fs[0], nil
	}
	return nil, http.ErrMissingFile
}

// FormFiles returns the files for the provided form key.
func (c *Context) FormFiles(key string) ([]*multipart.FileHeader, error) {
	mf, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}

	if fs := mf.File[key]; len(fs) > 0 {
		return fs, nil
	}
	return nil, http.ErrMissingFile
}

// MultipartForm is the parsed multipart form, including file uploads.
func (c *Context) MultipartForm() (*multipart.Form, error) {
	if c.Request.MultipartForm == nil {
		if err := c.Request.ParseMultipartForm(c.engine.MaxMultipartMemory); err != nil {
			return nil, err
		}
	}
	return c.Request.MultipartForm, nil
}

// SaveUploadedFile save uploaded file to the specific local file dst.
func (c *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	return httpx.SaveMultipartFile(file, dst)
}

// CopyUploadedFile read uploaded file to the specific buffer 'dst'.
func (c *Context) CopyUploadedFile(file *multipart.FileHeader, dst io.Writer) error {
	return httpx.CopyMultipartFile(file, dst)
}

// ReadUploadedFile read uploaded file and return it's content []byte.
func (c *Context) ReadUploadedFile(file *multipart.FileHeader) ([]byte, error) {
	return httpx.ReadMultipartFile(file)
}

// MustBind checks the Content-Type to select a binding engine automatically,
// Depending on the "Content-Type" header different bindings are used:
//
//	"application/json" --> JSON binding
//	"application/xml"  --> XML binding
//
// otherwise --> returns an error.
// It parses the request's body as JSON if Content-Type == "application/json" using JSON or XML as a JSON input.
// It decodes the json payload into the struct specified as a pointer.
// It writes a 400 error and sets Content-Type header "text/plain" in the response if input is not valid.
func (c *Context) MustBind(obj any) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.MustBindWith(obj, b)
}

// MustBindJSON is a shortcut for c.MustBindWith(obj, binding.JSON).
func (c *Context) MustBindJSON(obj any) error {
	return c.MustBindWith(obj, binding.JSON)
}

// MustBindXML is a shortcut for c.MustBindWith(obj, binding.BindXML).
func (c *Context) MustBindXML(obj any) error {
	return c.MustBindWith(obj, binding.XML)
}

// MustBindQuery is a shortcut for c.MustBindWith(obj, binding.Query).
func (c *Context) MustBindQuery(obj any) error {
	return c.MustBindWith(obj, binding.Query)
}

// MustBindPlain is a shortcut for c.MustBindWith(obj, binding.Plain).
func (c *Context) MustBindPlain(obj any) error {
	return c.MustBindWith(obj, binding.Plain)
}

// MustBindHeader is a shortcut for c.MustBindWith(obj, binding.Header).
func (c *Context) MustBindHeader(obj any) error {
	return c.MustBindWith(obj, binding.Header)
}

// MustBindURI binds the passed struct pointer using binding.Uri.
// It will abort the request with HTTP 400 if any error occurs.
func (c *Context) MustBindURI(obj any) error {
	if err := c.BindURI(obj); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return err
	}
	return nil
}

// MustBindWith binds the passed struct pointer using the specified binding engine.
// It will abort the request with HTTP 400 if any error occurs.
// See the binding package.
func (c *Context) MustBindWith(obj any, b binding.Binding) error {
	if err := c.BindWith(obj, b); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return err
	}
	return nil
}

// Bind checks the Content-Type to select a binding engine automatically,
// Depending on the "Content-Type" header different bindings are used:
//
//	"application/json" --> JSON binding
//	"application/xml"  --> XML binding
//
// otherwise --> returns an error
// It parses the request's body as JSON if Content-Type == "application/json" using JSON or XML as a JSON input.
// It decodes the json payload into the struct specified as a pointer.
// Like c.Bind() but this method does not set the response status code to 400 and abort if the json is not valid.
func (c *Context) Bind(obj any) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.BindWith(obj, b)
}

// BindJSON is a shortcut for c.BindWith(obj, binding.JSON).
func (c *Context) BindJSON(obj any) error {
	return c.BindWith(obj, binding.JSON)
}

// BindXML is a shortcut for c.BindWith(obj, binding.XML).
func (c *Context) BindXML(obj any) error {
	return c.BindWith(obj, binding.XML)
}

// BindQuery is a shortcut for c.BindWith(obj, binding.Query).
func (c *Context) BindQuery(obj any) error {
	return c.BindWith(obj, binding.Query)
}

// BindPlain is a shortcut for c.BindWith(obj, binding.Plain).
func (c *Context) ShouldBindPlain(obj any) error {
	return c.BindWith(obj, binding.Plain)
}

// BindHeader is a shortcut for c.BindWith(obj, binding.Header).
func (c *Context) BindHeader(obj any) error {
	return c.BindWith(obj, binding.Header)
}

// BindURI binds the passed struct pointer using the specified binding engine.
func (c *Context) BindURI(obj any) (err error) {
	m := make(map[string][]string, len(c.Params))
	for _, v := range c.Params {
		m[v.Key] = []string{v.Value}
	}

	err = binding.URI.BindURI(m, obj)
	if err != nil {
		return err
	}

	err = c.engine.Validator.ValidateStruct(obj)
	return
}

// BindWith binds the passed struct pointer using the specified binding engine.
// See the binding package.
func (c *Context) BindWith(obj any, b binding.Binding) (err error) {
	err = b.Bind(c.Request, obj)
	if err != nil {
		return err
	}

	err = c.engine.Validator.ValidateStruct(obj)
	return
}

// BindBodyWith is similar with BindWith, but it stores the request
// body into the context, and reuse when it is called again.
//
// NOTE: This method reads the body before binding. So you should use
// BindWith for better performance if you need to call only once.
func (c *Context) BindBodyWith(obj any, bb binding.BodyBinding) (err error) {
	var body []byte
	if cb, ok := c.Get(BodyBytesKey); ok {
		if cbb, ok := cb.([]byte); ok {
			body = cbb
		}
	}
	if body == nil {
		body, err = io.ReadAll(c.Request.Body)
		if err != nil {
			return
		}
		c.Set(BodyBytesKey, body)
	}

	err = bb.BindBody(body, obj)
	if err != nil {
		return
	}

	err = c.engine.Validator.ValidateStruct(obj)
	return
}

// BindBodyWithJSON is a shortcut for c.BindBodyWith(obj, binding.JSON).
func (c *Context) BindBodyWithJSON(obj any) error {
	return c.BindBodyWith(obj, binding.JSON)
}

// BindBodyWithXML is a shortcut for c.BindBodyWith(obj, binding.XML).
func (c *Context) BindBodyWithXML(obj any) error {
	return c.BindBodyWith(obj, binding.XML)
}

// BindBodyWithPlain is a shortcut for c.BindBodyWith(obj, binding.Plain).
func (c *Context) BindBodyWithPlain(obj any) error {
	return c.BindBodyWith(obj, binding.Plain)
}

// RequestHostname use net.SplitHostPort() to get host name only.
// If net.SplitHostPort(c.Request.Host) has error, return c.Request.Host .
func (c *Context) RequestHostname() string {
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err == nil {
		return host
	}
	return c.Request.Host
}

// IsSecure implements one best effort algorithm to check the https request.
// It checks the Request.URL.Scheme is 'https' or Request.TLS != nil.
// It called c.RemoteIP() under the hood, to check if the remote IP is a trusted proxy or not,
// then checks the SSLProxyHeaders (default "X-Forwarded-Proto": "https") to determine the https request.
func (c *Context) IsSecure() bool {
	r := c.Request
	if strings.EqualFold(r.URL.Scheme, "https") || r.TLS != nil {
		return true
	}

	// It also checks if the remoteIP is a trusted proxy or not.
	// In order to perform this validation, it will see if the IP is contained within at least one of the CIDR blocks
	// defined by Engine.SetTrustedProxies()
	remoteIP := net.ParseIP(c.RemoteIP())
	if remoteIP == nil || !c.engine.isTrustedProxy(remoteIP) {
		return false
	}

	sslProxyHeaders := c.engine.SSLProxyHeaders
	for k, v := range sslProxyHeaders {
		if c.GetHeader(k) == v {
			return true
		}
	}

	return false
}

// ClientIP implements one best effort algorithm to return the real client IP.
// It called c.RemoteIP() under the hood, to check if the remote IP is a trusted proxy or not.
// If it is it will then try to parse the headers defined in Engine.RemoteIPHeaders (defaulting to [X-Forwarded-For, X-Real-Ip]).
// If the headers are not syntactically valid OR the remote IP does not correspond to a trusted proxy,
// the remote IP (coming from Request.RemoteAddr) is returned.
func (c *Context) ClientIP() string {
	// Check if we're running on a trusted platform, continue running backwards if error
	if c.engine.TrustedIPHeader != "" {
		// Developers can define their own header of Trusted Platform or use predefined constants
		if addr := c.GetHeader(c.engine.TrustedIPHeader); addr != "" {
			return addr
		}
	}

	// It also checks if the remoteIP is a trusted proxy or not.
	// In order to perform this validation, it will see if the IP is contained within at least one of the CIDR blocks
	// defined by Engine.SetTrustedProxies()
	remoteIP := net.ParseIP(c.RemoteIP())
	if remoteIP == nil {
		return ""
	}

	if len(c.engine.RemoteIPHeaders) > 0 && c.engine.isTrustedProxy(remoteIP) {
		for _, headerName := range c.engine.RemoteIPHeaders {
			ip, valid := c.engine.validateClientIP(c.GetHeader(headerName))
			if valid {
				return ip
			}
		}
	}
	return remoteIP.String()
}

// RemoteIP parses the IP from Request.RemoteAddr, normalizes and returns the IP (without the port).
func (c *Context) RemoteIP() string {
	ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr))
	if err != nil {
		return ""
	}
	return ip
}

// ContentType returns the Content-Type header of the request.
func (c *Context) ContentType() string {
	return filterFlags(c.GetHeader("Content-Type"))
}

// IsWebsocket returns true if the request headers indicate that a websocket
// handshake is being initiated by the client.
func (c *Context) IsWebsocket() bool {
	if strings.Contains(strings.ToLower(c.GetHeader("Connection")), "upgrade") &&
		strings.EqualFold(c.GetHeader("Upgrade"), "websocket") {
		return true
	}
	return false
}

/************************************/
/******** RESPONSE RENDERING ********/
/************************************/

// bodyAllowedForStatus is a copy of http.bodyAllowedForStatus non-exported function.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}

// Status sets the HTTP response code.
func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

// Header is an intelligent shortcut for c.Writer.Header().Set(key, value).
// It writes a header in the response.
// If value == "", this method removes the header `c.Writer.Header().Del(key)`
func (c *Context) Header(key, value string) {
	if value == "" {
		c.Writer.Header().Del(key)
		return
	}
	c.Writer.Header().Set(key, value)
}

// GetHeader returns value from request headers.
func (c *Context) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

// GetRawData returns stream data.
func (c *Context) GetRawData() ([]byte, error) {
	if c.Request.Body == nil {
		return nil, errors.New("cannot read nil body")
	}
	return io.ReadAll(c.Request.Body)
}

// SetCookie adds a Set-Cookie header to the ResponseWriter's headers.
// The provided cookie must have a valid Name. Invalid cookies may be
// silently dropped.
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Writer, cookie)
}

// Cookie returns the named cookie provided in the request or
// ErrNoCookie if not found. And return the named cookie is unescaped.
// If multiple cookies match the given name, only one cookie will
// be returned.
func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// Render writes the response headers and calls render.Render to render data.
func (c *Context) Render(code int, r render.Render) {
	c.Status(code)

	if !bodyAllowedForStatus(code) {
		r.WriteContentType(c.Writer)
		c.Writer.WriteHeaderNow()
		return
	}

	if err := r.Render(c.Writer); err != nil {
		panic(err)
	}
}

// HTML renders the HTTP template specified by its file name.
// It also updates the HTTP code and sets the Content-Type as "text/html".
// See http://golang.org/doc/articles/wiki/
func (c *Context) HTML(code int, name string, obj any) {
	instance := c.engine.HTMLTemplates.Instance(name, obj)
	c.Render(code, instance)
}

// IndentedJSON serializes the given struct as pretty JSON (indented + endlines) into the response body.
// It also sets the Content-Type as "application/json".
// WARNING: we recommend using this only for development purposes since printing pretty JSON is
// more CPU and bandwidth consuming. Use Context.JSON() instead.
func (c *Context) IndentedJSON(code int, obj any) {
	c.Render(code, render.IndentedJSON{Data: obj})
}

// SecureJSON serializes the given struct as Secure JSON into the response body.
// Default prepends "while(1)," to response body if the given struct is array values.
// It also sets the Content-Type as "application/json".
func (c *Context) SecureJSON(code int, obj any) {
	c.Render(code, render.SecureJSON{Prefix: c.engine.secureJSONPrefix, Data: obj})
}

// JSONP serializes the given struct as JSON into the response body.
// It adds padding to response body to request data from a server residing in a different domain than the client.
// It also sets the Content-Type as "application/javascript".
func (c *Context) JSONP(code int, obj any) {
	callback := c.Query("callback")
	if callback == "" {
		c.Render(code, render.JSON{Data: obj})
		return
	}
	c.Render(code, render.JsonpJSON{Callback: callback, Data: obj})
}

// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func (c *Context) JSON(code int, obj any) {
	c.Render(code, render.JSON{Data: obj})
}

// AsciiJSON serializes the given struct as JSON into the response body with unicode to ASCII string.
// It also sets the Content-Type as "application/json".
func (c *Context) AsciiJSON(code int, obj any) {
	c.Render(code, render.AsciiJSON{Data: obj})
}

// PureJSON serializes the given struct as JSON into the response body.
// PureJSON, unlike JSON, does not replace special html characters with their unicode entities.
func (c *Context) PureJSON(code int, obj any) {
	c.Render(code, render.PureJSON{Data: obj})
}

// XML serializes the given struct as XML into the response body.
// It also sets the Content-Type as "application/xml".
func (c *Context) XML(code int, obj any) {
	c.Render(code, render.XML{Data: obj})
}

// String writes the given string into the response body.
func (c *Context) String(code int, format string, values ...any) {
	c.Render(code, render.String{Format: format, Data: values})
}

// Redirect returns an HTTP redirect to the specific location.
func (c *Context) Redirect(code int, location string) {
	c.Render(-1, render.Redirect{
		Code:     code,
		Location: location,
		Request:  c.Request,
	})
}

// Data writes some data into the body stream and updates the HTTP code.
func (c *Context) Data(code int, contentType string, data []byte) {
	c.Render(code, render.Data{
		ContentType: contentType,
		Data:        data,
	})
}

// DataFromReader writes the specified reader into the body stream and updates the HTTP code.
func (c *Context) DataFromReader(code int, contentLength int64, contentType string, reader io.Reader, extraHeaders map[string]string) {
	c.Render(code, render.Reader{
		Headers:       extraHeaders,
		ContentType:   contentType,
		ContentLength: contentLength,
		Reader:        reader,
	})
}

// File writes the specified file into the body stream in an efficient way.
func (c *Context) File(filepath string) {
	http.ServeFile(c.Writer, c.Request, filepath)
}

// FileFromFS writes the specified file from http.FileSystem into the body stream in an efficient way.
func (c *Context) FileFromFS(filepath string, fs http.FileSystem) {
	defer func(old string) {
		c.Request.URL.Path = old
	}(c.Request.URL.Path)

	c.Request.URL.Path = filepath

	http.FileServer(fs).ServeHTTP(c.Writer, c.Request)
}

// SetAttachmentHeader set response header Content-Disposition: attachment; filename=...
func (c *Context) SetAttachmentHeader(filename string) {
	httpx.SetAttachmentHeader(c.Writer.Header(), filename)
}

// FileAttachment writes the specified file into the body stream in an efficient way
// On the client side, the file will typically be downloaded with the given filename
func (c *Context) FileAttachment(filepath, filename string) {
	c.SetAttachmentHeader(filename)
	http.ServeFile(c.Writer, c.Request, filepath)
}

// SSEvent writes a Server-Sent Event into the body stream.
func (c *Context) SSEvent(name string, message any) {
	c.Render(-1, sse.Event{
		Event: name,
		Data:  message,
	})
}

// Stream sends a streaming response and returns a boolean
// indicates "Is client disconnected in middle of stream"
func (c *Context) Stream(step func(w io.Writer) bool) bool {
	w := c.Writer
	clientGone := w.CloseNotify()
	for {
		select {
		case <-clientGone:
			return true
		default:
			keepOpen := step(w)
			w.Flush()
			if !keepOpen {
				return false
			}
		}
	}
}

/************************************/
/******** CONTENT NEGOTIATION *******/
/************************************/

// Negotiate contains all negotiations data.
type Negotiate struct {
	Offered  []string
	HTMLName string
	HTMLData any
	JSONData any
	XMLData  any
	Data     any
}

// Negotiate calls different Render according to acceptable Accept format.
func (c *Context) Negotiate(code int, config Negotiate) {
	switch c.NegotiateFormat(config.Offered...) {
	case binding.MIMEJSON:
		data := chooseData(config.JSONData, config.Data)
		c.JSON(code, data)

	case binding.MIMEHTML:
		data := chooseData(config.HTMLData, config.Data)
		c.HTML(code, config.HTMLName, data)

	case binding.MIMEXML:
		data := chooseData(config.XMLData, config.Data)
		c.XML(code, data)

	default:
		c.AbortWithError(http.StatusNotAcceptable, errors.New("the accepted formats are not offered by the server"))
	}
}

// NegotiateFormat returns an acceptable Accept format.
func (c *Context) NegotiateFormat(offered ...string) string {
	if len(offered) == 0 {
		panic("you must provide at least one offer")
	}

	if c.Accepted == nil {
		c.Accepted = parseAccept(c.GetHeader("Accept"))
	}
	if len(c.Accepted) == 0 {
		return offered[0]
	}
	for _, accepted := range c.Accepted {
		for _, offer := range offered {
			// According to RFC 2616 and RFC 2396, non-ASCII characters are not allowed in headers,
			// therefore we can just iterate over the string without casting it into []rune
			i := 0
			for ; i < len(accepted) && i < len(offer); i++ {
				if accepted[i] == '*' || offer[i] == '*' {
					return offer
				}
				if accepted[i] != offer[i] {
					break
				}
			}
			if i == len(accepted) {
				return offer
			}
		}
	}
	return ""
}

// SetAccepted sets Accept header data.
func (c *Context) SetAccepted(formats ...string) {
	c.Accepted = formats
}

/************************************/
/***** GOLANG.ORG/X/NET/CONTEXT *****/
/************************************/

// hasRequestContext returns whether c.Request has Context and fallback.
func (c *Context) hasRequestContext() bool {
	hasFallback := c.engine != nil && c.engine.ContextWithFallback
	hasRequestContext := c.Request != nil && c.Request.Context() != nil
	return hasFallback && hasRequestContext
}

// Deadline returns that there is no deadline (ok==false) when c.Request has no Context.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	if !c.hasRequestContext() {
		return
	}
	return c.Request.Context().Deadline()
}

// Done returns nil (chan which will wait forever) when c.Request has no Context.
func (c *Context) Done() <-chan struct{} {
	if !c.hasRequestContext() {
		return nil
	}
	return c.Request.Context().Done()
}

// Err returns nil when c.Request has no Context.
func (c *Context) Err() error {
	if !c.hasRequestContext() {
		return nil
	}
	return c.Request.Context().Err()
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
func (c *Context) Value(key any) any {
	if key == RequestKey {
		return c.Request
	}
	if key == ContextKey {
		return c
	}
	if keyAsString, ok := key.(string); ok {
		if val, exists := c.Get(keyAsString); exists {
			return val
		}
	}
	if !c.hasRequestContext() {
		return nil
	}
	return c.Request.Context().Value(key)
}
