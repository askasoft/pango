package cal

import (
	"fmt"
	"reflect"
	"time"

	"github.com/askasoft/pango/cas"
	"github.com/askasoft/pango/ref"
)

var priorities = map[reflect.Type]int{
	ref.TypeString:   11,
	ref.TypeDuration: 21,
	ref.TypeFloat64:  31,
	ref.TypeFloat32:  32,
	ref.TypeInt64:    41,
	ref.TypeUint64:   42,
	ref.TypeInt32:    43,
	ref.TypeUint32:   44,
	ref.TypeInt16:    45,
	ref.TypeUint16:   46,
	ref.TypeInt8:     47,
	ref.TypeUint8:    48,
}

func init() {
	intSize := 32 << (^uint(0) >> 63) // 32 or 64

	if intSize == 64 {
		priorities[ref.TypeInt] = priorities[ref.TypeInt64]
		priorities[ref.TypeUint] = priorities[ref.TypeUint64]
	} else {
		priorities[ref.TypeInt] = priorities[ref.TypeInt32]
		priorities[ref.TypeUint] = priorities[ref.TypeUint32]
	}
}

func cast(a, b any) (any, error) {
	at, bt := reflect.TypeOf(a), reflect.TypeOf(b)
	ap, bp := priorities[at], priorities[bt]

	if ap == 0 || bp == 0 {
		return a, nil
	}

	if ap > bp {
		return ref.CastTo(a, bt)
	}

	return a, nil
}

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
	v, err := cast(a, b)
	if err != nil {
		return a, err
	}

	switch na := v.(type) {
	case string:
		nb, err := cas.ToString(b)
		return na + nb, err
	case time.Duration:
		nb, err := cas.ToDuration(b)
		return na + nb, err
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
		return a, fmt.Errorf("add: unsupported type for '%T'", v)
	}
}

// Subs returns the result of a - b[0] - b[1] ...
func Subs(a any, b ...any) (r any, err error) {
	r = a
	for _, v := range b {
		r, err = Sub(r, v)
		if err != nil {
			return
		}
	}
	return
}

// Sub subtract returns the result of a - b
func Sub(a, b any) (any, error) {
	v, err := cast(a, b)
	if err != nil {
		return a, err
	}

	switch na := v.(type) {
	case time.Duration:
		nb, err := cas.ToDuration(b)
		return na - nb, err
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
		return a, fmt.Errorf("subtract: unsupported type for '%T'", v)
	}
}

// Negate returns the result of -a
func Negate(a any) (any, error) {
	switch na := a.(type) {
	case int:
		return -na, nil
	case int8:
		return -na, nil
	case int16:
		return -na, nil
	case int32:
		return -na, nil
	case int64:
		return -na, nil
	case uint:
		return -int(na), nil
	case uint8:
		return -int8(na), nil
	case uint16:
		return -int16(na), nil
	case uint32:
		return -int32(na), nil
	case uint64:
		return -int64(na), nil
	case float32:
		return -na, nil
	case float64:
		return -na, nil
	default:
		return a, fmt.Errorf("negate: unsupported type for '%T'", a)
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
	v, err := cast(a, b)
	if err != nil {
		return a, err
	}

	switch na := v.(type) {
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
		return a, fmt.Errorf("multiply: unsupported type for '%T'", v)
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
	v, err := cast(a, b)
	if err != nil {
		return a, err
	}

	switch na := v.(type) {
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
		return a, fmt.Errorf("divide: unsupported type for '%T'", v)
	}
}

// Mods returns the result of a % b[0] % b[1] ...
func Mods(a any, b ...any) (r any, err error) {
	r = a
	for _, v := range b {
		r, err = Mod(r, v)
		if err != nil {
			return
		}
	}
	return
}

// Mod returns the result of a % b
func Mod(a, b any) (any, error) {
	v, err := cast(a, b)
	if err != nil {
		return a, err
	}

	switch na := v.(type) {
	case int:
		nb, err := cas.ToInt(b)
		return na % nb, err
	case int8:
		nb, err := cas.ToInt8(b)
		return na % nb, err
	case int16:
		nb, err := cas.ToInt16(b)
		return na % nb, err
	case int32:
		nb, err := cas.ToInt32(b)
		return na % nb, err
	case int64:
		nb, err := cas.ToInt64(b)
		return na % nb, err
	case uint:
		nb, err := cas.ToUint(b)
		return na % nb, err
	case uint8:
		nb, err := cas.ToUint8(b)
		return na % nb, err
	case uint16:
		nb, err := cas.ToUint16(b)
		return na % nb, err
	case uint32:
		nb, err := cas.ToUint32(b)
		return na % nb, err
	case uint64:
		nb, err := cas.ToUint64(b)
		return na % nb, err
	case float32:
		nb, err := cas.ToInt64(b)
		return int64(na) % nb, err
	case float64:
		nb, err := cas.ToInt64(b)
		return int64(na) % nb, err
	default:
		return a, fmt.Errorf("mod: unsupported type for '%T'", v)
	}
}
