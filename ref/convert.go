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
	return nil, fmt.Errorf("cannot convert '%T' to '%s'", v, t.String())
}

// CastTo cast the value v to the specified Type t
func CastTo(v any, t reflect.Type) (cv any, err error) {
	switch t {
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
			cv, err = ConvertTo(v, t)
		default:
			// case reflect.Complex64, reflect.Complex128:
			// case reflect.Pointer:
			// case reflect.Map:
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

			rv := reflect.ValueOf(v)
			if rv.Type().ConvertibleTo(t) {
				return rv.Convert(t).Interface(), nil
			}

			err = fmt.Errorf("cannot convert %T to %s", v, t.String())
		}
	}

	return
}
