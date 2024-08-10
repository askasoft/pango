package cpt

import (
	"testing"

	"github.com/askasoft/pango/str"
)

func TestAesCBCEncrypt(t *testing.T) {
	bits := []int{128, 192, 256}
	for _, bit := range bits {
		ac := NewAesCBC("1234567890abcde", bit)
		for i := 64; i <= 128; i++ {
			rs := str.RandLetterNumbers(i)
			es, err := ac.EncryptString(rs)
			if err != nil {
				t.Fatal(i, err)
			}

			// fmt.Printf("[%d] len(es) = %d\n", i, len(es))

			ds, err := ac.DecryptString(es)
			if err != nil {
				t.Fatal(i, err)
			}

			if ds != rs {
				t.Fatalf("[%d] want %q, but %q", i, rs, ds)
			}
		}
	}
}
