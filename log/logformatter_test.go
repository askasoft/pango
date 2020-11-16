package log

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTextFormatSimple(t *testing.T) {
	tf := TextFmtSimple
	le := newEvent(&logger{}, LevelInfo, "simple")
	le.When = time.Time{}
	assert.Equal(t, `[INFO ] simple`+eol, tf.Format(le))
}

func TestTextFormatDefault(t *testing.T) {
	tf := TextFmtDefault
	le := newEvent(&logger{}, LevelInfo, "default")
	le.When = time.Time{}
	le.Caller(2, false)
	assert.Equal(t, `0001-01-01T00:00:00.000 INFO  logformatter_test.go:`+strconv.Itoa(le.Line)+` log.TestTextFormatDefault() - default`+eol, tf.Format(le))
}

func TestTextFormatDate(t *testing.T) {
	tf := NewTextFormatter("%d - %m")
	le := newEvent(&logger{}, LevelInfo, "date")
	le.When = time.Time{}
	assert.Equal(t, `0001-01-01 00:00:00 +0000 UTC - date`, tf.Format(le))
}

func TestTextFormatProp(t *testing.T) {
	tf := NewTextFormatter("%x{a} %x{-}")
	lg := &logger{}
	lg.SetProp("a", "av")
	le := newEvent(lg, LevelInfo, "prop")
	le.When = time.Time{}
	assert.Equal(t, `av <nil>`, tf.Format(le))
}

func TestTextFormatProps1(t *testing.T) {
	tf := NewTextFormatter("%X")
	lg := &logger{}
	lg.SetProp("a", "av")
	lg.SetProp("b", "bv")
	lg.SetProp("c", "cv")
	lg.SetProp("n", 11)
	lg.SetProp("x", nil)
	le := newEvent(lg, LevelInfo, "props")
	le.When = time.Time{}
	assert.Equal(t, `a=av b=bv c=cv n=11 x=<nil>`, tf.Format(le))
}

func TestTextFormatProps2(t *testing.T) {
	tf := NewTextFormatter("%X{=|,}")
	lg := &logger{}
	lg.SetProp("a", "av")
	lg.SetProp("b", "bv")
	lg.SetProp("c", "cv")
	lg.SetProp("n", 11)
	lg.SetProp("x", nil)
	le := newEvent(lg, LevelInfo, "props")
	le.When = time.Time{}
	assert.Equal(t, `a=av,b=bv,c=cv,n=11,x=<nil>`, tf.Format(le))
}

func TestJSONFormatDefault(t *testing.T) {
	jf := JSONFmtDefault
	le := newEvent(&logger{}, LevelInfo, "default")
	le.When = time.Time{}
	le.Caller(2, false)
	assert.Equal(t, `{"when":"0001-01-01T00:00:00.000Z", "level":"INFO ", "file":"logformatter_test.go", "line":`+strconv.Itoa(le.Line)+`, "func":"log.TestJSONFormatDefault", "msg": "default"}`+eol, jf.Format(le))
}

func TestJSONFormatProp(t *testing.T) {
	jf := NewJSONFormatter(`{"a":%x{a}, "n":%x{n}, "-":%x{-}}`)
	lg := &logger{}
	lg.SetProp("a", "av")
	lg.SetProp("n", 11)
	le := newEvent(lg, LevelInfo, "prop")
	le.When = time.Time{}
	assert.Equal(t, `{"a":"av", "n":11, "-":null}`, jf.Format(le))
}

func TestJSONFormatProps(t *testing.T) {
	jf := NewJSONFormatter(`{"m":%X}`)
	lg := &logger{}
	lg.SetProp("a", "av")
	lg.SetProp("b", "bv")
	lg.SetProp("c", "cv")
	lg.SetProp("n", 11)
	lg.SetProp("x", nil)
	le := newEvent(lg, LevelInfo, "props")
	le.When = time.Time{}
	assert.Equal(t, `{"m":{"a":"av","b":"bv","c":"cv","n":11,"x":null}}`, jf.Format(le))
}
