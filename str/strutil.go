package str

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"
)

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

// Compare returns an integer comparing two strings lexicographically.
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
func Compare(a, b string) int {
	return bytes.Compare(UnsafeBytes(a), UnsafeBytes(b))
}
