package str

import (
	crand "crypto/rand"
	mrand "math/rand"
)

// RandNumbers create a random number string
func RandNumbers(size int) string {
	return RandString(size, Numbers)
}

// RandLetterNumbers create a random letter number string
func RandLetterNumbers(size int) string {
	return RandString(size, LetterNumbers)
}

// RandLetters create a random letter string
func RandLetters(size int) string {
	return RandString(size, Letters)
}

// RandUpperLetters create a random upper letter string
func RandUpperLetters(size int) string {
	return RandString(size, UpperLetters)
}

// RandLowerLetters create a random lower letter string
func RandLowerLetters(size int) string {
	return RandString(size, LowerLetters)
}

// RandSymbols create a random letter string
func RandSymbols(size int) string {
	return RandString(size, Symbols)
}

// RandString create a random string by the input chars
// if chars is omitted, the LetterNumberSymbols is used
func RandString(size int, chars ...string) string {
	cs := LetterDigitSymbols
	if len(chars) > 0 {
		cs = chars[0]
	}

	n := len(cs)

	bs := make([]byte, size)
	if _, err := crand.Read(bs); err != nil {
		_, _ = mrand.Read(bs) //nolint: gosec
	}

	for i := 0; i < size; i++ {
		bs[i] = cs[int(bs[i])%n]
	}

	return UnsafeString(bs)
}
