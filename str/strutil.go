package str

import (
	"bytes"
	"regexp"
	"strings"
	"unicode/utf8"
)

// StartsWith Tests if the string s starts with the specified prefix b.
func StartsWith(s string, b string) bool {
	if IsEmpty(b) {
		return true
	}
	if IsEmpty(s) {
		return false
	}
	if len(s) < len(b) {
		return false
	}

	a := s[0:len(b)]
	return a == b
}

// EndsWith Tests if the string s ends with the specified suffix b.
func EndsWith(s string, b string) bool {
	if IsEmpty(b) {
		return true
	}
	if IsEmpty(s) {
		return false
	}
	if len(s) < len(b) {
		return false
	}

	a := s[len(s)-len(b):]
	return a == b
}

// StartsWithByte Tests if the byte slice s starts with the specified prefix b.
func StartsWithByte(s string, b byte) bool {
	if IsEmpty(s) {
		return false
	}

	a := s[0]
	return a == b
}

// EndsWithByte Tests if the byte slice bs ends with the specified suffix b.
func EndsWithByte(s string, b byte) bool {
	if IsEmpty(s) {
		return false
	}

	a := s[len(s)-1]
	return a == b
}

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

// CountAny counts the number of non-overlapping instances of any character of chars in s.
// If chars is an empty string, Count returns 1 + the number of Unicode code points in s.
func CountAny(s, chars string) int {
	// special case
	if len(chars) < 2 {
		return strings.Count(s, chars)
	}

	n := 0
	for _, c := range s {
		if strings.ContainsRune(chars, c) {
			n++
		}
	}
	return n
}

// SplitAny split string into string slice by any rune in chars
func SplitAny(s, chars string) []string {
	if len(chars) < 2 {
		return strings.Split(s, chars)
	}
	if s == "" {
		return []string{s}
	}

	n := CountAny(s, chars)
	a := make([]string, 0, n)
	b := 0
	for i, c := range s {
		if strings.ContainsRune(chars, c) {
			a = append(a, s[b:i])
			b = i + utf8.RuneLen(c)
		}
	}

	a = append(a, s[b:])
	return a
}

// SplitAnyNoEmpty split string (exclude empty string) into string slice by any rune in chars
func SplitAnyNoEmpty(s, chars string) []string {
	if len(chars) < 1 {
		return strings.Split(s, chars)
	}
	if s == "" {
		return nil
	}

	n := CountAny(s, chars)
	a := make([]string, 0, n)
	b := 0
	for i, c := range s {
		if strings.ContainsRune(chars, c) {
			if i > b {
				a = append(a, s[b:i])
			}
			b = i + utf8.RuneLen(c)
		}
	}

	if b < len(s) {
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

// TrimSpaces trim every string in the string array.
func TrimSpaces(ss []string) []string {
	for i := 0; i < len(ss); i++ {
		ss[i] = strings.TrimSpace(ss[i])
	}
	return ss
}

// RemoveEmptys remove empty string in the string array ss, and returns the new string array
func RemoveEmptys(ss []string) []string {
	ds := make([]string, 0, len(ss))
	for _, s := range ss {
		if s != "" {
			ds = append(ds, s)
		}
	}
	return ds
}
