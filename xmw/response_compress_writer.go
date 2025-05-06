package xmw

import (
	"bytes"
	"io"

	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin"
)

const (
	rcwStateNone = iota
	rcwStateSkip // skip
	rcwStateBuff // buffering
	rcwStateComp // compressing
)

type responseCompressWriter struct {
	xin.ResponseWriter

	enc string
	xrc *ResponseCompressor
	ctx *xin.Context
	buf *bytes.Buffer
	cw  Compressor

	state int
}

func (rcw *responseCompressWriter) Close() {
	if rcw.buf != nil {
		if rcw.buf.Len() > 0 {
			rcw.ResponseWriter.Write(rcw.buf.Bytes()) //nolint: errcheck
			rcw.buf.Reset()
		}
		rcw.xrc.putBuffer(rcw.buf)
		rcw.buf = nil
	}

	if rcw.cw != nil {
		rcw.cw.Close()
		rcw.cw.Reset(io.Discard)
		if _, ok := rcw.cw.(*passCompressor); !ok {
			rcw.xrc.putCompressor(rcw.enc, rcw.cw)
		}
		rcw.cw = nil
	}

	rcw.ctx.Writer = rcw.ResponseWriter

	rcw.xrc = nil
	rcw.ResponseWriter = nil
	rcw.ctx = nil
	rcw.state = rcwStateNone
}

func (rcw *responseCompressWriter) checkHeader() {
	if rcw.state != rcwStateNone {
		return
	}

	h := rcw.Header()
	if h.Get("Content-Encoding") != "" {
		rcw.state = rcwStateSkip
		return
	}

	mts := rcw.xrc.mimeTypes
	if mts != nil {
		ct := str.SubstrBeforeByte(h.Get("Content-Type"), ';')
		if !mts.Contains(ct) {
			rcw.state = rcwStateSkip
			return
		}
	}

	if rcw.ctx.Request.Header.Get("Via") != "" {
		if rcw.xrc.proxied&ProxiedExpired == ProxiedExpired {
			if h.Get("Expires") == "" {
				rcw.state = rcwStateSkip
				return
			}
		}
		if rcw.xrc.proxied&ProxiedNoCache == ProxiedNoCache {
			if !str.ContainsFold(h.Get("Cache-Control"), "no-cache") {
				rcw.state = rcwStateSkip
				return
			}
		}
		if rcw.xrc.proxied&ProxiedNoStore == ProxiedNoStore {
			if !str.ContainsFold(h.Get("Cache-Control"), "no-store") {
				rcw.state = rcwStateSkip
				return
			}
		}
		if rcw.xrc.proxied&ProxiedPrivate == ProxiedPrivate {
			if !str.ContainsFold(h.Get("Cache-Control"), "private") {
				rcw.state = rcwStateSkip
				return
			}
		}
		if rcw.xrc.proxied&ProxiedNoLastModified == ProxiedNoLastModified {
			if h.Get("Last-Modified") != "" {
				rcw.state = rcwStateSkip
				return
			}
		}
		if rcw.xrc.proxied&ProxiedNoETag == ProxiedNoETag {
			if h.Get("ETag") != "" {
				rcw.state = rcwStateSkip
				return
			}
		}
	}

	rcw.state = rcwStateBuff
	rcw.buf = rcw.xrc.getBuffer()
}

func (rcw *responseCompressWriter) checkBuffer(data []byte) (err error) {
	if rcw.state != rcwStateBuff {
		return
	}

	if rcw.buf.Len()+len(data) < rcw.xrc.minLength {
		return
	}

	rcw.state = rcwStateComp

	rcw.cw = rcw.xrc.getCompressor(rcw.enc)
	if rcw.cw == nil {
		rcw.cw = &passCompressor{rcw.ResponseWriter}
	} else {
		rcw.modifyHeader()
		rcw.cw.Reset(rcw.ResponseWriter)
	}

	if rcw.buf.Len() > 0 {
		_, err = rcw.cw.Write(rcw.buf.Bytes())
		rcw.buf.Reset()
	}
	return
}

func (rcw *responseCompressWriter) modifyHeader() {
	h := rcw.Header()

	h.Del("Content-Length")
	h.Set("Content-Encoding", rcw.enc)
	if rcw.xrc.vary {
		h.Set("Vary", "Accept-Encoding")
	}
}

// implements xin.ResponseWriter
func (rcw *responseCompressWriter) WriteString(s string) (int, error) {
	return rcw.Write(str.UnsafeBytes(s))
}

// implements http.ResponseWriter
func (rcw *responseCompressWriter) Write(data []byte) (int, error) {
	rcw.checkHeader()

	if err := rcw.checkBuffer(data); err != nil {
		return 0, err
	}

	switch rcw.state {
	case rcwStateBuff:
		return rcw.buf.Write(data)
	case rcwStateComp:
		return rcw.cw.Write(data)
	case rcwStateSkip:
		return rcw.ResponseWriter.Write(data)
	}

	panic("ResponseCompressor: Invalid State")
}

// Flush implements the http.Flush interface.
func (rcw *responseCompressWriter) Flush() {
	switch rcw.state {
	case rcwStateComp:
		rcw.cw.Flush()
		rcw.ResponseWriter.Flush()
	case rcwStateSkip:
		rcw.ResponseWriter.Flush()
	}
}
