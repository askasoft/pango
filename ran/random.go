package ran

import (
	"crypto/rand"
	"math"
	"math/big"

	"github.com/askasoft/pango/str"
)

func Read(bs []byte) {
	_, _ = rand.Read(bs)
}

func RandInt() int {
	return RandIntn(math.MaxInt)
}

func RandIntn(n int) int {
	if n <= math.MaxInt32 {
		return int(RandInt31n(int32(n)))
	}
	return int(RandInt63n(int64(n)))
}

func RandInt31() int32 {
	return RandInt31n(math.MaxInt32)
}

func RandInt31n(max int32) int32 {
	val, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int32(val.Int64())
}

func RandInt63() int64 {
	return RandInt63n(math.MaxInt64)
}

func RandInt63n(max int64) int64 {
	val, _ := rand.Int(rand.Reader, big.NewInt(max))
	return val.Int64()
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
	Read(bs)

	for i, b := range bs {
		bs[i] = cs[int(b)%n]
	}

	return str.UnsafeString(bs)
}
