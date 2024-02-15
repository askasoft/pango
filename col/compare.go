package col

import (
	"github.com/askasoft/pango/str"
)

// Compare will make type assertion (see CompareString(a,b) for example),
// which will panic if a or b are not of the asserted type.
//
// Should return a int:
//
//	negative , if a < b
//	zero     , if a == b
//	positive , if a > b
type Compare func(a, b T) int

// CompareString provides a basic comparison on string
func CompareString(a, b T) int {
	return str.Compare(a.(string), b.(string))
}

// CompareStringFold provides a basic case-insensitive comparison on string
func CompareStringFold(a, b T) int {
	return str.CompareFold(a.(string), b.(string))
}

// CompareInt provides a basic comparison on int
func CompareInt(a, b T) int {
	x := a.(int)
	y := b.(int)

	switch {
	case x > y:
		return 1
	case x < y:
		return -1
	default:
		return 0
	}
}

// CompareInt8 provides a basic comparison on int8
func CompareInt8(a, b T) int {
	x := a.(int8)
	y := b.(int8)

	switch {
	case x > y:
		return 1
	case x < y:
		return -1
	default:
		return 0
	}
}

// CompareInt16 provides a basic comparison on int16
func CompareInt16(a, b T) int {
	x := a.(int16)
	y := b.(int16)

	switch {
	case x > y:
		return 1
	case x < y:
		return -1
	default:
		return 0
	}
}

// CompareInt32 provides a basic comparison on int32
func CompareInt32(a, b T) int {
	x := a.(int32)
	y := b.(int32)

	switch {
	case x > y:
		return 1
	case x < y:
		return -1
	default:
		return 0
	}
}

// CompareInt64 provides a basic comparison on int64
func CompareInt64(a, b T) int {
	x := a.(int64)
	y := b.(int64)

	switch {
	case x > y:
		return 1
	case x < y:
		return -1
	default:
		return 0
	}
}

// CompareUInt provides a basic comparison on uint
func CompareUInt(a, b T) int {
	x := a.(uint)
	y := b.(uint)

	switch {
	case x > y:
		return 1
	case x < y:
		return -1
	default:
		return 0
	}
}

// CompareUInt8 provides a basic comparison on uint8
func CompareUInt8(a, b T) int {
	x := a.(uint8)
	y := b.(uint8)

	switch {
	case x > y:
		return 1
	case x < y:
		return -1
	default:
		return 0
	}
}

// CompareUInt16 provides a basic comparison on uint16
func CompareUInt16(a, b T) int {
	x := a.(uint16)
	y := b.(uint16)

	switch {
	case x > y:
		return 1
	case x < y:
		return -1
	default:
		return 0
	}
}

// CompareUInt32 provides a basic comparison on uint32
func CompareUInt32(a, b T) int {
	x := a.(uint32)
	y := b.(uint32)

	switch {
	case x > y:
		return 1
	case x < y:
		return -1
	default:
		return 0
	}
}

// CompareUInt64 provides a basic comparison on uint64
func CompareUInt64(a, b T) int {
	x := a.(uint64)
	y := b.(uint64)

	switch {
	case x > y:
		return 1
	case x < y:
		return -1
	default:
		return 0
	}
}

// CompareFloat32 provides a basic comparison on float32
func CompareFloat32(a, b T) int {
	x := a.(float32)
	y := b.(float32)

	xNaN := x != x
	yNaN := y != y

	switch {
	case xNaN && yNaN:
		return 0
	case yNaN || x > y:
		return 1
	case xNaN || x < y:
		return -1
	default:
		return 0
	}
}

// CompareFloat64 provides a basic comparison on float64
func CompareFloat64(a, b T) int {
	x := a.(float64)
	y := b.(float64)

	xNaN := x != x
	yNaN := y != y

	switch {
	case xNaN && yNaN:
		return 0
	case yNaN || x > y:
		return 1
	case xNaN || x < y:
		return -1
	default:
		return 0
	}
}

// CompareByte provides a basic comparison on byte
func CompareByte(a, b T) int {
	x := a.(byte)
	y := b.(byte)

	switch {
	case x > y:
		return 1
	case x < y:
		return -1
	default:
		return 0
	}
}

// CompareRune provides a basic comparison on rune
func CompareRune(a, b T) int {
	x := a.(rune)
	y := b.(rune)

	switch {
	case x > y:
		return 1
	case x < y:
		return -1
	default:
		return 0
	}
}
