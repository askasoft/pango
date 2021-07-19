package log

import (
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

func testFormatEvent(lf Formatter, le *Event) string {
	sb := &strings.Builder{}
	lf.Write(sb, le)
	return sb.String()
}

func assertFormatEvent(t *testing.T, lf Formatter, le *Event, want string) {
	a := testFormatEvent(lf, le)

	if a != want {
		t.Errorf("\nactual: %v\nexcept: %v", a, want)
	}
}

func TestTextFormatSimple(t *testing.T) {
	tf := TextFmtSimple
	le := newEvent(&logger{}, LevelInfo, "simple")
	le.when = time.Time{}

	assertFormatEvent(t, tf, le, `[I] simple`+eol)
}

func TestTextFormatDefault(t *testing.T) {
	tf := TextFmtDefault
	le := newEvent(&logger{}, LevelInfo, "default")
	le.when = time.Time{}
	le.Caller(2, false)

	assertFormatEvent(t, tf, le, `0001-01-01T00:00:00.000 INFO  logformatter_test.go:`+
		strconv.Itoa(le.Line())+` log.TestTextFormatDefault() - default`+eol)
}

func TestTextFormatDate(t *testing.T) {
	tf := NewTextFormatter("%t - %m")
	le := newEvent(&logger{}, LevelInfo, "date")
	le.when = time.Time{}

	assertFormatEvent(t, tf, le, `0001-01-01T00:00:00.000 - date`)
}

func TestTextFormatProp(t *testing.T) {
	tf := NewTextFormatter("%x{a} %x{-}")
	lg := NewLog().GetLogger("")
	lg.SetProp("a", "av")
	le := newEvent(lg, LevelInfo, "prop")
	le.when = time.Time{}

	assertFormatEvent(t, tf, le, `av <nil>`)
}

func TestTextFormatProps1(t *testing.T) {
	tf := NewTextFormatter("%X")
	lg := NewLog().GetLogger("")
	lg.SetProp("a", "av")
	lg.SetProp("b", "bv")
	lg.SetProp("c", "cv")
	lg.SetProp("n", 11)
	lg.SetProp("x", nil)
	le := newEvent(lg, LevelInfo, "props")
	le.when = time.Time{}

	exp := strings.Split(`a=av b=bv c=cv n=11 x=<nil>`, " ")
	sort.Strings(exp)
	act := strings.Split(testFormatEvent(tf, le), " ")
	sort.Strings(act)

	if reflect.DeepEqual(exp, act) {
		t.Errorf("\nactual: %v\nexcept: %v", act, exp)
	}
}

func TestTextFormatProps2(t *testing.T) {
	tf := NewTextFormatter("%X{=|,}")
	lg := NewLog().GetLogger("")
	lg.SetProp("a", "av")
	lg.SetProp("b", "bv")
	lg.SetProp("c", "cv")
	lg.SetProp("n", 11)
	lg.SetProp("x", nil)
	le := newEvent(lg, LevelInfo, "props")
	le.when = time.Time{}

	exp := strings.Split(`a=av,b=bv,c=cv,n=11,x=<nil>`, ",")
	sort.Strings(exp)
	act := strings.Split(testFormatEvent(tf, le), ",")
	sort.Strings(act)

	if reflect.DeepEqual(exp, act) {
		t.Errorf("\nactual: %v\nexcept: %v", act, exp)
	}
}

func TestNewTextFormatSimple(t *testing.T) {
	tf := NewTextFormatter("SIMPLE")
	le := newEvent(&logger{}, LevelInfo, "simple")
	le.when = time.Time{}

	assertFormatEvent(t, tf, le, `[I] simple`+eol)
}

func TestNewTextFormatSubject(t *testing.T) {
	tf := NewTextFormatter("SUBJECT")
	le := newEvent(&logger{}, LevelInfo, "subject")
	le.when = time.Time{}

	assertFormatEvent(t, tf, le, `[INFO] subject`)
}

func TestNewTextFormatDefault(t *testing.T) {
	tf := NewTextFormatter("DEFAULT")
	le := newEvent(&logger{}, LevelInfo, "default")
	le.when = time.Time{}
	le.Caller(2, false)

	assertFormatEvent(t, tf, le, `0001-01-01T00:00:00.000 INFO  logformatter_test.go:`+
		strconv.Itoa(le.Line())+` log.TestNewTextFormatDefault() - default`+eol)
}

func TestNewLogFormatTextDefault(t *testing.T) {
	tf := NewLogFormatter("text:DEFAULT")
	le := newEvent(&logger{}, LevelInfo, "default")
	le.when = time.Time{}
	le.Caller(2, false)

	assertFormatEvent(t, tf, le, `0001-01-01T00:00:00.000 INFO  logformatter_test.go:`+
		strconv.Itoa(le.Line())+` log.TestNewLogFormatTextDefault() - default`+eol)
}

func TestJSONFormatDefault(t *testing.T) {
	jf := JSONFmtDefault
	le := newEvent(&logger{}, LevelInfo, "default")
	le.when = time.Now()
	le.Caller(2, false)

	assertFormatEvent(t, jf, le, `{"when": "`+le.when.Format(defaultTimeFormat)+
		`", "level": "INFO", "file": "logformatter_test.go", "line": `+strconv.Itoa(le.Line())+
		`, "func": "log.TestJSONFormatDefault", "msg": "default", "trace": ""}`+eol)
}

func TestJSONFormatProp(t *testing.T) {
	jf := NewJSONFormatter(`{"a":%x{a}, "n":%x{n}, "-":%x{-}}`)
	log := NewLog()
	lg := log.GetLogger("")
	lg.SetProp("a", "av")
	log.SetProp("n", 11)
	le := newEvent(lg, LevelInfo, "prop")
	le.when = time.Time{}

	assertFormatEvent(t, jf, le, `{"a":"av", "n":11, "-":null}`)
}

func TestJSONFormatProps(t *testing.T) {
	jf := NewJSONFormatter(`{"m":%X}`)
	log := NewLog()
	lg := log.GetLogger("")
	lg.SetProp("a", "av")
	lg.SetProp("b", "bv")
	lg.SetProp("c", "cv")
	log.SetProp("n", 11)
	log.SetProp("x", nil)
	le := newEvent(lg, LevelInfo, "props")
	le.when = time.Time{}

	assertFormatEvent(t, jf, le, `{"m":{"a":"av","b":"bv","c":"cv","n":11,"x":null}}`)
}

func TestNewJSONFormatDefault(t *testing.T) {
	jf := NewJSONFormatter("DEFAULT")
	le := newEvent(&logger{}, LevelInfo, "default")
	le.when = time.Now()
	le.Caller(2, false)

	assertFormatEvent(t, jf, le, `{"when": "`+le.when.Format(defaultTimeFormat)+
		`", "level": "INFO", "file": "logformatter_test.go", "line": `+strconv.Itoa(le.Line())+
		`, "func": "log.TestNewJSONFormatDefault", "msg": "default", "trace": ""}`+eol)
}

func TestNewLogFormatJSONDefault(t *testing.T) {
	jf := NewLogFormatter("json:DEFAULT")
	le := newEvent(&logger{}, LevelInfo, "default")
	le.when = time.Now()
	le.Caller(2, false)

	assertFormatEvent(t, jf, le, `{"when": "`+le.when.Format(defaultTimeFormat)+
		`", "level": "INFO", "file": "logformatter_test.go", "line": `+strconv.Itoa(le.Line())+
		`, "func": "log.TestNewLogFormatJSONDefault", "msg": "default", "trace": ""}`+eol)
}
