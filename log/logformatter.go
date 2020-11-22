package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

var eol = geteol()
var lvlPrefixs = [LevelTrace + 1]string{"N", "F", "E", "W", "I", "D", "T"}
var lvlStrings = [LevelTrace + 1]string{"NONE ", "FATAL", "ERROR", "WARN ", "INFO ", "DEBUG", "TRACE"}

// Formatter log formater interface
type Formatter interface {
	Write(bb *bytes.Buffer, le *Event)
	Format(le *Event) string
}

// TextFmtSubject subject log format "[%l] %m"
var TextFmtSubject = NewTextFormatter("[%l] %m")

// TextFmtSimple simple log format "[%p] %m%n"
var TextFmtSimple = NewTextFormatter("[%p] %m%n")

// TextFmtDefault default log format "%d{2006-01-02T15:04:05.000} %l %S:%L %F() - %m%n%T"
var TextFmtDefault = NewTextFormatter("%d{2006-01-02T15:04:05.000} %l %S:%L %F() - %m%n%T")

// JSONFmtDefault default log format `{"when":%d{2006-01-02T15:04:05.000Z07:00}, "level":%l, "file":%S, "line":%L, "func":%F, "msg": %m}%n`
var JSONFmtDefault = NewJSONFormatter(`{"when":%d{2006-01-02T15:04:05.000Z07:00}, "level":%l, "file":%S, "line":%L, "func":%F, "msg": %m}%n`)

// NewTextFormatter create a Text Formatter instance
// Text Format
// %c: logger name
// %d{format}: date
// %m: message
// %n: EOL('\n')
// %p: log level prefix
// %l: log level string
// %S: caller source file name (!!SLOW!!)
// %L: caller source line number (!!SLOW!!)
// %F: caller function name (!!SLOW!!)
// %T: caller stack trace (!!SLOW!!)
// %x{key}: logger property
// %X{=| }: logger propertys (operator|separator)
func NewTextFormatter(format string) *TextFormatter {
	tf := &TextFormatter{}
	tf.Init(format)
	return tf
}

// NewJSONFormatter create a Json Formatter instance
// JSON Format
// %c: logger name
// %d{format}: date
// %m: message
// %n: EOL('\n')
// %p: log level prefix
// %l: log level string
// %S: caller source file name (!!SLOW!!)
// %L: caller source line number (!!SLOW!!)
// %F: caller function name (!!SLOW!!)
// %T: caller stack trace (!!SLOW!!)
// %X{key}: logger value
// %X: logger propertys
func NewJSONFormatter(format string) *JSONFormatter {
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
			fmts = append(fmts, strfmt(format[s:i]))
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
			fmt = namefmt
		case 'p':
			fmt = lvlpfmt
		case 'l':
			fmt = lvlsfmt
		case 'm':
			fmt = msgfmt
		case 'n':
			fmt = eolfmt
		case 'F':
			fmt = funcfmt
		case 'S':
			fmt = filefmt
		case 'L':
			fmt = linefmt
		case 'T':
			fmt = tracefmt
		case 'd':
			p := format[i+1:]
			if len(p) > 0 && p[0] == '{' {
				e := strings.IndexByte(p, '}')
				if e > 0 {
					fmt = timefmt(p[1:e])
					i += e + 1
					break
				}
			}
			fmt = datefmt
		case 'x':
			p := format[i+1:]
			if len(p) > 0 && p[0] == '{' {
				e := strings.IndexByte(p, '}')
				if e > 0 {
					fmt = propfmt(p[1:e])
					i += e + 1
				}
			}
		case 'X':
			p := format[i+1:]
			if len(p) > 0 && p[0] == '{' {
				e := strings.IndexByte(p, '}')
				if e > 0 {
					fmt = propsfmt(p[1:e])
					i += e + 1
					break
				}
			}
			fmt = propsfmt("=| ")
		}

		if fmt != nil {
			fmts = append(fmts, fmt)
			s = i + 1
		}
	}

	if s < len(format) {
		fmts = append(fmts, strfmt(format[s:]))
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
			fmts = append(fmts, strfmt(format[s:i]))
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
			fmt = quotefmt(namefmt)
		case 'p':
			fmt = quotefmt(lvlpfmt)
		case 'l':
			fmt = quotefmt(lvlsfmt)
		case 'm':
			fmt = quotefmt(msgfmt)
		case 'n':
			fmt = eolfmt
		case 'F':
			fmt = quotefmt(funcfmt)
		case 'S':
			fmt = quotefmt(filefmt)
		case 'L':
			fmt = linefmt
		case 'T':
			fmt = quotefmt(tracefmt)
		case 'd':
			p := format[i+1:]
			if len(p) > 0 && p[0] == '{' {
				e := strings.IndexByte(p, '}')
				if e > 0 {
					fmt = quotefmt(timefmt(p[1:e]))
					i += e + 1
					break
				}
			}
			fmt = quotefmt(datefmt)
		case 'x':
			p := format[i+1:]
			if len(p) > 0 && p[0] == '{' {
				e := strings.IndexByte(p, '}')
				if e > 0 {
					fmt = jpropfmt(p[1:e])
					i += e + 1
				}
			}
		case 'X':
			fmt = jpropsfmt
		}

		if fmt != nil {
			fmts = append(fmts, fmt)
			s = i + 1
		}
	}

	if s < len(format) {
		fmts = append(fmts, strfmt(format[s:]))
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
	ss := strings.Builder{}
	for _, f := range fmts {
		s := f(le)
		ss.WriteString(s)
	}
	return ss.String()
}

func geteol() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

func quotefmt(ff fmtfunc) fmtfunc {
	return func(le *Event) string {
		return fmt.Sprintf("%q", ff(le))
	}
}

func strfmt(s string) fmtfunc {
	return func(le *Event) string {
		return s
	}
}

func lvlpfmt(le *Event) string {
	return lvlPrefixs[le.Level]
}

func lvlsfmt(le *Event) string {
	return lvlStrings[le.Level]
}

func namefmt(le *Event) string {
	return le.Logger.GetName()
}

func msgfmt(le *Event) string {
	return le.Msg
}

func eolfmt(le *Event) string {
	return eol
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

func datefmt(le *Event) string {
	return fmt.Sprint(le.When)
}

func timefmt(layout string) fmtfunc {
	return func(le *Event) string {
		return le.When.Format(layout)
	}
}

func propfmt(key string) fmtfunc {
	return func(le *Event) string {
		return fmt.Sprint(le.Logger.GetProp(key))
	}
}

func propsfmt(f string) fmtfunc {
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

func jpropfmt(key string) fmtfunc {
	return func(le *Event) string {
		b, _ := json.Marshal(le.Logger.GetProp(key))
		return string(b)
	}
}

func jpropsfmt(le *Event) string {
	b, _ := json.Marshal(le.Logger.GetProps())
	return string(b)
}
