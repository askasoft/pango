package ref

import (
	"fmt"
	"reflect"

	"github.com/askasoft/pango/cas"
)

// ConvertTo convert the value v to the specified Type t, if v is convertible.
func ConvertTo(v any, t reflect.Type) (any, error) {
	rv := reflect.ValueOf(v)
	if rv.Type().ConvertibleTo(t) {
		return rv.Convert(t).Interface(), nil
	}
	return nil, fmt.Errorf("ref: cannot convert '%T' to '%s'", v, t.String())
}

// CastTo cast the value v to the specified Type t
func CastTo(v any, t reflect.Type) (any, error) {
	switch t {
	case nil:
		return v, nil
	case TypeDuration:
		return cas.ToDuration(v)
	case TypeTime:
		return cas.ToTime(v)
	default:
		switch t.Kind() {
		case reflect.Bool:
			return cas.ToBool(v)
		case reflect.Int:
			return cas.ToInt(v)
		case reflect.Int8:
			return cas.ToInt8(v)
		case reflect.Int16:
			return cas.ToInt16(v)
		case reflect.Int32:
			return cas.ToInt32(v)
		case reflect.Int64:
			return cas.ToInt64(v)
		case reflect.Uint:
			return cas.ToUint(v)
		case reflect.Uint8:
			return cas.ToUint8(v)
		case reflect.Uint16:
			return cas.ToUint16(v)
		case reflect.Uint32:
			return cas.ToUint32(v)
		case reflect.Uint64:
			return cas.ToUint64(v)
		case reflect.Float32:
			return cas.ToFloat32(v)
		case reflect.Float64:
			return cas.ToFloat64(v)
		case reflect.String:
			return cas.ToString(v)
		case reflect.Slice:
			if v == nil {
				return reflect.MakeSlice(t, 0, 0).Interface(), nil
			}

			rv := reflect.ValueOf(v)
			if rv.Type().ConvertibleTo(t) {
				return rv.Convert(t).Interface(), nil
			}

			switch rv.Kind() {
			case reflect.Slice, reflect.Array:
			default:
				return nil, fmt.Errorf("ref: cannot convert '%T' to '%s'", v, t.String())
			}

			et := t.Elem()

			sv := reflect.MakeSlice(t, 0, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				ev := rv.Index(i)
				if ev.Type() != et {
					cv, err := CastTo(ev.Interface(), et)
					if err != nil {
						return nil, err
					}
					ev = reflect.ValueOf(cv)
				}
				sv = reflect.Append(sv, ev)
			}
			return sv.Interface(), nil
		case reflect.Map:
			if v == nil {
				return reflect.MakeMap(t).Interface(), nil
			}

			rv := reflect.ValueOf(v)
			if rv.Type().ConvertibleTo(t) {
				return rv.Convert(t).Interface(), nil
			}

			if rv.Kind() != reflect.Map {
				return nil, fmt.Errorf("ref: cannot convert '%T' to '%s'", v, t.String())
			}

			kt, et := t.Key(), t.Elem()

			mv := reflect.MakeMapWithSize(t, rv.Len())
			for _, kv := range rv.MapKeys() {
				vv := rv.MapIndex(kv)

				if kv.Type() != kt {
					ck, err := CastTo(kv.Interface(), kt)
					if err != nil {
						return nil, err
					}
					kv = reflect.ValueOf(ck)
				}

				if vv.Type() != et {
					cv, err := CastTo(vv.Interface(), et)
					if err != nil {
						return nil, err
					}
					vv = reflect.ValueOf(cv)
				}

				mv.SetMapIndex(kv, vv)
			}
			return mv.Interface(), nil
		default:
			// case reflect.Complex64, reflect.Complex128:
			// case reflect.Pointer:
			// case reflect.Array:
			// case reflect.Chan:
			// case reflect.Func:
			// case reflect.Interface:
			// case reflect.Pointer:
			// case reflect.Struct:
			// case reflect.UnsafePointer:
			if v == nil {
				return reflect.New(t).Interface(), nil
			}
			return ConvertTo(v, t)
		}
	}
}
