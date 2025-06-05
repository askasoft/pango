package cal

import (
	"fmt"
	"time"

	"github.com/askasoft/pango/cas"
	"github.com/askasoft/pango/ref"
)

func LogicAnd(a any, vs ...any) bool {
	if ref.IsZero(a) {
		return false
	}

	for _, v := range vs {
		if ref.IsZero(v) {
			return false
		}
	}

	return true
}

func LogicOr(a any, vs ...any) bool {
	if !ref.IsZero(a) {
		return true
	}

	for _, v := range vs {
		if !ref.IsZero(v) {
			return true
		}
	}

	return false
}

// LogicEq returns the result of a == b
func LogicEq(a, b any) (bool, error) {
	if a == b {
		return true, nil
	}

	v, err := cast(a, b)
	if err != nil {
		return false, err
	}

	switch na := v.(type) {
	case time.Time:
		nb, err := cas.ToTime(b)
		return na.Equal(nb), err
	case time.Duration:
		nb, err := cas.ToDuration(b)
		return na == nb, err
	case string:
		nb, err := cas.ToString(b)
		return na == nb, err
	case int:
		nb, err := cas.ToInt(b)
		return na == nb, err
	case int8:
		nb, err := cas.ToInt8(b)
		return na == nb, err
	case int16:
		nb, err := cas.ToInt16(b)
		return na == nb, err
	case int32:
		nb, err := cas.ToInt32(b)
		return na == nb, err
	case int64:
		nb, err := cas.ToInt64(b)
		return na == nb, err
	case uint:
		nb, err := cas.ToUint(b)
		return na == nb, err
	case uint8:
		nb, err := cas.ToUint8(b)
		return na == nb, err
	case uint16:
		nb, err := cas.ToUint16(b)
		return na == nb, err
	case uint32:
		nb, err := cas.ToUint32(b)
		return na == nb, err
	case uint64:
		nb, err := cas.ToUint64(b)
		return na == nb, err
	case float32:
		nb, err := cas.ToFloat32(b)
		return na == nb, err
	case float64:
		nb, err := cas.ToFloat64(b)
		return na == nb, err
	default:
		return false, fmt.Errorf("LogicEq: unknown type for '%T'", a)
	}
}

// LogicNeq returns the result of a != b
func LogicNeq(a, b any) (r bool, err error) {
	r, err = LogicEq(a, b)
	r = !r
	return
}

// LogicGt returns the result of a > b
func LogicGt(a, b any) (bool, error) {
	v, err := cast(a, b)
	if err != nil {
		return false, err
	}

	switch na := v.(type) {
	case time.Time:
		nb, err := cas.ToTime(b)
		return na.After(nb), err
	case time.Duration:
		nb, err := cas.ToDuration(b)
		return na > nb, err
	case string:
		nb, err := cas.ToString(b)
		return na > nb, err
	case int:
		nb, err := cas.ToInt(b)
		return na > nb, err
	case int8:
		nb, err := cas.ToInt8(b)
		return na > nb, err
	case int16:
		nb, err := cas.ToInt16(b)
		return na > nb, err
	case int32:
		nb, err := cas.ToInt32(b)
		return na > nb, err
	case int64:
		nb, err := cas.ToInt64(b)
		return na > nb, err
	case uint:
		nb, err := cas.ToUint(b)
		return na > nb, err
	case uint8:
		nb, err := cas.ToUint8(b)
		return na > nb, err
	case uint16:
		nb, err := cas.ToUint16(b)
		return na > nb, err
	case uint32:
		nb, err := cas.ToUint32(b)
		return na > nb, err
	case uint64:
		nb, err := cas.ToUint64(b)
		return na > nb, err
	case float32:
		nb, err := cas.ToFloat32(b)
		return na > nb, err
	case float64:
		nb, err := cas.ToFloat64(b)
		return na > nb, err
	default:
		return false, fmt.Errorf("LogicGt: unknown type for '%T'", a)
	}
}

// LogicGte returns the result of a >= b
func LogicGte(a, b any) (bool, error) {
	v, err := cast(a, b)
	if err != nil {
		return false, err
	}

	switch na := v.(type) {
	case time.Time:
		nb, err := cas.ToTime(b)
		return nb.Before(na), err
	case time.Duration:
		nb, err := cas.ToDuration(b)
		return na >= nb, err
	case string:
		nb, err := cas.ToString(b)
		return na >= nb, err
	case int:
		nb, err := cas.ToInt(b)
		return na >= nb, err
	case int8:
		nb, err := cas.ToInt8(b)
		return na >= nb, err
	case int16:
		nb, err := cas.ToInt16(b)
		return na >= nb, err
	case int32:
		nb, err := cas.ToInt32(b)
		return na >= nb, err
	case int64:
		nb, err := cas.ToInt64(b)
		return na >= nb, err
	case uint:
		nb, err := cas.ToUint(b)
		return na >= nb, err
	case uint8:
		nb, err := cas.ToUint8(b)
		return na >= nb, err
	case uint16:
		nb, err := cas.ToUint16(b)
		return na >= nb, err
	case uint32:
		nb, err := cas.ToUint32(b)
		return na >= nb, err
	case uint64:
		nb, err := cas.ToUint64(b)
		return na >= nb, err
	case float32:
		nb, err := cas.ToFloat32(b)
		return na >= nb, err
	case float64:
		nb, err := cas.ToFloat64(b)
		return na >= nb, err
	default:
		return false, fmt.Errorf("LoginGte: unknown type for '%T'", a)
	}
}

// LogicLt returns the result of a < b
func LogicLt(a, b any) (bool, error) {
	v, err := cast(a, b)
	if err != nil {
		return false, err
	}

	switch na := v.(type) {
	case time.Time:
		nb, err := cas.ToTime(b)
		return na.Before(nb), err
	case time.Duration:
		nb, err := cas.ToDuration(b)
		return na < nb, err
	case string:
		nb, err := cas.ToString(b)
		return na < nb, err
	case int:
		nb, err := cas.ToInt(b)
		return na < nb, err
	case int8:
		nb, err := cas.ToInt8(b)
		return na < nb, err
	case int16:
		nb, err := cas.ToInt16(b)
		return na < nb, err
	case int32:
		nb, err := cas.ToInt32(b)
		return na < nb, err
	case int64:
		nb, err := cas.ToInt64(b)
		return na < nb, err
	case uint:
		nb, err := cas.ToUint(b)
		return na < nb, err
	case uint8:
		nb, err := cas.ToUint8(b)
		return na < nb, err
	case uint16:
		nb, err := cas.ToUint16(b)
		return na < nb, err
	case uint32:
		nb, err := cas.ToUint32(b)
		return na < nb, err
	case uint64:
		nb, err := cas.ToUint64(b)
		return na < nb, err
	case float32:
		nb, err := cas.ToFloat32(b)
		return na < nb, err
	case float64:
		nb, err := cas.ToFloat64(b)
		return na < nb, err
	default:
		return false, fmt.Errorf("LoginLt: unknown type for '%T'", a)
	}
}

// LogicLte returns the result of a <= b
func LogicLte(a, b any) (bool, error) {
	v, err := cast(a, b)
	if err != nil {
		return false, err
	}

	switch na := v.(type) {
	case time.Time:
		nb, err := cas.ToTime(b)
		return nb.After(na), err
	case time.Duration:
		nb, err := cas.ToDuration(b)
		return na <= nb, err
	case string:
		nb, err := cas.ToString(b)
		return na <= nb, err
	case int:
		nb, err := cas.ToInt(b)
		return na <= nb, err
	case int8:
		nb, err := cas.ToInt8(b)
		return na <= nb, err
	case int16:
		nb, err := cas.ToInt16(b)
		return na <= nb, err
	case int32:
		nb, err := cas.ToInt32(b)
		return na <= nb, err
	case int64:
		nb, err := cas.ToInt64(b)
		return na <= nb, err
	case uint:
		nb, err := cas.ToUint(b)
		return na <= nb, err
	case uint8:
		nb, err := cas.ToUint8(b)
		return na <= nb, err
	case uint16:
		nb, err := cas.ToUint16(b)
		return na <= nb, err
	case uint32:
		nb, err := cas.ToUint32(b)
		return na <= nb, err
	case uint64:
		nb, err := cas.ToUint64(b)
		return na <= nb, err
	case float32:
		nb, err := cas.ToFloat32(b)
		return na <= nb, err
	case float64:
		nb, err := cas.ToFloat64(b)
		return na <= nb, err
	default:
		return false, fmt.Errorf("LoginLte: unknown type for '%T'", a)
	}
}
