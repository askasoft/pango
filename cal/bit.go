package cal

import (
	"fmt"

	"github.com/askasoft/pango/cas"
)

// BitNot returns the result of ~a
func BitNot(a any) (any, error) {
	switch na := a.(type) {
	case int:
		return ^na, nil
	case int8:
		return ^na, nil
	case int16:
		return ^na, nil
	case int32:
		return ^na, nil
	case int64:
		return ^na, nil
	case uint:
		return ^na, nil
	case uint8:
		return ^na, nil
	case uint16:
		return ^na, nil
	case uint32:
		return ^na, nil
	case uint64:
		return ^na, nil
	default:
		return a, fmt.Errorf("cal: BitNot(^) unknown type for '%T'", a)
	}
}

// BitAnds returns the result of a & b[0] & b[1] ...
func BitAnds(a any, b ...any) (r any, err error) {
	r = a
	for _, v := range b {
		r, err = BitAnd(r, v)
		if err != nil {
			return
		}
	}
	return
}

// BitAnd returns the result of a & b
func BitAnd(a, b any) (any, error) {
	switch na := a.(type) {
	case int:
		nb, err := cas.ToInt(b)
		return na & nb, err
	case int8:
		nb, err := cas.ToInt8(b)
		return na & nb, err
	case int16:
		nb, err := cas.ToInt16(b)
		return na & nb, err
	case int32:
		nb, err := cas.ToInt32(b)
		return na & nb, err
	case int64:
		nb, err := cas.ToInt64(b)
		return na & nb, err
	case uint:
		nb, err := cas.ToUint(b)
		return na & nb, err
	case uint8:
		nb, err := cas.ToUint8(b)
		return na & nb, err
	case uint16:
		nb, err := cas.ToUint16(b)
		return na & nb, err
	case uint32:
		nb, err := cas.ToUint32(b)
		return na & nb, err
	case uint64:
		nb, err := cas.ToUint64(b)
		return na & nb, err
	default:
		return a, fmt.Errorf("cal: BitAnd(&) unknown type for '%T'", a)
	}
}

// BitOrs returns the result of a | b[0] | b[1] ...
func BitOrs(a any, b ...any) (r any, err error) {
	r = a
	for _, v := range b {
		r, err = BitOr(r, v)
		if err != nil {
			return
		}
	}
	return
}

// BitOr returns the result of a | b
func BitOr(a, b any) (any, error) {
	switch na := a.(type) {
	case int:
		nb, err := cas.ToInt(b)
		return na | nb, err
	case int8:
		nb, err := cas.ToInt8(b)
		return na | nb, err
	case int16:
		nb, err := cas.ToInt16(b)
		return na | nb, err
	case int32:
		nb, err := cas.ToInt32(b)
		return na | nb, err
	case int64:
		nb, err := cas.ToInt64(b)
		return na | nb, err
	case uint:
		nb, err := cas.ToUint(b)
		return na | nb, err
	case uint8:
		nb, err := cas.ToUint8(b)
		return na | nb, err
	case uint16:
		nb, err := cas.ToUint16(b)
		return na | nb, err
	case uint32:
		nb, err := cas.ToUint32(b)
		return na | nb, err
	case uint64:
		nb, err := cas.ToUint64(b)
		return na | nb, err
	default:
		return a, fmt.Errorf("cal: BitOr(|) unknown type for '%T'", a)
	}
}

// BitXors returns the result of a ^ b[0] ^ b[1] ...
func BitXors(a any, b ...any) (r any, err error) {
	r = a
	for _, v := range b {
		r, err = BitXor(r, v)
		if err != nil {
			return
		}
	}
	return
}

// BitXor returns the result of a ^ b
func BitXor(a, b any) (any, error) {
	switch na := a.(type) {
	case int:
		nb, err := cas.ToInt(b)
		return na ^ nb, err
	case int8:
		nb, err := cas.ToInt8(b)
		return na ^ nb, err
	case int16:
		nb, err := cas.ToInt16(b)
		return na ^ nb, err
	case int32:
		nb, err := cas.ToInt32(b)
		return na ^ nb, err
	case int64:
		nb, err := cas.ToInt64(b)
		return na ^ nb, err
	case uint:
		nb, err := cas.ToUint(b)
		return na ^ nb, err
	case uint8:
		nb, err := cas.ToUint8(b)
		return na ^ nb, err
	case uint16:
		nb, err := cas.ToUint16(b)
		return na ^ nb, err
	case uint32:
		nb, err := cas.ToUint32(b)
		return na ^ nb, err
	case uint64:
		nb, err := cas.ToUint64(b)
		return na ^ nb, err
	default:
		return a, fmt.Errorf("cal: BitXor(^) unknown type for '%T'", a)
	}
}

// BitLeft returns the result of a << b
func BitLeft(a, b any) (any, error) {
	switch na := a.(type) {
	case int:
		nb, err := cas.ToInt(b)
		return na << nb, err
	case int8:
		nb, err := cas.ToInt8(b)
		return na << nb, err
	case int16:
		nb, err := cas.ToInt16(b)
		return na << nb, err
	case int32:
		nb, err := cas.ToInt32(b)
		return na << nb, err
	case int64:
		nb, err := cas.ToInt64(b)
		return na << nb, err
	case uint:
		nb, err := cas.ToUint(b)
		return na << nb, err
	case uint8:
		nb, err := cas.ToUint8(b)
		return na << nb, err
	case uint16:
		nb, err := cas.ToUint16(b)
		return na << nb, err
	case uint32:
		nb, err := cas.ToUint32(b)
		return na << nb, err
	case uint64:
		nb, err := cas.ToUint64(b)
		return na << nb, err
	default:
		return a, fmt.Errorf("cal: BitLeft(<<) unknown type for '%T'", a)
	}
}

// BitRight returns the result of a >> b
func BitRight(a, b any) (any, error) {
	switch na := a.(type) {
	case int:
		nb, err := cas.ToInt(b)
		return na >> nb, err
	case int8:
		nb, err := cas.ToInt8(b)
		return na >> nb, err
	case int16:
		nb, err := cas.ToInt16(b)
		return na >> nb, err
	case int32:
		nb, err := cas.ToInt32(b)
		return na >> nb, err
	case int64:
		nb, err := cas.ToInt64(b)
		return na >> nb, err
	case uint:
		nb, err := cas.ToUint(b)
		return na >> nb, err
	case uint8:
		nb, err := cas.ToUint8(b)
		return na >> nb, err
	case uint16:
		nb, err := cas.ToUint16(b)
		return na >> nb, err
	case uint32:
		nb, err := cas.ToUint32(b)
		return na >> nb, err
	case uint64:
		nb, err := cas.ToUint64(b)
		return na >> nb, err
	default:
		return a, fmt.Errorf("cal: BitRight(>>) unknown type for '%T'", a)
	}
}
