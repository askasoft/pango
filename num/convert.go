package num

import (
	"fmt"
	"strconv"
)

// Atoi use strconv.ParseInt(s, 0, strconv.IntSize) to parse string 's' to int, return n[0] if error.
func Atoi(s string, n ...int) int {
	if s == "" {
		if len(n) > 0 {
			return n[0]
		}
		return 0
	}

	i, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil && len(n) > 0 {
		return n[0]
	}
	return int(i)
}

func Itoa(i int) string {
	return strconv.Itoa(i)
}

// Atol use strconv.ParseInt(s, 0, 64) to parse string 's' to int64, return n[0] if error.
func Atol(s string, n ...int64) int64 {
	if s == "" {
		if len(n) > 0 {
			return n[0]
		}
		return 0
	}

	i, err := strconv.ParseInt(s, 0, 64)
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
	if s == "" {
		if len(n) > 0 {
			return n[0]
		}
		return 0
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil && len(n) > 0 {
		return n[0]
	}
	return f
}

// Ftoa converts a float to a string with no trailing zeros.
func Ftoa(f float64) string {
	return StripTrailingZeros(strconv.FormatFloat(f, 'f', 6, 64))
}

// FtoaWithDigits converts a float to a string but limits the resulting string
// to the given number of decimal places, and no trailing zeros.
func FtoaWithDigits(f float64, digits int) string {
	switch {
	case digits < 0:
		return Ftoa(f)
	case digits > 0:
		return StripTrailingZeros(fmt.Sprintf("%.*f", digits, f))
	default:
		return fmt.Sprintf("%.0f", f)
	}
}

func StripTrailingZeros(s string) string {
	i := len(s) - 1
	for i > 0 {
		if s[i] == '.' {
			i--
			break
		}
		if s[i] != '0' {
			break
		}
		i--
	}
	return s[:i+1]
}
