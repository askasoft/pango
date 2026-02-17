package ref

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/askasoft/pango/str"
)

var (
	errEmptyPropName = errors.New("ref: empty property name")
	errInvalidObject = errors.New("ref: invalid object")
)

type MissingFieldError struct {
	Type  reflect.Type
	Field string
}

func (mfe *MissingFieldError) Error() string {
	return fmt.Sprintf("ref: missing field %q of %v", mfe.Field, mfe.Type)
}

func AsMissingFieldError(err error) (mfe *MissingFieldError, ok bool) {
	ok = errors.As(err, &mfe)
	return
}

func IsMissingFieldError(err error) bool {
	_, ok := AsMissingFieldError(err)
	return ok
}

func GetProperty(o any, k string) (any, error) {
	if k == "" {
		return nil, errEmptyPropName
	}

	rv := reflect.ValueOf(o)
	switch rv.Kind() {
	case reflect.Invalid:
		return nil, errInvalidObject
	case reflect.Pointer:
		re := rv.Elem()
		switch re.Kind() {
		case reflect.Map:
			return mapGet(re, k)
		default:
			return getProperty(rv, k)
		}
	case reflect.Map:
		return mapGet(rv, k)
	default:
		return getProperty(rv, k)
	}
}

func getProperty(rv reflect.Value, k string) (ret any, err error) {
	defer func() {
		if er := recover(); er != nil {
			err = fmt.Errorf("ref: GetProperty(%v, %q): %v", rv.Type(), k, er)
		}
	}()

	p := str.Capitalize(k)

	// method(function) first
	{
		mv := rv.MethodByName(p)
		if mv.IsValid() {
			ret = mv.Interface()
			return
		}
	}

	// use getter method (java-like)
	{
		fn := "Get" + p
		mv := rv.MethodByName(fn)
		if mv.IsValid() {
			mt := mv.Type()
			if mt.NumIn() == 0 && (mt.NumOut() == 1 || mt.NumOut() == 2) {
				rs := mv.Call(nil)
				ret = rs[0].Interface()
				if len(rs) == 2 {
					if er, ok := rs[1].Interface().(error); ok {
						err = er
					}
				}
				return
			}
		}
	}

	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		err = fmt.Errorf("ref: expected %s but got %v (%v)", reflect.Struct, rv.Kind(), rv.Type())
		return
	}

	fv := rv.FieldByName(p)
	if !fv.IsValid() {
		err = &MissingFieldError{rv.Type(), k}
		return
	}

	ret = fv.Interface()
	return
}

func SetProperty(o any, k string, v any) error {
	if k == "" {
		return errEmptyPropName
	}

	rv := reflect.ValueOf(o)
	switch rv.Kind() {
	case reflect.Invalid:
		return errInvalidObject
	case reflect.Pointer:
		re := rv.Elem()
		switch re.Kind() {
		case reflect.Map:
			return mapSet(re, k, v)
		default:
			return setProperty(rv, k, v)
		}
	case reflect.Map:
		return mapSet(rv, k, v)
	default:
		return setProperty(rv, k, v)
	}
}

func setProperty(rv reflect.Value, k string, v any) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = fmt.Errorf("ref: SetProperty(%v, %q, %T): %v", rv.Type(), k, v, er)
		}
	}()

	p := str.Capitalize(k)

	// use setter method (java-like)
	{
		fn := "Set" + p
		mv := rv.MethodByName(fn)
		if mv.IsValid() {
			mt := mv.Type()
			if mt.NumIn() == 1 {
				av, er := CastTo(v, mv.Type().In(0))
				if er != nil {
					err = er
					return
				}

				rs := mv.Call([]reflect.Value{reflect.ValueOf(av)})
				for _, r := range rs {
					if er, ok := r.Interface().(error); ok {
						err = er
						return
					}
				}
				return
			}
		}
	}

	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		err = fmt.Errorf("ref: expected %s but got %v (%v)", reflect.Struct, rv.Kind(), rv.Type())
		return
	}

	f := rv.FieldByName(p)
	if f.IsValid() && f.CanSet() {
		t := f.Type()
		cv, er := CastTo(v, t)
		if er != nil {
			err = er
			return
		}

		f.Set(reflect.ValueOf(cv))
		return
	}

	err = &MissingFieldError{rv.Type(), k}
	return
}
