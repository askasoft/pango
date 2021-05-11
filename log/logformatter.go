package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pandafw/pango/iox"
)

var eol = iox.EOL

const defaultTimeFormat = "2006-01-02T15:04:05.000"

// Formatter log formater interface
type Formatter interface {
	Write(bb *bytes.Buffer, le *Event)
	Format(le *Event) string
}

// TextFmtSubject subject log format "[%l] %m"
var TextFmtSubject = newTextFormatter("[%l] %m")

// TextFmtSimple simple log format "[%p] %m%n"
var TextFmtSimple = newTextFormatter("[%p] %m%n")

// TextFmtDefault default log format "%t %l{-5s} %S:%L %F() - %m%n%T"
var TextFmtDefault = newTextFormatter("%t %l{-5s} %S:%L %F() - %m%n%T")

// JSONFmtDefault default log format `{"when": %t, "level": %l, "file": %S, "line": %L, "func": %F, "msg": %m, "trace": %T}%n`
var JSONFmtDefault = newJSONFormatter(`{"when": %t, "level": %l, "file": %S, "line": %L, "func": %F, "msg": %m, "trace": %T}%n`)

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
// %S: caller source file name (!!SLOW!!)
// %L: caller source line number (!!SLOW!!)
// %F: caller function name (!!SLOW!!)
// %T: caller stack trace (!!SLOW!!)
// %m: message
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
	tf.Init(format)
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
// %S: caller source file name (!!SLOW!!)
// %L: caller source line number (!!SLOW!!)
// %F: caller function name (!!SLOW!!)
// %T: caller stack trace (!!SLOW!!)
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
	jf.Init(format)
	return jf
}

// TextFormatter text formatter
type TextFormatter struct {
	fmts []fmtfunc
}

// Format format the log event to the buffer 'bb'
func (tf *TextFormatter) Write(bb *bytes.Buffer, le *Event) {
	write(bb, le, tf.fmts)
}

// Format format the log event to a string
func (tf *TextFormatter) Format(le *Event) string {
	return format(le, tf.fmts)
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

// Init initialize the text formatter
func (tf *TextFormatter) Init(format string) {
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
		case 't':
			p := getFormatOption(format, &i)
			if p == "" {
				p = defaultTimeFormat
			}
			fmt = timefmtc(p)
		case 'c':
			p := getFormatOption(format, &i)
			if p == "" {
				fmt = namefmt
			} else {
				fmt = namefmtc("%" + p)
			}
		case 'p':
			p := getFormatOption(format, &i)
			if p == "" {
				fmt = lvlpfmt
			} else {
				fmt = lvlpfmtc("%" + p)
			}
		case 'l':
			p := getFormatOption(format, &i)
			if p == "" {
				fmt = lvlsfmt
			} else {
				fmt = lvlsfmtc("%" + p)
			}
		case 'x':
			p := getFormatOption(format, &i)
			if p != "" {
				fmt = propfmtc(p)
			}
		case 'X':
			p := getFormatOption(format, &i)
			if p == "" {
				p = "=| "
			}
			fmt = propsfmtc(p)
		case 'S':
			fmt = filefmt
		case 'L':
			fmt = linefmt
		case 'F':
			fmt = funcfmt
		case 'T':
			fmt = tracefmt
		case 'm':
			fmt = msgfmt
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

	tf.fmts = fmts
}

// JSONFormatter json formatter
type JSONFormatter struct {
	fmts []fmtfunc
}

// Write format the log event as a json string to the buffer 'bb'
func (jf *JSONFormatter) Write(bb *bytes.Buffer, le *Event) {
	write(bb, le, jf.fmts)
}

// Format format the log event to a json string
func (jf *JSONFormatter) Format(le *Event) string {
	return format(le, jf.fmts)
}

// Init initialize the json formatter
func (jf *JSONFormatter) Init(format string) {
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
		case 't':
			p := getFormatOption(format, &i)
			if p == "" {
				p = defaultTimeFormat
			}
			fmt = quotefmtc(timefmtc(p))
		case 'c':
			p := getFormatOption(format, &i)
			if p == "" {
				fmt = namefmt
			} else {
				fmt = namefmtc("%" + p)
			}
			fmt = quotefmtc(fmt)
		case 'p':
			p := getFormatOption(format, &i)
			if p == "" {
				fmt = lvlpfmt
			} else {
				fmt = lvlpfmtc("%" + p)
			}
			fmt = quotefmtc(fmt)
		case 'l':
			p := getFormatOption(format, &i)
			if p == "" {
				fmt = lvlsfmt
			} else {
				fmt = lvlsfmtc("%" + p)
			}
			fmt = quotefmtc(fmt)
		case 'x':
			p := getFormatOption(format, &i)
			if p != "" {
				fmt = jpropfmtc(p)
			}
		case 'X':
			fmt = jpropsfmt
		case 'S':
			fmt = quotefmtc(filefmt)
		case 'L':
			fmt = linefmt
		case 'F':
			fmt = quotefmtc(funcfmt)
		case 'T':
			fmt = quotefmtc(tracefmt)
		case 'm':
			fmt = quotefmtc(msgfmt)
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
	jf.fmts = fmts
}

//-------------------------------------------------

type fmtfunc func(le *Event) string

func write(bb *bytes.Buffer, le *Event, fmts []fmtfunc) {
	for _, f := range fmts {
		s := f(le)
		bb.WriteString(s)
	}
}

func format(le *Event, fmts []fmtfunc) string {
	ss := &strings.Builder{}
	for _, f := range fmts {
		s := f(le)
		ss.WriteString(s)
	}
	return ss.String()
}

func quotefmtc(ff fmtfunc) fmtfunc {
	return func(le *Event) string {
		return fmt.Sprintf("%q", ff(le))
	}
}

func strfmtc(s string) fmtfunc {
	return func(le *Event) string {
		return s
	}
}

func timefmtc(layout string) fmtfunc {
	return func(le *Event) string {
		return le.When.Format(layout)
	}
}

func namefmt(le *Event) string {
	return le.Logger.GetName()
}

func namefmtc(f string) fmtfunc {
	return func(le *Event) string {
		return fmt.Sprintf(f, le.Logger.GetName())
	}
}

func lvlpfmt(le *Event) string {
	return le.Level.Prefix()
}

func lvlpfmtc(f string) fmtfunc {
	return func(le *Event) string {
		return fmt.Sprintf(f, le.Level.Prefix())
	}
}

func lvlsfmt(le *Event) string {
	return le.Level.String()
}

func lvlsfmtc(f string) fmtfunc {
	return func(le *Event) string {
		return fmt.Sprintf(f, le.Level.String())
	}
}

func propfmtc(key string) fmtfunc {
	return func(le *Event) string {
		return fmt.Sprint(le.Logger.GetProp(key))
	}
}

func propsfmtc(f string) fmtfunc {
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
			a = append(a, fmt.Sprintf("%s%s%v", k, d, v))
		}
		return strings.Join(a, j)
	}
}

func jpropfmtc(key string) fmtfunc {
	return func(le *Event) string {
		b, _ := json.Marshal(le.Logger.GetProp(key))
		return string(b)
	}
}

func jpropsfmt(le *Event) string {
	b, _ := json.Marshal(le.Logger.GetProps())
	return string(b)
}

func funcfmt(le *Event) string {
	return le.Func
}

func filefmt(le *Event) string {
	return le.File
}

func linefmt(le *Event) string {
	return strconv.Itoa(le.Line)
}

func tracefmt(le *Event) string {
	return le.Trace
}

func msgfmt(le *Event) string {
	return le.Msg
}

func eolfmt(le *Event) string {
	return eol
}
