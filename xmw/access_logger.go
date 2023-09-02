package xmw

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/xin"
)

// AccessLogTimeFormat default log time format
const AccessLogTimeFormat = "2006-01-02T15:04:05.000"

// AccessLogTextFormat default log format
// TIME STATUS LATENCY LENGTH CLIENT_IP REMOTE_ADDR LISTEN METHOD HOST URL
const AccessLogTextFormat = "text:%t\t%S\t%T\t%L\t%c\t%r\t%A\t%m\t%h\t%u%n"

// AccessLogJSONFormat default log format
const AccessLogJSONFormat = `json:{"when": %t, "status": %S, "latency": %T, "length": %L, "clientIP": %c, "remoteAddr": %r, "listen": %A, "method": %m, "host": %h, "url": %u}%n`

// AccessLogger access loger for XIN
type AccessLogger struct {
	outputer io.Writer
	formats  []fmtfunc
	disabled bool
}

type logevt struct {
	Start time.Time
	End   time.Time
	Ctx   *xin.Context
}

type fmtfunc func(p *logevt) string

// DefaultAccessLogger create a log middleware for xin access logger
// Equals: NewAccessLogger(xin.Logger.Outputer("XINA", log.LevelTrace), xmw.AccessLogTextFormat)
func DefaultAccessLogger(xin *xin.Engine) *AccessLogger {
	return NewAccessLogger(xin.Logger.GetOutputer("XINA", log.LevelTrace), AccessLogTextFormat)
}

// NewAccessLogger create a log middleware for xin access logger
// Access Log Format:
// text:...     json:...
//
//	%t{format} - Request start time, if {format} is omitted, '2006-01-02T15:04:05.000' is used.
//	%c - Client IP ([X-Forwarded-For, X-Real-Ip] or RemoteIP())
//	%r - Remote IP:Port
//	%u - Request URL
//	%p - Request protocol
//	%m - Request method (GET, POST, etc.)
//	%q - Query string (prepended with a '?' if it exists)
//	%h - Request host
//	%h{name} - Request header
//	%A - Server listen address
//	%T - Time taken to process the request, in milliseconds
//	%S - HTTP status code of the response
//	%L - Response body length
//	%H{name} - Response header
//	%n: EOL(Windows: "\r\n", Other: "\n")
func NewAccessLogger(outputer io.Writer, format string) *AccessLogger {
	return &AccessLogger{outputer: outputer, formats: parseFormat(format)}
}

// Disable disable the logger or not
func (al *AccessLogger) Disable(disabled bool) {
	al.disabled = disabled
}

// Handler returns the HandlerFunc
func (al *AccessLogger) Handler() xin.HandlerFunc {
	return al.Handle
}

// Handle process xin request
func (al *AccessLogger) Handle(c *xin.Context) {
	w := al.outputer
	if w == nil || al.disabled {
		c.Next()
		return
	}

	p := &logevt{Start: time.Now(), Ctx: c}

	// process request
	c.Next()

	p.End = time.Now()

	// al.formats can be modified concurrently
	fmts := al.formats

	// write access al
	bb := &bytes.Buffer{}
	for _, f := range fmts {
		s := f(p)
		bb.WriteString(s)
	}
	w.Write(bb.Bytes()) //nolint: errcheck
}

// SetOutput set the access al output writer
func (al *AccessLogger) SetOutput(w io.Writer) {
	al.outputer = w
}

// SetFormat set the access al format
func (al *AccessLogger) SetFormat(format string) {
	al.formats = parseFormat(format)
}

func parseFormat(format string) []fmtfunc {
	if strings.HasPrefix(format, "text:") {
		return parseTextFormat(format[5:])
	}
	if strings.HasPrefix(format, "json:") {
		return parseJSONFormat(format[5:])
	}
	return parseTextFormat(format)
}

func parseTextFormat(format string) []fmtfunc {
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
		case 'r':
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
		case 'L':
			fmt = responseBodyLen
		case 'H':
			p := getFormatOption(format, &i)
			if p != "" {
				fmt = responseHeader(p)
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

	return fmts
}

func parseJSONFormat(format string) []fmtfunc {
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
		case 'r':
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
		case 'L':
			fmt = responseBodyLen
		case 'H':
			p := getFormatOption(format, &i)
			if p != "" {
				fmt = quotefmtc(responseHeader(p))
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

	return fmts
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
	return func(p *logevt) string {
		return fmt.Sprintf("%q", ff(p))
	}
}

func strfmtc(s string) fmtfunc {
	return func(p *logevt) string {
		return s
	}
}

func timefmtc(layout string) fmtfunc {
	return func(p *logevt) string {
		return p.Start.Format(layout)
	}
}

func eolfmt(p *logevt) string {
	return iox.EOL
}

func latency(p *logevt) string {
	return strconv.FormatInt(p.End.Sub(p.Start).Milliseconds(), 10)
}

func clientIP(p *logevt) string {
	return p.Ctx.ClientIP()
}

func remoteAddr(p *logevt) string {
	return p.Ctx.Request.RemoteAddr
}

func listenAddr(p *logevt) string {
	ctx := p.Ctx.Request.Context()
	addr, ok := ctx.Value(http.LocalAddrContextKey).(net.Addr)
	if ok {
		return addr.String()
	}
	return ""
}

func requestURL(p *logevt) string {
	return p.Ctx.Request.URL.String()
}

func requestHost(p *logevt) string {
	return p.Ctx.Request.Host
}

func requestProto(p *logevt) string {
	return p.Ctx.Request.Proto
}

func requestMethod(p *logevt) string {
	return p.Ctx.Request.Method
}

func requestQuery(p *logevt) string {
	return p.Ctx.Request.URL.RawQuery
}

func requestHeader(name string) fmtfunc {
	return func(p *logevt) string {
		return p.Ctx.Request.Header.Get(name)
	}
}

func statusCode(p *logevt) string {
	return strconv.Itoa(p.Ctx.Writer.Status())
}

func responseBodyLen(p *logevt) string {
	return strconv.Itoa(p.Ctx.Writer.Size())
}

func responseHeader(name string) fmtfunc {
	return func(p *logevt) string {
		return p.Ctx.Writer.Header().Get(name)
	}
}
