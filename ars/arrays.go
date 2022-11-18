package ars

import (
	"bytes"

	"github.com/pandafw/pango/bye"
)

// Contains reports whether the c is contained in the slice a.
func Contains(a []any, c any) bool {
	return Index(a, c) >= 0
}

// ContainsByte reports whether the c is contained in the slice a.
func ContainsByte(a []byte, c byte) bool {
	return bye.ContainsByte(a, c)
}

// ContainsInt reports whether the c is contained in the slice a.
func ContainsInt(a []int, c int) bool {
	return IndexInt(a, c) >= 0
}

// ContainsInt8 reports whether the c is contained in the slice a.
func ContainsInt8(a []int8, c int8) bool {
	return IndexInt8(a, c) >= 0
}

// ContainsInt16 reports whether the c is contained in the slice a.
func ContainsInt16(a []int16, c int16) bool {
	return IndexInt16(a, c) >= 0
}

// ContainsInt32 reports whether the c is contained in the slice a.
func ContainsInt32(a []int32, c int32) bool {
	return IndexInt32(a, c) >= 0
}

// ContainsInt64 reports whether the c is contained in the slice a.
func ContainsInt64(a []int64, c int64) bool {
	return IndexInt64(a, c) >= 0
}

// ContainsUint reports whether the c is contained in the slice a.
func ContainsUint(a []uint, c uint) bool {
	return IndexUint(a, c) >= 0
}

// ContainsUint8 reports whether the c is contained in the slice a.
func ContainsUint8(a []uint8, c uint8) bool {
	return IndexUint8(a, c) >= 0
}

// ContainsUint16 reports whether the c is contained in the slice a.
func ContainsUint16(a []uint16, c uint16) bool {
	return IndexUint16(a, c) >= 0
}

// ContainsUint32 reports whether the c is contained in the slice a.
func ContainsUint32(a []uint32, c uint32) bool {
	return IndexUint32(a, c) >= 0
}

// ContainsUint64 reports whether the c is contained in the slice a.
func ContainsUint64(a []uint64, c uint64) bool {
	return IndexUint64(a, c) >= 0
}

// ContainsFloat32 reports whether the c is contained in the slice a.
func ContainsFloat32(a []float32, c float32) bool {
	return IndexFloat32(a, c) >= 0
}

// ContainsFloat64 reports whether the c is contained in the slice a.
func ContainsFloat64(a []float64, c float64) bool {
	return IndexFloat64(a, c) >= 0
}

// ContainsRune reports whether the c is contained in the slice a.
func ContainsRune(a []rune, c rune) bool {
	return IndexRune(a, c) >= 0
}

// ContainsString reports whether the c is contained in the slice a.
func ContainsString(a []string, c string) bool {
	return IndexString(a, c) >= 0
}

// Equal reports whether a and b
// are the same length and contain the same element.
// A nil argument is equivalent to an empty slice.
func Equal(a []any, b []any) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// EqualBytes reports whether a and b
// are the same length and contain the same element.
// A nil argument is equivalent to an empty slice.
func EqualBytes(a []byte, b []byte) bool {
	return bytes.Equal(a, b)
}

// EqualInts reports whether a and b
// are the same length and contain the same element.
// A nil argument is equivalent to an empty slice.
func EqualInts(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// EqualInt8s reports whether a and b
// are the same length and contain the same element.
// A nil argument is equivalent to an empty slice.
func EqualInt8s(a []int8, b []int8) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// EqualInt16s reports whether a and b
// are the same length and contain the same element.
// A nil argument is equivalent to an empty slice.
func EqualInt16s(a []int16, b []int16) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// EqualInt32s reports whether a and b
// are the same length and contain the same element.
// A nil argument is equivalent to an empty slice.
func EqualInt32s(a []int32, b []int32) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// EqualInt64s reports whether a and b
// are the same length and contain the same element.
// A nil argument is equivalent to an empty slice.
func EqualInt64s(a []int64, b []int64) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// EqualUints reports whether a and b
// are the same length and contain the same element.
// A nil argument is equivalent to an empty slice.
func EqualUints(a []uint, b []uint) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// EqualUint8s reports whether a and b
// are the same length and contain the same element.
// A nil argument is equivalent to an empty slice.
func EqualUint8s(a []uint8, b []uint8) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// EqualUint16s reports whether a and b
// are the same length and contain the same element.
// A nil argument is equivalent to an empty slice.
func EqualUint16s(a []uint16, b []uint16) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// EqualUint32s reports whether a and b
// are the same length and contain the same element.
// A nil argument is equivalent to an empty slice.
func EqualUint32s(a []uint32, b []uint32) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// EqualUint64s reports whether a and b
// are the same length and contain the same element.
// A nil argument is equivalent to an empty slice.
func EqualUint64s(a []uint64, b []uint64) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// EqualFloat32s reports whether a and b
// are the same length and contain the same element.
// A nil argument is equivalent to an empty slice.
func EqualFloat32s(a []float32, b []float32) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// EqualFloat64s reports whether a and b
// are the same length and contain the same element.
// A nil argument is equivalent to an empty slice.
func EqualFloat64s(a []float64, b []float64) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// EqualRunes reports whether a and b
// are the same length and contain the same element.
// A nil argument is equivalent to an empty slice.
func EqualRunes(a []rune, b []rune) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// EqualStrings reports whether a and b
// are the same length and contain the same element.
// A nil argument is equivalent to an empty slice.
func EqualStrings(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// Index returns the index of the first instance of c in a, or -1 if c is not present in a.
func Index(a []any, c any) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexByte returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexByte(a []byte, c byte) int {
	return bye.IndexByte(a, c)
}

// IndexInt returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexInt(a []int, c int) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexInt8 returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexInt8(a []int8, c int8) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexInt16 returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexInt16(a []int16, c int16) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexInt32 returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexInt32(a []int32, c int32) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexInt64 returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexInt64(a []int64, c int64) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexUint returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexUint(a []uint, c uint) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexUint8 returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexUint8(a []uint8, c uint8) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexUint16 returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexUint16(a []uint16, c uint16) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexUint32 returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexUint32(a []uint32, c uint32) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexUint64 returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexUint64(a []uint64, c uint64) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexFloat32 returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexFloat32(a []float32, c float32) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexFloat64 returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexFloat64(a []float64, c float64) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexRune returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexRune(a []rune, c rune) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}

// IndexString returns the index of the first instance of c in a, or -1 if c is not present in a.
func IndexString(a []string, c string) int {
	for i, e := range a {
		if e == c {
			return i
		}
	}
	return -1
}
