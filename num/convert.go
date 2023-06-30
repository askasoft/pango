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

// Atoi use strconv.Atoi(s) to parse string 's' to int, return n[0] if error.
func Atoi(s string, n ...int) int {
	i, err := strconv.Atoi(s)
	if err != nil && len(n) > 0 {
		return n[0]
	}
	return i
}

func Itoa(i int) string {
	return strconv.Itoa(i)
}

// Atol use strconv.ParseInt(s, 10, 64) to parse string 's' to int64, return n[0] if error.
func Atol(s string, n ...int64) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil && len(n) > 0 {
		return n[0]
	}
	return i
}

func Ltoa(i int64) string {
	return strconv.FormatInt(i, 10)
}

// Atol use strconv.ParseFloat(s, 64) to parse string 's' to float64, return n[0] if error.
func Atof(s string, n ...float64) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil && len(n) > 0 {
		return n[0]
	}
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
