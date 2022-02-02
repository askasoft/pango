package str

import "unicode/utf8"

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
	if i < 0 || i == len(s)-len(b) {
		return ""
	}
	return s[i+len(b):]
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
	if i < 0 || i == len(s)-1 {
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
	if i < 0 || i == len(s)-1 {
		return ""
	}
	return s[i+utf8.RuneLen(r):]
}

// SubstrBefore Gets the substring before the first occurrence of a separator b.
// The separator b is not returned.
// If nothing is found, the empty string is returned.
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

// SubstrBeforeByte Gets the substring before the first occurrence of a separator b.
// The separator b is not returned.
// If nothing is found, the empty string is returned.
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
// If nothing is found, the empty string is returned.
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
// If nothing is found, the empty string is returned.
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

// SubstrBeforeLastByte Gets the substring before the last occurrence of a separator b.
// The separator b is not returned.
// If nothing is found, the empty string is returned.
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
// If nothing is found, the empty string is returned.
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
