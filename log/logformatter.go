package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
)

// EOL windows: "\r\n" other: "\n"
var EOL = iox.EOL

const defaultTimeFormat = "2006-01-02T15:04:05.000"

// Formatter log formatter interface
type Formatter interface {
	Write(w io.Writer, le *Event)
}

// TextFmtSubject subject log format "[%l] %m"
var TextFmtSubject = newTextFormatter("[%l] %m")

// TextFmtSimple simple log format "[%p] %m%n"
var TextFmtSimple = newTextFormatter("[%p] %m%n")

// TextFmtDefault default log format "%t %l{-5s} %c %S:%L %F() - %m%n%T"
var TextFmtDefault = newTextFormatter("%t %l{-5s} %c %S:%L %F() - %m%n%T")

// JSONFmtDefault default log format `{"time": %t, "level": %l, "name": %c, "host": %h, "file": %S, "line": %L, "func": %F, "msg": %m, "trace": %T}%n`
var JSONFmtDefault = newJSONFormatter(`{"time": %t, "level": %l, "name": %c, "host": %h, "file": %S, "line": %L, "func": %F, "msg": %m, "trace": %T}%n`)

// NewLogFormatter create a text or json formatter
// text:[%p] %m%n -> TextFormatter
// json:{"level":%l, "msg": %m}%n  -> JSONFormatter
func NewLogFormatter(format string) Formatter {
	if strings.HasPrefix(format, "text:") {
		return NewTextFormatter(format[5:])
	}
	if strings.HasPrefix(format, "json:") {
		return NewJSONFormatter(format[5:])
	}
	return NewTextFormatter(format)
}

// NewTextFormatter create a Text Formatter instance
// Text Format
// %t{format}: time, if {format} is omitted, '2006-01-02T15:04:05.000' will be used
// %c{format}: logger name
// %p{format}: log level prefix
// %l{format}: log level string
// %h{format}: hostname
// %e{key}: os environment variable
// %x{key}: logger property
// %X{=| }: logger properties (operator|separator)
// %S: caller source file name
// %L: caller source line number
// %F: caller function name
// %T: caller stack trace
// %m: message
// %q: quoted message
// %n: EOL(Windows: "\r\n", Other: "\n")
func NewTextFormatter(format string) *TextFormatter {
	switch format {
	case "DEFAULT":
		return TextFmtDefault
	case "SIMPLE":
		return TextFmtSimple
	case "SUBJECT":
		return TextFmtSubject
	default:
		return newTextFormatter(format)
	}
}

func newTextFormatter(format string) *TextFormatter {
	tf := &TextFormatter{}
	tf.SetFormat(format)
	return tf
}

// NewJSONFormatter create a Json Formatter instance
// JSON Format
// %t{format}: time, if {format} is omitted, '2006-01-02T15:04:05.000' will be used
// %c{format}: logger name
// %p{format}: log level prefix
// %l{format}: log level string
// %h{format}: hostname
// %e{key}: os environment variable
// %x{key}: logger property
// %X: logger properties (json format)
// %S: caller source file name
// %L: caller source line number
// %F: caller function name
// %T: caller stack trace
// %m: message
// %n: EOL(Windows: "\r\n", Other: "\n")
func NewJSONFormatter(format string) *JSONFormatter {
	switch format {
	case "DEFAULT":
		return JSONFmtDefault
	default:
		return newJSONFormatter(format)
	}
}

func newJSONFormatter(format string) *JSONFormatter {
	jf := &JSONFormatter{}
	jf.SetFormat(format)
	return jf
}

// TextFormatter text formatter
type TextFormatter struct {
	fmts []fmtfunc
}

// Format format the log event to the writer w
func (tf *TextFormatter) Write(w io.Writer, le *Event) {
	write(w, le, tf.fmts)
}

// SetFormat initialize the text formatter
func (tf *TextFormatter) SetFormat(format string) {
	fmts := make([]fmtfunc, 0, 10)

	s := 0
	for i := 0; i < len(format); i++ {
		c := format[i]
		if c != '%' {
			continue
		}

		// string
		if s < i {
			fmts = append(fmts, fcString(format[s:i]))
		}

		i++
		s = i
		if i >= len(format) {
			break
		}

		// symbol
		var ff fmtfunc
		switch format[i] {
		case 't':
			p := getFormatOption(format, &i)
			if p == "" {
				p = defaultTimeFormat
			}
			switch p {
			case "unix":
				ff = ffTimeUnix
			case "unixmilli":
				ff = ffTimeUnixMilli
			case "unixmicro":
				ff = ffTimeUnixMicro
			case "unixnano":
				ff = ffTimeUnixNano
			default:
				ff = fcTimeFormat(p)
			}
		case 'c':
			p := getFormatOption(format, &i)
			if p == "" {
				ff = ffName
			} else {
				ff = fcNameFormat("%" + p)
			}
		case 'p':
			p := getFormatOption(format, &i)
			if p == "" {
				ff = ffLevelPrefix
			} else {
				ff = fcLevelPrefix("%" + p)
			}
		case 'l':
			p := getFormatOption(format, &i)
			if p == "" {
				ff = ffLevelString
			} else {
				ff = fcLevelString("%" + p)
			}
		case 'h':
			p := getFormatOption(format, &i)
			h, _ := os.Hostname()
			if p == "" {
				ff = fcString(h)
			} else {
				ff = fcStringFormat(p, h)
			}
		case 'e':
			p := getFormatOption(format, &i)
			if p != "" {
				ff = fcGetenv(p)
			}
		case 'x':
			p := getFormatOption(format, &i)
			if p != "" {
				ff = fcPropText(p)
			}
		case 'X':
			p := getFormatOption(format, &i)
			if p == "" {
				p = "=| "
			}
			ff = fcPropsText(p)
		case 'S':
			ff = ffFile
		case 'L':
			ff = ffLine
		case 'F':
			ff = ffFunc
		case 'T':
			ff = ffTrace
		case 'm':
			ff = ffMsg
		case 'q':
			ff = ffQmsg
		case 'n':
			ff = ffEol
		}

		if ff != nil {
			fmts = append(fmts, ff)
			s = i + 1
		}
	}

	if s < len(format) {
		fmts = append(fmts, fcString(format[s:]))
	}

	tf.fmts = fmts
}

// JSONFormatter json formatter
type JSONFormatter struct {
	fmts []fmtfunc
}

// Write format the log event as a json string to the writer w
func (jf *JSONFormatter) Write(w io.Writer, le *Event) {
	write(w, le, jf.fmts)
}

// SetFormat initialize the json formatter
func (jf *JSONFormatter) SetFormat(format string) {
	fmts := make([]fmtfunc, 0, 10)

	s := 0
	for i := 0; i < len(format); i++ {
		c := format[i]
		if c != '%' {
			continue
		}

		// string
		if s < i {
			fmts = append(fmts, fcString(format[s:i]))
		}

		i++
		s = i
		if i >= len(format) {
			break
		}

		// symbol
		var ff fmtfunc
		switch format[i] {
		case 't':
			p := getFormatOption(format, &i)
			if p == "" {
				p = defaultTimeFormat
			}
			switch p {
			case "unix":
				ff = ffTimeUnix
			case "unixmilli":
				ff = ffTimeUnixMilli
			case "unixmicro":
				ff = ffTimeUnixMicro
			case "unixnano":
				ff = ffTimeUnixNano
			default:
				ff = fcQuote(fcTimeFormat(p))
			}
		case 'c':
			p := getFormatOption(format, &i)
			if p == "" {
				ff = ffName
			} else {
				ff = fcNameFormat("%" + p)
			}
			ff = fcQuote(ff)
		case 'p':
			p := getFormatOption(format, &i)
			if p == "" {
				ff = ffLevelPrefix
			} else {
				ff = fcLevelPrefix("%" + p)
			}
			ff = fcQuote(ff)
		case 'l':
			p := getFormatOption(format, &i)
			if p == "" {
				ff = ffLevelString
			} else {
				ff = fcLevelString("%" + p)
			}
			ff = fcQuote(ff)
		case 'h':
			p := getFormatOption(format, &i)
			h, _ := os.Hostname()
			if p == "" {
				ff = fcString(h)
			} else {
				ff = fcStringFormat(p, h)
			}
			ff = fcQuote(ff)
		case 'e':
			p := getFormatOption(format, &i)
			if p != "" {
				ff = fcQuote(fcGetenv(p))
			}
		case 'x':
			p := getFormatOption(format, &i)
			if p != "" {
				ff = fcPropJSON(p)
			}
		case 'X':
			ff = fcPropsJSON
		case 'S':
			ff = fcQuote(ffFile)
		case 'L':
			ff = ffLine
		case 'F':
			ff = fcQuote(ffFunc)
		case 'T':
			ff = fcQuote(ffTrace)
		case 'm':
			ff = fcQuote(ffMsg)
		case 'n':
			ff = ffEol
		}

		if ff != nil {
			fmts = append(fmts, ff)
			s = i + 1
		}
	}

	if s < len(format) {
		fmts = append(fmts, fcString(format[s:]))
	}
	jf.fmts = fmts
}

//-------------------------------------------------

type fmtfunc func(le *Event) string

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

func write(w io.Writer, le *Event, fmts []fmtfunc) {
	for _, f := range fmts {
		s := f(le)
		iox.WriteString(w, s) //nolint: errcheck
	}
}

func ffName(le *Event) string {
	return le.Name
}

func ffLevelPrefix(le *Event) string {
	return le.Level.Prefix()
}

func ffLevelString(le *Event) string {
	return le.Level.String()
}

func ffFunc(le *Event) string {
	return le.Func
}

func ffFile(le *Event) string {
	return le.File
}

func ffLine(le *Event) string {
	return strconv.Itoa(le.Line)
}

func ffTrace(le *Event) string {
	return le.Trace
}

func ffMsg(le *Event) string {
	return le.Msg
}

func ffQmsg(le *Event) string {
	return strconv.Quote(le.Msg)
}

func ffEol(le *Event) string {
	return EOL
}

func fcQuote(ff fmtfunc) fmtfunc {
	return func(le *Event) string {
		return strconv.Quote(ff(le))
	}
}

func fcString(s string) fmtfunc {
	return func(le *Event) string {
		return s
	}
}

func fcStringFormat(f, s string) fmtfunc {
	return func(le *Event) string {
		return fmt.Sprintf(f, s)
	}
}

func ffTimeUnix(le *Event) string {
	return num.Ltoa(le.Time.Unix())
}

func ffTimeUnixMilli(le *Event) string {
	return num.Ltoa(le.Time.UnixMilli())
}

func ffTimeUnixMicro(le *Event) string {
	return num.Ltoa(le.Time.UnixMicro())
}

func ffTimeUnixNano(le *Event) string {
	return num.Ltoa(le.Time.UnixNano())
}

func fcTimeFormat(layout string) fmtfunc {
	return func(le *Event) string {
		return le.Time.Format(layout)
	}
}

func fcNameFormat(f string) fmtfunc {
	return func(le *Event) string {
		return fmt.Sprintf(f, le.Name)
	}
}

func fcLevelPrefix(f string) fmtfunc {
	return func(le *Event) string {
		return fmt.Sprintf(f, le.Level.Prefix())
	}
}

func fcLevelString(f string) fmtfunc {
	return func(le *Event) string {
		return fmt.Sprintf(f, le.Level.String())
	}
}

func fcGetenv(key string) fmtfunc {
	return func(le *Event) string {
		return os.Getenv(key)
	}
}

func fcPropText(key string) fmtfunc {
	return func(le *Event) string {
		if v, ok := le.Props[key]; ok {
			return fmt.Sprint(v)
		}
		return ""
	}
}

func fcPropsText(f string) fmtfunc {
	ss := strings.Split(f, "|")
	d := ss[0]
	j := ""
	if len(ss) > 1 {
		j = ss[1]
	}
	return func(le *Event) string {
		if len(le.Props) == 0 {
			return ""
		}

		a := make([]string, 0, len(le.Props))
		for k, v := range le.Props {
			if v == nil {
				v = ""
			}
			a = append(a, fmt.Sprintf("%s%s%v", k, d, v))
		}
		return strings.Join(a, j)
	}
}

func fcPropJSON(key string) fmtfunc {
	return func(le *Event) string {
		if v, ok := le.Props[key]; ok {
			b, _ := json.Marshal(v)
			return str.UnsafeString(b)
		}
		return `null`
	}
}

func fcPropsJSON(le *Event) string {
	b, _ := json.Marshal(le.Props)
	return str.UnsafeString(b)
}

//-------------------------------------------------

type FormatSupport struct {
	Formatter Formatter    // log formatter
	Buffer    bytes.Buffer // log buffer
}

// SetFormat set the log formatter
func (fs *FormatSupport) SetFormat(format string) {
	fs.Formatter = NewLogFormatter(format)
}

// Format format the log event
func (fs *FormatSupport) GetFormatter(le *Event, df ...Formatter) Formatter {
	f := fs.Formatter
	if f == nil {
		if len(df) > 0 {
			f = df[0]
		} else {
			f = TextFmtDefault
		}
	}
	return f
}

// Format format the log event
func (fs *FormatSupport) Format(le *Event, df ...Formatter) []byte {
	f := fs.GetFormatter(le, df...)
	fs.Buffer.Reset()
	f.Write(&fs.Buffer, le)
	return fs.Buffer.Bytes()
}

// Append format the log event and append to buffer
func (fs *FormatSupport) Append(le *Event, df ...Formatter) {
	f := fs.GetFormatter(le, df...)
	f.Write(&fs.Buffer, le)
}

type SubjectSuport struct {
	Subjecter Formatter    // log formatter
	SubBuffer bytes.Buffer // log buffer
}

// SetSubject set the subject formatter
func (ss *SubjectSuport) SetSubject(format string) {
	ss.Subjecter = NewLogFormatter(format)
}

// GetFormatter get Formatter
func (ss *SubjectSuport) SubFormat(le *Event) []byte {
	f := ss.Subjecter
	if f == nil {
		f = TextFmtSubject
	}

	ss.SubBuffer.Reset()
	f.Write(&ss.SubBuffer, le)
	return ss.SubBuffer.Bytes()
}
