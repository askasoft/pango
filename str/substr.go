package str

import (
	"fmt"
	"unicode/utf8"
)

// BOM "\uFEFF"
const BOM = "\uFEFF"

// SkipBOM Returns a string without BOM.
// internal call TrimPrefix(s, "\uFEFF")
func SkipBOM(s string) string {
	return TrimPrefix(s, BOM)
}

// SubstrAfter Gets the substring after the first occurrence of a separator b.
// The separator b is not returned.
// If nothing is found, the empty string is returned.
// SubstrAfter("", *)        = ""
// SubstrAfter("abc", "a")   = "bc"
// SubstrAfter("abcba", "b") = "cba"
// SubstrAfter("abc", "c")   = ""
// SubstrAfter("abc", "d")   = ""
func SubstrAfter(s string, b string) string {
	if s == "" {
		return s
	}

	i := Index(s, b)
	if i < 0 {
		return ""
	}
	return s[i+len(b):]
}

// SubstrAfterAny Gets the substring after the first occurrence of any char in separator b.
// The separator b is not returned.
// If nothing is found, the empty string is returned.
// SubstrAfterAny("", *)        = ""
// SubstrAfterAny("abc", "az")   = "bc"
// SubstrAfterAny("abcba", "b") = "cba"
// SubstrAfterAny("abc", "c")   = ""
// SubstrAfterAny("abc", "d")   = ""
func SubstrAfterAny(s string, b string) string {
	if s == "" {
		return s
	}

	i := IndexAny(s, b)
	if i < 0 {
		return ""
	}

	_, z := utf8.DecodeRuneInString(s[i:])
	return s[i+z:]
}

// SubstrAfterByte Gets the substring after the first occurrence of a separator b.
// The separator b is not returned.
// If nothing is found, the empty string is returned.
// SubstrAfterByte("", *)        = ""
// SubstrAfterByte("abc", 'a')   = "bc"
// SubstrAfterByte("abcba", 'b') = "cba"
// SubstrAfterByte("abc", 'c')   = ""
// SubstrAfterByte("abc", 'd')   = ""
func SubstrAfterByte(s string, b byte) string {
	if s == "" {
		return s
	}

	i := IndexByte(s, b)
	if i < 0 {
		return ""
	}
	return s[i+1:]
}

// SubstrAfterRune Gets the substring after the first occurrence of a separator r.
// The separator r is not returned.
// If nothing is found, the empty string is returned.
// SubstrAfterRune("", *)        = ""
// SubstrAfterRune("abc", 'a')   = "bc"
// SubstrAfterRune("abcba", 'b') = "cba"
// SubstrAfterRune("abc", 'c')   = ""
// SubstrAfterRune("abc", 'd')   = ""
func SubstrAfterRune(s string, r rune) string {
	if s == "" {
		return s
	}

	i := IndexRune(s, r)
	if i < 0 {
		return ""
	}
	return s[i+utf8.RuneLen(r):]
}

// SubstrAfterLast Gets the substring after the last occurrence of a separator b.
// The separator b is not returned.
// If nothing is found, the empty string is returned.
//
// SubstrAfterLast("", *)        = ""
// SubstrAfterLast("abc", "a")   = "bc"
// SubstrAfterLast("abcba", "b") = "a"
// SubstrAfterLast("abc", "c")   = ""
// SubstrAfterLast("a", "a")     = ""
// SubstrAfterLast("a", "z")     = ""
func SubstrAfterLast(s string, b string) string {
	if s == "" {
		return s
	}

	i := LastIndex(s, b)
	if i < 0 {
		return ""
	}
	return s[i+len(b):]
}

// SubstrAfterLastAny Gets the substring after the last occurrence of any char in separator b.
// The separator b is not returned.
// If nothing is found, the empty string is returned.
//
// SubstrAfterLastAny("", *)        = ""
// SubstrAfterLastAny("abc", "az")   = "bc"
// SubstrAfterLastAny("abcba", "b") = "a"
// SubstrAfterLastAny("abc", "c")   = ""
// SubstrAfterLastAny("a", "a")     = ""
// SubstrAfterLastAny("a", "z")     = ""
func SubstrAfterLastAny(s string, b string) string {
	if s == "" {
		return s
	}

	i := LastIndexAny(s, b)
	if i < 0 {
		return ""
	}

	_, z := utf8.DecodeRuneInString(s[i:])
	return s[i+z:]
}

// SubstrAfterLastByte Gets the substring after the last occurrence of a separator b.
// The separator b is not returned.
// If nothing is found, the empty string is returned.
//
// SubstrAfterLastByte("", *)        = ""
// SubstrAfterLastByte("abc", 'a')   = "bc"
// SubstrAfterLastByte("abcba", 'b') = "a"
// SubstrAfterLastByte("abc", 'c')   = ""
// SubstrAfterLastByte("a", 'a')     = ""
// SubstrAfterLastByte("a", 'z')     = ""
func SubstrAfterLastByte(s string, b byte) string {
	if s == "" {
		return s
	}

	i := LastIndexByte(s, b)
	if i < 0 {
		return ""
	}
	return s[i+1:]
}

// SubstrAfterLastRune Gets the substring after the last occurrence of a separator r.
// The separator r is not returned.
// If nothing is found, the empty string is returned.
//
// SubstrAfterLastRune("", *)        = ""
// SubstrAfterLastRune("abc", 'a')   = "bc"
// SubstrAfterLastRune("abcba", 'b') = "a"
// SubstrAfterLastRune("abc", 'c')   = ""
// SubstrAfterLastRune("a", 'a')     = ""
// SubstrAfterLastRune("a", 'z')     = ""
func SubstrAfterLastRune(s string, r rune) string {
	if s == "" {
		return s
	}

	i := LastIndexRune(s, r)
	if i < 0 {
		return ""
	}
	return s[i+utf8.RuneLen(r):]
}

// SubstrBefore Gets the substring before the first occurrence of a separator b.
// The separator b is not returned.
// If nothing is found, the input string is returned.
// SubstrBefore("", *)        = ""
// SubstrBefore("abc", "a")   = ""
// SubstrBefore("abcba", "b") = "a"
// SubstrBefore("abc", "c")   = "ab"
// SubstrBefore("abc", "d")   = "abc"
func SubstrBefore(s string, b string) string {
	if s == "" {
		return s
	}

	i := Index(s, b)
	if i < 0 {
		return s
	}
	return s[:i]
}

// SubstrBeforeAny Gets the substring before the first occurrence of any chat in separator b.
// The separator b is not returned.
// If nothing is found, the input string is returned.
// SubstrBeforeAny("", *)        = ""
// SubstrBeforeAny("abc", "a")   = ""
// SubstrBeforeAny("abcba", "cb") = "a"
// SubstrBeforeAny("abc", "zc")   = "ab"
// SubstrBeforeAny("abc", "d")   = "abc"
func SubstrBeforeAny(s string, b string) string {
	if s == "" {
		return s
	}

	i := IndexAny(s, b)
	if i < 0 {
		return s
	}
	return s[:i]
}

// SubstrBeforeByte Gets the substring before the first occurrence of a separator b.
// The separator b is not returned.
// If nothing is found, the input string is returned.
// SubstrBeforeByte("", *)        = ""
// SubstrBeforeByte("abc", 'a')   = ""
// SubstrBeforeByte("abcba", 'b') = "a"
// SubstrBeforeByte("abc", 'c')   = "ab"
// SubstrBeforeByte("abc", 'd')   = "abc"
func SubstrBeforeByte(s string, b byte) string {
	if s == "" {
		return s
	}

	i := IndexByte(s, b)
	if i < 0 {
		return s
	}
	return s[:i]
}

// SubstrBeforeRune Gets the substring before the first occurrence of a separator r.
// The separator r is not returned.
// If nothing is found, the input string is returned.
// SubstrBeforeRune("", *)        = ""
// SubstrBeforeRune("abc", 'a')   = ""
// SubstrBeforeRune("abcba", 'b') = "a"
// SubstrBeforeRune("abc", 'c')   = "ab"
// SubstrBeforeRune("abc", 'd')   = "abc"
func SubstrBeforeRune(s string, r rune) string {
	if s == "" {
		return s
	}

	i := IndexRune(s, r)
	if i < 0 {
		return s
	}
	return s[:i]
}

// SubstrBeforeLast Gets the substring before the last occurrence of a separator b.
// The separator b is not returned.
// If nothing is found, the input string is returned.
//
// SubstrBeforeLast("", *)        = ""
// SubstrBeforeLast("abc", "a")   = ""
// SubstrBeforeLast("abcba", "b") = "a"
// SubstrBeforeLast("abc", "c")   = "ab"
// SubstrBeforeLast("a", "a")     = ""
// SubstrBeforeLast("a", "z")     = "a"
// SubstrBeforeLast("a", "")      = "a"
func SubstrBeforeLast(s string, b string) string {
	if s == "" {
		return s
	}

	i := LastIndex(s, b)
	if i < 0 {
		return s
	}
	return s[:i]
}

// SubstrBeforeLastAny Gets the substring before the last occurrence of any char in separator b.
// The separator b is not returned.
// If nothing is found, the input string is returned.
//
// SubstrBeforeLastAny("", *)        = ""
// SubstrBeforeLastAny("abc", "a")   = ""
// SubstrBeforeLastAny("abcba", "bc") = "a"
// SubstrBeforeLastAny("abc", "zc")   = "ab"
// SubstrBeforeLastAny("a", "a")     = ""
// SubstrBeforeLastAny("a", "z")     = "a"
// SubstrBeforeLastAny("a", "")      = "a"
func SubstrBeforeLastAny(s string, b string) string {
	if s == "" {
		return s
	}

	i := LastIndexAny(s, b)
	if i < 0 {
		return s
	}
	return s[:i]
}

// SubstrBeforeLastByte Gets the substring before the last occurrence of a separator b.
// The separator b is not returned.
// If nothing is found, the input string is returned.
//
// SubstrBeforeLastByte("", *)        = ""
// SubstrBeforeLastByte("abc", 'a')   = ""
// SubstrBeforeLastByte("abcba", 'b') = "abc"
// SubstrBeforeLastByte("abc", 'c')   = "ab"
// SubstrBeforeLastByte("a", 'a')     = ""
// SubstrBeforeLastByte("a", 'z')     = "a"
func SubstrBeforeLastByte(s string, b byte) string {
	if s == "" {
		return s
	}

	i := LastIndexByte(s, b)
	if i < 0 {
		return s
	}
	return s[:i]
}

// SubstrBeforeLastRune Gets the substring before the last occurrence of a separator r.
// The separator r is not returned.
// If nothing is found, the input string is returned.
//
// SubstrBeforeLastRune("", *)        = ""
// SubstrBeforeLastRune("abc", 'a')   = ""
// SubstrBeforeLastRune("abcba", 'b') = "abc"
// SubstrBeforeLastRune("abc", 'c')   = "ab"
// SubstrBeforeLastRune("a", 'a')     = ""
// SubstrBeforeLastRune("a", 'z')     = "a"
func SubstrBeforeLastRune(s string, r rune) string {
	if s == "" {
		return s
	}

	i := LastIndexRune(s, r)
	if i < 0 {
		return s
	}
	return s[:i]
}

// CutPrefix returns s without the provided leading prefix string
// and reports whether it found the prefix.
// If s doesn't start with prefix, CutPrefix returns s, false.
// If prefix is the empty string, CutPrefix returns s, true.
func CutPrefix(s, prefix string) (after string, found bool) {
	if !HasPrefix(s, prefix) {
		return s, false
	}
	return s[len(prefix):], true
}

// CutSuffix returns s without the provided ending suffix string
// and reports whether it found the suffix.
// If s doesn't end with suffix, CutSuffix returns s, false.
// If suffix is the empty string, CutSuffix returns s, true.
func CutSuffix(s, suffix string) (before string, found bool) {
	if !HasSuffix(s, suffix) {
		return s, false
	}
	return s[:len(s)-len(suffix)], true
}

// Cut slices s around the first instance of sep,
// returning the text before and after sep.
// The found result reports whether sep appears in s.
// If sep does not appear in s, cut returns s, "", false.
func Cut(s, sep string) (before, after string, found bool) {
	if i := Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

// CutByte slices s around the first instance of sep,
// returning the text before and after sep.
// The found result reports whether sep appears in s.
// If sep does not appear in s, cut returns s, "", false.
func CutByte(s string, sep byte) (before, after string, found bool) {
	if i := IndexByte(s, sep); i >= 0 {
		return s[:i], s[i+1:], true
	}
	return s, "", false
}

// CutRune slices s around the first instance of sep,
// returning the text before and after sep.
// The found result reports whether sep appears in s.
// If sep does not appear in s, cut returns s, "", false.
func CutRune(s string, sep rune) (before, after string, found bool) {
	if i := IndexRune(s, sep); i >= 0 {
		return s[:i], s[i+RuneLen(sep):], true
	}
	return s, "", false
}

// CutAny slices s around the first instance of any Unicode code point
// from chars in s, returning the text before and after separator.
// The found result reports whether any chars's rune appears in s.
// If not found, cut returns s, "", false.
func CutAny(s string, chars string) (before, after string, found bool) {
	if i := IndexAny(s, chars); i >= 0 {
		_, z := utf8.DecodeRuneInString(s[i:])
		return s[:i], s[i+z:], true
	}
	return s, "", false
}

// CutFunc slices s around the first instance of the Unicode
// code point satisfying f(c), returning the text before and after f(c).
// The found result reports whether f(c) appears in s.
// If f(c) does not appear in s, cut returns s, "", false.
func CutFunc(s string, f func(rune) bool) (before, after string, found bool) {
	if i := IndexFunc(s, f); i >= 0 {
		_, z := utf8.DecodeRuneInString(s[i:])
		return s[:i], s[i+z:], true
	}
	return s, "", false
}

// LastCut slices s around the last instance of sep,
// returning the text before and after sep.
// The found result reports whether sep appears in s.
// If sep does not appear in s, cut returns s, "", false.
func LastCut(s, sep string) (before, after string, found bool) {
	if i := LastIndex(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

// LastCutByte slices s around the last instance of sep,
// returning the text before and after sep.
// The found result reports whether sep appears in s.
// If sep does not appear in s, cut returns s, "", false.
func LastCutByte(s string, sep byte) (before, after string, found bool) {
	if i := LastIndexByte(s, sep); i >= 0 {
		return s[:i], s[i+1:], true
	}
	return s, "", false
}

// LastCutRune slices s around the last instance of sep,
// returning the text before and after sep.
// The found result reports whether sep appears in s.
// If sep does not appear in s, cut returns s, "", false.
func LastCutRune(s string, sep rune) (before, after string, found bool) {
	if i := LastIndexRune(s, sep); i >= 0 {
		return s[:i], s[i+RuneLen(sep):], true
	}
	return s, "", false
}

// LastCutAny slices s around the last instance of any Unicode code point
// from chars in s, returning the text before and after separator.
// The found result reports whether any chars's rune appears in s.
// If not found, cut returns s, "", false.
func LastCutAny(s string, chars string) (before, after string, found bool) {
	if i := LastIndexAny(s, chars); i >= 0 {
		_, z := utf8.DecodeRuneInString(s[i:])
		return s[:i], s[i+z:], true
	}
	return s, "", false
}

// LastCutFunc slices s around the last instance of the Unicode
// code point satisfying f(c), returning the text before and after f(c).
// The found result reports whether f(c) appears in s.
// If f(c) does not appear in s, cut returns s, "", false.
func LastCutFunc(s string, f func(rune) bool) (before, after string, found bool) {
	if i := LastIndexFunc(s, f); i >= 0 {
		_, z := utf8.DecodeRuneInString(s[i:])
		return s[:i], s[i+z:], true
	}
	return s, "", false
}

// CutAt slices s around the rune position p,
// returning the text before and start position p.
// panic if p < 0.
// if p == 0, return "", s.
// if p > RuneCount(s), return s, "".
func CutAt(s string, p int) (before, after string) {
	if p < 0 {
		panic(fmt.Sprintf("invalid argument: position %d must not be negative", p))
	}
	if p == 0 {
		return "", s
	}

	if p > len(s) {
		return s, ""
	}

	z := RuneCount(s)
	if z < p {
		return s, ""
	}

	i, c, a := 0, 0, s
	for a != "" {
		_, z = utf8.DecodeRuneInString(a)
		i += z
		c++
		if c >= p {
			break
		}
		a = a[z:]
	}

	return s[:i], s[i:]
}

// Left return the leftmost n rune string.
func Left(s string, n int) string {
	if n <= 0 {
		return ""
	}

	if len(s) <= n {
		return s
	}

	z := RuneCount(s)
	if z <= n {
		return s
	}

	i, c, a := 0, 0, s
	for a != "" {
		_, z = utf8.DecodeRuneInString(a)
		i += z
		c++
		if c >= n {
			break
		}
		a = a[z:]
	}

	return s[:i]
}

// Right return the rightmost n rune string.
func Right(s string, n int) string {
	if n <= 0 {
		return ""
	}

	if len(s) <= n {
		return s
	}

	z := RuneCount(s)
	if z <= n {
		return s
	}

	i, c, a := len(s), 0, s
	for a != "" {
		_, z = utf8.DecodeLastRuneInString(a)
		i -= z
		c++
		if c >= n {
			break
		}
		a = s[:i]
	}

	return s[i:]
}

// Mid slices s around the rune position p,
// returning the max n rune string start from the position p.
// panic if p < 0.
func Mid(s string, p, n int) string {
	_, a := CutAt(s, p)
	return Left(a, n)
}

// Ellipsis abbreviates a string with max rune count `n` using ellipses (default "...").
// This will turn "Now is the time for all good men" into "Now is the time for..."
func Ellipsis(s string, n int, ellipses ...string) string {
	if n <= 0 || s == "" {
		return ""
	}

	sc := RuneCount(s)
	if sc <= n {
		return s
	}

	e := "..."
	for _, d := range ellipses {
		if d != "" {
			e = d
			break
		}
	}

	ec := RuneCount(e)
	if n <= ec {
		return Left(e, n)
	}

	return Left(s, n-ec) + e
}

// Ellipsiz abbreviates a string with max rune count `n` using ellipses (default "…").
// The non ascii rune (> 0x80) in string `s` will count as 2.
// This will turn "こんにちは世界" into "こんにちは…"
func Ellipsiz(s string, n int, ellipses ...string) string {
	if n <= 0 || s == "" {
		return ""
	}

	e := "…"
	for _, d := range ellipses {
		if d != "" {
			e = d
			break
		}
	}

	ec := RuneCount(e)

	cut := n - ec
	sz, j := 0, -1
	for i, r := range s {
		sz++
		if r > utf8.RuneSelf {
			sz++
		}
		if sz > cut && j < 0 {
			j = i
		}
		if sz > n {
			if n <= ec {
				return Left(e, n)
			}
			return s[0:j] + e
		}
	}
	return s
}
