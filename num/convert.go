package num

import (
	"strconv"
	"strings"
)

func stripTrailingZeros(s string) string {
	offset := len(s) - 1
	for offset > 0 {
		if s[offset] == '.' {
			offset--
			break
		}
		if s[offset] != '0' {
			break
		}
		offset--
	}
	return s[:offset+1]
}

func stripTrailingDigits(s string, digits int) string {
	if i := strings.Index(s, "."); i >= 0 {
		if digits <= 0 {
			return s[:i]
		}
		i++
		if i+digits >= len(s) {
			return s
		}
		return s[:i+digits]
	}
	return s
}

func Atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func Itoa(i int) string {
	return strconv.Itoa(i)
}

func Atol(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func Ltoa(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Atof(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// Ftoa converts a float to a string with no trailing zeros.
func Ftoa(f float64) string {
	return stripTrailingZeros(strconv.FormatFloat(f, 'f', 6, 64))
}

// FtoaWithDigits converts a float to a string but limits the resulting string
// to the given number of decimal places, and no trailing zeros.
func FtoaWithDigits(f float64, digits int) string {
	return stripTrailingZeros(stripTrailingDigits(strconv.FormatFloat(f, 'f', 6, 64), digits))
}
