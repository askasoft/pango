package str

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"
)

// If return (b ? t : f)
func If(b bool, t, f string) string {
	if b {
		return t
	}
	return f
}

// IfEmpty return (a == "" ? b : a)
func IfEmpty(a, b string) string {
	if a == "" {
		return b
	}
	return a
}

// Compare returns an integer comparing two strings lexicographically.
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
func Compare(a, b string) int {
	return bytes.Compare(UnsafeBytes(a), UnsafeBytes(b))
}

// CompareFold returns an integer comparing two strings case-insensitive.
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
func CompareFold(s, t string) int {
	// ASCII fast path
	i := 0
	for ; i < len(s) && i < len(t); i++ {
		sr := s[i]
		tr := t[i]
		if sr|tr >= utf8.RuneSelf {
			goto hasUnicode
		}

		if tr == sr {
			continue
		}

		// ASCII only, sr/tr must be upper/lower case
		if 'A' <= sr && sr <= 'Z' {
			sr += ('a' - 'A')
		}
		if 'A' <= tr && tr <= 'Z' {
			tr += ('a' - 'A')
		}

		switch {
		case sr < tr:
			return -1
		case sr > tr:
			return 1
		}
	}

	// Check if we've exhausted both strings.
	{
		r := len(s) - len(t)
		switch {
		case r < 0:
			return -1
		case r > 0:
			return 1
		default:
			return 0
		}
	}

hasUnicode:
	s = s[i:]
	t = t[i:]
	for _, sr := range s {
		// If t is exhausted the strings are not equal.
		if len(t) == 0 {
			return 1
		}

		// Extract first rune from second string.
		var tr rune
		if t[0] < utf8.RuneSelf {
			tr, t = rune(t[0]), t[1:]
		} else {
			r, size := utf8.DecodeRuneInString(t)
			tr, t = r, t[size:]
		}

		// If they match, keep going;
		if tr == sr {
			continue
		}

		// Fast check for ASCII.
		if sr < utf8.RuneSelf && tr < utf8.RuneSelf {
			// ASCII only, sr/tr must be upper/lower case
			if 'A' <= sr && sr <= 'Z' {
				sr += ('a' - 'A')
			}
			if 'A' <= tr && tr <= 'Z' {
				tr += ('a' - 'A')
			}

			switch {
			case sr < tr:
				return -1
			case sr > tr:
				return 1
			default:
				continue
			}
		}

		sr = unicode.ToLower(sr)
		tr = unicode.ToLower(tr)
		switch {
		case sr < tr:
			return -1
		case sr > tr:
			return 1
		}
	}

	// First string is empty, so check if the second one is also empty.
	if len(t) == 0 {
		return 0
	}

	return -1
}

// Capitalize returns a copy of the string s that the start letter
// mapped to their Unicode upper case.
func Capitalize(s string) string {
	if s == "" {
		return s
	}

	if s[0] < utf8.RuneSelf {
		if unicode.IsUpper(rune(s[0])) {
			return s
		}
		return ToUpper(s[:1]) + s[1:]
	}

	r, size := utf8.DecodeRuneInString(s)
	if unicode.IsUpper(r) {
		return s
	}
	return string(unicode.ToUpper(r)) + s[size:]
}

// CamelCase returns a copy of the string s with camel case.
func CamelCase(s string) string {
	var sb strings.Builder
	var uc bool
	for _, r := range s {
		if r == '_' || r == '-' {
			uc = true
			continue
		}

		if uc {
			sb.WriteRune(unicode.ToUpper(r))
			uc = false
		} else {
			sb.WriteRune(unicode.ToLower(r))
		}
	}
	return sb.String()
}

// PascalCase returns a copy of the string s with pascal case.
func PascalCase(s string) string {
	var sb strings.Builder
	var uc bool
	for _, r := range s {
		if r == '_' || r == '-' {
			uc = true
			continue
		}

		if uc {
			sb.WriteRune(unicode.ToUpper(r))
			uc = false
		} else {
			if sb.Len() == 0 {
				sb.WriteRune(unicode.ToUpper(r))
			} else {
				sb.WriteRune(unicode.ToLower(r))
			}
		}
	}
	return sb.String()
}

// SnakeCase returns a copy of the string s with snake case c (default _).
func SnakeCase(s string, c ...rune) string {
	d := '_'
	if len(c) > 0 {
		d = c[0]
	}

	var sb strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				sb.WriteRune(d)
			}
			sb.WriteRune(unicode.ToLower(r))
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// Strip returns a slice of the string s, with all leading
// and trailing white space removed, as defined by Unicode.
func Strip(s string) string {
	return strings.TrimSpace(s)
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

// StripLeft returns a slice of the string s, with all leading
// white space removed, as defined by Unicode.
func StripLeft(s string) string {
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

// StripRight returns a slice of the string s, with all
// trailing white space removed, as defined by Unicode.
func StripRight(s string) string {
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

// RuneLen returns the number of runes in s.
func RuneLen(r rune) int {
	return utf8.RuneLen(r)
}

// RuneCount returns the number of runes in s.
func RuneCount(s string) int {
	return utf8.RuneCountInString(s)
}

// RuneEqualFold reports whether sr and tr,
// are equal under Unicode case-folding, which is a more general
// form of case-insensitivity.
func RuneEqualFold(sr, tr rune) bool {
	// Easy case.
	if tr == sr {
		return true
	}

	// Make sr < tr to simplify what follows.
	if tr < sr {
		tr, sr = sr, tr
	}

	// Fast check for ASCII.
	if tr < utf8.RuneSelf {
		// ASCII only, sr/tr must be upper/lower case
		return 'A' <= sr && sr <= 'Z' && tr == sr+'a'-'A'
	}

	// General case. SimpleFold(x) returns the next equivalent rune > x
	// or wraps around to smaller values.
	r := unicode.SimpleFold(sr)
	for r != sr && r < tr {
		r = unicode.SimpleFold(r)
	}
	return r == tr
}

// RemoveByte Removes all occurrences of the byte b from the source string s.
func RemoveByte(s string, b byte) string {
	if s == "" {
		return s
	}

	sb := Builder{}
	for {
		i := strings.IndexByte(s, b)
		if i < 0 {
			if sb.Len() == 0 {
				return s
			}
			sb.WriteString(s)
			return sb.String()
		}

		sb.WriteString(s[:i])
		s = s[i+1:]
	}
}

// RemoveRune Removes all occurrences of the rune r from the source string s.
func RemoveRune(s string, r rune) string {
	if s == "" {
		return s
	}

	n := utf8.RuneLen(r)

	sb := Builder{}
	for {
		i := strings.IndexRune(s, r)
		if i < 0 {
			if sb.Len() == 0 {
				return s
			}
			sb.WriteString(s)
			return sb.String()
		}

		sb.WriteString(s[:i])
		s = s[i+n:]
	}
}

// RemoveAny Removes all occurrences of characters from within the source string s.
func RemoveAny(s string, r string) string {
	if s == "" || r == "" {
		return s
	}

	sb := Builder{}
	for {
		i := strings.IndexAny(s, r)
		if i < 0 {
			if sb.Len() == 0 {
				return s
			}
			sb.WriteString(s)
			return sb.String()
		}

		sb.WriteString(s[:i])

		_, n := utf8.DecodeRuneInString(s[i:])
		s = s[i+n:]
	}
}

// RemoveFunc Removes all occurrences of characters from within the source string which satisfy f(c).
func RemoveFunc(s string, f func(r rune) bool) string {
	if s == "" {
		return s
	}

	sb := Builder{}
	for {
		i := strings.IndexFunc(s, f)
		if i < 0 {
			if sb.Len() == 0 {
				return s
			}
			sb.WriteString(s)
			return sb.String()
		}

		sb.WriteString(s[:i])

		_, n := utf8.DecodeRuneInString(s[i:])
		s = s[i+n:]
	}
}

// Remove Removes all substring r from within the source string s.
func Remove(s string, r string) string {
	if s == "" || r == "" {
		return s
	}

	n := len(r)
	sb := Builder{}
	for {
		i := strings.Index(s, r)
		if i < 0 {
			if sb.Len() == 0 {
				return s
			}
			sb.WriteString(s)
			return sb.String()
		}

		sb.WriteString(s[:i])
		s = s[i+n:]
	}
}
