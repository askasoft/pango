package strutil

import (
	"regexp"
	"strings"
)

// ContainsByte reports whether b is within s.
func ContainsByte(s string, b byte) bool {
	return strings.IndexByte(s, b) >= 0
}

// ContainsAnyByte reports whether any byte in sbs is within str.
func ContainsAnyByte(str string, sbs string) bool {
	if IsEmpty(sbs) {
		return true
	}

	if IsEmpty(str) {
		return false
	}

	l := len(str)
	for i := 0; i < l; i++ {
		b := str[i]
		if ContainsByte(sbs, b) {
			return true
		}
	}

	return false
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

// RemoveAnyByte Removes all occurrences of characters from within the source string.
func RemoveAnyByte(str string, rbs string) string {
	if IsEmpty(str) || IsEmpty(rbs) {
		return str
	}

	l := len(str)
	bs := make([]byte, l)

	p := 0
	for i := 0; i < l; i++ {
		b := str[i]
		if !ContainsByte(rbs, b) {
			bs[p] = b
			p++
		}
	}
	return string(bs[0:p])
}

// SplitAnyByte split string into string slice by any byte in chars
func SplitAnyByte(s, chars string) []string {
	if len(chars) == 0 {
		a := [1]string{s}
		return a[:]
	}

	a := make([]string, 0, 2)
	l := len(s)
	b := 0
	for i := 0; i < l; i++ {
		c := s[i]
		if strings.IndexByte(chars, c) >= 0 {
			if i > b {
				a = append(a, s[b:i])
			}
			b = i + 1
		}
	}
	if b < l {
		a = append(a, s[b:l])
	}
	return a
}

// Matches checks if string matches the pattern (pattern is regular expression)
// In case of error return false
func Matches(str, pattern string) bool {
	match, _ := regexp.MatchString(pattern, str)
	return match
}
