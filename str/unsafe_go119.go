//go:build !go1.20

package str

import (
	"unsafe"
)

// UnsafeString returns a string from []byte with no memory allocations
func UnsafeString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

// UnsafeBytes returns a []byte from string with no memory allocations
func UnsafeBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
