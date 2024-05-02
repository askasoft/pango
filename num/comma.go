package num

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// Comma produces a string form of the given number in base 10 with
// commas after every three orders of magnitude.
//
// e.g. Comma(834142) -> 834,142
// e.g. Comma(834142, "_") -> 834_142
func Comma(n any, c ...string) string {
	v := reflect.ValueOf(n)

	switch v.Kind() {
	case reflect.Int8, reflect.Uint8:
		return strconv.Itoa(int(v.Int()))
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return CommaInt(v.Int(), c...)
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return CommaUint(v.Uint(), c...)
	case reflect.Float32, reflect.Float64:
		return CommaFloat(v.Float(), c...)
	default:
		panic(fmt.Sprintf("Comma: unknown type for '%v' (%T)", n, n))
	}
}

// CommaInt produces a string form of the given number in base 10 with
// commas after every three orders of magnitude.
//
// e.g. CommaInt(834142) -> 834,142
// e.g. CommaInt(834142, "_") -> 834_142
func CommaInt(v int64, c ...string) string {
	// Min int64 can't be negated to a usable value, so it has to be special cased.
	if v == math.MinInt64 {
		return "-9,223,372,036,854,775,808"
	}

	sign := ""
	if v < 0 {
		sign = "-"
		v = 0 - v
	}

	return sign + CommaUint((uint64)(v), c...)
}

func comma(c ...string) string {
	sep := ","
	if len(c) > 0 {
		sep = c[0]
	}
	return sep
}

// CommaUint produces a string form of the given number in base 10 with
// commas after every three orders of magnitude.
//
// e.g. CommaUint(834142) -> 834,142
// e.g. CommaUint(834142, "_") -> 834_142
func CommaUint(v uint64, c ...string) string {
	parts := []string{"", "", "", "", "", "", ""}
	j := len(parts) - 1

	for v > 999 {
		parts[j] = strconv.FormatUint(v%1000, 10)
		switch len(parts[j]) {
		case 2:
			parts[j] = "0" + parts[j]
		case 1:
			parts[j] = "00" + parts[j]
		}
		v = v / 1000
		j--
	}
	parts[j] = strconv.Itoa(int(v))

	sep := comma(c...)
	return strings.Join(parts[j:], sep)
}

// CommaFloat produces a string form of the given number in base 10 with
// commas after every three orders of magnitude.
//
// e.g. CommaFloat(834142.32) -> 834,142.32
// e.g. CommaFloat(834142.32, "_") -> 834_142.32
func CommaFloat(v float64, c ...string) string {
	s := strconv.FormatFloat(v, 'f', -1, 64)
	return CommaString(s, c...)
}

// CommaFloatWithDigits works like the Commaf but limits the resulting
// string to the given number of decimal places.
//
// e.g. CommaFloatWithDigits(834142.32, 1) -> 834,142.3
// e.g. CommaFloatWithDigits(834142.32, 1, "_") -> 834_142.3
func CommaFloatWithDigits(f float64, digits int, c ...string) string {
	s := FtoaWithDigits(f, digits)
	return CommaString(s, c...)
}

// CommaString produces a string form of the given number string in base 10 with
// commas after every three orders of magnitude.
//
// e.g. CommaString("834142.32") -> 834,142.32
// e.g. CommaString("834142.32", "_") -> 834_142.32
func CommaString(s string, c ...string) string {
	if s == "" {
		return s
	}

	sep := comma(c...)

	var buf bytes.Buffer

	if s[0] == '-' || s[0] == '+' {
		buf.WriteByte(s[0])
		s = s[1:]
	}

	parts := strings.Split(s, ".")
	pos := 0
	if len(parts[0])%3 != 0 {
		pos += len(parts[0]) % 3
		buf.WriteString(parts[0][:pos])
		buf.WriteString(sep)
	}
	for ; pos < len(parts[0]); pos += 3 {
		buf.WriteString(parts[0][pos : pos+3])
		buf.WriteString(sep)
	}
	buf.Truncate(buf.Len() - 1)

	if len(parts) > 1 {
		buf.Write([]byte{'.'})
		buf.WriteString(parts[1])
	}
	return buf.String()
}
