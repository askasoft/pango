package str

import (
	"bytes"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
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

// RuneCount returns the number of runes in s.
func RuneCount(s string) int {
	return utf8.RuneCountInString(s)
}

// CountRune counts the number of non-overlapping instances of rune c in s.
func CountRune(s string, c rune) int {
	n := 0
	for _, r := range s {
		if r == c {
			n++
		}
	}
	return n
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

// ContainsByte reports whether b is within s.
func ContainsByte(s string, b byte) bool {
	return strings.IndexByte(s, b) >= 0
}

// LastIndexRune returns the index of the last instance of the Unicode code point
// r, or -1 if rune is not present in s.
// If r is utf8.RuneError, it returns the last instance of any
// invalid UTF-8 byte sequence.
func LastIndexRune(s string, r rune) int {
	switch {
	case 0 <= r && r < utf8.RuneSelf:
		return strings.LastIndexByte(s, byte(r))
	case r == utf8.RuneError:
		n := -1
		for i, r := range s {
			if r == utf8.RuneError {
				n = i
			}
		}
		return n
	case !utf8.ValidRune(r):
		return -1
	default:
		return strings.LastIndex(s, string(r))
	}
}

// SplitFunc splits the string s at each run of Unicode code points c satisfying f(c)
// and returns an array of slices of s.
// If s does not satisfying f(c), Split returns a
// slice of length 1 whose only element is s.
func SplitFunc(s string, f func(rune) bool) []string {
	if s == "" {
		return []string{s}
	}

	a := make([]string, 0, 32)
	b := 0
	for i, c := range s {
		if f(c) {
			a = append(a, s[b:i])
			b = i + utf8.RuneLen(c)
		}
	}

	a = append(a, s[b:])
	return a
}

// SplitAny split string into string slice by any rune in chars
func SplitAny(s, chars string) []string {
	if s == "" {
		return []string{s}
	}

	if len(chars) < 2 {
		return strings.Split(s, chars)
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

// FieldsRune split string (exclude empty string) into string slice by rune c
func FieldsRune(s string, c rune) []string {
	if s == "" {
		return []string{}
	}

	n := CountRune(s, c)
	a := make([]string, 0, n)
	b := 0
	for i, r := range s {
		if r == c {
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

// FieldsAny split string (exclude empty string) into string slice by any rune in chars
func FieldsAny(s, chars string) []string {
	if s == "" {
		return []string{}
	}

	if len(chars) < 1 {
		return strings.Split(s, chars)
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

// StartsWith Tests if the string s starts with the specified prefix b.
func StartsWith(s string, b string) bool {
	return strings.HasPrefix(s, b)
}

// EndsWith Tests if the string s ends with the specified suffix b.
func EndsWith(s string, b string) bool {
	return strings.HasSuffix(s, b)
}

// StartsWithByte Tests if the byte slice s starts with the specified prefix b.
func StartsWithByte(s string, b byte) bool {
	if s == "" {
		return false
	}

	a := s[0]
	return a == b
}

// EndsWithByte Tests if the byte slice bs ends with the specified suffix b.
func EndsWithByte(s string, b byte) bool {
	if s == "" {
		return false
	}

	a := s[len(s)-1]
	return a == b
}

// RemoveByte Removes all occurrences of the byte b from the source string str.
func RemoveByte(s string, b byte) string {
	if s == "" || strings.IndexByte(s, b) < 0 {
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
	if str == "" || rcs == "" {
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

// RemoveAnyByte Removes all occurrences of bytes from within the source string.
func RemoveAnyByte(s string, rbs string) string {
	if s == "" {
		return s
	}

	l := len(s)
	bs := make([]byte, l)
	p := 0
	for i := 0; i < l; i++ {
		if !ContainsByte(rbs, s[i]) {
			bs[p] = s[i]
			p++
		}
	}
	return string(bs[0:p])
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
