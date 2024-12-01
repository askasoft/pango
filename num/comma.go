package num

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Comma produces a string form of the given number in base 10 with
// commas after every three orders of magnitude.
//
// e.g. Comma(834142) -> 834,142
// e.g. Comma(834142.1234, 3) -> 834,142.123
// e.g. Comma(834142, "_") -> 834_142
// e.g. Comma(834142.1234, "_", 3) -> 834_142.123
func Comma(n any, args ...any) string {
	s, err := CommaAny(n, args...)
	if err != nil {
		panic(err)
	}
	return s
}

// CommaAny produces a string form of the given number in base 10 with
// commas after every three orders of magnitude.
//
// e.g. CommaAny(834142) -> 834,142
// e.g. CommaAny(834142.1234, 3) -> 834,142.123
// e.g. CommaAny(834142, "_") -> 834_142
// e.g. CommaAny(834142.1234, "_", 3) -> 834_142.123
func CommaAny(n any, args ...any) (string, error) {
	c, d := ",", 0
	for _, a := range args {
		switch p := a.(type) {
		case string:
			c = p
		case int:
			d = p
		}
	}

	switch v := n.(type) {
	case int8:
		return Itoa(int(v)), nil
	case int16:
		return CommaInt(int64(v), c), nil
	case int32:
		return CommaInt(int64(v), c), nil
	case int64:
		return CommaInt(v, c), nil
	case int:
		return CommaInt(int64(v), c), nil
	case uint8:
		return Itoa(int(v)), nil
	case uint16:
		return CommaUint(uint64(v), c), nil
	case uint32:
		return CommaUint(uint64(v), c), nil
	case uint64:
		return CommaUint(v, c), nil
	case uint:
		return CommaUint(uint64(v), c), nil
	case float32:
		if d > 0 {
			return CommaFloatWithDigits(float64(v), d, c), nil
		}
		return CommaFloat(float64(v), c), nil
	case float64:
		if d > 0 {
			return CommaFloatWithDigits(float64(v), d, c), nil
		}
		return CommaFloat(float64(v), c), nil
	default:
		return "", fmt.Errorf("Comma: unknown type for '%v' (%T)", n, n)
	}
}

// CommaInt produces a string form of the given number in base 10 with
// commas after every three orders of magnitude.
//
// e.g. CommaInt(834142) -> 834,142
// e.g. CommaInt(834142, "_") -> 834_142
func CommaInt(v int64, c string) string {
	// Min int64 can't be negated to a usable value, so it has to be special cased.
	if v == math.MinInt64 {
		return "-9,223,372,036,854,775,808"
	}

	sign := ""
	if v < 0 {
		sign = "-"
		v = 0 - v
	}

	return sign + CommaUint((uint64)(v), c)
}

// CommaUint produces a string form of the given number in base 10 with
// commas after every three orders of magnitude.
//
// e.g. CommaUint(834142) -> 834,142
// e.g. CommaUint(834142, "_") -> 834_142
func CommaUint(v uint64, c string) string {
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

	return strings.Join(parts[j:], c)
}

// CommaFloat produces a string form of the given number in base 10 with
// commas after every three orders of magnitude.
//
// e.g. CommaFloat(834142.32) -> 834,142.32
// e.g. CommaFloat(834142.32, "_") -> 834_142.32
func CommaFloat(v float64, c string) string {
	s := strconv.FormatFloat(v, 'f', -1, 64)
	return CommaString(s, c)
}

// CommaFloatWithDigits works like the Commaf but limits the resulting
// string to the given number of decimal places.
//
// e.g. CommaFloatWithDigits(834142.32, 1) -> 834,142.3
// e.g. CommaFloatWithDigits(834142.32, 1, "_") -> 834_142.3
func CommaFloatWithDigits(f float64, digits int, c string) string {
	s := FtoaWithDigits(f, digits)
	return CommaString(s, c)
}

// CommaString produces a string form of the given number string in base 10 with
// commas after every three orders of magnitude.
//
// e.g. CommaString("834142.32") -> 834,142.32
// e.g. CommaString("834142.32", "_") -> 834_142.32
func CommaString(s string, c string) string {
	if s == "" {
		return s
	}

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
		buf.WriteString(c)
	}
	for ; pos < len(parts[0]); pos += 3 {
		buf.WriteString(parts[0][pos : pos+3])
		buf.WriteString(c)
	}
	buf.Truncate(buf.Len() - 1)

	if len(parts) > 1 {
		buf.Write([]byte{'.'})
		buf.WriteString(parts[1])
	}
	return buf.String()
}
