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
func StringAfterByte(s string, c byte) string {
	if len(s) == 0 {
		return s
	}

	i := strings.IndexByte(s, c)
	if i < 0 {
		return ""
	}
	return s[i+1:]
}
