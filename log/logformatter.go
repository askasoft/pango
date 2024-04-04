package log

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/str"
)

// EOL windows: "\r\n" other: "\n"
var EOL = iox.EOL

const defaultTimeFormat = "2006-01-02T15:04:05.000"

// Formatter log formater interface
type Formatter interface {
	Write(w io.Writer, le *Event)
}

// TextFmtSubject subject log format "[%l] %m"
var TextFmtSubject = newTextFormatter("[%l] %m")

// TextFmtSimple simple log format "[%p] %m%n"
var TextFmtSimple = newTextFormatter("[%p] %m%n")

// TextFmtDefault default log format "%t %l{-5s} %c %S:%L %F() - %m%n%T"
var TextFmtDefault = newTextFormatter("%t %l{-5s} %c %S:%L %F() - %m%n%T")

// JSONFmtDefault default log format `{"time": %t, "level": %l, "name": %c, "file": %S, "line": %L, "func": %F, "msg": %m, "trace": %T}%n`
var JSONFmtDefault = newJSONFormatter(`{"time": %t, "level": %l, "name": %c, "file": %S, "line": %L, "func": %F, "msg": %m, "trace": %T}%n`)

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
		var fmt fmtfunc
		switch format[i] {
		case 't':
			p := getFormatOption(format, &i)
			if p == "" {
				p = defaultTimeFormat
			}
			fmt = fcTime(p)
		case 'c':
			p := getFormatOption(format, &i)
			if p == "" {
				fmt = ffName
			} else {
				fmt = fcName("%" + p)
			}
		case 'p':
			p := getFormatOption(format, &i)
			if p == "" {
				fmt = ffLevelPrefix
			} else {
				fmt = fcLevelPrefix("%" + p)
			}
		case 'l':
			p := getFormatOption(format, &i)
			if p == "" {
				fmt = ffLevelString
			} else {
				fmt = fcLevelString("%" + p)
			}
		case 'x':
			p := getFormatOption(format, &i)
			if p != "" {
				fmt = fcProp(p)
			}
		case 'X':
			p := getFormatOption(format, &i)
			if p == "" {
				p = "=| "
			}
			fmt = fcProps(p)
		case 'S':
			fmt = ffFile
		case 'L':
			fmt = ffLine
		case 'F':
			fmt = ffFunc
		case 'T':
			fmt = ffTrace
		case 'm':
			fmt = ffMsg
		case 'q':
			fmt = ffQmsg
		case 'n':
			fmt = ffEol
		}

		if fmt != nil {
			fmts = append(fmts, fmt)
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
		var fmt fmtfunc
		switch format[i] {
		case 't':
			p := getFormatOption(format, &i)
			if p == "" {
				p = defaultTimeFormat
			}
			fmt = fcQuote(fcTime(p))
		case 'c':
			p := getFormatOption(format, &i)
			if p == "" {
				fmt = ffName
			} else {
				fmt = fcName("%" + p)
			}
			fmt = fcQuote(fmt)
		case 'p':
			p := getFormatOption(format, &i)
			if p == "" {
				fmt = ffLevelPrefix
			} else {
				fmt = fcLevelPrefix("%" + p)
			}
			fmt = fcQuote(fmt)
		case 'l':
			p := getFormatOption(format, &i)
			if p == "" {
				fmt = ffLevelString
			} else {
				fmt = fcLevelString("%" + p)
			}
			fmt = fcQuote(fmt)
		case 'x':
			p := getFormatOption(format, &i)
			if p != "" {
				fmt = fcPropJSON(p)
			}
		case 'X':
			fmt = fcPropsJSON
		case 'S':
			fmt = fcQuote(ffFile)
		case 'L':
			fmt = ffLine
		case 'F':
			fmt = fcQuote(ffFunc)
		case 'T':
			fmt = fcQuote(ffTrace)
		case 'm':
			fmt = fcQuote(ffMsg)
		case 'n':
			fmt = ffEol
		}

		if fmt != nil {
			fmts = append(fmts, fmt)
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
	return le.Logger.GetName()
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

func fcTime(layout string) fmtfunc {
	return func(le *Event) string {
		return le.Time.Format(layout)
	}
}

func fcName(f string) fmtfunc {
	return func(le *Event) string {
		return fmt.Sprintf(f, le.Logger.GetName())
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

func fcProp(key string) fmtfunc {
	return func(le *Event) string {
		v := le.Logger.GetProp(key)
		if v == nil {
			return ""
		}
		return fmt.Sprint(v)
	}
}

func fcProps(f string) fmtfunc {
	ss := strings.Split(f, "|")
	d := ss[0]
	j := ""
	if len(ss) > 1 {
		j = ss[1]
	}
	return func(le *Event) string {
		m := le.Logger.GetProps()
		if m == nil {
			return ""
		}

		a := make([]string, 0, len(m))
		for k, v := range m {
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
		b, _ := json.Marshal(le.Logger.GetProp(key))
		return str.UnsafeString(b)
	}
}

func fcPropsJSON(le *Event) string {
	b, _ := json.Marshal(le.Logger.GetProps())
	return str.UnsafeString(b)
}
