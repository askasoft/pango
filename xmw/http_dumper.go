package xmw

//nolint: gosec
import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/xin"
)

const dumpTimeFormat = "2006-01-02T15:04:05.000"

// HTTPDumper dump http request and response
type HTTPDumper struct {
	outputer  io.Writer
	disabled  bool
	maxlength int64
}

// DefaultHTTPDumper create a middleware for xin http dumper
// Equals: NewHTTPDumper(xin.Logger.Outputer("XHD", log.LevelInfo))
func DefaultHTTPDumper(xin *xin.Engine) *HTTPDumper {
	return NewHTTPDumper(xin.Logger.GetOutputer("XHD", log.LevelInfo))
}

// NewHTTPDumper create a middleware for xin http dumper
func NewHTTPDumper(outputer io.Writer) *HTTPDumper {
	return &HTTPDumper{outputer: outputer, maxlength: 1 << 20}
}

// Disable disable the dumper or not
func (hd *HTTPDumper) Disable(disabled bool) {
	hd.disabled = disabled
}

// SetMaxlength set the maxlength of request/resonse that the dumper should dump
func (hd *HTTPDumper) SetMaxlength(maxlength int64) {
	hd.maxlength = maxlength
}

// Handler returns the xin.HandlerFunc
func (hd *HTTPDumper) Handler() xin.HandlerFunc {
	return hd.Handle
}

// Handle process xin request
func (hd *HTTPDumper) Handle(c *xin.Context) {
	w := hd.outputer
	if w == nil || hd.disabled {
		c.Next()
		return
	}

	id := hd.dumpRequest(w, c.Request)

	cw := c.Writer
	dw := newHttpDumpWriter(c.Writer, hd.maxlength)
	c.Writer = dw

	c.Next()

	hd.dumpResponse(w, id, c.Request, dw)
	c.Writer = cw
}

// SetOutput set the access log output writer
func (hd *HTTPDumper) SetOutput(w io.Writer) {
	hd.outputer = w
}

const eol = "\r\n"

func (hd *HTTPDumper) dumpRequest(w io.Writer, req *http.Request) string {
	db := (req.ContentLength >= 0 && req.ContentLength <= hd.maxlength)
	bs, _ := httputil.DumpRequest(req, db)

	id := fmt.Sprintf("%x", md5.Sum(bs)) //nolint: gosec

	bb := &bytes.Buffer{}

	bb.WriteString(fmt.Sprintf(">>>>>>>> %s %s >>>>>>>>", time.Now().Format(dumpTimeFormat), id))
	bb.WriteString(eol)
	if len(bs) > 0 {
		bb.Write(bs)
	}
	bb.WriteString(eol)
	bb.WriteString(eol)

	w.Write(bb.Bytes()) //nolint: errcheck

	return id
}

func (hd *HTTPDumper) dumpResponse(w io.Writer, id string, req *http.Request, hdw *httpDumpWriter) {
	bb := &bytes.Buffer{}

	bb.WriteString(fmt.Sprintf("<<<<<<<< %s %s <<<<<<<<", time.Now().Format(dumpTimeFormat), id))
	bb.WriteString(eol)

	res := &http.Response{
		Proto:      req.Proto,
		ProtoMajor: req.ProtoMajor,
		ProtoMinor: req.ProtoMinor,
	}
	res.StatusCode = hdw.ResponseWriter.Status()
	res.Header = hdw.ResponseWriter.Header()
	res.Body = io.NopCloser(hdw.bb)
	res.Write(bb) //nolint: errcheck

	bb.WriteString(eol)
	bb.WriteString(eol)

	w.Write(bb.Bytes()) //nolint: errcheck
}

type httpDumpWriter struct {
	xin.ResponseWriter
	bb *bytes.Buffer
	lw io.Writer
}

func newHttpDumpWriter(xrw xin.ResponseWriter, maxlength int64) *httpDumpWriter {
	bb := &bytes.Buffer{}
	hdw := &httpDumpWriter{xrw, bb, iox.LimitWriter(bb, maxlength)}
	return hdw
}

func (hdw *httpDumpWriter) Write(data []byte) (int, error) {
	hdw.lw.Write(data) //nolint: errcheck
	return hdw.ResponseWriter.Write(data)
}
