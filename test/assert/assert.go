package assert

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/askasoft/pango/ref"
)

type TestingT interface {
	Errorf(format string, args ...any)
}

// isNil checks if a specified object is nil or not, without Failing.
func isNil(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case
		reflect.Chan, reflect.Func,
		reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return rv.IsNil()
	default:
		return false
	}
}

func isZero(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)

	return !rv.IsValid() || rv.IsZero()
}

func caller() (caller string) {
	caller = "UNKNOWN"

	rpc := make([]uintptr, 30)
	n := runtime.Callers(2, rpc)
	if n > 0 {
		found := false
		frames := runtime.CallersFrames(rpc)
		for frame, next := frames.Next(); next; frame, next = frames.Next() {
			if strings.Contains(frame.File, "/pango/assert/") {
				found = true
				continue
			}

			if found {
				caller = fmt.Sprintf("%s:%d", frame.File, frame.Line)
				return
			}
		}
	}
	return
}

func Fail(t TestingT, format string, args ...any) bool {
	msgs := append([]any{caller()}, args...)
	t.Errorf("%s\n"+format, msgs...)
	return false
}

func Nil(t TestingT, a any, msgs ...string) bool {
	if !isNil(a) {
		Fail(t, "%v not nil\nactual: %v", msgs, a)
		return false
	}
	return true
}

func NotNil(t TestingT, a any, msgs ...string) bool {
	if isNil(a) {
		Fail(t, "%v nil\nactual: %v", msgs, a)
		return false
	}
	return true
}

func Zero(t TestingT, a any, msgs ...string) bool {
	if !isZero(a) {
		Fail(t, "%v not zero\nactual: %v", msgs, a)
		return false
	}
	return true
}

func Error(t TestingT, err error, msgs ...any) bool {
	if err == nil {
		return Fail(t, "An error is expected but got nil.", msgs...)
	}
	return true
}

func NoError(t TestingT, err error, msgs ...any) bool {
	if err != nil {
		return Fail(t, fmt.Sprintf("Received unexpected error:\n%+v", err), msgs...)
	}

	return true
}

func EqualError(t TestingT, err error, want string, msgs ...any) bool {
	if !Error(t, err, msgs...) {
		return false
	}

	if err.Error() != want {
		return Fail(t, "%v error messsage not equal\nactual: %v,    want: %v", msgs, err.Error(), want)
	}
	return true
}

// getLen tries to get the length of an object.
// It returns (0, false) if impossible.
func getLen(x any) (length int, ok bool) {
	v := reflect.ValueOf(x)
	defer func() {
		ok = recover() == nil
	}()
	return v.Len(), true
}

func Len(t TestingT, a any, length int, msgs ...any) bool {
	l, ok := getLen(a)
	if !ok {
		return Fail(t, fmt.Sprintf("\"%v\" could not be applied builtin len()", a), msgs...)
	}

	if l != length {
		return Fail(t, fmt.Sprintf("\"%v\" should have %d item(s), but has %d", a, length, l), msgs...)
	}
	return true
}

func Empty(t TestingT, a any, msgs ...any) bool {
	l, ok := getLen(a)
	if !ok {
		return Fail(t, fmt.Sprintf("\"%v\" could not be applied builtin len()", a), msgs...)
	}

	if l != 0 {
		return Fail(t, fmt.Sprintf("\"%v\" is not empty, but has %d", a, l), msgs...)
	}
	return true
}

func True(t TestingT, value bool, msgs ...any) bool {
	if !value {
		return Fail(t, "Should be true", msgs...)
	}
	return true
}

func False(t TestingT, value bool, msgs ...any) bool {
	if value {
		return Fail(t, "Should be false", msgs...)
	}
	return true

}

func Equal(t TestingT, w, a any, msgs ...string) bool {
	if !reflect.DeepEqual(a, w) {
		return Fail(t, "%v not equal\nactual: %v,    want: %v", msgs, a, w)
	}
	return true
}

func NotEqual(t TestingT, w, a any, msgs ...string) bool {
	if !reflect.DeepEqual(a, w) {
		return Fail(t, "%v but equal\nactual: %v,    want: %v", msgs, a, w)
	}
	return true
}

func Contains(t TestingT, s, contains string, msgs ...any) bool {
	ok := strings.Contains(s, contains)
	if !ok {
		return Fail(t, fmt.Sprintf("%#v does not contain %#v", s, contains), msgs...)
	}
	return true
}

func NotContains(t TestingT, s, contains string, msgs ...any) bool {
	ok := strings.Contains(s, contains)
	if ok {
		return Fail(t, fmt.Sprintf("%#v contains %#v", s, contains), msgs...)
	}
	return true
}

func Regexp(t TestingT, p string, a string) bool {
	re := regexp.MustCompile(p)
	if !re.Match([]byte(a)) {
		return Fail(t, "%q not match %q", p, a)
	}
	return true
}

func didPanic(f func()) (did bool, err any) {
	defer func() {
		err = recover()
		did = err != nil
	}()

	f()

	return
}

func Panics(t TestingT, f func()) bool {
	if did, _ := didPanic(f); !did {
		return Fail(t, "not panic: %s()", ref.NameOfFunc(f))
	}
	return true
}

func PanicsWithValue(t TestingT, msg string, f func()) bool {
	if did, _ := didPanic(f); !did {
		return Fail(t, "%s not panic: %s()", msg, ref.NameOfFunc(f))
	}
	return true
}

func NotPanics(t TestingT, f func()) bool {
	if did, err := didPanic(f); did {
		return Fail(t, "panic: %v - %s()", err, ref.NameOfFunc(f))
	}
	return true
}
