package bye

import "bytes"

// IsEmpty checks if the byte slice is null.
func IsEmpty(bs []byte) bool {
	return len(bs) == 0
}

// StartsWith Tests if the byte slice s starts with the specified prefix b.
func StartsWith(s []byte, b []byte) bool {
	return bytes.HasPrefix(s, b)
}

// EndsWith Tests if the byte slice bs ends with the specified suffix b.
func EndsWith(s []byte, b []byte) bool {
	return bytes.HasSuffix(s, b)
}

// StartsWithByte Tests if the byte slice s starts with the specified prefix b.
func StartsWithByte(s []byte, b byte) bool {
	if len(s) == 0 {
		return false
	}

	a := s[0]
	return a == b
}

// EndsWithByte Tests if the byte slice bs ends with the specified suffix b.
func EndsWithByte(s []byte, b byte) bool {
	if len(s) == 0 {
		return false
	}

	a := s[len(s)-1]
	return a == b
}
