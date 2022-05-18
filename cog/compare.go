package cog

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
type Compare[T any] func(a, b T) int

// CompareString provides a fast comparison on strings
func CompareString(a, b string) int {
	return str.Compare(a, b)
}

// CompareInt provides a basic comparison on int
func CompareInt(a, b int) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// CompareInt8 provides a basic comparison on int8
func CompareInt8(a, b int8) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// CompareInt16 provides a basic comparison on int16
func CompareInt16(a, b int16) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// CompareInt32 provides a basic comparison on int32
func CompareInt32(a, b int32) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// CompareInt64 provides a basic comparison on int64
func CompareInt64(a, b int64) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// CompareUInt provides a basic comparison on uint
func CompareUInt(a, b uint) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// CompareUInt8 provides a basic comparison on uint8
func CompareUInt8(a, b uint8) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// CompareUInt16 provides a basic comparison on uint16
func CompareUInt16(a, b uint16) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// CompareUInt32 provides a basic comparison on uint32
func CompareUInt32(a, b uint32) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// CompareUInt64 provides a basic comparison on uint64
func CompareUInt64(a, b uint64) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// CompareFloat32 provides a basic comparison on float32
func CompareFloat32(a, b float32) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// CompareFloat64 provides a basic comparison on float64
func CompareFloat64(a, b float64) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// CompareByte provides a basic comparison on byte
func CompareByte(a, b byte) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// CompareRune provides a basic comparison on rune
func CompareRune(a, b rune) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}
