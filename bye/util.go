package bye

import (
	"bytes"
	"unicode"
	"unicode/utf8"
)

// IsEmpty checks if the byte slice is null.
func IsEmpty(s []byte) bool {
	return len(s) == 0
}

// ContainsByte reports whether the byte is contained in the slice s.
func ContainsByte(s []byte, b byte) bool {
	return bytes.IndexByte(s, b) >= 0
}

// StartsWith Tests if the byte slice s starts with the specified prefix b.
func StartsWith(s []byte, b []byte) bool {
	return bytes.HasPrefix(s, b)
}

// EndsWith Tests if the byte slice s ends with the specified suffix b.
func EndsWith(s []byte, b []byte) bool {
	return bytes.HasSuffix(s, b)
}

// StartsWithByte Tests if the byte slice s starts with the specified prefix b.
func StartsWithByte(s []byte, b byte) bool {
	if len(s) == 0 {
		return false
	}

	return s[0] == b
}

// EndsWithByte Tests if the byte slice s ends with the specified suffix b.
func EndsWithByte(s []byte, b byte) bool {
	if len(s) == 0 {
		return false
	}

	return s[len(s)-1] == b
}

// CountByte counts the number of b in s.
func CountByte(s []byte, b byte) int {
	n := 0
	for {
		i := IndexByte(s, b)
		if i < 0 {
			return n
		}
		n++
		s = s[i+1:]
	}
}

// Strip returns a slice of the bytes s, with all leading
// and trailing white space removed, as defined by Unicode.
func Strip(s []byte) []byte {
	return bytes.TrimSpace(s)
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

// StripLeft returns a slice of the bytes s, with all leading
// white space removed, as defined by Unicode.
func StripLeft(s []byte) []byte {
	// Fast path for ASCII: look for the first ASCII non-space byte
	start := 0
	for ; start < len(s); start++ {
		c := s[start]
		if c >= utf8.RuneSelf {
			// If we run into a non-ASCII byte, fall back to the
			// slower unicode-aware method on the remaining bytes
			return TrimLeftFunc(s[start:], unicode.IsSpace)
		}
		if asciiSpace[c] == 0 {
			break
		}
	}

	// At this point s[start:] starts with an ASCII
	// non-space bytes, so we're done. Non-ASCII cases have already
	// been handled above.
	return s[start:]
}

// StripRight returns a slice of the bytes s, with all
// trailing white space removed, as defined by Unicode.
func StripRight(s []byte) []byte {
	// Now look for the first ASCII non-space byte from the end
	stop := len(s)
	for ; stop > 0; stop-- {
		c := s[stop-1]
		if c >= utf8.RuneSelf {
			return TrimRightFunc(s[:stop], unicode.IsSpace)
		}
		if asciiSpace[c] == 0 {
			break
		}
	}

	// At this point s[:stop] ends with an ASCII
	// non-space bytes, so we're done. Non-ASCII cases have already
	// been handled above.
	return s[:stop]
}
