package xmw

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/xin"
)

// AccessLogTimeFormat default log time format
const AccessLogTimeFormat = "2006-01-02T15:04:05.000"

// AccessLogTextFormat default text log format
// TIME STATUS LATENCY SIZE CLIENT_IP REMOTE_ADDR METHOD HOST URL HEADER(User-Agent)
const AccessLogTextFormat = "text:%t\t%S\t%T\t%B\t%c\t%r\t%m\t%h\t%u\t%h{User-Agent}%n"

// AccessLogJSONFormat default json log format
const AccessLogJSONFormat = `json:{"when": %t, "server": %H, "status": %S, "latency": %T, "size": %B, "client_ip": %c, "remote_addr": %r, "method": %m, "host": %h, "url": %u, "user_agent": %h{User-Agent}}%n`

type AccessLogCtx struct {
	Start time.Time
	End   time.Time
	Ctx   *xin.Context
}

// AccessLogWriter access log writer for XIN
//
//	%t{format} - Request start time, if {format} is omitted, '2006-01-02T15:04:05.000' is used.
//	%c - Client IP ([X-Forwarded-For, X-Real-Ip] or RemoteIP())
//	%r - Remote IP:Port (%a)
//	%u - Request URL
//	%p - Request protocol
//	%m - Request method (GET, POST, etc.)
//	%q - Query string (prepended with a '?' if it exists)
//	%h - Request host
//	%h{name} - Request header
//	%A - Server listen address
//	%T - Time taken to process the request, in milliseconds
//	%S - HTTP status code of the response
//	%B - Response body length (%L)
//	%H - Local hostname
//	%H{name} - Response header
//	%n: EOL(Windows: "\r\n", Other: "\n")
type AccessLogWriter interface {
	Write(*AccessLogCtx)
}

// NewAccessLogMultiWriter create a multi writer
func NewAccessLogMultiWriter(ws ...AccessLogWriter) *AccessLogMultiWriter {
	return &AccessLogMultiWriter{Writers: ws}
}

// AccessLogMultiWriter write log to multiple writers.
type AccessLogMultiWriter struct {
	Writers []AccessLogWriter
}

// Write write the access log to multiple writers.
func (mw *AccessLogMultiWriter) Write(alc *AccessLogCtx) {
	for _, w := range mw.Writers {
		w.Write(alc)
	}
}

// NewAccessLogWriter create a text or json access log writer
// text:... -> AccessLogTextWriter
// json:... -> AccessLogJSONWriter
func NewAccessLogWriter(writer io.Writer, format string) AccessLogWriter {
	if strings.HasPrefix(format, "text:") {
		return NewAccessLogTextWriter(writer, format[5:])
	}
	if strings.HasPrefix(format, "json:") {
		return NewAccessLogJSONWriter(writer, format[5:])
	}
	return NewAccessLogTextWriter(writer, format)
}

type fmtfunc func(p *AccessLogCtx) string

// NewAccessLogTextWriter create text style writer for AccessLogger
func NewAccessLogTextWriter(writer io.Writer, format string) *AccessLogTextWriter {
	altw := &AccessLogTextWriter{writer: writer}
	altw.SetFormat(format)
	return altw
}

// AccessLogTextWriter format(text) and write access log
type AccessLogTextWriter struct {
	writer  io.Writer
	formats []fmtfunc
}

// SetOutput set the access log writer
func (altw *AccessLogTextWriter) SetOutput(w io.Writer) {
	altw.writer = w
}

// Write write the access log
func (altw *AccessLogTextWriter) Write(alc *AccessLogCtx) {
	writeAccessLog(altw.writer, alc, altw.formats)
}

// SetFormat set the access alw format
func (altw *AccessLogTextWriter) SetFormat(format string) {
	fmts := make([]fmtfunc, 0, 10)

	s := 0
	for i := 0; i < len(format); i++ {
		c := format[i]
		if c != '%' {
			continue
		}

		// string
		if s < i {
			fmts = append(fmts, strfmtc(format[s:i]))
		}

		i++
		s = i
		if i >= len(format) {
			break
		}

		// symbol
		var fmt fmtfunc
		switch format[i] {
		case 'c':
			fmt = clientIP
		case 'r', 'a':
			fmt = remoteAddr
		case 'u':
			fmt = requestURL
		case 'p':
			fmt = requestProto
		case 'm':
			fmt = requestMethod
		case 'q':
			fmt = requestQuery
		case 'h':
			p := getFormatOption(format, &i)
			if p != "" {
				fmt = requestHeader(p)
			} else {
				fmt = requestHost
			}
		case 't':
			p := getFormatOption(format, &i)
			if p == "" {
				p = AccessLogTimeFormat
			}
			fmt = timefmtc(p)
		case 'A':
			fmt = listenAddr
		case 'S':
			fmt = statusCode
		case 'T':
			fmt = latency
		case 'B', 'L':
			fmt = responseBodyLen
		case 'H':
			p := getFormatOption(format, &i)
			if p != "" {
				fmt = responseHeader(p)
			} else {
				s, _ := os.Hostname()
				fmt = strfmtc(s)
			}
		case 'n':
			fmt = eolfmt
		}

		if fmt != nil {
			fmts = append(fmts, fmt)
			s = i + 1
		}
	}

	if s < len(format) {
		fmts = append(fmts, strfmtc(format[s:]))
	}

	altw.formats = fmts
}

// NewAccessLogJSONWriter create json style writer for AccessLogger
func NewAccessLogJSONWriter(writer io.Writer, format string) *AccessLogJSONWriter {
	aljw := &AccessLogJSONWriter{writer: writer}
	aljw.SetFormat(format)
	return aljw
}

// AccessLogJSONWriter format(json-style) and write access log
type AccessLogJSONWriter struct {
	writer  io.Writer
	formats []fmtfunc
}

// SetOutput set the access log writer
func (aljw *AccessLogJSONWriter) SetOutput(w io.Writer) {
	aljw.writer = w
}

// Write write the access log
func (aljw *AccessLogJSONWriter) Write(alc *AccessLogCtx) {
	writeAccessLog(aljw.writer, alc, aljw.formats)
}

// SetFormat set the access alw format
func (aljw *AccessLogJSONWriter) SetFormat(format string) {
	fmts := make([]fmtfunc, 0, 10)

	s := 0
	for i := 0; i < len(format); i++ {
		c := format[i]
		if c != '%' {
			continue
		}

		// string
		if s < i {
			fmts = append(fmts, strfmtc(format[s:i]))
		}

		i++
		s = i
		if i >= len(format) {
			break
		}

		// symbol
		var fmt fmtfunc
		switch format[i] {
		case 'c':
			fmt = quotefmtc(clientIP)
		case 'r', 'a':
			fmt = quotefmtc(remoteAddr)
		case 'u':
			fmt = quotefmtc(requestURL)
		case 'p':
			fmt = quotefmtc(requestProto)
		case 'm':
			fmt = quotefmtc(requestMethod)
		case 'q':
			fmt = quotefmtc(requestQuery)
		case 'h':
			p := getFormatOption(format, &i)
			if p != "" {
				fmt = requestHeader(p)
			} else {
				fmt = requestHost
			}
			fmt = quotefmtc(fmt)
		case 't':
			p := getFormatOption(format, &i)
			if p == "" {
				p = AccessLogTimeFormat
			}
			fmt = quotefmtc(timefmtc(p))
		case 'A':
			fmt = quotefmtc(listenAddr)
		case 'S':
			fmt = statusCode
		case 'T':
			fmt = latency
		case 'B', 'L':
			fmt = responseBodyLen
		case 'H':
			p := getFormatOption(format, &i)
			if p != "" {
				fmt = quotefmtc(responseHeader(p))
			} else {
				s, _ := os.Hostname()
				fmt = quotefmtc(strfmtc(s))
			}
		case 'n':
			fmt = eolfmt
		}

		if fmt != nil {
			fmts = append(fmts, fmt)
			s = i + 1
		}
	}

	if s < len(format) {
		fmts = append(fmts, strfmtc(format[s:]))
	}

	aljw.formats = fmts
}

// -------------------------------------------------

func writeAccessLog(w io.Writer, alc *AccessLogCtx, fmts []fmtfunc) {
	bb := &bytes.Buffer{}
	for _, f := range fmts {
		s := f(alc)
		bb.WriteString(s)
	}
	w.Write(bb.Bytes()) //nolint: errcheck
}

func getFormatOption(format string, i *int) string {
	p := format[*i+1:]
	if len(p) > 0 && p[0] == '{' {
		e := strings.IndexByte(p, '}')
		if e > 0 {
			*i += e + 1
			return p[1:e]
		}
	}
	return ""
}

// -------------------------------------------------

func quotefmtc(ff fmtfunc) fmtfunc {
	return func(p *AccessLogCtx) string {
		return fmt.Sprintf("%q", ff(p))
	}
}

func strfmtc(s string) fmtfunc {
	return func(p *AccessLogCtx) string {
		return s
	}
}

func timefmtc(layout string) fmtfunc {
	return func(p *AccessLogCtx) string {
		return p.Start.Format(layout)
	}
}

func eolfmt(p *AccessLogCtx) string {
	return iox.EOL
}

func latency(p *AccessLogCtx) string {
	return strconv.FormatInt(p.End.Sub(p.Start).Milliseconds(), 10)
}

func clientIP(p *AccessLogCtx) string {
	return p.Ctx.ClientIP()
}

func remoteAddr(p *AccessLogCtx) string {
	return p.Ctx.Request.RemoteAddr
}

func listenAddr(p *AccessLogCtx) string {
	ctx := p.Ctx.Request.Context()
	addr, ok := ctx.Value(http.LocalAddrContextKey).(net.Addr)
	if ok {
		return addr.String()
	}
	return ""
}

func requestURL(p *AccessLogCtx) string {
	return p.Ctx.Request.URL.String()
}

func requestHost(p *AccessLogCtx) string {
	return p.Ctx.Request.Host
}

func requestProto(p *AccessLogCtx) string {
	return p.Ctx.Request.Proto
}

func requestMethod(p *AccessLogCtx) string {
	return p.Ctx.Request.Method
}

func requestQuery(p *AccessLogCtx) string {
	return p.Ctx.Request.URL.RawQuery
}

func requestHeader(name string) fmtfunc {
	return func(p *AccessLogCtx) string {
		return p.Ctx.Request.Header.Get(name)
	}
}

func statusCode(p *AccessLogCtx) string {
	return strconv.Itoa(p.Ctx.Writer.Status())
}

func responseBodyLen(p *AccessLogCtx) string {
	return strconv.Itoa(p.Ctx.Writer.Size())
}

func responseHeader(name string) fmtfunc {
	return func(p *AccessLogCtx) string {
		return p.Ctx.Writer.Header().Get(name)
	}
}
