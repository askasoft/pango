package cpt

import (
	"testing"

	"github.com/askasoft/pango/str"
)

func TestAesCBCEncryptString(t *testing.T) {
	ac := NewAesCBC("1234567890123456")
	for i := 1; i < 100; i++ {
		c := str.RandLetterNumbers(i)
		o, err := ac.EncryptString(c)
		if err != nil {
			t.Fatal(i, err)
		}

		s, err := ac.DecryptString(o)
		if err != nil {
			t.Fatal(i, err)
		}

		if s != c {
			t.Fatalf("[%d] want %q, but %q", i, c, s)
		}
	}
}

func TestAesCBCEncryptData(t *testing.T) {
	ac := NewAesCBC("1234567890123456", "0987654321654321")
	for i := 1; i < 100; i++ {
		c := str.RandLetterNumbers(i)
		o, err := ac.EncryptData(str.UnsafeBytes(c))
		if err != nil {
			t.Fatal(i, err)
		}

		s, err := ac.DecryptData(o)
		if err != nil {
			t.Fatal(i, err)
		}

		if string(s) != c {
			t.Fatalf("[%d] want %q, but %q", i, c, string(s))
		}
	}
}
