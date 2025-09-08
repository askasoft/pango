package ref

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/askasoft/pango/str"
)

func GetProperty(o any, k string) (any, error) {
	if k == "" {
		return nil, errors.New("ref: empty property name")
	}

	rv := reflect.ValueOf(o)
	switch rv.Kind() {
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

	for _, fn := range []string{"Get" + p, p} {
		mv := rv.MethodByName(fn)
		if mv.IsValid() {
			mt := mv.Type()
			if mt.NumIn() != 0 || (mt.NumOut() != 1 && mt.NumOut() != 2) {
				return nil, fmt.Errorf("ref: invalid getter method %q of %v", fn, rv.Type())
			}

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

	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ref: missing field %q of %v", p, rv.Type())
	}

	fv := rv.FieldByName(k)
	if !fv.IsValid() {
		return nil, fmt.Errorf("ref: missing field %q of %v", p, rv.Type())
	}

	ret = fv.Interface()
	return
}

func SetProperty(o any, k string, v any) error {
	if k == "" {
		return errors.New("ref: empty property name")
	}

	rv := reflect.ValueOf(o)
	switch rv.Kind() {
	case reflect.Pointer:
		re := rv.Elem()
		switch re.Kind() {
		case reflect.Map:
			_, err := mapSet(re, k, v)
			return err
		default:
			return setProperty(rv, k, v)
		}
	case reflect.Map:
		_, err := mapSet(rv, k, v)
		return err
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

	{
		fn := "Set" + p
		mv := rv.MethodByName(fn)
		if mv.IsValid() {
			mt := mv.Type()
			if mt.NumIn() != 1 {
				return fmt.Errorf("ref: invalid setter method %q of %v", fn, rv.Type())
			}

			av, err := CastTo(v, mv.Type().In(0))
			if err != nil {
				return err
			}

			rs := mv.Call([]reflect.Value{reflect.ValueOf(av)})
			for _, r := range rs {
				if err, ok := r.Interface().(error); ok {
					return err
				}
			}
			return nil
		}
	}

	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	f := rv.FieldByName(p)
	if f.IsValid() && f.CanSet() {
		t := f.Type()
		cv, err := CastTo(v, t)
		if err != nil {
			return err
		}

		f.Set(reflect.ValueOf(cv))
		return nil
	}

	return fmt.Errorf("ref: missing field %q of %v", k, rv.Type())
}
