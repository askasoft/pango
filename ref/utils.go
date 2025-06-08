package ref

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
)

func NameOfFunc(f any) string {
	return NameOfFuncValue(reflect.ValueOf(f))
}

func NameOfFuncValue(fv reflect.Value) string {
	return runtime.FuncForPC(fv.Pointer()).Name()
}

// IsNil checks if a specified object is nil or not, without Failing.
func IsNil(v any) bool {
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

func IsZero(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)

	return !rv.IsValid() || rv.IsZero()
}

func InvokeMethod(obj any, name string, args ...any) ([]any, error) {
	if name == "" {
		return nil, errors.New("ref: empty function name")
	}

	rv := reflect.ValueOf(obj)

	mv := rv.MethodByName(name)
	if !mv.IsValid() {
		if rv.Kind() == reflect.Pointer {
			rv = rv.Elem()
		}

		mv = rv.FieldByName(name)
		if !mv.IsValid() || mv.Kind() != reflect.Func {
			return nil, fmt.Errorf("ref: missing function %q of %T", name, obj)
		}
	}

	mt := mv.Type()
	if mt.NumIn() != len(args) {
		return nil, fmt.Errorf("ref: %s(): invalid argument count, want %d, got %d", NameOfFuncValue(mv), mt.NumIn(), len(args))
	}

	var avs []reflect.Value
	for i, a := range args {
		t := mt.In(i)

		v, err := CastTo(a, t)
		if err != nil {
			return nil, fmt.Errorf("ref: method %T.%q(): invalid argument #%d - %w", obj, name, i, err)
		}

		avs = append(avs, reflect.ValueOf(v))
	}

	rvs := mv.Call(avs)

	var rets []any
	for _, rv := range rvs {
		rets = append(rets, rv.Interface())
	}
	return rets, nil
}
