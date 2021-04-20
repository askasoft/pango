package ref

import (
	"fmt"
	"reflect"
	"strconv"
)

// Convert convert the value v to the specified Type t
func Convert(v interface{}, t reflect.Type) (reflect.Value, error) {
	src := reflect.ValueOf(v)

	if src.Type().ConvertibleTo(t) {
		return src.Convert(t), nil
	}

	if src.Kind() == reflect.String {
		s := v.(string)
		switch t.Kind() {
		case reflect.Int:
			i, err := strconv.Atoi(s)
			return reflect.ValueOf(i), err
		case reflect.Int8:
			i, err := strconv.ParseInt(s, 10, 8)
			return reflect.ValueOf(int8(i)), err
		case reflect.Int16:
			i, err := strconv.ParseInt(s, 10, 16)
			return reflect.ValueOf(int16(i)), err
		case reflect.Int32:
			i, err := strconv.ParseInt(s, 10, 32)
			return reflect.ValueOf(int32(i)), err
		case reflect.Int64:
			i, err := strconv.ParseInt(s, 10, 64)
			return reflect.ValueOf(int64(i)), err
		case reflect.Uint:
			i, err := strconv.Atoi(s)
			return reflect.ValueOf(uint(i)), err
		case reflect.Uint8:
			i, err := strconv.ParseUint(s, 10, 8)
			return reflect.ValueOf(uint8(i)), err
		case reflect.Uint16:
			i, err := strconv.ParseUint(s, 10, 16)
			return reflect.ValueOf(uint16(i)), err
		case reflect.Uint32:
			i, err := strconv.ParseUint(s, 10, 32)
			return reflect.ValueOf(uint32(i)), err
		case reflect.Uint64:
			i, err := strconv.ParseUint(s, 10, 64)
			return reflect.ValueOf(uint64(i)), err

		case reflect.Float32:
			f, err := strconv.ParseFloat(s, 32)
			return reflect.ValueOf(float32(f)), err

		case reflect.Float64:
			f, err := strconv.ParseFloat(s, 64)
			return reflect.ValueOf(f), err

		case reflect.Bool:
			b, err := strconv.ParseBool(s)
			return reflect.ValueOf(b), err
		}
	}

	return reflect.ValueOf(nil), fmt.Errorf("value of type " + src.Type().String() + " cannot be converted to type " + t.String())
}
