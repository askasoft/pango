package str

import (
	"math/rand"
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

// RandString create a random string by the input chars
// if chars is omitted, the LetterNumberSymbols is used
func RandString(size int, chars ...string) string {
	seed := LetterNumberSymbols
	if len(chars) > 0 {
		seed = chars[0]
	}

	n := len(seed)
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = seed[rand.Intn(n)] //nolint: gosec
	}

	return string(buf)
}
