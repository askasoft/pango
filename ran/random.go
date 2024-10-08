package ran

import (
	crand "crypto/rand"
	"math/big"
	mrand "math/rand"

	"github.com/askasoft/pango/str"
)

func RandInt() int {
	return int(RandInt63())
}

func RandInt31() int32 {
	val, err := crand.Int(crand.Reader, big.NewInt(2147483647))
	if err == nil {
		return int32(val.Int64())
	}
	return mrand.Int31()
}

func RandInt63() int64 {
	val, err := crand.Int(crand.Reader, big.NewInt(9223372036854775807))
	if err == nil {
		return val.Int64()
	}
	return mrand.Int63()
}

// RandNumbers create a random number string
func RandNumbers(size int) string {
	return RandString(size, str.Numbers)
}

// RandLetterNumbers create a random letter number string
func RandLetterNumbers(size int) string {
	return RandString(size, str.LetterNumbers)
}

// RandLetters create a random letter string
func RandLetters(size int) string {
	return RandString(size, str.Letters)
}

// RandUpperLetters create a random upper letter string
func RandUpperLetters(size int) string {
	return RandString(size, str.UpperLetters)
}

// RandLowerLetters create a random lower letter string
func RandLowerLetters(size int) string {
	return RandString(size, str.LowerLetters)
}

// RandSymbols create a random letter string
func RandSymbols(size int) string {
	return RandString(size, str.Symbols)
}

// RandString create a random string by the input chars
// if chars is omitted, the LetterNumberSymbols is used
func RandString(size int, chars ...string) string {
	cs := str.LetterDigitSymbols
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

	return str.UnsafeString(bs)
}
