package sql

import "strings"

// EscapeLike escape sql like string
func EscapeLike(str string) string {
	sb := strings.Builder{}

	for _, c := range str {
		if c == '~' {
			sb.WriteRune('~')
			sb.WriteRune('~')
			continue
		}

		if c == '%' || c == '_' {
			sb.WriteRune('~')
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
func EscapeString(str string) string {
	return strings.ReplaceAll(str, "'", "''")
}

// StringLike build a string for like '%' + str + '%'
func StringLike(str string) string {
	return "%" + EscapeLike(str) + "%"
}

// StartsLike build a string for like prefix str + '%'
func StartsLike(str string) string {
	return EscapeLike(str) + "%"
}

// StartsLike build a string for like suffix '%' + str
func EndsLike(str string) string {
	return "%" + EscapeLike(str)
}
