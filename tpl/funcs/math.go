package funcs

import (
	"fmt"

	"github.com/askasoft/pango/cas"
)

// Add returns the result of a + b[0] + b[1] ...
func Add(a any, b ...any) (r any, err error) {
	r = a
	for _, v := range b {
		r, err = add(r, v)
		if err != nil {
			return
		}
	}
	return
}

func add(a, b any) (any, error) {
	switch na := a.(type) {
	case int:
		nb, err := cas.ToInt(b)
		return na + nb, err
	case int8:
		nb, err := cas.ToInt8(b)
		return na + nb, err
	case int16:
		nb, err := cas.ToInt16(b)
		return na + nb, err
	case int32:
		nb, err := cas.ToInt32(b)
		return na + nb, err
	case int64:
		nb, err := cas.ToInt64(b)
		return na + nb, err
	case uint:
		nb, err := cas.ToUint(b)
		return na + nb, err
	case uint8:
		nb, err := cas.ToUint8(b)
		return na + nb, err
	case uint16:
		nb, err := cas.ToUint16(b)
		return na + nb, err
	case uint32:
		nb, err := cas.ToUint32(b)
		return na + nb, err
	case uint64:
		nb, err := cas.ToUint64(b)
		return na + nb, err
	case float32:
		nb, err := cas.ToFloat32(b)
		return na + nb, err
	case float64:
		nb, err := cas.ToFloat64(b)
		return na + nb, err
	default:
		return nil, fmt.Errorf("Add: unknown type for '%T'", a)
	}
}

// Subtract returns the result of a - b[0] - b[1] ...
func Subtract(a any, b ...any) (r any, err error) {
	r = a
	for _, v := range b {
		r, err = subtract(r, v)
		if err != nil {
			return
		}
	}
	return
}

func subtract(a, b any) (any, error) {
	switch na := a.(type) {
	case int:
		nb, err := cas.ToInt(b)
		return na - nb, err
	case int8:
		nb, err := cas.ToInt8(b)
		return na - nb, err
	case int16:
		nb, err := cas.ToInt16(b)
		return na - nb, err
	case int32:
		nb, err := cas.ToInt32(b)
		return na - nb, err
	case int64:
		nb, err := cas.ToInt64(b)
		return na - nb, err
	case uint:
		nb, err := cas.ToUint(b)
		return na - nb, err
	case uint8:
		nb, err := cas.ToUint8(b)
		return na - nb, err
	case uint16:
		nb, err := cas.ToUint16(b)
		return na - nb, err
	case uint32:
		nb, err := cas.ToUint32(b)
		return na - nb, err
	case uint64:
		nb, err := cas.ToUint64(b)
		return na - nb, err
	case float32:
		nb, err := cas.ToFloat32(b)
		return na - nb, err
	case float64:
		nb, err := cas.ToFloat64(b)
		return na - nb, err
	default:
		return nil, fmt.Errorf("Subtract: unknown type for '%T'", a)
	}
}

// Multiply returns the result of a * b[0] * b[1] ...
func Multiply(a any, b ...any) (r any, err error) {
	r = a
	for _, v := range b {
		r, err = multiply(r, v)
		if err != nil {
			return
		}
	}
	return
}

func multiply(a, b any) (any, error) {
	switch na := a.(type) {
	case int:
		nb, err := cas.ToInt(b)
		return na * nb, err
	case int8:
		nb, err := cas.ToInt8(b)
		return na * nb, err
	case int16:
		nb, err := cas.ToInt16(b)
		return na * nb, err
	case int32:
		nb, err := cas.ToInt32(b)
		return na * nb, err
	case int64:
		nb, err := cas.ToInt64(b)
		return na * nb, err
	case uint:
		nb, err := cas.ToUint(b)
		return na * nb, err
	case uint8:
		nb, err := cas.ToUint8(b)
		return na * nb, err
	case uint16:
		nb, err := cas.ToUint16(b)
		return na * nb, err
	case uint32:
		nb, err := cas.ToUint32(b)
		return na * nb, err
	case uint64:
		nb, err := cas.ToUint64(b)
		return na * nb, err
	case float32:
		nb, err := cas.ToFloat32(b)
		return na * nb, err
	case float64:
		nb, err := cas.ToFloat64(b)
		return na * nb, err
	default:
		return nil, fmt.Errorf("Multiply: unknown type for '%T'", a)
	}
}

// Divide returns the result of a / b[0] / b[1] ...
func Divide(a any, b ...any) (r any, err error) {
	r = a
	for _, v := range b {
		r, err = divide(r, v)
		if err != nil {
			return
		}
	}
	return
}

func divide(a, b any) (any, error) {
	switch na := a.(type) {
	case int:
		nb, err := cas.ToInt(b)
		return na / nb, err
	case int8:
		nb, err := cas.ToInt8(b)
		return na / nb, err
	case int16:
		nb, err := cas.ToInt16(b)
		return na / nb, err
	case int32:
		nb, err := cas.ToInt32(b)
		return na / nb, err
	case int64:
		nb, err := cas.ToInt64(b)
		return na / nb, err
	case uint:
		nb, err := cas.ToUint(b)
		return na / nb, err
	case uint8:
		nb, err := cas.ToUint8(b)
		return na / nb, err
	case uint16:
		nb, err := cas.ToUint16(b)
		return na / nb, err
	case uint32:
		nb, err := cas.ToUint32(b)
		return na / nb, err
	case uint64:
		nb, err := cas.ToUint64(b)
		return na / nb, err
	case float32:
		nb, err := cas.ToFloat32(b)
		return na / nb, err
	case float64:
		nb, err := cas.ToFloat64(b)
		return na / nb, err
	default:
		return nil, fmt.Errorf("Divide: unknown type for '%T'", a)
	}
}
