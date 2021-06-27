package gindump

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pandafw/pango/log"
)

const defaultTimeFormat = "2006-01-02T15:04:05.000"

// Dumper dump http request and response
type Dumper struct {
	outputer io.Writer
	disabled bool
}

// Default create a default http dumper
// Equals to: New(log.Outputer("HTTP", log.LevelTrace))
func Default() *Dumper {
	return New(log.Outputer("HTTP", log.LevelTrace))
}

// New create a log middleware for gin http dumper
func New(outputer io.Writer) *Dumper {
	return &Dumper{outputer: outputer}
}

// Disable disable the dumper or not
func (d *Dumper) Disable(disabled bool) {
	d.disabled = disabled
}

// Handler returns the gin.HandlerFunc
func (d *Dumper) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		d.handle(c)
	}
}

// handle process gin request
func (d *Dumper) handle(c *gin.Context) {
	w := d.outputer
	if w == nil || d.disabled {
		c.Next()
		return
	}

	// dump request
	id := dumpRequest(w, c.Request)

	dw := &dumpWriter{c.Writer, &http.Response{
		Proto:      c.Request.Proto,
		ProtoMajor: c.Request.ProtoMajor,
		ProtoMinor: c.Request.ProtoMinor,
	}, &bytes.Buffer{}}
	c.Writer = dw

	// process request
	c.Next()

	// dump response
	dumpResponse(w, id, dw)
}

// SetOutput set the access log output writer
func (d *Dumper) SetOutput(w io.Writer) {
	d.outputer = w
}

const eol = "\r\n"

func dumpRequest(w io.Writer, req *http.Request) string {
	bs, _ := httputil.DumpRequest(req, true)

	id := fmt.Sprintf("%x", sha1.Sum(bs))

	bb := &bytes.Buffer{}

	// Seperate Line
	bb.WriteString(fmt.Sprintf(">>>>>>>> %s %s >>>>>>>>", time.Now().Format(defaultTimeFormat), id))
	bb.WriteString(eol)
	if len(bs) > 0 {
		bb.Write(bs)
	}
	bb.WriteString(eol)
	bb.WriteString(eol)

	// dump
	w.Write(bb.Bytes())

	return id
}

func dumpResponse(w io.Writer, id string, dw *dumpWriter) {
	bb := &bytes.Buffer{}

	// Seperate Line
	bb.WriteString(fmt.Sprintf("<<<<<<<< %s %s <<<<<<<<", time.Now().Format(defaultTimeFormat), id))
	bb.WriteString(eol)

	// http response
	dw.res.Header = dw.ResponseWriter.Header()
	dw.res.Body = ioutil.NopCloser(dw.bb)
	dw.res.Write(bb)
	bb.WriteString(eol)
	bb.WriteString(eol)

	// dump
	w.Write(bb.Bytes())
}

type dumpWriter struct {
	gin.ResponseWriter
	res *http.Response
	bb  *bytes.Buffer
}

func (dw *dumpWriter) WriteHeader(statusCode int) {
	dw.res.StatusCode = statusCode
	dw.ResponseWriter.WriteHeader(statusCode)
}

func (dw *dumpWriter) Write(data []byte) (int, error) {
	dw.bb.Write(data)
	return dw.ResponseWriter.Write(data)
}
