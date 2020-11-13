package log

import (
	"fmt"
	"strconv"
	"strings"
)

// Formatter log formater interface
type Formatter interface {
	Format(le *Event) string
}

// TextFormatter text formatter implement struct
type TextFormatter struct {
	fmts []formatter
}

type formatter interface {
	format(sb *strings.Builder, le *Event)
}

// NewFormatter create a Formatter instance
// Log Format
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
func NewFormatter(format string) Formatter {
	tf := &TextFormatter{}
	tf.Init(format)
	return tf
}

// FormatterSimple simple log format "[%l] %m"
var FormatterSimple = NewFormatter("[%l] %m")

// FormatterDefault default log format "%d{2006-01-02T15:04:05.000} %l %S:%L %F() - %m%T%n"
var FormatterDefault = NewFormatter("%d{2006-01-02T15:04:05.000} %l %S:%L %F() - %m%T%n")

var lvlPrefixs = [LevelTrace + 1]string{"N", "F", "E", "W", "I", "D", "T"}
var lvlStrings = [LevelTrace + 1]string{"NONE ", "FATAL", "ERROR", "WARN ", "INFO ", "DEBUG", "TRACE"}

// Format format the log event to a string
func (tf *TextFormatter) Format(le *Event) string {
	sb := strings.Builder{}
	for _, f := range tf.fmts {
		f.format(&sb, le)
	}
	return sb.String()
}

// Init initialize the formatter
func (tf *TextFormatter) Init(format string) {
	s := 0
	for i := 0; i < len(format); i++ {
		c := format[i]
		if c != '%' {
			continue
		}

		// string
		if s < i {
			tf.fmts = append(tf.fmts, &strfmt{format[s:i]})
		}

		i++
		s = i
		if i >= len(format) {
			break
		}

		// symbol
		var fmt formatter
		switch format[i] {
		case 'c':
			fmt = &namefmt{}
		case 'p':
			fmt = &lvlpfmt{}
		case 'l':
			fmt = &lvlsfmt{}
		case 'm':
			fmt = &msgfmt{}
		case 'n':
			fmt = &eolfmt{}
		case 'F':
			fmt = &funcfmt{}
		case 'S':
			fmt = &filefmt{}
		case 'L':
			fmt = &linefmt{}
		case 'T':
			fmt = &tracefmt{}
		case 'd':
			p := format[i+1:]
			if len(p) > 0 && p[0] == '{' {
				e := strings.IndexByte(p, '}')
				if e > 0 {
					fmt = &timefmt{p[1:e]}
					i += e + 1
				}
			}
		case 'X':
			p := format[i+1:]
			if len(p) > 0 && p[0] == '{' {
				e := strings.IndexByte(p, '}')
				if e > 0 {
					fmt = &paramfmt{p[1:e]}
					i += e + 1
				}
			}
		}

		if fmt != nil {
			tf.fmts = append(tf.fmts, fmt)
			s = i + 1
		}
	}

	if s < len(format) {
		tf.fmts = append(tf.fmts, &strfmt{format[s:]})
	}
}

type strfmt struct {
	s string
}

func (f *strfmt) format(sb *strings.Builder, le *Event) {
	sb.WriteString(f.s)
}

type lvlpfmt struct {
}

func (f *lvlpfmt) format(sb *strings.Builder, le *Event) {
	sb.WriteString(lvlPrefixs[le.Level])
}

type lvlsfmt struct {
}

func (f *lvlsfmt) format(sb *strings.Builder, le *Event) {
	sb.WriteString(lvlStrings[le.Level])
}

type namefmt struct {
}

func (f *namefmt) format(sb *strings.Builder, le *Event) {
	sb.WriteString(le.Logger.GetName())
}

type msgfmt struct {
}

func (f *msgfmt) format(sb *strings.Builder, le *Event) {
	sb.WriteString(le.Msg)
}

type eolfmt struct {
}

func (f *eolfmt) format(sb *strings.Builder, le *Event) {
	sb.WriteString("\n")
}

type funcfmt struct {
}

func (f *funcfmt) format(sb *strings.Builder, le *Event) {
	sb.WriteString(le.Func)
}

type filefmt struct {
}

func (f *filefmt) format(sb *strings.Builder, le *Event) {
	sb.WriteString(le.File)
}

type linefmt struct {
}

func (f *linefmt) format(sb *strings.Builder, le *Event) {
	sb.WriteString(strconv.Itoa(le.Line))
}

type tracefmt struct {
}

func (f *tracefmt) format(sb *strings.Builder, le *Event) {
	sb.WriteString(le.Trace)
}

type timefmt struct {
	layout string
}

func (f *timefmt) format(sb *strings.Builder, le *Event) {
	sb.WriteString(le.When.Format(f.layout))
}

type paramfmt struct {
	key string
}

func (f *paramfmt) format(sb *strings.Builder, le *Event) {
	sb.WriteString(fmt.Sprint(le.Logger.GetParam(f.key)))
}
