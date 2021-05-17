package log

import (
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTextFormatSimple(t *testing.T) {
	tf := TextFmtSimple
	le := newEvent(&logger{}, LevelInfo, "simple")
	le.When = time.Time{}
	assert.Equal(t, `[I] simple`+eol, tf.Format(le))
}

func TestTextFormatDefault(t *testing.T) {
	tf := TextFmtDefault
	le := newEvent(&logger{}, LevelInfo, "default")
	le.When = time.Time{}
	le.Caller(2, false)
	assert.Equal(t, `0001-01-01T00:00:00.000 INFO  logformatter_test.go:`+strconv.Itoa(le.Line)+` log.TestTextFormatDefault() - default`+eol, tf.Format(le))
}

func TestTextFormatDate(t *testing.T) {
	tf := NewTextFormatter("%t - %m")
	le := newEvent(&logger{}, LevelInfo, "date")
	le.When = time.Time{}
	assert.Equal(t, `0001-01-01T00:00:00.000 - date`, tf.Format(le))
}

func TestTextFormatProp(t *testing.T) {
	tf := NewTextFormatter("%x{a} %x{-}")
	lg := NewLog().GetLogger("")
	lg.SetProp("a", "av")
	le := newEvent(lg, LevelInfo, "prop")
	le.When = time.Time{}
	assert.Equal(t, `av <nil>`, tf.Format(le))
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

	exp := strings.Split(`a=av b=bv c=cv n=11 x=<nil>`, " ")
	sort.Strings(exp)
	act := strings.Split(tf.Format(le), " ")
	sort.Strings(act)
	assert.Equal(t, exp, act)
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

	exp := strings.Split(`a=av,b=bv,c=cv,n=11,x=<nil>`, ",")
	sort.Strings(exp)
	act := strings.Split(tf.Format(le), ",")
	sort.Strings(act)
	assert.Equal(t, exp, act)
}

func TestNewTextFormatSimple(t *testing.T) {
	tf := NewTextFormatter("SIMPLE")
	le := newEvent(&logger{}, LevelInfo, "simple")
	le.When = time.Time{}
	assert.Equal(t, `[I] simple`+eol, tf.Format(le))
}

func TestNewTextFormatSubject(t *testing.T) {
	tf := NewTextFormatter("SUBJECT")
	le := newEvent(&logger{}, LevelInfo, "subject")
	le.When = time.Time{}
	assert.Equal(t, `[INFO] subject`, tf.Format(le))
}

func TestNewTextFormatDefault(t *testing.T) {
	tf := NewTextFormatter("DEFAULT")
	le := newEvent(&logger{}, LevelInfo, "default")
	le.When = time.Time{}
	le.Caller(2, false)
	assert.Equal(t, `0001-01-01T00:00:00.000 INFO  logformatter_test.go:`+strconv.Itoa(le.Line)+` log.TestNewTextFormatDefault() - default`+eol, tf.Format(le))
}

func TestNewLogFormatTextDefault(t *testing.T) {
	tf := NewLogFormatter("text:DEFAULT")
	le := newEvent(&logger{}, LevelInfo, "default")
	le.When = time.Time{}
	le.Caller(2, false)
	assert.Equal(t, `0001-01-01T00:00:00.000 INFO  logformatter_test.go:`+strconv.Itoa(le.Line)+` log.TestNewLogFormatTextDefault() - default`+eol, tf.Format(le))
}

func TestJSONFormatDefault(t *testing.T) {
	jf := JSONFmtDefault
	le := newEvent(&logger{}, LevelInfo, "default")
	le.When = time.Now()
	le.Caller(2, false)
	assert.Equal(t, `{"when": "`+le.When.Format(defaultTimeFormat)+`", "level": "INFO", "file": "logformatter_test.go", "line": `+strconv.Itoa(le.Line)+`, "func": "log.TestJSONFormatDefault", "msg": "default", "trace": ""}`+eol, jf.Format(le))
}

func TestJSONFormatProp(t *testing.T) {
	jf := NewJSONFormatter(`{"a":%x{a}, "n":%x{n}, "-":%x{-}}`)
	log := NewLog()
	lg := log.GetLogger("")
	lg.SetProp("a", "av")
	log.SetProp("n", 11)
	le := newEvent(lg, LevelInfo, "prop")
	le.When = time.Time{}
	assert.Equal(t, `{"a":"av", "n":11, "-":null}`, jf.Format(le))
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
	assert.Equal(t, `{"m":{"a":"av","b":"bv","c":"cv","n":11,"x":null}}`, jf.Format(le))
}

func TestNewJSONFormatDefault(t *testing.T) {
	jf := NewJSONFormatter("DEFAULT")
	le := newEvent(&logger{}, LevelInfo, "default")
	le.When = time.Now()
	le.Caller(2, false)
	assert.Equal(t, `{"when": "`+le.When.Format(defaultTimeFormat)+`", "level": "INFO", "file": "logformatter_test.go", "line": `+strconv.Itoa(le.Line)+`, "func": "log.TestNewJSONFormatDefault", "msg": "default", "trace": ""}`+eol, jf.Format(le))
}

func TestNewLogFormatJSONDefault(t *testing.T) {
	jf := NewLogFormatter("json:DEFAULT")
	le := newEvent(&logger{}, LevelInfo, "default")
	le.When = time.Now()
	le.Caller(2, false)
	assert.Equal(t, `{"when": "`+le.When.Format(defaultTimeFormat)+`", "level": "INFO", "file": "logformatter_test.go", "line": `+strconv.Itoa(le.Line)+`, "func": "log.TestNewLogFormatJSONDefault", "msg": "default", "trace": ""}`+eol, jf.Format(le))
}
