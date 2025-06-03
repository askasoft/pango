package ref

import (
	"fmt"
	"reflect"

	"github.com/askasoft/pango/str"
)

func GetProperty(o any, k string) (any, error) {
	rv := reflect.ValueOf(o)

	switch rv.Kind() {
	case reflect.Map:
		return mapGet(rv, k)
	case reflect.Ptr:
		return getProperty(rv, k)
	default:
		return nil, fmt.Errorf("ref: %T is not a pointer", o)
	}
}

func getProperty(rv reflect.Value, k string) (v any, err error) {
	defer func() {
		if er := recover(); er != nil {
			err = fmt.Errorf("ref: GetProperty(%v, %q): %v", rv.Type(), k, er)
		}
	}()

	p := str.Capitalize(k)

	m := rv.MethodByName("Get" + p)
	if m.IsValid() && m.Type().NumIn() == 0 && (m.Type().NumOut() == 1 || m.Type().NumOut() == 2) {
		rs := m.Call(nil)
		v = rs[0].Interface()
		if len(rs) == 2 {
			if er, ok := rs[1].Interface().(error); ok {
				err = er
			}
		}
		return
	}

	f := rv.Elem().FieldByName(p)
	if f.IsValid() {
		v = f.Interface()
		return
	}

	return nil, fmt.Errorf("ref: missing property %q of %v", k, rv.Type())
}

func SetProperty(o any, k string, v any) error {
	rv := reflect.ValueOf(o)

	switch rv.Kind() {
	case reflect.Map:
		_, err := mapSet(rv, k, v)
		return err
	case reflect.Ptr:
		return setProperty(rv, k, v)
	default:
		return fmt.Errorf("ref: %T is not a pointer", o)
	}
}

func setProperty(rv reflect.Value, k string, v any) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = fmt.Errorf("ref: SetProperty(%v, %q, %T): %v", rv.Type(), k, v, er)
		}
	}()

	p := str.Capitalize(k)

	m := rv.MethodByName("Set" + p)
	if m.IsValid() && m.Type().NumIn() == 1 {
		t := m.Type().In(0)

		i, err := CastTo(v, t)
		if err != nil {
			return err
		}

		rs := m.Call([]reflect.Value{reflect.ValueOf(i)})
		for _, r := range rs {
			if err, ok := r.Interface().(error); ok {
				return err
			}
		}
		return nil
	}

	f := rv.Elem().FieldByName(p)
	if f.IsValid() && f.CanSet() {
		t := f.Type()
		cv, err := CastTo(v, t)
		if err != nil {
			return err
		}

		f.Set(reflect.ValueOf(cv))
		return nil
	}

	return fmt.Errorf("ref: missing property %q of %v", k, rv.Type())
}
