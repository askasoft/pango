package xmw

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"
	"regexp"
	"strings"
	"sync"

	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin"
)

// http://nginx.org/en/docs/http/ngx_http_gzip_module.html

const (
	BestCompression    = flate.BestCompression
	BestSpeed          = flate.BestSpeed
	DefaultCompression = flate.DefaultCompression
	NoCompression      = flate.NoCompression
)

// ResponseCompressor Compresses responses using the “gzip” method
type ResponseCompressor struct {
	compressLevel int

	// protoMajor Sets the minimum HTTP Major version of a request required to compress a response.
	// Default: 1
	protoMajor int

	// protoMinor Sets the minimum HTTP Minor version of a request required to compress a response.
	// Default: 1
	protoMinor int

	// Proxied Enables or disables compressing of responses for proxied requests depending on the request and response.
	// The fact that the request is proxied is determined by the presence of the “Via” request header field.
	// Default: any
	proxied ProxiedFlag

	// Vary Enables or disables inserting the “Vary: Accept-Encoding” response header field.
	// Default: true
	vary bool

	// the minimum length of a response that will be compressed.
	// Default: 1024
	minLength int

	// mimeTypes Enables compressing of responses for the specified MIME types.
	mimeTypes *stringSet

	// ignorePathPrefixs Ignored URL Path Prefixs
	ignorePathPrefixs prefixs

	// ignoreRegexps Ignored URL Path Regexp
	ignorePathRegexps regexps

	// disabled Disable the compressor
	// Default: false
	disabled bool

	bufPool *sync.Pool
	gzwPool *sync.Pool
	zlwPool *sync.Pool
}

// DefaultResponseCompressor create a default zipper
// = NewResponseCompressor(DefaultCompression, 1024)
func DefaultResponseCompressor() *ResponseCompressor {
	return NewResponseCompressor(DefaultCompression, 1024)
}

// NewResponseCompressor create a http response compressor
// proxied: ProxiedAny
// vary: true
// minLength: 1024
func NewResponseCompressor(compressLevel, minLength int) *ResponseCompressor {
	rc := &ResponseCompressor{
		compressLevel: compressLevel,
		protoMajor:    1,
		protoMinor:    1,
		proxied:       ProxiedAny,
		vary:          true,
		minLength:     minLength,
	}
	rc.bufPool = &sync.Pool{
		New: func() any {
			return &bytes.Buffer{}
		},
	}
	rc.gzwPool = &sync.Pool{
		New: rc.newGzipWriter,
	}
	rc.zlwPool = &sync.Pool{
		New: rc.newZlibWriter,
	}

	rc.SetMimeTypes(
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
	return rc
}

// SetHTTPVersion Sets the minimum HTTP Proto version of a request required to compress a response.
func (rc *ResponseCompressor) SetHTTPVersion(major, minor int) {
	rc.protoMajor = major
	rc.protoMinor = minor
}

// SetProxied Enables or disables compressing of responses for proxied requests depending on the request and response.
// The fact that the request is proxied is determined by the presence of the “Via” request header field.
// The directive accepts multiple parameters:
// off
//
//	disables compression for all proxied requests, ignoring other parameters;
//
// any (Default)
//
//	enables compression for all proxied requests.
//
// auth
//
//	enables compression if a request header includes the “Authorization” field;
//
// expired
//
//	enables compression if a response header includes the “Expires” field with a value that disables caching;
//
// no-cache
//
//	enables compression if a response header includes the “Cache-Control” field with the “no-cache” parameter;
//
// no-store
//
//	enables compression if a response header includes the “Cache-Control” field with the “no-store” parameter;
//
// private
//
//	enables compression if a response header includes the “Cache-Control” field with the “private” parameter;
//
// no_last_modified
//
//	enables compression if a response header does not include the “Last-Modified” field;
//
// no_etag
//
//	enables compression if a response header does not include the “ETag” field;
func (rc *ResponseCompressor) SetProxied(ps ...string) {
	rc.proxied = toProxiedFlag(ps...)
}

// Vary Enables or disables inserting the “Vary: Accept-Encoding” response header field.
// Default: true
func (rc *ResponseCompressor) Vary(vary bool) {
	rc.vary = vary
}

// SetMimeTypes Enables compressing of responses for the specified MIME types.
// The special value "*" matches any MIME type.
// Default:
//
//	text/html
//	text/plain
//	text/xml
//	text/css
//	text/javascript
//	text/json
//	text/comma-separated-values
//	text/tab-separated-values
//	application/xml
//	application/xhtml+xml
//	application/rss+xml
//	application/atom_xml
//	application/json
//	application/javascript
//	application/x-javascript
func (rc *ResponseCompressor) SetMimeTypes(mts ...string) {
	if len(mts) == 0 {
		rc.mimeTypes = nil
		return
	}

	hs := newStringSet(mts...)
	if hs.Contains("*") {
		hs = nil
	}
	rc.mimeTypes = hs
}

// IgnorePathPrefix ignore URL path prefix
func (rc *ResponseCompressor) IgnorePathPrefix(ps ...string) {
	rc.ignorePathPrefixs = ps
}

// IgnorePathRegexp ignore URL path regexp
func (rc *ResponseCompressor) IgnorePathRegexp(ps ...string) {
	rs := make([]*regexp.Regexp, len(ps))
	for i, p := range ps {
		rs[i] = regexp.MustCompile(p)
	}
	rc.ignorePathRegexps = rs
}

// Disable disable the gzip compress or not
func (rc *ResponseCompressor) Disable(disabled bool) {
	rc.disabled = disabled
}

// Handle process xin request
func (rc *ResponseCompressor) Handle(ctx *xin.Context) {
	ae := rc.shouldCompress(ctx)
	if ae == acceptEncodingNone {
		ctx.Next()
		return
	}

	rcw := &responseCompressWriter{
		ae:             ae,
		rc:             rc,
		ctx:            ctx,
		ResponseWriter: ctx.Writer,
	}

	ctx.Writer = rcw
	defer rcw.Close()

	ctx.Next()
}

func (rc *ResponseCompressor) getBuffer() *bytes.Buffer {
	return rc.bufPool.Get().(*bytes.Buffer)
}

func (rc *ResponseCompressor) putBuffer(buf *bytes.Buffer) {
	rc.bufPool.Put(buf)
}

func (rc *ResponseCompressor) newGzipWriter() any {
	w, err := gzip.NewWriterLevel(io.Discard, rc.compressLevel)
	if err != nil {
		panic(err)
	}
	return w
}

func (rc *ResponseCompressor) newZlibWriter() any {
	w, err := zlib.NewWriterLevel(io.Discard, rc.compressLevel)
	if err != nil {
		panic(err)
	}
	return w
}

func (rc *ResponseCompressor) getCompressor(ac acceptEncoding) compressor {
	switch ac {
	case acceptEncodingGzip:
		return rc.gzwPool.Get().(*gzip.Writer)
	case acceptEncodingDeflate:
		return rc.zlwPool.Get().(*zlib.Writer)
	default:
		panic("ResponseCompressor: Invalid Encoding")
	}
}

func (rc *ResponseCompressor) putCompressor(ac acceptEncoding, cw compressor) {
	switch ac {
	case acceptEncodingGzip:
		rc.gzwPool.Put(cw)
	case acceptEncodingDeflate:
		rc.zlwPool.Put(cw)
	default:
		panic("ResponseCompressor: Invalid Encoding")
	}
}

func (rc *ResponseCompressor) shouldCompress(c *xin.Context) (ae acceptEncoding) {
	req := c.Request

	if rc.disabled {
		return
	}

	if !req.ProtoAtLeast(rc.protoMajor, rc.protoMinor) {
		return
	}

	if str.ContainsFold(req.Header.Get("Connection"), "Upgrade") {
		return
	}

	if str.ContainsFold(req.Header.Get("Content-Type"), "text/event-stream") {
		return
	}

	if req.Header.Get("Via") != "" {
		if rc.proxied == ProxiedOff {
			return
		}
		if rc.proxied&ProxiedAuth == ProxiedAuth {
			if req.Header.Get("Authorization") == "" {
				return
			}
		}
	}

	if rc.ignorePathPrefixs.Contains(req.URL.Path) {
		return
	}
	if rc.ignorePathRegexps.Contains(req.URL.Path) {
		return
	}

	ac := req.Header.Get("Accept-Encoding")
	if str.ContainsFold(ac, "gzip") {
		ae = acceptEncodingGzip
	} else if str.ContainsFold(ac, "deflate") {
		ae = acceptEncodingDeflate
	}

	return
}

type acceptEncoding int

const (
	acceptEncodingNone acceptEncoding = iota
	acceptEncodingGzip
	acceptEncodingDeflate
)

func (ae acceptEncoding) String() string {
	switch ae {
	case acceptEncodingGzip:
		return "gzip"
	case acceptEncodingDeflate:
		return "deflate"
	default:
		return ""
	}
}

// ProxiedFlag Proxied flag
type ProxiedFlag int

// Proxied option flags
const (
	ProxiedAny ProxiedFlag = 1 << iota
	ProxiedAuth
	ProxiedExpired
	ProxiedNoCache
	ProxiedNoStore
	ProxiedPrivate
	ProxiedNoLastModified
	ProxiedNoETag
	ProxiedOff = 0
)

// String return level string
func (pf ProxiedFlag) String() string {
	if pf == ProxiedOff {
		return "off"
	}

	fs := make([]string, 0, 9)
	if pf&ProxiedAny == ProxiedAny {
		fs = append(fs, "any")
	}
	if pf&ProxiedAuth == ProxiedAuth {
		fs = append(fs, "auth")
	}
	if pf&ProxiedExpired == ProxiedExpired {
		fs = append(fs, "expired")
	}
	if pf&ProxiedNoCache == ProxiedNoCache {
		fs = append(fs, "no-cache")
	}
	if pf&ProxiedNoStore == ProxiedNoStore {
		fs = append(fs, "no-store")
	}
	if pf&ProxiedPrivate == ProxiedPrivate {
		fs = append(fs, "private")
	}
	if pf&ProxiedNoLastModified == ProxiedNoLastModified {
		fs = append(fs, "no_last_modified")
	}
	if pf&ProxiedNoETag == ProxiedNoETag {
		fs = append(fs, "no_etag")
	}

	return strings.Join(fs, " ")
}

// toProxiedFlag parse proxied flag from string
func toProxiedFlag(ps ...string) (pf ProxiedFlag) {
	for _, s := range ps {
		s = strings.ToLower(s)
		switch s {
		case "off":
			return ProxiedOff
		case "any":
			pf |= ProxiedAny
		case "auth":
			pf |= ProxiedAuth
		case "expired":
			pf |= ProxiedExpired
		case "no-cache":
			pf |= ProxiedNoCache
		case "no-store":
			pf |= ProxiedNoStore
		case "private":
			pf |= ProxiedPrivate
		case "no_last_modified":
			pf |= ProxiedNoLastModified
		case "no_etag":
			pf |= ProxiedNoETag
		}
	}

	return
}

type prefixs []string

func (ps prefixs) Contains(uri string) bool {
	for _, path := range ps {
		if strings.HasPrefix(uri, path) {
			return true
		}
	}
	return false
}

type regexps []*regexp.Regexp

func (rs regexps) Contains(uri string) bool {
	for _, re := range rs {
		if re.MatchString(uri) {
			return true
		}
	}
	return false
}
