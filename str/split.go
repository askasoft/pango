package str

import (
	"strings"
	"unicode/utf8"
)

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
