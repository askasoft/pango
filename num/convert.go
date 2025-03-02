package num

import (
	"fmt"
	"strconv"
)

// NonZeroi return the first non-zero value of defs if error.
func NonZeroi(ns ...int) int {
	for _, n := range ns {
		if n != 0 {
			return n
		}
	}
	return 0
}

// Atoi use strconv.ParseInt(s, 0, strconv.IntSize) to parse string 's' to int,
// return the first non-zero value of defs if error.
func Atoi(s string, defs ...int) int {
	if s == "" {
		return NonZeroi(defs...)
	}

	if i, err := strconv.ParseInt(s, 0, strconv.IntSize); err == nil {
		return int(i)
	}
	return NonZeroi(defs...)
}

// Itoa is equivalent to strconv.FormatInt(int64(i), 10).
func Itoa(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

// NonZerol return the first non-zero value of defs if error.
func NonZerol(ns ...int64) int64 {
	for _, n := range ns {
		if n != 0 {
			return n
		}
	}
	return 0
}

// Atol use strconv.ParseInt(s, 0, 64) to parse string 's' to int64,
// return the first non-zero value of defs if error.
func Atol(s string, defs ...int64) int64 {
	if s == "" {
		return NonZerol(defs...)
	}

	if i, err := strconv.ParseInt(s, 0, 64); err == nil {
		return i
	}
	return NonZerol(defs...)
}

// Itoa is equivalent to strconv.FormatInt(i, 10).
func Ltoa(i int64) string {
	return strconv.FormatInt(i, 10)
}

// NonZerof return the first non-zero value of defs if error.
func NonZerof(ns ...float64) float64 {
	for _, n := range ns {
		if n != 0 {
			return n
		}
	}
	return 0
}

// Atol use strconv.ParseFloat(s, 64) to parse string 's' to float64,
// return the first non-zero value of defs if error.
func Atof(s string, defs ...float64) float64 {
	if s == "" {
		return NonZerof(defs...)
	}

	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return NonZerof(defs...)
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
