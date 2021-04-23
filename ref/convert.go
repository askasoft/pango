package ref

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var errNotSupport = errors.New("not support")

// Convert convert the value v to the specified Type t
func Convert(v interface{}, t reflect.Type) (ti interface{}, err error) {
	switch t.Kind() {
	case reflect.Int:
		ti, err = toInt(v)
	case reflect.Int8:
		ti, err = toInt8(v)
	case reflect.Int16:
		ti, err = toInt16(v)
	case reflect.Int32:
		ti, err = toInt32(v)
	case reflect.Int64:
		ti, err = toInt64(v)
	case reflect.Uint:
		ti, err = toUint(v)
	case reflect.Uint8:
		ti, err = toUint8(v)
	case reflect.Uint16:
		ti, err = toUint16(v)
	case reflect.Uint32:
		ti, err = toUint32(v)
	case reflect.Uint64:
		ti, err = toUint64(v)
	case reflect.Float32:
		ti, err = toFloat32(v)
	case reflect.Float64:
		ti, err = toFloat64(v)
	case reflect.Bool:
		ti, err = toBool(v)
	case reflect.String:
		ti, err = toString(v)
	default:
		sv := reflect.ValueOf(v)
		if sv.IsValid() {
			if sv.Type().ConvertibleTo(t) {
				return sv.Convert(t).Interface(), nil
			}
			if sv.IsZero() {
				return reflect.New(t).Interface(), nil
			}
			if sv.IsNil() {
				return nil, nil
			}
		}
		err = errNotSupport
	}

	if err != nil {
		err = fmt.Errorf("cannot convert value %v to type %s:%s", v, t.String(), err.Error())
	}
	return
}

func toInt(si interface{}) (interface{}, error) {
	if si == nil {
		return int(0), nil
	}

	switch s := si.(type) {
	case string:
		if s == "" {
			return 0, nil
		}
		i, err := strconv.Atoi(s)
		return i, err
	case bool:
		if s {
			return int(1), nil
		}
		return int(0), nil
	case int8:
		return int(s), nil
	case int16:
		return int(s), nil
	case int32:
		return int(s), nil
	case int64:
		return int(s), nil
	case int:
		return int(s), nil
	case uint8:
		return int(s), nil
	case uint16:
		return int(s), nil
	case uint32:
		return int(s), nil
	case uint64:
		return int(s), nil
	case uint:
		return int(s), nil
	case float32:
		return int(s), nil
	case float64:
		return int(s), nil
	}
	return nil, errNotSupport
}

func toInt8(si interface{}) (interface{}, error) {
	if si == nil {
		return int8(0), nil
	}

	switch s := si.(type) {
	case string:
		if s == "" {
			return int8(0), nil
		}
		i, err := strconv.ParseInt(s, 0, 8)
		return int8(i), err
	case bool:
		if s {
			return int8(1), nil
		}
		return int8(0), nil
	case int8:
		return int8(s), nil
	case int16:
		return int8(s), nil
	case int32:
		return int8(s), nil
	case int64:
		return int8(s), nil
	case int:
		return int8(s), nil
	case uint8:
		return int8(s), nil
	case uint16:
		return int8(s), nil
	case uint32:
		return int8(s), nil
	case uint64:
		return int8(s), nil
	case uint:
		return int8(s), nil
	case float32:
		return int8(s), nil
	case float64:
		return int8(s), nil
	}
	return nil, errNotSupport
}

func toInt16(si interface{}) (interface{}, error) {
	if si == nil {
		return int16(0), nil
	}

	switch s := si.(type) {
	case string:
		if s == "" {
			return int16(0), nil
		}
		i, err := strconv.ParseInt(s, 0, 16)
		return int16(i), err
	case bool:
		if s {
			return int16(1), nil
		}
		return int16(0), nil
	case int8:
		return int16(s), nil
	case int16:
		return int16(s), nil
	case int32:
		return int16(s), nil
	case int64:
		return int16(s), nil
	case int:
		return int16(s), nil
	case uint8:
		return int16(s), nil
	case uint16:
		return int16(s), nil
	case uint32:
		return int16(s), nil
	case uint64:
		return int16(s), nil
	case uint:
		return int16(s), nil
	case float32:
		return int16(s), nil
	case float64:
		return int16(s), nil
	}
	return nil, errNotSupport
}

func toInt32(si interface{}) (interface{}, error) {
	if si == nil {
		return int32(0), nil
	}

	switch s := si.(type) {
	case string:
		if s == "" {
			return int32(0), nil
		}
		i, err := strconv.ParseInt(s, 0, 32)
		return int32(i), err
	case bool:
		if s {
			return int32(1), nil
		}
		return int32(0), nil
	case int8:
		return int32(s), nil
	case int16:
		return int32(s), nil
	case int32:
		return int32(s), nil
	case int64:
		return int32(s), nil
	case int:
		return int32(s), nil
	case uint8:
		return int32(s), nil
	case uint16:
		return int32(s), nil
	case uint32:
		return int32(s), nil
	case uint64:
		return int32(s), nil
	case uint:
		return int32(s), nil
	case float32:
		return int32(s), nil
	case float64:
		return int32(s), nil
	}
	return nil, errNotSupport
}

func toInt64(si interface{}) (interface{}, error) {
	if si == nil {
		return int64(0), nil
	}

	switch s := si.(type) {
	case string:
		if s == "" {
			return int64(0), nil
		}
		i, err := strconv.ParseInt(s, 0, 64)
		return int64(i), err
	case bool:
		if s {
			return int64(1), nil
		}
		return int64(0), nil
	case int8:
		return int64(s), nil
	case int16:
		return int64(s), nil
	case int32:
		return int64(s), nil
	case int64:
		return int64(s), nil
	case int:
		return int64(s), nil
	case uint8:
		return int64(s), nil
	case uint16:
		return int64(s), nil
	case uint32:
		return int64(s), nil
	case uint64:
		return int64(s), nil
	case uint:
		return int64(s), nil
	case float32:
		return int64(s), nil
	case float64:
		return int64(s), nil
	}
	return nil, errNotSupport
}

func toUint(si interface{}) (interface{}, error) {
	if si == nil {
		return uint(0), nil
	}

	switch s := si.(type) {
	case string:
		if s == "" {
			return uint(0), nil
		}
		i, err := strconv.ParseUint(s, 0, strconv.IntSize)
		return uint(i), err
	case bool:
		if s {
			return uint(1), nil
		}
		return uint(0), nil
	case int8:
		return uint(s), nil
	case int16:
		return uint(s), nil
	case int32:
		return uint(s), nil
	case int64:
		return uint(s), nil
	case int:
		return uint(s), nil
	case uint8:
		return uint(s), nil
	case uint16:
		return uint(s), nil
	case uint32:
		return uint(s), nil
	case uint64:
		return uint(s), nil
	case uint:
		return uint(s), nil
	case float32:
		return uint(s), nil
	case float64:
		return uint(s), nil
	}
	return nil, errNotSupport
}

func toUint8(si interface{}) (interface{}, error) {
	if si == nil {
		return uint8(0), nil
	}

	switch s := si.(type) {
	case string:
		if s == "" {
			return uint8(0), nil
		}
		i, err := strconv.ParseUint(s, 0, 8)
		return uint8(i), err
	case bool:
		if s {
			return uint8(1), nil
		}
		return uint8(0), nil
	case int8:
		return uint8(s), nil
	case int16:
		return uint8(s), nil
	case int32:
		return uint8(s), nil
	case int64:
		return uint8(s), nil
	case int:
		return uint8(s), nil
	case uint8:
		return uint8(s), nil
	case uint16:
		return uint8(s), nil
	case uint32:
		return uint8(s), nil
	case uint64:
		return uint8(s), nil
	case uint:
		return uint8(s), nil
	case float32:
		return uint8(s), nil
	case float64:
		return uint8(s), nil
	}
	return nil, errNotSupport
}

func toUint16(si interface{}) (interface{}, error) {
	if si == nil {
		return uint16(0), nil
	}

	switch s := si.(type) {
	case string:
		if s == "" {
			return uint16(0), nil
		}
		i, err := strconv.ParseUint(s, 0, 16)
		return uint16(i), err
	case bool:
		if s {
			return uint16(1), nil
		}
		return uint16(0), nil
	case int8:
		return uint16(s), nil
	case int16:
		return uint16(s), nil
	case int32:
		return uint16(s), nil
	case int64:
		return uint16(s), nil
	case int:
		return uint16(s), nil
	case uint8:
		return uint16(s), nil
	case uint16:
		return uint16(s), nil
	case uint32:
		return uint16(s), nil
	case uint64:
		return uint16(s), nil
	case uint:
		return uint16(s), nil
	case float32:
		return uint16(s), nil
	case float64:
		return uint16(s), nil
	}
	return nil, errNotSupport
}

func toUint32(si interface{}) (interface{}, error) {
	if si == nil {
		return uint32(0), nil
	}

	switch s := si.(type) {
	case string:
		if s == "" {
			return uint32(0), nil
		}
		i, err := strconv.ParseUint(s, 0, 32)
		return uint32(i), err
	case bool:
		if s {
			return uint32(1), nil
		}
		return uint32(0), nil
	case int8:
		return uint32(s), nil
	case int16:
		return uint32(s), nil
	case int32:
		return uint32(s), nil
	case int64:
		return uint32(s), nil
	case int:
		return uint32(s), nil
	case uint8:
		return uint32(s), nil
	case uint16:
		return uint32(s), nil
	case uint32:
		return uint32(s), nil
	case uint64:
		return uint32(s), nil
	case uint:
		return uint32(s), nil
	case float32:
		return uint32(s), nil
	case float64:
		return uint32(s), nil
	}
	return nil, errNotSupport
}

func toUint64(si interface{}) (interface{}, error) {
	if si == nil {
		return uint64(0), nil
	}

	switch s := si.(type) {
	case string:
		if s == "" {
			return uint64(0), nil
		}
		i, err := strconv.ParseUint(s, 0, 64)
		return uint64(i), err
	case bool:
		if s {
			return uint64(1), nil
		}
		return uint64(0), nil
	case int8:
		return uint64(s), nil
	case int16:
		return uint64(s), nil
	case int32:
		return uint64(s), nil
	case int64:
		return uint64(s), nil
	case int:
		return uint64(s), nil
	case uint8:
		return uint64(s), nil
	case uint16:
		return uint64(s), nil
	case uint32:
		return uint64(s), nil
	case uint64:
		return uint64(s), nil
	case uint:
		return uint64(s), nil
	case float32:
		return uint64(s), nil
	case float64:
		return uint64(s), nil
	}
	return nil, errNotSupport
}

func toFloat32(si interface{}) (interface{}, error) {
	if si == nil {
		return float32(0), nil
	}

	switch s := si.(type) {
	case string:
		if s == "" {
			return float32(0), nil
		}
		i, err := strconv.ParseFloat(s, 32)
		return float32(i), err
	case bool:
		if s {
			return float32(1), nil
		}
		return float32(0), nil
	case int8:
		return float32(s), nil
	case int16:
		return float32(s), nil
	case int32:
		return float32(s), nil
	case int64:
		return float32(s), nil
	case int:
		return float32(s), nil
	case uint8:
		return float32(s), nil
	case uint16:
		return float32(s), nil
	case uint32:
		return float32(s), nil
	case uint64:
		return float32(s), nil
	case uint:
		return float32(s), nil
	case float32:
		return float32(s), nil
	case float64:
		return float32(s), nil
	}
	return nil, errNotSupport
}

func toFloat64(si interface{}) (interface{}, error) {
	if si == nil {
		return float64(0), nil
	}

	switch s := si.(type) {
	case string:
		if s == "" {
			return float64(0), nil
		}
		i, err := strconv.ParseFloat(s, 64)
		return float64(i), err
	case bool:
		if s {
			return float64(1), nil
		}
		return float64(0), nil
	case int8:
		return float64(s), nil
	case int16:
		return float64(s), nil
	case int32:
		return float64(s), nil
	case int64:
		return float64(s), nil
	case int:
		return float64(s), nil
	case uint8:
		return float64(s), nil
	case uint16:
		return float64(s), nil
	case uint32:
		return float64(s), nil
	case uint64:
		return float64(s), nil
	case uint:
		return float64(s), nil
	case float32:
		return float64(s), nil
	case float64:
		return float64(s), nil
	}
	return nil, errNotSupport
}

func toBool(si interface{}) (interface{}, error) {
	if si == nil {
		return false, nil
	}

	switch s := si.(type) {
	case string:
		if s == "" {
			return false, nil
		}
		return strconv.ParseBool(s)
	case bool:
		return s, nil
	case int8:
		return s != 0, nil
	case int16:
		return s != 0, nil
	case int32:
		return s != 0, nil
	case int64:
		return s != 0, nil
	case int:
		return s != 0, nil
	case uint8:
		return s != 0, nil
	case uint16:
		return s != 0, nil
	case uint32:
		return s != 0, nil
	case uint64:
		return s != 0, nil
	case uint:
		return s != 0, nil
	case float32:
		return s != 0, nil
	case float64:
		return s != 0, nil
	}
	return nil, errNotSupport
}

func toString(si interface{}) (interface{}, error) {
	if si == nil {
		return "", nil
	}

	switch s := si.(type) {
	case string:
		return s, nil
	case bool:
		if s {
			return "true", nil
		}
		return "false", nil
	case int8:
		return strconv.FormatInt(int64(s), 10), nil
	case int16:
		return strconv.FormatInt(int64(s), 10), nil
	case int32:
		return strconv.FormatInt(int64(s), 10), nil
	case int64:
		return strconv.FormatInt(int64(s), 10), nil
	case int:
		return strconv.FormatInt(int64(s), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint64:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint:
		return strconv.FormatUint(uint64(s), 10), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	}
	return nil, errNotSupport
}
