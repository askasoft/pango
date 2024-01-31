package str

import (
	"strings"
	"unicode/utf8"
)

// SplitLength splits the string s by the length 'n' to a string slice.
// Each string 's' in the result slice satisfying len(s) < n.
func SplitLength(s string, n int) []string {
	if len(s) <= n || n < 1 {
		return []string{s}
	}

	a := make([]string, 0, len(s)/n+1)

	b := 0
	for {
		_, z := utf8.DecodeRuneInString(s[b:])
		if b+z <= n {
			b += z
			continue
		}

		x := b
		if b == 0 {
			x = z
		}

		a = append(a, s[:x])
		s = s[x:]
		if s == "" {
			return a
		}
		if len(s) <= n {
			a = append(a, s)
			return a
		}
		if b != 0 {
			b = z
		}
	}
}

// SplitFunc splits the string s at each rune of Unicode code points c satisfying f(c)
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
func FieldsRune(s string, r rune) []string {
	if s == "" {
		return []string{}
	}

	n := CountRune(s, r)
	a := make([]string, 0, n)

	b := 0
	z := utf8.RuneLen(r)
	for i, c := range s {
		if r == c {
			if i > b {
				a = append(a, s[b:i])
			}
			b = i + z
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
