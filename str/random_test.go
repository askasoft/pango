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

func testRandString(t *testing.T, fn string, f func(int) string, n int, cs string) bool {
	s := f(n)
	if len(s) != n {
		return false
	}

	for _, c := range s {
		if !ContainsRune(cs, c) {
			t.Errorf("RandNumbers(%d) = %s", n, s)
			return false
		}
	}
	return true
}
