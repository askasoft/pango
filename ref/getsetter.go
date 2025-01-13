package ref

import (
	"fmt"
	"reflect"

	"github.com/askasoft/pango/str"
)

func GetProperty(o any, k string) (v any, err error) {
	if !IsPtrType(o) {
		return nil, fmt.Errorf("%T is not a pointer", o)
	}

	defer func() {
		if er := recover(); er != nil {
			err = fmt.Errorf("GetProperty(%T, %q): %v", o, k, er)
		}
	}()

	p := str.Capitalize(k)
	r := reflect.ValueOf(o)

	m := r.MethodByName("Get" + p)
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

	f := r.Elem().FieldByName(p)
	if f.IsValid() {
		v = f.Interface()
		return
	}

	return nil, fmt.Errorf("Missing property %q of %v", k, r.Type())
}

func SetProperty(o any, k string, v any) (err error) {
	if !IsPtrType(o) {
		return fmt.Errorf("%T is not a pointer", o)
	}

	defer func() {
		if er := recover(); er != nil {
			err = fmt.Errorf("SetProperty(%T, %q, %T): %v", o, k, v, er)
		}
	}()

	p := str.Capitalize(k)
	r := reflect.ValueOf(o)

	m := r.MethodByName("Set" + p)
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

	f := r.Elem().FieldByName(p)
	if f.IsValid() && f.CanSet() {
		t := f.Type()
		cv, err := CastTo(v, t)
		if err != nil {
			return err
		}

		f.Set(reflect.ValueOf(cv))
		return nil
	}

	return fmt.Errorf("Missing property %q of %v", k, r.Type())
}
