package str

import (
	"math/rand"
	"time"
)

var seed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

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
func RandString(size int, chars string) string {
	if chars == "" {
		chars = LetterNumberSymbols
	}

	n := len(chars)
	buf := make([]byte, size, size)
	for i := 0; i < size; i++ {
		buf[i] = chars[seed.Intn(n)]
	}

	return string(buf)
}
