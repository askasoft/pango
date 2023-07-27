package log

import (
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
		t.Errorf("\n actual: %v\n expect: %v", a, want)
	}
}

func TestTextFormatSimple(t *testing.T) {
	tf := TextFmtSimple
	le := newEvent(&logger{}, LevelInfo, "simple")
	le.When = time.Time{}

	assertFormatEvent(t, tf, le, `[I] simple`+EOL)
}

func TestTextFormatDefault(t *testing.T) {
	tf := TextFmtDefault
	le := newEvent(&logger{name: "TEXT"}, LevelInfo, "default")
	le.When = time.Time{}
	le.CallerDepth(2, false)

	assertFormatEvent(t, tf, le, `0001-01-01T00:00:00.000 INFO  TEXT logformatter_test.go:`+
		strconv.Itoa(le.Line)+` log.TestTextFormatDefault() - default`+EOL)
}

func TestTextFormatDate(t *testing.T) {
	tf := NewTextFormatter("%t - %m")
	le := newEvent(&logger{}, LevelInfo, "date")
	le.When = time.Time{}

	assertFormatEvent(t, tf, le, `0001-01-01T00:00:00.000 - date`)
}

func TestTextFormatProp(t *testing.T) {
	tf := NewTextFormatter("%x{a} %x{-}")
	lg := NewLog().GetLogger("")
	lg.SetProp("a", "av")
	le := newEvent(lg, LevelInfo, "prop")
	le.When = time.Time{}

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
	le.When = time.Time{}

	w := `a=av b=bv c=cv n=11 x=<nil>`
	as := strings.Split(testFormatEvent(tf, le), " ")
	sort.Strings(as)
	a := strings.Join(as, " ")
	if w != a {
		t.Errorf("\n actual: %v\n expect: %v", a, w)
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
	le.When = time.Time{}

	w := `a=av,b=bv,c=cv,n=11,x=<nil>`
	as := strings.Split(testFormatEvent(tf, le), ",")
	sort.Strings(as)
	a := strings.Join(as, ",")
	if w != a {
		t.Errorf("\n actual: %v\n expect: %v", a, w)
	}
}

func TestNewTextFormatSimple(t *testing.T) {
	tf := NewTextFormatter("SIMPLE")
	le := newEvent(&logger{}, LevelInfo, "simple")
	le.When = time.Time{}

	assertFormatEvent(t, tf, le, `[I] simple`+EOL)
}

func TestNewTextFormatSubject(t *testing.T) {
	tf := NewTextFormatter("SUBJECT")
	le := newEvent(&logger{}, LevelInfo, "subject")
	le.When = time.Time{}

	assertFormatEvent(t, tf, le, `[INFO] subject`)
}

func TestNewTextFormatDefault(t *testing.T) {
	tf := NewTextFormatter("DEFAULT")
	le := newEvent(&logger{name: "TEXT"}, LevelInfo, "default")
	le.When = time.Time{}
	le.CallerDepth(2, false)

	assertFormatEvent(t, tf, le, `0001-01-01T00:00:00.000 INFO  TEXT logformatter_test.go:`+
		strconv.Itoa(le.Line)+` log.TestNewTextFormatDefault() - default`+EOL)
}

func TestNewLogFormatTextDefault(t *testing.T) {
	tf := NewLogFormatter("text:DEFAULT")
	le := newEvent(&logger{name: "TEXT"}, LevelInfo, "default")
	le.When = time.Time{}
	le.CallerDepth(2, false)

	assertFormatEvent(t, tf, le, `0001-01-01T00:00:00.000 INFO  TEXT logformatter_test.go:`+
		strconv.Itoa(le.Line)+` log.TestNewLogFormatTextDefault() - default`+EOL)
}

func TestJSONFormatDefault(t *testing.T) {
	jf := JSONFmtDefault
	le := newEvent(&logger{name: "JSON"}, LevelInfo, "default")
	le.When = time.Now()
	le.CallerDepth(2, false)

	assertFormatEvent(t, jf, le, `{"when": "`+le.When.Format(defaultTimeFormat)+
		`", "level": "INFO", "name": "JSON", "file": "logformatter_test.go", "line": `+strconv.Itoa(le.Line)+
		`, "func": "log.TestJSONFormatDefault", "msg": "default", "trace": ""}`+EOL)
}

func TestJSONFormatProp(t *testing.T) {
	jf := NewJSONFormatter(`{"a":%x{a}, "n":%x{n}, "-":%x{-}}`)
	log := NewLog()
	lg := log.GetLogger("")
	lg.SetProp("a", "av")
	log.SetProp("n", 11)
	le := newEvent(lg, LevelInfo, "prop")
	le.When = time.Time{}

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
	le.When = time.Time{}

	assertFormatEvent(t, jf, le, `{"m":{"a":"av","b":"bv","c":"cv","n":11,"x":null}}`)
}

func TestNewJSONFormatDefault(t *testing.T) {
	jf := NewJSONFormatter("DEFAULT")
	le := newEvent(&logger{name: "JSON"}, LevelInfo, "default")
	le.When = time.Now()
	le.CallerDepth(2, false)

	assertFormatEvent(t, jf, le, `{"when": "`+le.When.Format(defaultTimeFormat)+
		`", "level": "INFO", "name": "JSON", "file": "logformatter_test.go", "line": `+strconv.Itoa(le.Line)+
		`, "func": "log.TestNewJSONFormatDefault", "msg": "default", "trace": ""}`+EOL)
}

func TestNewLogFormatJSONDefault(t *testing.T) {
	jf := NewLogFormatter("json:DEFAULT")
	le := newEvent(&logger{name: "JSON"}, LevelInfo, "default")
	le.When = time.Now()
	le.CallerDepth(2, false)

	assertFormatEvent(t, jf, le, `{"when": "`+le.When.Format(defaultTimeFormat)+
		`", "level": "INFO", "name": "JSON", "file": "logformatter_test.go", "line": `+strconv.Itoa(le.Line)+
		`, "func": "log.TestNewLogFormatJSONDefault", "msg": "default", "trace": ""}`+EOL)
}
