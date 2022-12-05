package cpt

import (
	"math/rand"
	"testing"

	"github.com/pandafw/pango/str"
)

func TestSalt(t *testing.T) {
	for i := 0; i < 100; i++ {
		salt := str.RandLetterNumbers(rand.Intn(100))
		src := str.RandLetterNumbers(i)

		salted := Salt(SecretChars, salt, src)
		if len(src) != len(salted) {
			t.Fatalf("[%d] len(src)=%d, len(salted)=%d", i, len(src), len(salted))
		}

		unsalted := Unsalt(SecretChars, salt, salted)
		if src != unsalted {
			t.Fatalf("[%d] unsalted = %q, want %q", i, unsalted, src)
		}
	}
}
