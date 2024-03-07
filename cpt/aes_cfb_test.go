package cpt

import (
	"fmt"
	"testing"

	"github.com/askasoft/pango/str"
)

func TestAesCFBEncryptString(t *testing.T) {
	ac := NewAesCFB("1234567890abedefg")
	for i := 1; i < 100; i++ {
		c := str.RandLetterNumbers(i)
		o, err := ac.EncryptString(c)
		if err != nil {
			t.Fatal(i, err)
		}
		fmt.Println(i, o)

		s, err := ac.DecryptString(o)
		if err != nil {
			t.Fatal(i, err)
		}

		if s != c {
			t.Fatalf("[%d] want %q, but %q", i, c, s)
		}
	}
}

func TestAesCFBEncryptData(t *testing.T) {
	ac := NewAesCFB("1234567890abcde", "0987654321654321")
	for i := 1; i < 100; i++ {
		c := str.RandLetterNumbers(i)
		o, err := ac.EncryptData(str.UnsafeBytes(c))
		if err != nil {
			t.Fatal(i, err)
		}
		fmt.Println(i, o)

		s, err := ac.DecryptData(o)
		if err != nil {
			t.Fatal(i, err)
		}

		if string(s) != c {
			t.Fatalf("[%d] want %q, but %q", i, c, string(s))
		}
	}
}
