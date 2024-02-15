package col

import "github.com/askasoft/pango/str"

// Less will make type assertion (see LessString(a,b) for example),
// which will panic if a or b are not of the asserted type.
//
// Should return a bool:
//
//	true , if a < b
//	false, if a >= b
type Less func(a, b T) bool

// LessString string less function
func LessString(a, b T) bool {
	return a.(string) < b.(string)
}

// LessStringFold string case-insensitive less function
func LessStringFold(a, b T) bool {
	return str.CompareFold(a.(string), b.(string)) < 0
}

// LessByte byte less function
func LessByte(a, b T) bool {
	return a.(byte) < b.(byte)
}

// LessRune rune less function
func LessRune(a, b T) bool {
	return a.(rune) < b.(rune)
}

// LessInt int less function
func LessInt(a, b T) bool {
	return a.(int) < b.(int)
}

// LessInt8 int8 less function
func LessInt8(a, b T) bool {
	return a.(int8) < b.(int8)
}

// LessInt16 int16 less function
func LessInt16(a, b T) bool {
	return a.(int16) < b.(int16)
}

// LessInt32 int32 less function
func LessInt32(a, b T) bool {
	return a.(int32) < b.(int32)
}

// LessInt64 int64 less function
func LessInt64(a, b T) bool {
	return a.(int64) < b.(int64)
}

// LessUint uint less function
func LessUint(a, b T) bool {
	return a.(uint) < b.(uint)
}

// LessUint8 uint8 less function
func LessUint8(a, b T) bool {
	return a.(uint8) < b.(uint8)
}

// LessUint16 uint16 less function
func LessUint16(a, b T) bool {
	return a.(uint16) < b.(uint16)
}

// LessUint32 uint32 less function
func LessUint32(a, b T) bool {
	return a.(uint32) < b.(uint32)
}

// LessUint64 uint64 less function
func LessUint64(a, b T) bool {
	return a.(uint64) < b.(uint64)
}

// LessFloat32 float32 less function
func LessFloat32(a, b T) bool {
	x := a.(float32)
	y := b.(float32)
	return (x != x && y == y) || x < y
}

// LessFloat64 float64 less function
func LessFloat64(a, b T) bool {
	x := a.(float64)
	y := b.(float64)
	return (x != x && y == y) || x < y
}
