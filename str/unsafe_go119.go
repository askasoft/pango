//go:build !go1.20
// +build !go1.20

package str

import (
	"unsafe"
)

// UnsafeString returns a string from []byte with no memory allocations
func UnsafeString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}
