package cmp

import (
	"github.com/askasoft/pango/str"
)

// CompareString provides a basic comparison on string
func CompareString(a, b string) int {
	return str.Compare(a, b)
}

// CompareStringFold provides a basic case-insensitive comparison on string
func CompareStringFold(a, b string) int {
	return str.CompareFold(a, b)
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
	aNaN := a != a
	bNaN := b != b

	switch {
	case aNaN && bNaN:
		return 0
	case bNaN || a > b:
		return 1
	case aNaN || a < b:
		return -1
	default:
		return 0
	}
}

// CompareFloat64 provides a basic comparison on float64
func CompareFloat64(a, b float64) int {
	aNaN := a != a
	bNaN := b != b

	switch {
	case aNaN && bNaN:
		return 0
	case bNaN || a > b:
		return 1
	case aNaN || a < b:
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
