package str

import (
	"math/rand"
	"time"
)

var seed = rand.New(rand.NewSource(time.Now().UnixNano())) //nolint: gosec

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
	cs := LetterNumberSymbols
	if len(chars) > 0 {
		cs = chars[0]
	}

	n := len(cs)
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = cs[seed.Intn(n)]
	}

	return string(buf)
}
