package ginlog

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pandafw/pango/iox"
	"github.com/pandafw/pango/log"
	"github.com/pandafw/pango/net/httpx"
)

// DefaultTimeFormat default log time format
const DefaultTimeFormat = "2006-01-02T15:04:05.000"

// DefaultTextLogFormat default log format
// TIME STATUS LATENCY LENGTH CLIENT_IP REMOTE_ADDR LISTEN METHOD HOST URL
const DefaultTextLogFormat = "text:%t\t%S\t%T\t%L\t%c\t%r\t%A\t%m\t%h\t%u%n"

// DefaultJSONLogFormat default log format
const DefaultJSONLogFormat = `json:{"when": %t, "status": %S, "latency": %T, "length": %L, "clientIP": %c, "remoteAddr": %r, "listen": %A, "method": %m, "host": %h, "url": %u}%n`

// Logger access loger for GIN
type Logger struct {
	outputer io.Writer
	formats  []fmtfunc
	disabled bool
}

type param struct {
	Start time.Time
	End   time.Time
	Ctx   *gin.Context
}

type fmtfunc func(p *param) string

// Default create a default log
// Equals to: New(log.Outputer("GIN", log.LevelTrace), DefaultTextLogFormt)
func Default() *Logger {
	return New(log.Outputer("GIN", log.LevelTrace), DefaultTextLogFormat)
}

// New create a log middleware for gin access log
// Access Log Format:
// text:...     json:...
//   %t{format} - Request start time, if {format} is omitted, '2006-01-02T15:04:05.000' is used.
//   %c - Client IP ([X-Forwarded-For, X-Real-Ip] or RemoteIP())
//   %r - Remote IP:Port
//   %u - Request URL
//   %p - Request protocol
//   %m - Request method (GET, POST, etc.)
//   %q - Query string (prepended with a '?' if it exists)
//   %h - Request host
//   %h{name} - Request header
//   %A - Server listen address
//   %T - Time taken to process the request, in milliseconds
//   %S - HTTP status code of the response
//   %L - Response body length
//   %H{name} - Response header
//   %n: EOL(Windows: "\r\n", Other: "\n")
func New(outputer io.Writer, format string) *Logger {
	return &Logger{outputer: outputer, formats: parseFormat(format)}
}

// Disable disable the logger or not
func (log *Logger) Disable(disabled bool) {
	log.disabled = disabled
}

// Handler returns the gin.HandlerFunc
func (log *Logger) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.handle(c)
	}
}

// handle process gin request
func (log *Logger) handle(c *gin.Context) {
	w := log.outputer
	if w == nil || log.disabled {
		c.Next()
		return
	}

	p := &param{Start: time.Now(), Ctx: c}

	// process request
	c.Next()

	p.End = time.Now()

	// log.formats can be modified concurrently
	fmts := log.formats

	// write access log
	bb := &bytes.Buffer{}
	for _, f := range fmts {
		s := f(p)
		bb.WriteString(s)
	}
	w.Write(bb.Bytes())
}

// SetOutput set the access log output writer
func (log *Logger) SetOutput(w io.Writer) {
	log.outputer = w
}

// SetFormat set the access log format
func (log *Logger) SetFormat(format string) {
	log.formats = parseFormat(format)
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
			}
			fmt = requestHost
		case 't':
			p := getFormatOption(format, &i)
			if p == "" {
				p = DefaultTimeFormat
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
				p = DefaultTimeFormat
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

//-------------------------------------------------
func quotefmtc(ff fmtfunc) fmtfunc {
	return func(p *param) string {
		return fmt.Sprintf("%q", ff(p))
	}
}

func strfmtc(s string) fmtfunc {
	return func(p *param) string {
		return s
	}
}

func timefmtc(layout string) fmtfunc {
	return func(p *param) string {
		return p.Start.Format(layout)
	}
}

func eolfmt(p *param) string {
	return iox.EOL
}

func latency(p *param) string {
	return strconv.FormatInt(p.End.Sub(p.Start).Milliseconds(), 10)
}

func clientIP(p *param) string {
	return httpx.GetClientIP(p.Ctx.Request)
}

func remoteAddr(p *param) string {
	return p.Ctx.Request.RemoteAddr
}

func listenAddr(p *param) string {
	ctx := p.Ctx.Request.Context()
	addr, ok := ctx.Value(http.LocalAddrContextKey).(net.Addr)
	if ok {
		return addr.String()
	}
	return ""
}

func requestURL(p *param) string {
	return p.Ctx.Request.URL.String()
}

func requestHost(p *param) string {
	return p.Ctx.Request.Host
}

func requestProto(p *param) string {
	return p.Ctx.Request.Proto
}

func requestMethod(p *param) string {
	return p.Ctx.Request.Method
}

func requestQuery(p *param) string {
	return p.Ctx.Request.URL.RawQuery
}

func requestPath(p *param) string {
	return p.Ctx.Request.URL.Path
}

func requestHeader(name string) fmtfunc {
	return func(p *param) string {
		return p.Ctx.Request.Header.Get(name)
	}
}

func statusCode(p *param) string {
	return strconv.Itoa(p.Ctx.Writer.Status())
}

func responseBodyLen(p *param) string {
	return strconv.Itoa(p.Ctx.Writer.Size())
}

func responseHeader(name string) fmtfunc {
	return func(p *param) string {
		return p.Ctx.Writer.Header().Get(name)
	}
}
