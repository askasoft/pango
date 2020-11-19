package str

import (
	"bytes"
	"regexp"
	"strings"
)

// ContainsByte reports whether b is within s.
func ContainsByte(s string, b byte) bool {
	return strings.IndexByte(s, b) >= 0
}

// RemoveByte Removes all occurrences of the byte b from the source string str.
func RemoveByte(s string, b byte) string {
	if IsEmpty(s) || strings.IndexByte(s, b) < 0 {
		return s
	}

	l := len(s)
	bs := make([]byte, l)
	p := 0
	for i := 0; i < l; i++ {
		if s[i] != b {
			bs[p] = s[i]
			p++
		}
	}
	return string(bs[0:p])
}

// RemoveAny Removes all occurrences of characters from within the source string.
func RemoveAny(str string, rcs string) string {
	if IsEmpty(str) || IsEmpty(rcs) {
		return str
	}

	bb := bytes.Buffer{}
	bb.Grow(len(str))

	for _, c := range str {
		if strings.ContainsRune(rcs, c) {
			continue
		}
		bb.WriteRune(c)
	}
	return bb.String()
}

// SplitAny split string into string slice by any byte in chars
func SplitAny(s, chars string) []string {
	if len(chars) == 0 {
		return []string{s}
	}

	a := make([]string, 0, 2)
	b := -1
	for i, c := range s {
		if b < 0 {
			b = i
		}
		if strings.ContainsRune(chars, c) {
			if i > b {
				a = append(a, s[b:i])
			}
			b = -1
		}
	}

	if b >= 0 && b < len(s) {
		a = append(a, s[b:])
	}
	return a
}

// Matches checks if string matches the pattern (pattern is regular expression)
// In case of error return false
func Matches(str, pattern string) bool {
	match, _ := regexp.MatchString(pattern, str)
	return match
}
