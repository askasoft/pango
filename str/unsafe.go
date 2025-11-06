package str

import (
	"unsafe"
)

// UnsafeBytes returns a []byte from string with no memory allocations
func UnsafeBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// UnsafeString returns a string from []byte with no memory allocations
func UnsafeString(bs []byte) string {
	return unsafe.String(unsafe.SliceData(bs), len(bs))
}
