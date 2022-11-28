package str

import (
	"unicode"
)

// IsEmpty checks if the string is null.
func IsEmpty(s string) bool {
	return len(s) == 0
}

// IsNotEmpty checks if the string is not null.
func IsNotEmpty(s string) bool {
	return len(s) > 0
}

// IsASCII checks if the string contains ASCII chars only.
func IsASCII(s string) bool {
	if s == "" {
		return false
	}

	l := len(s)
	for i := 0; i < l; i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// IsPrintableASCII checks if the string contains printable ASCII chars only.
func IsPrintableASCII(s string) bool {
	if s == "" {
		return false
	}

	l := len(s)
	for i := 0; i < l; i++ {
		b := s[i]
		if b < ' ' || b > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// IsAlpha checks if the string contains only letters (a-zA-Z).
func IsAlpha(s string) bool {
	if s == "" {
		return false
	}

	l := len(s)
	for i := 0; i < l; i++ {
		b := s[i]
		if !((b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')) {
			return false
		}
	}
	return true
}

// IsAlphaNumber checks if the string contains only letters and numbers.
func IsAlphaNumber(s string) bool {
	if s == "" {
		return false
	}

	l := len(s)
	for i := 0; i < l; i++ {
		b := s[i]
		if !((b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')) {
			return false
		}
	}
	return true
}

// IsUTFLetter checks if the string contains only unicode letter characters.
// Similar to IsAlpha but for all languages.
func IsUTFLetter(s string) bool {
	if s == "" {
		return false
	}

	for _, c := range s {
		if !unicode.IsLetter(c) {
			return false
		}
	}
	return true

}

// IsNumber checks if the string contains only numbers.
func IsNumber(s string) bool {
	if s == "" {
		return false
	}

	l := len(s)
	for i := 0; i < l; i++ {
		b := s[i]
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}

// IsNumeric checks if the string contains only numbers and prefix [+-].
func IsNumeric(s string) bool {
	if s == "" {
		return false
	}

	if len(s) > 1 && (s[0] == '+' || s[0] == '-') {
		s = s[1:]
	}

	l := len(s)
	for i := 0; i < l; i++ {
		b := s[i]
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}

// IsUTFNumeric checks if the string contains only unicode numbers of any kind.
// Numbers can be 0-9 but also Fractions ¾,Roman Ⅸ and Hangzhou 〩.
// Prefix +- are allowed.
func IsUTFNumeric(s string) bool {
	if s == "" {
		return false
	}

	if len(s) > 1 && (s[0] == '+' || s[0] == '-') {
		s = s[1:]
	}

	for _, c := range s {
		if !unicode.IsNumber(c) {
			return false
		}
	}
	return true

}

// IsUTFDigit checks if the string contains only unicode radix-10 decimal digits.
func IsUTFDigit(s string) bool {
	if s == "" {
		return false
	}

	if len(s) > 1 && (s[0] == '+' || s[0] == '-') {
		s = s[1:]
	}

	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

// IsUTFLetterNumber checks if the string contains only unicode letters and numbers.
func IsUTFLetterNumber(s string) bool {
	if s == "" {
		return false
	}

	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsNumber(c) { //letters && numbers are ok
			return false
		}
	}
	return true

}

// IsHexDecimal checks if the string is a hexadecimal number `^(0[xX])?[0-9a-fA-F]+$`.
func IsHexDecimal(s string) bool {
	if s == "" {
		return false
	}

	if len(s) > 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
		s = s[2:]
	}

	l := len(s)
	for i := 0; i < l; i++ {
		b := s[i]
		if !((b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F') || (b >= '0' && b <= '9')) {
			return false
		}
	}
	return true
}

// IsLowerCase checks if the string is lowercase.
func IsLowerCase(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

// HasLowerCase checks if the string contains at least 1 lowercase.
func HasLowerCase(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

// IsUpperCase checks if the string is uppercase.
func IsUpperCase(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if unicode.IsLower(r) {
			return false
		}
	}
	return true
}

// HasUpperCase checks if the string contains as least 1 uppercase.
func HasUpperCase(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

// IsWhitespace checks the string only contains whitespace
func IsWhitespace(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// HasWhitespace checks if the string contains any whitespace
func HasWhitespace(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if unicode.IsSpace(r) {
			return true
		}
	}
	return false
}

// HasMultibyte checks if the string contains one or more multibyte chars.
func HasMultibyte(s string) bool {
	if s == "" {
		return false
	}

	l := len(s)
	for i := 0; i < l; i++ {
		if s[i] > unicode.MaxASCII {
			return true
		}
	}
	return false
}
