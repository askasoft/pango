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

func title(msgs ...any) string {
	if len(msgs) == 0 {
		return ""
	}

	if s, ok := msgs[0].(string); ok {
		return fmt.Sprintf(s, msgs[1:]...)
	}

	sb := &strings.Builder{}
	fmt.Fprint(sb, msgs[0])
	for _, n := range msgs[1:] {
		sb.WriteString(",")
		fmt.Fprint(sb, n)
	}
	return sb.String()
}

func nameOfFunc(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func caller() (caller string) {
	caller = "UNKNOWN"

	rpc := make([]uintptr, 30)
	n := runtime.Callers(2, rpc)
	if n > 0 {
		found := false
		frames := runtime.CallersFrames(rpc)
		for frame, next := frames.Next(); next; frame, next = frames.Next() {
			if strings.Contains(frame.File, "/pango/test/") {
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

func Nil(t TestingT, a any, msgs ...any) bool {
	if !isNil(a) {
		Fail(t, "%s> not nil\nactual: %v", title(msgs...), a)
		return false
	}
	return true
}

func NotNil(t TestingT, a any, msgs ...any) bool {
	if isNil(a) {
		Fail(t, "%s> nil\nactual: %v", title(msgs...), a)
		return false
	}
	return true
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
		reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return rv.IsNil()
	default:
		return false
	}
}

func Zero(t TestingT, a any, msgs ...any) bool {
	if !isZero(a) {
		Fail(t, "%s> not zero\nactual: %v", title(msgs...), a)
		return false
	}
	return true
}

func isZero(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)

	return !rv.IsValid() || rv.IsZero()
}

func Error(t TestingT, err error, msgs ...any) bool {
	if err == nil {
		return Fail(t, "%s> an error is expected but got nil.", title(msgs...))
	}
	return true
}

func NoError(t TestingT, err error, msgs ...any) bool {
	if err != nil {
		return Fail(t, "%s received unexpected error:\n%+v", title(msgs...), err)
	}

	return true
}

func EqualError(t TestingT, err error, want string, msgs ...any) bool {
	if !Error(t, err, msgs...) {
		return false
	}

	if err.Error() != want {
		return Fail(t, "%s> error messsage not equal\nactual: %v\n  want: %v", title(msgs...), err.Error(), want)
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
		return Fail(t, "%s> \"%v\" could not be applied builtin len()", title(msgs...), a)
	}

	if l != length {
		return Fail(t, "%s> \"%v\" should have %d item(s), but has %d", title(msgs...), a, length, l)
	}
	return true
}

func Empty(t TestingT, a any, msgs ...any) bool {
	l, ok := getLen(a)
	if !ok {
		return Fail(t, "%s> \"%v\" could not be applied builtin len()", title(msgs...), a)
	}

	if l != 0 {
		return Fail(t, "%s> \"%v\" is not empty, but has %d", title(msgs...), a, l)
	}
	return true
}

func True(t TestingT, value bool, msgs ...any) bool {
	if !value {
		return Fail(t, "%s> Should be true", title(msgs...))
	}
	return true
}

func False(t TestingT, value bool, msgs ...any) bool {
	if value {
		return Fail(t, "%s> Should be false", title(msgs...))
	}
	return true

}

func Equal(t TestingT, w, a any, msgs ...any) bool {
	if !reflect.DeepEqual(a, w) {
		return Fail(t, "%s> not equal\nactual: %v\n  want: %v", title(msgs...), a, w)
	}
	return true
}

func NotEqual(t TestingT, w, a any, msgs ...any) bool {
	if reflect.DeepEqual(a, w) {
		return Fail(t, "%s> but equal\nactual: %v\n  want: %v", title(msgs...), a, w)
	}
	return true
}

func Contains(t TestingT, s, contains string, msgs ...any) bool {
	ok := strings.Contains(s, contains)
	if !ok {
		return Fail(t, "%s> %#v does not contain %#v", title(msgs...), s, contains)
	}
	return true
}

func NotContains(t TestingT, s, contains string, msgs ...any) bool {
	ok := strings.Contains(s, contains)
	if ok {
		return Fail(t, "%s> %#v contains %#v", title(msgs...), s, contains)
	}
	return true
}

func Regexp(t TestingT, p string, a string) bool {
	re := regexp.MustCompile(p)
	if !re.Match([]byte(a)) {
		return Fail(t, "> %q not match %q", p, a)
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
		return Fail(t, "> not panic: %s()", nameOfFunc(f))
	}
	return true
}

func PanicsWithValue(t TestingT, val any, f func()) bool {
	did, err := didPanic(f)
	if !did {
		return Fail(t, "> not panic: %s()", nameOfFunc(f))
	}
	if err != val {
		return Fail(t, "> panic: %s()\nactual: %v\n  want: %v", nameOfFunc(f), err, val)
	}
	return true
}

func NotPanics(t TestingT, f func()) bool {
	if did, err := didPanic(f); did {
		return Fail(t, "> panic: %v - %s()", err, ref.NameOfFunc(f))
	}
	return true
}
