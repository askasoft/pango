package num

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/askasoft/pango/asg"
)

// Atoi use strconv.ParseInt(s, 0, strconv.IntSize) to parse string 's' to int,
// returns the default defs[0] value if error.
func Atoi(s string, defs ...int) int {
	if s == "" {
		return asg.First(defs)
	}

	if n, err := strconv.ParseInt(s, 0, strconv.IntSize); err == nil {
		return int(n)
	}
	return asg.First(defs)
}

// Itoa is equivalent to strconv.FormatInt(int64(n), 10).
func Itoa(n int) string {
	return strconv.FormatInt(int64(n), 10)
}

// Atol use strconv.ParseInt(s, 0, 64) to parse string 's' to int64,
// returns the default defs[0] value if error.
func Atol(s string, defs ...int64) int64 {
	if s == "" {
		return asg.First(defs)
	}

	if n, err := strconv.ParseInt(s, 0, 64); err == nil {
		return n
	}
	return asg.First(defs)
}

// Itoa is equivalent to strconv.FormatInt(n, 10).
func Ltoa(n int64) string {
	return strconv.FormatInt(n, 10)
}

// Atol use strconv.ParseFloat(s, 64) to parse string 's' to float64,
// returns the default defs[0] value if error.
func Atof(s string, defs ...float64) float64 {
	if s == "" {
		return asg.First(defs)
	}

	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return asg.First(defs)
}

// Ftoa converts a float to a string with no trailing zeros.
func Ftoa(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
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

var int2roman = []struct {
	n int
	s string
}{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

// FormatRoman converts a positive integer to roman number string.
func FormatRoman(n int) (string, error) {
	if n < 1 || n > 3999 {
		// Roman numerals are conventionally defined for numbers 1 to 3999
		return "", fmt.Errorf("num.FormatRoman(%d): out of Roman numeral range (1-3999)", n)
	}

	var r strings.Builder
	for _, c := range int2roman {
		for n >= c.n {
			r.WriteString(c.s)
			n -= c.n
		}
	}
	return r.String(), nil
}

// IntToRoman converts a positive integer to roman number string.
func IntToRoman(n int) string {
	s, _ := FormatRoman(n)
	return s
}

var roman2int = map[byte]int{
	'I': 1,
	'V': 5,
	'X': 10,
	'L': 50,
	'C': 100,
	'D': 500,
	'M': 1000,
	'i': 1,
	'v': 5,
	'x': 10,
	'l': 50,
	'c': 100,
	'd': 500,
	'm': 1000,
}

// ParseRoman converts a roman numeric string to integer.
func ParseRoman(s string) (int, error) {
	if s == "" {
		return 0, errors.New(`num.ParseRoman(""): empty roman number`)
	}

	r, n := 0, len(s)
	for i := range n {
		cv, ok := roman2int[s[i]]
		if !ok {
			return 0, fmt.Errorf("num.ParseRoman(%q): invalid roman number", s)
		}

		// Check for subtractive cases
		if i+1 < n {
			nv := roman2int[s[i+1]]
			if cv < nv {
				r -= cv
			} else {
				r += cv
			}
		} else {
			// Last character, simply add its value
			r += cv
		}
	}
	return r, nil
}

// RomanToInt converts a roman numeric string to integer.
func RomanToInt(s string, defs ...int) int {
	if s == "" {
		return asg.First(defs)
	}

	if n, err := ParseRoman(s); err == nil {
		return n
	}

	return asg.First(defs)
}

// FormatAlpha converts a positive integer to Excel-style column string.
func FormatAlpha(n int) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("num.FormatAlpha(%d): invalid excel-style alpha number", n)
	}

	z := 0
	for x := n; x > 0; z++ {
		x = (x - 1) / 26
	}

	r := make([]rune, z)
	for i := len(r) - 1; i >= 0; i-- {
		// Excelの列番号は1始まりなので、1引いてから26で割る
		n-- // 0基準に変換
		r[i] = rune('A' + (n % 26))
		n /= 26
	}
	return string(r), nil
}

// IntToAlpha converts a positive integer to Excel-style column string.
func IntToAlpha(n int) string {
	s, _ := FormatAlpha(n)
	return s
}

// ParseAlpha converts a string like "A", "Z", "AA" to a number (Excel-style).
func ParseAlpha(s string) (int, error) {
	if s == "" {
		return 0, errors.New(`num.ParseAlpha(""): empty alpha number`)
	}

	r := 0
	for i := range len(s) {
		v := s[i]
		switch {
		case v >= 'A' && v <= 'Z':
			r = r*26 + int(v-'A'+1)
		case v >= 'a' && v <= 'z':
			r = r*26 + int(v-'a'+1)
		default:
			return 0, fmt.Errorf("num.ParseAlpha(%q): invalid alpha number", s)
		}
	}
	return r, nil
}

// AlphaToInt converts a string like "A", "Z", "AA" to a number (Excel-style).
func AlphaToInt(s string, defs ...int) int {
	if s == "" {
		return asg.First(defs)
	}

	if n, err := ParseAlpha(s); err == nil {
		return n
	}

	return asg.First(defs)
}
