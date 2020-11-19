package str

import (
	"math/rand"
	"time"
)

var seed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandDigits create a random digit string
func RandDigits(size int) string {
	return RandString(size, Digits)
}

// RandDigitLetters create a random digit letter string
func RandDigitLetters(size int) string {
	return RandString(size, DigitLetters)
}

// RandString create a random string by the input chars
func RandString(size int, chars ...string) string {
	cs := SymbolDigitLetters
	if len(chars) > 0 {
		cs = chars[0]
	}

	buf := make([]byte, size, size)
	for i := 0; i < size; i++ {
		buf[i] = cs[seed.Intn(len(cs))]
	}

	return string(buf)
}
