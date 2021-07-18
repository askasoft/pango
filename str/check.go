package str

import (
	"unicode"
)

// IsASCII checks if the string contains ASCII chars only.
func IsASCII(s string) bool {
	if s == "" {
		return false
	}

	l := len(s)
	for i := 0; i < l; i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// IsPrintableASCII checks if the string contains printable ASCII chars only.
func IsPrintableASCII(s string) bool {
	if s == "" {
		return false
	}

	l := len(s)
	for i := 0; i < l; i++ {
		b := s[i]
		if b < ' ' || b > unicode.MaxASCII {
			return false
		}
	}
	return true
}
