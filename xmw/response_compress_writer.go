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

type compressor interface {
	io.Writer
	Reset(w io.Writer)
	Flush() error
	Close() error
}

type responseCompressWriter struct {
	xin.ResponseWriter

	ae  acceptEncoding
	rc  *ResponseCompressor
	ctx *xin.Context
	buf *bytes.Buffer
	cw  compressor

	state int
}

func (rcw *responseCompressWriter) Close() {
	if rcw.buf != nil {
		if rcw.buf.Len() > 0 {
			rcw.ResponseWriter.Write(rcw.buf.Bytes()) //nolint: errcheck
			rcw.buf.Reset()
		}
		rcw.rc.putBuffer(rcw.buf)
		rcw.buf = nil
	}

	if rcw.cw != nil {
		rcw.cw.Close()
		rcw.cw.Reset(io.Discard)
		rcw.rc.putCompressor(rcw.ae, rcw.cw)
		rcw.cw = nil
	}

	rcw.ctx.Writer = rcw.ResponseWriter

	rcw.rc = nil
	rcw.ResponseWriter = nil
	rcw.ctx = nil
	rcw.state = rcwStateNone
}

func (rcw *responseCompressWriter) checkHeader() {
	if rcw.state != rcwStateNone {
		return
	}

	h := rcw.ResponseWriter.Header()
	if h.Get("Content-Encoding") != "" {
		rcw.state = rcwStateSkip
		return
	}

	mts := rcw.rc.mimeTypes
	if mts != nil {
		ct := str.SubstrBeforeByte(h.Get("Content-Type"), ';')
		if !mts.Contains(ct) {
			rcw.state = rcwStateSkip
			return
		}
	}

	if rcw.ctx.Request.Header.Get("Via") != "" {
		if rcw.rc.proxied&ProxiedExpired == ProxiedExpired {
			if h.Get("Expires") == "" {
				rcw.state = rcwStateSkip
				return
			}
		}
		if rcw.rc.proxied&ProxiedNoCache == ProxiedNoCache {
			if !str.ContainsFold(h.Get("Cache-Control"), "no-cache") {
				rcw.state = rcwStateSkip
				return
			}
		}
		if rcw.rc.proxied&ProxiedNoStore == ProxiedNoStore {
			if !str.ContainsFold(h.Get("Cache-Control"), "no-store") {
				rcw.state = rcwStateSkip
				return
			}
		}
		if rcw.rc.proxied&ProxiedPrivate == ProxiedPrivate {
			if !str.ContainsFold(h.Get("Cache-Control"), "private") {
				rcw.state = rcwStateSkip
				return
			}
		}
		if rcw.rc.proxied&ProxiedNoLastModified == ProxiedNoLastModified {
			if h.Get("Last-Modified") != "" {
				rcw.state = rcwStateSkip
				return
			}
		}
		if rcw.rc.proxied&ProxiedNoETag == ProxiedNoETag {
			if h.Get("ETag") != "" {
				rcw.state = rcwStateSkip
				return
			}
		}
	}

	rcw.state = rcwStateBuff
	rcw.buf = rcw.rc.getBuffer()
}

func (rcw *responseCompressWriter) checkBuffer(data []byte) (err error) {
	if rcw.state != rcwStateBuff {
		return
	}

	if rcw.buf.Len()+len(data) < rcw.rc.minLength {
		return
	}

	rcw.modifyHeader()

	rcw.state = rcwStateComp
	rcw.cw = rcw.rc.getCompressor(rcw.ae)
	rcw.cw.Reset(rcw.ResponseWriter)
	if rcw.buf.Len() > 0 {
		_, err = rcw.cw.Write(rcw.buf.Bytes())
		rcw.buf.Reset()
	}
	return
}

func (rcw *responseCompressWriter) modifyHeader() {
	h := rcw.ResponseWriter.Header()

	h.Del("Content-Length")
	h.Set("Content-Encoding", rcw.ae.String())
	if rcw.rc.vary {
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
