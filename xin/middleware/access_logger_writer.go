package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tmu"
	"github.com/askasoft/pango/xin"
)

// AccessLogTimeFormat default log time format
const AccessLogTimeFormat = "2006-01-02T15:04:05.000Z07:00"

// AccessLogTextFormat default text log format
// TIME STATUS LATENCY SIZE CLIENT_IP REMOTE_ADDR METHOD HOST URL HEADER(User-Agent)
const AccessLogTextFormat = "text:%t\t%S\t%D\t%B\t%c\t%r\t%m\t%s://%h%u\t%h{User-Agent}%n"

// AccessLogJSONFormat default json log format
const AccessLogJSONFormat = `json:{"when": %t, "server": %H, "status": %S, "latency": %T, "size": %B, "client_ip": %c, "remote_addr": %r, "method": %m, "scheme": %s, "host": %h, "url": %u, "user_agent": %h{User-Agent}}%n`

// AccessLogWriter access log writer for XIN
//
//	%t{format} - Request start time, if {format} is omitted, '2006-01-02T15:04:05.000Z07:00' is used.
//	%c - Client IP ([X-Forwarded-For, X-Real-Ip] or RemoteIP())
//	%r - Remote IP:Port (%a)
//	%u - Request URL
//	%p - Request protocol
//	%s - Request scheme (http, https)
//	%m - Request method (GET, POST, etc.)
//	%q - Query string (prepended with a '?' if it exists)
//	%h - Request host
//	%h{name} - Request header
//	%A - Server listen address
//	%D - Time taken to process the request, duration format string
//	%T - Time taken to process the request, number in milliseconds
//	%S - Response status code
//	%B - Response body length (%L)
//	%H - Local hostname
//	%H{name} - Response header
//	%n: EOL(Windows: "\r\n", Other: "\n")
type AccessLogWriter interface {
	Write(*xin.Context)
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
func (mw *AccessLogMultiWriter) Write(c *xin.Context) {
	for _, w := range mw.Writers {
		w.Write(c)
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

type fmtfunc func(c *xin.Context) string

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
func (altw *AccessLogTextWriter) Write(c *xin.Context) {
	writeAccessLog(altw.writer, c, altw.formats)
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
		case 's':
			fmt = requestScheme
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
		case 'D':
			fmt = latencyDuration
		case 'T':
			fmt = latencyMillis
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
func (aljw *AccessLogJSONWriter) Write(c *xin.Context) {
	writeAccessLog(aljw.writer, c, aljw.formats)
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
		case 's':
			fmt = quotefmtc(requestScheme)
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
		case 'D':
			fmt = quotefmtc(latencyDuration)
		case 'T':
			fmt = latencyMillis
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

func writeAccessLog(w io.Writer, c *xin.Context, fmts []fmtfunc) {
	bb := &bytes.Buffer{}
	for _, f := range fmts {
		s := f(c)
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
	return func(c *xin.Context) string {
		return fmt.Sprintf("%q", ff(c))
	}
}

func strfmtc(s string) fmtfunc {
	return func(c *xin.Context) string {
		return s
	}
}

func timefmtc(layout string) fmtfunc {
	return func(c *xin.Context) string {
		return c.GetTime(AccessLogStartKey).Format(layout)
	}
}

func eolfmt(c *xin.Context) string {
	return iox.EOL
}

func latencyDuration(c *xin.Context) string {
	s := c.GetTime(AccessLogStartKey)
	e := c.GetTime(AccessLogEndKey)
	return tmu.HumanDuration(e.Sub(s))
}

func latencyMillis(c *xin.Context) string {
	s := c.GetTime(AccessLogStartKey)
	e := c.GetTime(AccessLogEndKey)
	return strconv.FormatInt(e.Sub(s).Milliseconds(), 10)
}

func clientIP(c *xin.Context) string {
	return c.ClientIP()
}

func remoteAddr(c *xin.Context) string {
	return c.Request.RemoteAddr
}

func listenAddr(c *xin.Context) string {
	ctx := c.Request.Context()
	addr, ok := ctx.Value(http.LocalAddrContextKey).(net.Addr)
	if ok {
		return addr.String()
	}
	return ""
}

func requestURL(c *xin.Context) string {
	return c.Request.URL.String()
}

func requestHost(c *xin.Context) string {
	return c.Request.Host
}

func requestProto(c *xin.Context) string {
	return c.Request.Proto
}

func requestScheme(c *xin.Context) string {
	return str.If(c.IsSecure(), "https", "http")
}

func requestMethod(c *xin.Context) string {
	return c.Request.Method
}

func requestQuery(c *xin.Context) string {
	return c.Request.URL.RawQuery
}

func requestHeader(name string) fmtfunc {
	return func(c *xin.Context) string {
		return c.Request.Header.Get(name)
	}
}

func statusCode(c *xin.Context) string {
	return strconv.Itoa(c.Writer.Status())
}

func responseBodyLen(c *xin.Context) string {
	return strconv.Itoa(c.Writer.Size())
}

func responseHeader(name string) fmtfunc {
	return func(c *xin.Context) string {
		return c.Writer.Header().Get(name)
	}
}
