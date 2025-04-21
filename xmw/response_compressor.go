package xmw

import (
	"bytes"
	"regexp"
	"strings"
	"sync"

	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin"
)

// http://nginx.org/en/docs/http/ngx_http_gzip_module.html

// ResponseCompressor Compresses responses using the “gzip” method
type ResponseCompressor struct {
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

	Encodings map[string]CompressorProvider
}

// DefaultResponseCompressor create a default zipper
// = NewResponseCompressor(1024)
func DefaultResponseCompressor() *ResponseCompressor {
	return NewResponseCompressor(1024)
}

// NewResponseCompressor create a http response compressor
// proxied: ProxiedAny
// vary: true
// minLength: 1024
func NewResponseCompressor(minLength int) *ResponseCompressor {
	xrc := &ResponseCompressor{
		protoMajor: 1,
		protoMinor: 1,
		proxied:    ProxiedAny,
		vary:       true,
		minLength:  minLength,
	}
	xrc.bufPool = &sync.Pool{
		New: func() any {
			return &bytes.Buffer{}
		},
	}
	xrc.Encodings = map[string]CompressorProvider{
		"gzip":    NewGzipCompressorProvider(),
		"deflate": NewZlibCompressorProvider(),
	}

	xrc.SetMimeTypes(
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
	return xrc
}

// SetHTTPVersion Sets the minimum HTTP Proto version of a request required to compress a response.
func (xrc *ResponseCompressor) SetHTTPVersion(major, minor int) {
	xrc.protoMajor = major
	xrc.protoMinor = minor
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
func (xrc *ResponseCompressor) SetProxied(ps ...string) {
	xrc.proxied = toProxiedFlag(ps...)
}

// Vary Enables or disables inserting the “Vary: Accept-Encoding” response header field.
// Default: true
func (xrc *ResponseCompressor) Vary(vary bool) {
	xrc.vary = vary
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
func (xrc *ResponseCompressor) SetMimeTypes(mts ...string) {
	if len(mts) == 0 {
		xrc.mimeTypes = nil
		return
	}

	hs := newStringSet(mts...)
	if hs.Contains("*") {
		hs = nil
	}
	xrc.mimeTypes = hs
}

// IgnorePathPrefix ignore URL path prefix
func (xrc *ResponseCompressor) IgnorePathPrefix(ps ...string) {
	xrc.ignorePathPrefixs = ps
}

// IgnorePathRegexp ignore URL path regexp
func (xrc *ResponseCompressor) IgnorePathRegexp(ps ...string) {
	rs := make([]*regexp.Regexp, len(ps))
	for i, p := range ps {
		rs[i] = regexp.MustCompile(p)
	}
	xrc.ignorePathRegexps = rs
}

// Disable disable the gzip compress or not
func (xrc *ResponseCompressor) Disable(disabled bool) {
	xrc.disabled = disabled
}

// Handle process xin request
func (xrc *ResponseCompressor) Handle(ctx *xin.Context) {
	encoding := xrc.getEncoding(ctx)
	if encoding == "" {
		ctx.Next()
		return
	}

	rcw := &responseCompressWriter{
		enc:            encoding,
		xrc:            xrc,
		ctx:            ctx,
		ResponseWriter: ctx.Writer,
	}

	ctx.Writer = rcw
	defer rcw.Close()

	ctx.Next()
}

func (xrc *ResponseCompressor) getBuffer() *bytes.Buffer {
	return xrc.bufPool.Get().(*bytes.Buffer)
}

func (xrc *ResponseCompressor) putBuffer(buf *bytes.Buffer) {
	xrc.bufPool.Put(buf)
}

func (xrc *ResponseCompressor) getCompressor(encoding string) (c Compressor) {
	if p, ok := xrc.Encodings[encoding]; ok {
		return p.GetCompressor()
	}
	return nil
}

func (xrc *ResponseCompressor) putCompressor(encoding string, c Compressor) {
	if p, ok := xrc.Encodings[encoding]; ok {
		p.PutCompressor(c)
	}
}

func (xrc *ResponseCompressor) getEncoding(c *xin.Context) (encoding string) {
	if xrc.disabled {
		return
	}

	req := c.Request

	if !req.ProtoAtLeast(xrc.protoMajor, xrc.protoMinor) {
		return
	}
	if str.ContainsFold(req.Header.Get("Connection"), "Upgrade") {
		return
	}
	if str.ContainsFold(req.Header.Get("Content-Type"), "text/event-stream") {
		return
	}

	if req.Header.Get("Via") != "" {
		if xrc.proxied == ProxiedOff {
			return
		}
		if xrc.proxied&ProxiedAuth == ProxiedAuth {
			if req.Header.Get("Authorization") == "" {
				return
			}
		}
	}

	if xrc.ignorePathPrefixs.Contains(req.URL.Path) {
		return
	}
	if xrc.ignorePathRegexps.Contains(req.URL.Path) {
		return
	}

	acs := str.FieldsByte(req.Header.Get("Accept-Encoding"), ',')
	for _, ac := range acs {
		ac = str.ToLower(str.Strip(str.SubstrBeforeByte(ac, ';')))
		if _, ok := xrc.Encodings[ac]; ok {
			encoding = ac
			break
		}
	}
	return
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
