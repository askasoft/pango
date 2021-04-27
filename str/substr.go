package str

import (
	"strings"
)

// StringAfterByte Gets the substring after the first occurrence of a separator. The separator is not returned.
// If nothing is found, the empty string is returned.
// StringAfterByte("", *)        = ""
// StringAfterByte("abc", 'a')   = "bc"
// StringAfterByte("abcba", 'b') = "cba"
// StringAfterByte("abc", 'c')   = ""
// StringAfterByte("abc", 'd')   = ""
func StringAfterByte(s string, b byte) string {
	if s == "" {
		return s
	}

	i := strings.IndexByte(s, b)
	if i < 0 {
		return ""
	}
	return s[i+1:]
}

// StringAfterRune Gets the substring after the first occurrence of a separator. The separator is not returned.
// If nothing is found, the empty string is returned.
// StringAfterRune("", *)        = ""
// StringAfterRune("abc", 'a')   = "bc"
// StringAfterRune("abcba", 'b') = "cba"
// StringAfterRune("abc", 'c')   = ""
// StringAfterRune("abc", 'd')   = ""
func StringAfterRune(s string, r rune) string {
	if s == "" {
		return s
	}

	i := strings.IndexRune(s, r)
	if i < 0 {
		return ""
	}
	return s[i+1:]
}

// StringAfterLastByte Gets the substring after the last occurrence of a separator. The separator is not returned.
// If nothing is found, the empty string is returned.
//
// StringAfterLastByte("", *)        = ""
// StringAfterLastByte("abc", 'a')   = "bc"
// StringAfterLastByte("abcba", 'b') = "a"
// StringAfterLastByte("abc", 'c')   = ""
// StringAfterLastByte("a", 'a')     = ""
// StringAfterLastByte("a", 'z')     = ""
func StringAfterLastByte(s string, b byte) string {
	if s == "" {
		return s
	}

	i := strings.LastIndexByte(s, b)
	if i < 0 || i == len(s)-1 {
		return ""
	}
	return s[i+1:]
}

// StringAfterLast Gets the substring after the last occurrence of a separator. The separator is not returned.
// If nothing is found, the empty string is returned.
//
// StringAfterLast("", *)        = ""
// StringAfterLast("abc", "a")   = "bc"
// StringAfterLast("abcba", "b") = "a"
// StringAfterLast("abc", "c")   = ""
// StringAfterLast("a", "a")     = ""
// StringAfterLast("a", "z")     = ""
func StringAfterLast(s string, c string) string {
	if s == "" {
		return s
	}

	i := strings.LastIndex(s, c)
	if i < 0 || i == len(s)-len(c) {
		return ""
	}
	return s[i+len(c):]
}
