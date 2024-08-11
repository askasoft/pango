package cpt

import (
	"fmt"
	"testing"

	"github.com/askasoft/pango/str"
)

func TestAesCBCEncrypt(t *testing.T) {
	bits := []int{128, 192, 256}
	for _, bit := range bits {
		ae := NewAesCBCEncryptor("1234567890abcde", bit)
		ad := NewAesCBCDecryptor("1234567890abcde", bit)
		for i := 64; i <= 128; i++ {
			rs := str.RandLetterNumbers(i)
			es, err := ae.EncryptString(rs)
			if err != nil {
				t.Fatal(i, err)
			}

			// fmt.Printf("[%d] len(es) = %d\n", i, len(es))

			ds, err := ad.DecryptString(es)
			if err != nil {
				t.Fatal(i, err)
			}

			if ds != rs {
				t.Fatalf("[%d] want %q, but %q", i, rs, ds)
			}
		}
	}
}

func TestAesCFBEncrypt(t *testing.T) {
	bits := []int{128, 192, 256}
	for _, bit := range bits {
		ae := NewAesCFBEncryptor("1234567890abedefg", bit)
		ad := NewAesCFBDecryptor("1234567890abedefg", bit)
		for i := 1; i <= 128; i++ {
			c := str.RandLetterNumbers(i)
			o, err := ae.EncryptString(c)
			if err != nil {
				t.Fatal(i, err)
			}
			fmt.Println(i, o)

			s, err := ad.DecryptString(o)
			if err != nil {
				t.Fatal(i, err)
			}

			if s != c {
				t.Fatalf("[%d] want %q, but %q", i, c, s)
			}
		}
	}
}
