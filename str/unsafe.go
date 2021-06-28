package str

import (
	"unsafe"
)

// UnsafeBytes returns a []byte from string with no memory allocations
func UnsafeBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
