package ref

import (
	"fmt"
	"reflect"

	"github.com/askasoft/pango/str"
)

func SetProperty(o any, k string, v any) (err error) {
	if !IsPtrType(o) {
		return fmt.Errorf("%T is not a pointer", o)
	}

	defer func() {
		if er := recover(); er != nil {
			err = fmt.Errorf("SetProperty(%v, %q, %v): %v", o, k, v, er)
		}
	}()

	p := str.Capitalize(k)
	r := reflect.ValueOf(o)

	m := r.MethodByName("Set" + p)
	if m.IsValid() && m.Type().NumIn() == 1 {
		t := m.Type().In(0)

		i, err := Convert(v, t)
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
		cv, err := Convert(v, t)
		if err != nil {
			return err
		}

		f.Set(reflect.ValueOf(cv))
		return nil
	}

	return fmt.Errorf("Missing property %q of %v", k, r.Type())
}
