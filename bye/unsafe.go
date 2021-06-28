package bye

import (
	"unsafe"
)

// UnsafeString returns a string from []byte with no memory allocations
func UnsafeString(p []byte) string {
	return *(*string)(unsafe.Pointer(&p))
}
