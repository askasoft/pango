package mkdn

import "github.com/askasoft/pango/str"

var specials = "\\`*_~{}[]()<>#!|+-."

var escapes = "\\`*_~{}[]()<>#!|"

// EscapeString returns a string that escapes all markdown metacharacters inside the argument text.
func EscapeString(s string) string {
	i := str.IndexAny(s, escapes)
	if i < 0 {
		return s
	}

	var sb str.Builder

	for i >= 0 {
		sb.WriteString(s[0:i])
		sb.WriteByte('\\')
		sb.WriteByte(s[i])
		s = s[i+1:]
		i = str.IndexAny(s, escapes)
	}

	if s != "" {
		sb.WriteString(s)
	}

	return sb.String()
}

// UnescapeString returns a string unescaped for markdown metacharacters.
func UnescapeString(s string) string {
	i := str.IndexByte(s, '\\')
	if i < 0 || i >= len(s)-1 {
		return s
	}

	var sb str.Builder

	for i >= 0 && i < len(s)-1 {
		sb.WriteString(s[0:i])
		if str.ContainsByte(specials, s[i+1]) {
			sb.WriteByte(s[i+1])
			s = s[i+2:]
		} else {
			sb.WriteByte('\\')
			s = s[i+1:]
		}
		i = str.IndexByte(s, '\\')
	}

	if s != "" {
		sb.WriteString(s)
	}

	return sb.String()
}
