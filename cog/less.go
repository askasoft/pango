//go:build go1.18
// +build go1.18

package cog

import "github.com/askasoft/pango/str"

// Less will make type assertion (see LessString(a,b) for example),
// which will panic if a or b are not of the asserted type.
//
// Should return a bool:
//
//	true , if a < b
//	false, if a >= b
type Less[T any] func(a, b T) bool

// LessString string less function
func LessString(a, b string) bool {
	return a < b
}

// LessStringFold string case-insensitive less function
func LessStringFold(a, b string) bool {
	return str.CompareFold(a, b) < 0
}

// LessByte byte less function
func LessByte(a, b byte) bool {
	return a < b
}

// LessRune rune less function
func LessRune(a, b rune) bool {
	return a < b
}

// LessInt int less function
func LessInt(a, b int) bool {
	return a < b
}

// LessInt8 int8 less function
func LessInt8(a, b int8) bool {
	return a < b
}

// LessInt16 int16 less function
func LessInt16(a, b int16) bool {
	return a < b
}

// LessInt32 int32 less function
func LessInt32(a, b int32) bool {
	return a < b
}

// LessInt64 int64 less function
func LessInt64(a, b int64) bool {
	return a < b
}

// LessUint uint less function
func LessUint(a, b uint) bool {
	return a < b
}

// LessUint8 uint8 less function
func LessUint8(a, b uint8) bool {
	return a < b
}

// LessUint16 uint16 less function
func LessUint16(a, b uint16) bool {
	return a < b
}

// LessUint32 uint32 less function
func LessUint32(a, b uint32) bool {
	return a < b
}

// LessUint64 uint64 less function
func LessUint64(a, b uint64) bool {
	return a < b
}

// LessFloat32 float32 less function
func LessFloat32(a, b float32) bool {
	return a < b
}

// LessFloat64 float64 less function
func LessFloat64(a, b float64) bool {
	return a < b
}
