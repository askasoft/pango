package cmp

import (
	"github.com/pandafw/pango/str"
)

// Compare will make type assertion (see CompareString(a,b) for example),
// which will panic if a or b are not of the asserted type.
//
// Should return a int:
//    negative , if a < b
//    zero     , if a == b
//    positive , if a > b
type Compare func(a, b interface{}) int

// CompareString provides a fast comparison on strings
func CompareString(a, b interface{}) int {
	sa := a.(string)
	sb := b.(string)
	return str.Compare(sa, sb)
}

// CompareInt provides a basic comparison on int
func CompareInt(a, b interface{}) int {
	ia := a.(int)
	ib := b.(int)
	switch {
	case ia > ib:
		return 1
	case ia < ib:
		return -1
	default:
		return 0
	}
}

// CompareInt8 provides a basic comparison on int8
func CompareInt8(a, b interface{}) int {
	ia := a.(int8)
	ib := b.(int8)
	switch {
	case ia > ib:
		return 1
	case ia < ib:
		return -1
	default:
		return 0
	}
}

// CompareInt16 provides a basic comparison on int16
func CompareInt16(a, b interface{}) int {
	ia := a.(int16)
	ib := b.(int16)
	switch {
	case ia > ib:
		return 1
	case ia < ib:
		return -1
	default:
		return 0
	}
}

// CompareInt32 provides a basic comparison on int32
func CompareInt32(a, b interface{}) int {
	ia := a.(int32)
	ib := b.(int32)
	switch {
	case ia > ib:
		return 1
	case ia < ib:
		return -1
	default:
		return 0
	}
}

// CompareInt64 provides a basic comparison on int64
func CompareInt64(a, b interface{}) int {
	ia := a.(int64)
	ib := b.(int64)
	switch {
	case ia > ib:
		return 1
	case ia < ib:
		return -1
	default:
		return 0
	}
}

// CompareUInt provides a basic comparison on uint
func CompareUInt(a, b interface{}) int {
	ia := a.(uint)
	ib := b.(uint)
	switch {
	case ia > ib:
		return 1
	case ia < ib:
		return -1
	default:
		return 0
	}
}

// CompareUInt8 provides a basic comparison on uint8
func CompareUInt8(a, b interface{}) int {
	ia := a.(uint8)
	ib := b.(uint8)
	switch {
	case ia > ib:
		return 1
	case ia < ib:
		return -1
	default:
		return 0
	}
}

// CompareUInt16 provides a basic comparison on uint16
func CompareUInt16(a, b interface{}) int {
	ia := a.(uint16)
	ib := b.(uint16)
	switch {
	case ia > ib:
		return 1
	case ia < ib:
		return -1
	default:
		return 0
	}
}

// CompareUInt32 provides a basic comparison on uint32
func CompareUInt32(a, b interface{}) int {
	ia := a.(uint32)
	ib := b.(uint32)
	switch {
	case ia > ib:
		return 1
	case ia < ib:
		return -1
	default:
		return 0
	}
}

// CompareUInt64 provides a basic comparison on uint64
func CompareUInt64(a, b interface{}) int {
	ia := a.(uint64)
	ib := b.(uint64)
	switch {
	case ia > ib:
		return 1
	case ia < ib:
		return -1
	default:
		return 0
	}
}

// CompareFloat32 provides a basic comparison on float32
func CompareFloat32(a, b interface{}) int {
	ia := a.(float32)
	ib := b.(float32)
	switch {
	case ia > ib:
		return 1
	case ia < ib:
		return -1
	default:
		return 0
	}
}

// CompareFloat64 provides a basic comparison on float64
func CompareFloat64(a, b interface{}) int {
	ia := a.(float64)
	ib := b.(float64)
	switch {
	case ia > ib:
		return 1
	case ia < ib:
		return -1
	default:
		return 0
	}
}

// CompareByte provides a basic comparison on byte
func CompareByte(a, b interface{}) int {
	ia := a.(byte)
	ib := b.(byte)
	switch {
	case ia > ib:
		return 1
	case ia < ib:
		return -1
	default:
		return 0
	}
}

// CompareRune provides a basic comparison on rune
func CompareRune(a, b interface{}) int {
	ia := a.(rune)
	ib := b.(rune)
	switch {
	case ia > ib:
		return 1
	case ia < ib:
		return -1
	default:
		return 0
	}
}
