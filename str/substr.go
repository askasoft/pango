package str

import (
	"strings"
)

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

	i := strings.IndexByte(s, b)
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

	i := strings.IndexRune(s, r)
	if i < 0 {
		return ""
	}
	return s[i+1:]
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

	i := strings.LastIndexByte(s, b)
	if i < 0 || i == len(s)-1 {
		return ""
	}
	return s[i+1:]
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

	i := strings.LastIndex(s, b)
	if i < 0 || i == len(s)-len(b) {
		return ""
	}
	return s[i+len(b):]
}

// SubstrBeforeByte Gets the substring before the first occurrence of a separator b.
// The separator b is not returned.
// If nothing is found, the empty string is returned.
// SubstrBeforeByte("", *)        = ""
// SubstrBeforeByte("abc", 'a')   = ""
// SubstrBeforeByte("abcba", 'b') = "a"
// SubstrBeforeByte("abc", 'c')   = "ab"
// SubstrBeforeByte("abc", 'd')   = ""
func SubstrBeforeByte(s string, b byte) string {
	if s == "" {
		return s
	}

	i := strings.IndexByte(s, b)
	if i < 0 {
		return ""
	}
	return s[:i]
}

// SubstrBeforeRune Gets the substring before the first occurrence of a separator r.
// The separator r is not returned.
// If nothing is found, the empty string is returned.
// SubstrBeforeRune("", *)        = ""
// SubstrBeforeRune("abc", 'a')   = ""
// SubstrBeforeRune("abcba", 'b') = "a"
// SubstrBeforeRune("abc", 'c')   = "ab"
// SubstrBeforeRune("abc", 'd')   = ""
func SubstrBeforeRune(s string, r rune) string {
	if s == "" {
		return s
	}

	i := strings.IndexRune(s, r)
	if i <= 0 {
		return ""
	}
	return s[:i]
}

// SubstrBeforeLastByte Gets the substring before the last occurrence of a separator b.
// The separator b is not returned.
// If nothing is found, the empty string is returned.
//
// SubstrBeforeLastByte("", *)        = ""
// SubstrBeforeLastByte("abc", 'a')   = ""
// SubstrBeforeLastByte("abcba", 'b') = "abc"
// SubstrBeforeLastByte("abc", 'c')   = "ab"
// SubstrBeforeLastByte("a", 'a')     = ""
// SubstrBeforeLastByte("a", 'z')     = ""
func SubstrBeforeLastByte(s string, b byte) string {
	if s == "" {
		return s
	}

	i := strings.LastIndexByte(s, b)
	if i <= 0 {
		return ""
	}
	return s[:i]
}

// SubstrBeforeLast Gets the substring before the last occurrence of a separator b.
// The separator b is not returned.
// If nothing is found, the empty string is returned.
//
// SubstrBeforeLast("", *)        = ""
// SubstrBeforeLast("abc", "a")   = ""
// SubstrBeforeLast("abcba", "b") = "a"
// SubstrBeforeLast("abc", "c")   = "ab"
// SubstrBeforeLast("a", "a")     = ""
// SubstrBeforeLast("a", "z")     = ""
// SubstrBeforeLast("a", "")      = "a"
func SubstrBeforeLast(s string, b string) string {
	if s == "" {
		return s
	}

	i := strings.LastIndex(s, b)
	if i <= 0 {
		return ""
	}
	return s[:i]
}
