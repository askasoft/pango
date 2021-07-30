package gingzip

import (
	"compress/gzip"
	"io"
	"net/http"
	"regexp"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pandafw/pango/col"
	"github.com/pandafw/pango/str"
)

// http://nginx.org/en/docs/http/ngx_http_gzip_module.html

// These constants are copied from the gzip package.
const (
	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
)

// Zipper Compresses responses using the “gzip” method
type Zipper struct {
	// protoMajor Sets the minimum HTTP Major version of a request required to compress a response.
	// Default: 1
	protoMajor int

	// protoMinor Sets the minimum HTTP Minor version of a request required to compress a response.
	// Default: 1
	protoMinor int

	// Proxied Enables or disables gzipping of responses for proxied requests depending on the request and response.
	// The fact that the request is proxied is determined by the presence of the “Via” request header field.
	// Default: any
	proxied ProxiedFlag

	// Vary Enables or disables inserting the “Vary: Accept-Encoding” response header field.
	// Default: true
	vary bool

	// the minimum length of a response that will be gzipped.
	// Default: 1024
	minLength int

	// CompressLevel Sets a gzip compression level of a response.
	// Default: DefaultCompression
	compressLevel int

	// mimeTypes Enables gzipping of responses for the specified MIME types.
	mimeTypes *col.HashSet

	// ignorePathPrefixs Ignored URL Path Prefixs
	ignorePathPrefixs prefixs

	// ignoreRegexps Ignored URL Path Regexp
	ignorePathRegexps regexps

	// disabled Disable gzip
	// Default: false
	disabled bool

	// pool gzip writer pool
	pool *sync.Pool
}

// Default create a default zipper
// = NewZipper(DefaultCompression, 1024)
func Default() *Zipper {
	return NewZipper(DefaultCompression, 1024)
}

// NewZipper create a zipper
// proxied: ProxiedAny
// vary: true
// minLength: 1024
func NewZipper(compressLevel, minLength int) *Zipper {
	z := &Zipper{
		protoMajor:    1,
		protoMinor:    1,
		proxied:       ProxiedAny,
		vary:          true,
		compressLevel: compressLevel,
		minLength:     minLength,
	}
	z.pool = &sync.Pool{
		New: func() interface{} {
			gw := &gzipWriter{}
			w, err := gzip.NewWriterLevel(io.Discard, z.compressLevel)
			if err != nil {
				panic(err)
			}
			gw.gzw = w
			return gw
		},
	}

	z.SetMimeTypes(
		"text/html",
		"text/plain",
		"text/xml",
		"text/css",
		"text/javascript",
		"text/json",
		"text/comma-separated-values",
		"text/tab-separated-values",
		"application/xml",
		"application/xhtml+xml",
		"application/rss+xml",
		"application/atom_xml",
		"application/json",
		"application/javascript",
		"application/x-javascript",
	)
	return z
}

// SetHTTPVersion Sets the minimum HTTP Proto version of a request required to compress a response.
func (z *Zipper) SetHTTPVersion(major, minor int) {
	z.protoMajor = major
	z.protoMinor = minor
}

// SetProxied Enables or disables gzipping of responses for proxied requests depending on the request and response.
// The fact that the request is proxied is determined by the presence of the “Via” request header field.
// The directive accepts multiple parameters:
// off
//     disables compression for all proxied requests, ignoring other parameters;
// any (Default)
//     enables compression for all proxied requests.
// auth
//     enables compression if a request header includes the “Authorization” field;
// expired
//     enables compression if a response header includes the “Expires” field with a value that disables caching;
// no-cache
//     enables compression if a response header includes the “Cache-Control” field with the “no-cache” parameter;
// no-store
//     enables compression if a response header includes the “Cache-Control” field with the “no-store” parameter;
// private
//     enables compression if a response header includes the “Cache-Control” field with the “private” parameter;
// no_last_modified
//     enables compression if a response header does not include the “Last-Modified” field;
// no_etag
//     enables compression if a response header does not include the “ETag” field;
func (z *Zipper) SetProxied(ps ...string) {
	z.proxied = toProxiedFlag(ps...)
}

// Vary Enables or disables inserting the “Vary: Accept-Encoding” response header field.
// Default: true
func (z *Zipper) Vary(vary bool) {
	z.vary = vary
}

// SetMimeTypes Enables gzipping of responses for the specified MIME types.
// The special value "*" matches any MIME type.
// Default:
//   text/html
//   text/plain
//   text/xml
//   text/css
//   text/javascript
//   text/json
//   text/comma-separated-values
//   text/tab-separated-values
//   application/xml
//   application/xhtml+xml
//   application/rss+xml
//   application/atom_xml
//   application/json
//   application/javascript
//   application/x-javascript
func (z *Zipper) SetMimeTypes(mts ...string) {
	if len(mts) == 0 {
		z.mimeTypes = nil
		return
	}

	hs := col.NewStringHashSet(mts...)
	if hs.Contains("*") {
		hs = nil
	}
	z.mimeTypes = hs
}

// IgnorePathPrefix ignore URL path prefix
func (z *Zipper) IgnorePathPrefix(ps ...string) {
	z.ignorePathPrefixs = ps
}

// IgnorePathRegexp ignore URL path regexp
func (z *Zipper) IgnorePathRegexp(ps ...string) {
	rs := make([]*regexp.Regexp, len(ps), len(ps))
	for i, p := range ps {
		rs[i] = regexp.MustCompile(p)
	}
	z.ignorePathRegexps = rs
}

// Disable disable the gzip compress or not
func (z *Zipper) Disable(disabled bool) {
	z.disabled = disabled
}

// Handler returns the gin.HandlerFunc
func (z *Zipper) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		z.handle(c)
	}
}

// handle process gin request
func (z *Zipper) handle(c *gin.Context) {
	if !z.shouldCompress(c.Request) {
		c.Next()
		return
	}

	gw := z.pool.Get().(*gzipWriter)
	gw.zipper = z
	gw.ctx = c
	gw.ResponseWriter = c.Writer

	c.Writer = gw
	defer func() {
		gw.Close()
		z.pool.Put(gw)
	}()
	c.Next()
}

func (z *Zipper) shouldCompress(req *http.Request) bool {
	if z.disabled ||
		!req.ProtoAtLeast(z.protoMajor, z.protoMinor) ||
		!str.ContainsFold(req.Header.Get("Accept-Encoding"), "gzip") ||
		str.ContainsFold(req.Header.Get("Connection"), "Upgrade") ||
		str.ContainsFold(req.Header.Get("Content-Type"), "text/event-stream") {

		return false
	}

	if req.Header.Get("Via") != "" {
		if z.proxied == ProxiedOff {
			return false
		}
		if z.proxied&ProxiedAuth == ProxiedAuth {
			if req.Header.Get("Authorization") == "" {
				return false
			}
		}
	}

	if z.ignorePathPrefixs.Contains(req.URL.Path) {
		return false
	}
	if z.ignorePathRegexps.Contains(req.URL.Path) {
		return false
	}

	return true
}
