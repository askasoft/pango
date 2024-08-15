package str

import (
	"testing"
)

func TestRandNumbers(t *testing.T) {
	for i := 10; i < 100; i++ {
		testRandString(t, "RandNumbers", RandNumbers, i, Numbers)
	}
}

func TestRandLetterNumbers(t *testing.T) {
	for i := 10; i < 100; i++ {
		testRandString(t, "RandLetterNumbers", RandLetterNumbers, i, LetterNumbers)
	}
}

func TestRandLetters(t *testing.T) {
	for i := 10; i < 100; i++ {
		testRandString(t, "RandLetters", RandLetters, i, Letters)
	}
}
func TestRandUpperLetters(t *testing.T) {
	for i := 10; i < 100; i++ {
		testRandString(t, "RandUpperLetters", RandUpperLetters, i, UpperLetters)
	}
}

func TestRandLowerLetters(t *testing.T) {
	for i := 10; i < 100; i++ {
		testRandString(t, "RandLowerLetters", RandLowerLetters, i, LowerLetters)
	}
}

func TestRandSymbols(t *testing.T) {
	for i := 10; i < 100; i++ {
		testRandString(t, "RandSymbols", RandSymbols, i, Symbols)
	}
}

func TestRandString(t *testing.T) {
	for i := 10; i < 100; i++ {
		testRandString(t, "RandString", func(z int) string { return RandString(z) }, i, LetterDigitSymbols)
	}
}

func testRandString(t *testing.T, fn string, f func(int) string, n int, cs string) bool {
	s := f(n)
	if len(s) != n {
		return false
	}

	for _, c := range s {
		if !ContainsRune(cs, c) {
			t.Errorf("%s(%d) = %s", fn, n, s)
			return false
		}
	}
	return true
}
