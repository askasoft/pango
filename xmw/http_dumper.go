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

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/xin"
)

const dumpTimeFormat = "2006-01-02T15:04:05.000"

// HTTPDumper dump http request and response
type HTTPDumper struct {
	outputer io.Writer
	disabled bool
}

// DefaultHTTPDumper create a middleware for xin http dumper
// Equals: NewHTTPDumper(xin.Logger.Outputer("XIND", log.LevelTrace))
func DefaultHTTPDumper(xin *xin.Engine) *HTTPDumper {
	return NewHTTPDumper(xin.Logger.GetOutputer("XIND", log.LevelTrace))
}

// NewHTTPDumper create a middleware for xin http dumper
func NewHTTPDumper(outputer io.Writer) *HTTPDumper {
	return &HTTPDumper{outputer: outputer}
}

// Disable disable the dumper or not
func (hd *HTTPDumper) Disable(disabled bool) {
	hd.disabled = disabled
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

	// dump request
	id := hd.dumpRequest(w, c.Request)

	dw := &httpDumpWriter{c.Writer, &http.Response{
		Proto:      c.Request.Proto,
		ProtoMajor: c.Request.ProtoMajor,
		ProtoMinor: c.Request.ProtoMinor,
	}, &bytes.Buffer{}}
	c.Writer = dw

	// process request
	c.Next()

	// dump response
	hd.dumpResponse(w, id, dw)
}

// SetOutput set the access log output writer
func (hd *HTTPDumper) SetOutput(w io.Writer) {
	hd.outputer = w
}

const eol = "\r\n"

func (hd *HTTPDumper) dumpRequest(w io.Writer, req *http.Request) string {
	bs, _ := httputil.DumpRequest(req, true)

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

func (hd *HTTPDumper) dumpResponse(w io.Writer, id string, dw *httpDumpWriter) {
	bb := &bytes.Buffer{}

	bb.WriteString(fmt.Sprintf("<<<<<<<< %s %s <<<<<<<<", time.Now().Format(dumpTimeFormat), id))
	bb.WriteString(eol)

	dw.res.StatusCode = dw.ResponseWriter.Status()
	dw.res.Header = dw.ResponseWriter.Header()
	dw.res.Body = io.NopCloser(dw.bb)
	dw.res.Write(bb) //nolint: errcheck
	bb.WriteString(eol)
	bb.WriteString(eol)

	w.Write(bb.Bytes()) //nolint: errcheck
}

type httpDumpWriter struct {
	xin.ResponseWriter
	res *http.Response
	bb  *bytes.Buffer
}

func (dw *httpDumpWriter) Write(data []byte) (int, error) {
	dw.bb.Write(data)
	return dw.ResponseWriter.Write(data)
}
