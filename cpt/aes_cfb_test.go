package cpt

import (
	"fmt"
	"testing"

	"github.com/askasoft/pango/str"
)

func TestAes128CFBEncrypt(t *testing.T) {
	bits := []int{128, 192, 256}
	for _, bit := range bits {
		ac := NewAesCFB("1234567890abedefg", bit)
		for i := 1; i <= 128; i++ {
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
}
