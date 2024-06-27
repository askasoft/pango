package str

import (
	"strings"
	"unicode/utf8"
)

// SplitLength splits the string s by the length 'n' to a string slice.
// Each string 's' in the result slice satisfying len(s) <= n.
func SplitLength(s string, n int) []string {
	z := len(s)
	if z <= n || n < 1 {
		return []string{s}
	}

	a := make([]string, 0, z/n+1)

	x := 0
	b := 0
	for i, r := range s {
		d := utf8.RuneLen(r)
		b += d
		if b > n && i > x {
			a = append(a, s[x:i])
			x = i
			if z-x <= n {
				break
			}
			b = d
		}
	}

	if x < z {
		a = append(a, s[x:])
	}
	return a
}

// SplitCount splits the string s by the length 'n' to a string slice.
// Each string 's' in the result slice satisfying utf8.RuneCountInString(s) <= n.
func SplitCount(s string, n int) []string {
	z := RuneCount(s)
	if z <= n || n < 1 {
		return []string{s}
	}

	a := make([]string, 0, z/n+1)

	x := 0
	b := 0
	for i := range s {
		b++
		if b > n {
			a = append(a, s[x:i])
			x = i
			b = 1
		}
	}

	if x < len(s) {
		a = append(a, s[x:])
	}
	return a
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

// FieldsByte split string (exclude empty string) into string slice by byte c.
func FieldsByte(s string, c byte) []string {
	if s == "" {
		return []string{}
	}

	n := CountByte(s, c)
	a := make([]string, 0, n)

	b := 0
	for i := 0; i < len(s); i++ {
		if c == s[i] {
			if i > b {
				a = append(a, s[b:i])
			}
			b = i + 1
		}
	}

	if b < len(s) {
		a = append(a, s[b:])
	}
	return a
}

// FieldsRune split string (exclude empty string) into string slice by rune r.
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

// FieldsFuncIter splits the string s at each rune of Unicode code points c satisfying f(c)
// and call iter(s). stop split if iter(s) returns false.
func FieldsFuncIter(s string, f func(rune) bool, iter func(string) bool) {
	if s == "" {
		return
	}

	b := 0
	for i, c := range s {
		if f(c) {
			if i > b {
				if !iter(s[b:i]) {
					return
				}
			}
			b = i + utf8.RuneLen(c)
		}
	}

	if b < len(s) {
		iter(s[b:])
	}
}
