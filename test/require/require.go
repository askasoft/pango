package require

import (
	"github.com/askasoft/pango/test/assert"
)

type TestingT interface {
	assert.TestingT

	FailNow()
}

func Fail(t TestingT, format string, args ...any) {
	assert.Fail(t, format)
	t.FailNow()
}

func Nil(t TestingT, a any, msgs ...any) {
	if !assert.Nil(t, a, msgs...) {
		t.FailNow()
	}
}

func NotNil(t TestingT, a any, msgs ...any) {
	if !assert.NotNil(t, a, msgs...) {
		t.FailNow()
	}
}

func Zero(t TestingT, a any, msgs ...any) {
	if !assert.Zero(t, a, msgs...) {
		t.FailNow()
	}
}

func Error(t TestingT, err error, msgs ...any) {
	if !assert.Error(t, err, msgs...) {
		t.FailNow()
	}
}

func NoError(t TestingT, err error, msgs ...any) {
	if !assert.NoError(t, err, msgs...) {
		t.FailNow()
	}
}

func EqualError(t TestingT, err error, want string, msgs ...any) {
	if !assert.EqualError(t, err, want, msgs...) {
		t.FailNow()
	}
}

func Len(t TestingT, a any, length int, msgs ...any) {
	if !assert.Len(t, a, length, msgs...) {
		t.FailNow()
	}
}

func Empty(t TestingT, a any, msgs ...any) {
	if !assert.Empty(t, a, msgs...) {
		t.FailNow()
	}
}

func True(t TestingT, value bool, msgs ...any) {
	if !assert.True(t, value, msgs...) {
		t.FailNow()
	}
}

func False(t TestingT, value bool, msgs ...any) {
	if !assert.False(t, value, msgs...) {
		t.FailNow()
	}
}

func Equal(t TestingT, w, a any, msgs ...any) {
	if !assert.Equal(t, w, a, msgs...) {
		t.FailNow()
	}
}

func NotEqual(t TestingT, w, a any, msgs ...any) {
	if !assert.NotEqual(t, w, a, msgs...) {
		t.FailNow()
	}
}

func Contains(t TestingT, s, contains string, msgs ...any) {
	if !assert.Contains(t, s, contains, msgs...) {
		t.FailNow()
	}
}

func NotContains(t TestingT, s, contains string, msgs ...any) {
	if !assert.NotContains(t, s, contains, msgs...) {
		t.FailNow()
	}
}

func Regexp(t TestingT, p string, a string) {
	if !assert.Regexp(t, p, a) {
		t.FailNow()
	}
}

func Panics(t TestingT, f func()) {
	if !assert.Panics(t, f) {
		t.FailNow()
	}
}

func PanicsWithValue(t TestingT, val any, f func()) {
	if !assert.PanicsWithValue(t, val, f) {
		t.FailNow()
	}
}

func NotPanics(t TestingT, f func()) {
	if !assert.NotPanics(t, f) {
		t.FailNow()
	}
}
