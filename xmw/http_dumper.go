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

// HTTPDumpKey is the key for dump http saved in context
const HTTPDumpKey = "X_DUMP"

const dumpTimeFormat = "2006-01-02T15:04:05.000"

// HTTPDumper dump http request and response
type HTTPDumper struct {
	Outputer   io.Writer
	Maxlength  int64
	DumpCtxKey string // use xin.Context.Set(key, true | false) to enable/disable dump
	disabled   bool
}

// DefaultHTTPDumper create a middleware for xin http dumper
// Equals: NewHTTPDumper(xin.Logger.Outputer("XHD", log.LevelTrace))
func DefaultHTTPDumper(xin *xin.Engine) *HTTPDumper {
	return NewHTTPDumper(xin.Logger.GetOutputer("XHD", log.LevelTrace))
}

// NewHTTPDumper create a middleware for xin http dumper
func NewHTTPDumper(outputer io.Writer) *HTTPDumper {
	return &HTTPDumper{Outputer: outputer, Maxlength: 1 << 20, DumpCtxKey: HTTPDumpKey}
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
	w := hd.Outputer
	if w == nil {
		c.Next()
		return
	}

	d := hd.disabled
	if hd.DumpCtxKey != "" {
		if v, ok := c.Get(hd.DumpCtxKey); ok {
			d = !v.(bool)
		}
	}
	if d {
		c.Next()
		return
	}

	id := hd.dumpRequest(w, c.Request)

	cw := c.Writer
	dw := newHttpDumpWriter(c.Writer, hd.Maxlength)
	c.Writer = dw

	c.Next()

	hd.dumpResponse(w, id, c.Request, dw)
	c.Writer = cw
}

func (hd *HTTPDumper) dumpRequest(w io.Writer, req *http.Request) string {
	db := (req.ContentLength >= 0 && req.ContentLength <= hd.Maxlength)
	bs, _ := httputil.DumpRequest(req, db)

	id := fmt.Sprintf("%x", md5.Sum(bs)) //nolint: gosec

	bb := &bytes.Buffer{}

	bb.WriteString(fmt.Sprintf(">>>>>>>> %s %s >>>>>>>>\r\n", time.Now().Format(dumpTimeFormat), id))
	if len(bs) > 0 {
		bb.Write(bs)
	}
	bb.WriteString("\r\n\r\n")

	w.Write(bb.Bytes()) //nolint: errcheck

	return id
}

func (hd *HTTPDumper) dumpResponse(w io.Writer, id string, req *http.Request, hdw *httpDumpWriter) {
	bb := &bytes.Buffer{}

	bb.WriteString(fmt.Sprintf("<<<<<<<< %s %s <<<<<<<<\r\n", time.Now().Format(dumpTimeFormat), id))

	res := &http.Response{
		Proto:      req.Proto,
		ProtoMajor: req.ProtoMajor,
		ProtoMinor: req.ProtoMinor,
	}
	res.StatusCode = hdw.ResponseWriter.Status()
	res.Header = hdw.ResponseWriter.Header()
	res.Body = io.NopCloser(hdw.bb)
	res.Write(bb) //nolint: errcheck

	bb.WriteString("\r\n\r\n")

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
