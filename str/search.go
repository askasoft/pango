package str

import (
	"strings"
	"unicode/utf8"

	"github.com/pandafw/pango/str/rabinkarp"
)

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

// ContainsFold reports whether substr is within s (case insensitive).
func ContainsFold(s, substr string) bool {
	return IndexFold(s, substr) >= 0
}

// ContainsByte reports whether b is within s.
func ContainsByte(s string, b byte) bool {
	return strings.IndexByte(s, b) >= 0
}

// IndexFold returns the index of the first instance of substr in s (case insensitive), or -1 if substr is not present in s.
func IndexFold(s, substr string) int {
	ns := len(s)
	nb := len(substr)
	if ns < nb {
		return -1
	}
	if nb == 0 {
		return 0
	}
	if ns == nb {
		if strings.EqualFold(s, substr) {
			return 0
		}
		return -1
	}

	l := ns - nb
	for i := 0; i <= l; {
		src := s[i : i+nb]
		if strings.EqualFold(src, substr) {
			return i
		}
		_, z := utf8.DecodeRuneInString(src)
		i += z
	}
	return -1
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

// LastIndexFold returns the index of the last instance of substr in s (case insensitive), or -1 if substr is not present in s.
func LastIndexFold(s, substr string) int {
	n := len(substr)
	switch {
	case n == 0:
		return len(s)
	case n == 1:
		lc := ToLower(substr)
		uc := ToUpper(substr)
		if lc == uc {
			return LastIndexByte(s, substr[0])
		}
		return LastIndexAny(s, lc+uc)
	case n == len(s):
		if substr == s {
			return 0
		}
		return -1
	case n > len(s):
		return -1
	}

	// Rabin-Karp search from the end of the string
	hashss, pow := rabinkarp.HashStrRev(substr)
	last := len(s) - n
	var h uint32
	for i := len(s) - 1; i >= last; i-- {
		h = h*rabinkarp.PrimeRK + uint32(s[i])
	}
	if h == hashss && EqualFold(s[last:], substr) {
		return last
	}

	for i := last - 1; i >= 0; i-- {
		h *= rabinkarp.PrimeRK
		h += uint32(s[i])
		h -= pow * uint32(s[i+n])
		if h == hashss && EqualFold(s[i:i+n], substr) {
			return i
		}
	}
	return -1
}

// HasPrefixFold Tests if the string s starts with the specified prefix (case insensitive).
func HasPrefixFold(s, prefix string) bool {
	return len(s) >= len(prefix) && EqualFold(s[:len(prefix)], prefix)
}

// HasSuffixFold Tests if the string s ends with the specified suffix (case insensitive).
func HasSuffixFold(s, suffix string) bool {
	return len(s) >= len(suffix) && EqualFold(s[len(s)-len(suffix):], suffix)
}

// StartsWith Tests if the string s starts with the specified prefix.
func StartsWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// EndsWith Tests if the string s ends with the specified suffix.
func EndsWith(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// StartsWithFold Tests if the string s starts with the specified prefix (case insensitive).
func StartsWithFold(s, prefix string) bool {
	return HasPrefixFold(s, prefix)
}

// EndsWithFold Tests if the string s ends with the specified suffix (case insensitive).
func EndsWithFold(s, suffix string) bool {
	return HasSuffixFold(s, suffix)
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
