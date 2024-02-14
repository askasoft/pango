package sqx

import "strings"

// EscapeLike escape sql like string
func EscapeLike(s string) string {
	sb := strings.Builder{}

	for _, c := range s {
		if c == '~' {
			sb.WriteString("~~")
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
func EscapeString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// StringLike build a string for like '%' + s + '%'
func StringLike(s string) string {
	return "%" + EscapeLike(s) + "%"
}

// StartsLike build a string for like prefix s + '%'
func StartsLike(s string) string {
	return EscapeLike(s) + "%"
}

// StartsLike build a string for like suffix '%' + s
func EndsLike(s string) string {
	return "%" + EscapeLike(s)
}

// Quote quote string 's' with quote string 'q', return (q + s + q)
func Quote(s string, q string) string {
	return q + s + q
}

// QuoteFunc return a quote function with quote string 'q'
func QuoteFunc(q string) func(string) string {
	return func(s string) string {
		return Quote(s, q)
	}
}
