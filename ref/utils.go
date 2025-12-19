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
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Slice,
		reflect.Interface, reflect.Pointer, reflect.UnsafePointer:
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

func CallMethod(obj any, name string, args ...any) ([]any, error) {
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

	return CallFunction(mv, args)
}

func CallFunction(fv reflect.Value, args []any) ([]any, error) {
	ft := fv.Type()

	nin := ft.NumIn()

	var vars []any
	if ft.IsVariadic() {
		if nin-1 > len(args) {
			return nil, fmt.Errorf("ref: %q too few arguments, want %d~, got %d", ft, nin-1, len(args))
		}
		vars = args[nin-1:]
		args = args[:nin-1]
	} else {
		if nin != len(args) {
			return nil, fmt.Errorf("ref: %q invalid argument count, want %d, got %d", ft, nin, len(args))
		}
	}

	avs := make([]reflect.Value, 0, len(args)+len(vars))
	for i, a := range args {
		v, err := CastTo(a, ft.In(i))
		if err != nil {
			return nil, fmt.Errorf("ref: %q invalid argument #%d - %w", ft, i, err)
		}
		avs = append(avs, reflect.ValueOf(v))
	}

	if len(vars) > 0 {
		t := ft.In(nin - 1).Elem()

		for i, a := range vars {
			v, err := CastTo(a, t)
			if err != nil {
				return nil, fmt.Errorf("ref: %q invalid argument #%d - %w", ft, i+len(args), err)
			}
			avs = append(avs, reflect.ValueOf(v))
		}
	}

	rvs := fv.Call(avs)

	var rets []any
	for _, rv := range rvs {
		rets = append(rets, rv.Interface())
	}
	return rets, nil
}

func StructFieldsToMap(obj any) (map[string]any, error) {
	rv := reflect.Indirect(reflect.ValueOf(obj))

	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ref: %T is not a struct", obj)
	}

	rt := rv.Type()
	fn := rt.NumField()
	m := make(map[string]any, fn)

	for i := range fn {
		ft := rt.Field(i)
		if ft.IsExported() {
			m[ft.Name] = rv.Field(i).Interface()
		}
	}
	return m, nil
}

func IterStructFields(obj any, itf func(reflect.StructField, reflect.Value)) error {
	rv := reflect.Indirect(reflect.ValueOf(obj))

	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("ref: %T is not a struct", obj)
	}

	rt := rv.Type()
	for i := range rt.NumField() {
		itf(rt.Field(i), rv.Field(i))
	}
	return nil
}
