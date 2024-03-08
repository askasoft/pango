package xmw

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin"
)

const (
	gzwStateNone = iota
	gzwStateSkip
	gzwStateBuff
	gzwStateGzip
)

type gzipWriter struct {
	xin.ResponseWriter

	hgz *HTTPGziper
	ctx *xin.Context
	gzw *gzip.Writer
	buf bytes.Buffer

	state int
}

func (g *gzipWriter) Close() {
	if g.buf.Len() > 0 {
		g.ResponseWriter.Write(g.buf.Bytes()) //nolint: errcheck
		g.buf.Reset()
	} else {
		g.gzw.Flush()
	}

	g.ctx.Writer = g.ResponseWriter

	g.reset()
}

func (g *gzipWriter) reset() {
	g.hgz = nil
	g.ResponseWriter = nil
	g.ctx = nil
	g.gzw.Reset(io.Discard)
	g.buf.Reset()
	g.state = gzwStateNone
}

func (g *gzipWriter) checkHeader() {
	if g.state != gzwStateNone {
		return
	}

	h := g.ResponseWriter.Header()
	if str.ContainsFold(h.Get("Content-Encoding"), "gzip") {
		g.state = gzwStateSkip
		return
	}

	mts := g.hgz.mimeTypes
	if mts != nil {
		ct := str.SubstrBeforeByte(h.Get("Content-Type"), ';')
		if !mts.Contains(ct) {
			g.state = gzwStateSkip
			return
		}
	}

	if g.ctx.Request.Header.Get("Via") != "" {
		if g.hgz.proxied&ProxiedExpired == ProxiedExpired {
			if h.Get("Expires") == "" {
				g.state = gzwStateSkip
				return
			}
		}
		if g.hgz.proxied&ProxiedNoCache == ProxiedNoCache {
			if !str.ContainsFold(h.Get("Cache-Control"), "no-cache") {
				g.state = gzwStateSkip
				return
			}
		}
		if g.hgz.proxied&ProxiedNoStore == ProxiedNoStore {
			if !str.ContainsFold(h.Get("Cache-Control"), "no-store") {
				g.state = gzwStateSkip
				return
			}
		}
		if g.hgz.proxied&ProxiedPrivate == ProxiedPrivate {
			if !str.ContainsFold(h.Get("Cache-Control"), "private") {
				g.state = gzwStateSkip
				return
			}
		}
		if g.hgz.proxied&ProxiedNoLastModified == ProxiedNoLastModified {
			if h.Get("Last-Modified") != "" {
				g.state = gzwStateSkip
				return
			}
		}
		if g.hgz.proxied&ProxiedNoETag == ProxiedNoETag {
			if h.Get("ETag") != "" {
				g.state = gzwStateSkip
				return
			}
		}
	}

	g.state = gzwStateBuff
}

func (g *gzipWriter) checkBuffer(data []byte) (err error) {
	if g.state != gzwStateBuff {
		return
	}

	if g.buf.Len()+len(data) < g.hgz.minLength {
		return
	}

	g.modifyHeader()
	g.gzw.Reset(g.ResponseWriter)
	if g.buf.Len() > 0 {
		_, err = g.gzw.Write(g.buf.Bytes())
		g.buf.Reset()
	}
	g.state = gzwStateGzip
	return
}

func (g *gzipWriter) modifyHeader() {
	h := g.ResponseWriter.Header()

	h.Del("Content-Length")
	h.Set("Content-Encoding", "gzip")
	if g.hgz.vary {
		h.Set("Vary", "Accept-Encoding")
	}
}

// implements xin.ResponseWriter
func (g *gzipWriter) WriteString(s string) (int, error) {
	return g.Write(str.UnsafeBytes(s))
}

// implements http.ResponseWriter
func (g *gzipWriter) Write(data []byte) (int, error) {
	g.checkHeader()

	if err := g.checkBuffer(data); err != nil {
		return 0, err
	}

	switch g.state {
	case gzwStateBuff:
		return g.buf.Write(data)
	case gzwStateGzip:
		return g.gzw.Write(data)
	case gzwStateSkip:
		return g.ResponseWriter.Write(data)
	}

	panic("hgz: invalid state")
}

// Flush implements the http.Flush interface.
func (g *gzipWriter) Flush() {
	switch g.state {
	case gzwStateGzip:
		g.gzw.Flush()
		g.ResponseWriter.Flush()
	case gzwStateSkip:
		g.ResponseWriter.Flush()
	}
}
