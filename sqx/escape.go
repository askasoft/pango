package sqx

import "strings"

// EscapeLike escape sql like string.
// The default escape char is backslach '\\'.
func EscapeLike(s string, escape ...rune) string {
	esc := '\\'
	if len(escape) > 0 {
		esc = escape[0]
	}

	sb := strings.Builder{}

	for _, c := range s {
		if c == esc {
			sb.WriteRune(esc)
			sb.WriteRune(esc)
			continue
		}

		if c == '%' || c == '_' {
			sb.WriteRune(esc)
			sb.WriteRune(c)
			continue
		}

		sb.WriteRune(c)
	}

	return sb.String()
}

// At present, this method only turns single-quotes into doubled single-quotes (
// <code>"McHale's Navy"</code> => <code>"McHale”s Navy"</code>). It does not handle the cases
// of percent (%) or underscore (_) for use in LIKE clauses.
// </p>
// see http://www.jguru.com/faq/view.jsp?EID=8881
func EscapeString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// StringLike build a string for like '%' + s + '%'
// The default escape char is backslach '\\'.
func StringLike(s string, escape ...rune) string {
	return "%" + EscapeLike(s, escape...) + "%"
}

// StartsLike build a string for like prefix s + '%'
// The default escape char is backslach '\\'.
func StartsLike(s string, escape ...rune) string {
	return EscapeLike(s, escape...) + "%"
}

// StartsLike build a string for like suffix '%' + s
// The default escape char is backslach '\\'.
func EndsLike(s string, escape ...rune) string {
	return "%" + EscapeLike(s, escape...)
}
