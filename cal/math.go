package cal

import (
	"fmt"

	"github.com/askasoft/pango/cas"
)

// Adds returns the result of a + b[0] + b[1] ...
func Adds(a any, b ...any) (r any, err error) {
	r = a
	for _, v := range b {
		r, err = Add(r, v)
		if err != nil {
			return
		}
	}
	return
}

// Add returns the result of a + b
func Add(a, b any) (any, error) {
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
		return a, fmt.Errorf("Add: unknown type for '%T'", a)
	}
}

// Subtracts returns the result of a - b[0] - b[1] ...
func Subtracts(a any, b ...any) (r any, err error) {
	r = a
	for _, v := range b {
		r, err = Subtract(r, v)
		if err != nil {
			return
		}
	}
	return
}

// Subtract returns the result of a - b
func Subtract(a, b any) (any, error) {
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
		return a, fmt.Errorf("Subtract: unknown type for '%T'", a)
	}
}

// Multiplys returns the result of a * b[0] * b[1] ...
func Multiplys(a any, b ...any) (r any, err error) {
	r = a
	for _, v := range b {
		r, err = Multiply(r, v)
		if err != nil {
			return
		}
	}
	return
}

// Multiply returns the result of a * b
func Multiply(a, b any) (any, error) {
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
		return a, fmt.Errorf("Multiply: unknown type for '%T'", a)
	}
}

// Divides returns the result of a / b[0] / b[1] ...
func Divides(a any, b ...any) (r any, err error) {
	r = a
	for _, v := range b {
		r, err = Divide(r, v)
		if err != nil {
			return
		}
	}
	return
}

// Divide returns the result of a / b
func Divide(a, b any) (any, error) {
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
		return a, fmt.Errorf("Divide: unknown type for '%T'", a)
	}
}
