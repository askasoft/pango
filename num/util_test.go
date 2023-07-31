package num

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"testing/quick"
)

func TestStripTrailingDigits(t *testing.T) {
	err := quick.Check(func(s string, digits int) bool {
		stripped := stripTrailingDigits(s, digits)

		// A stripped string will always be a prefix of its original string
		if !strings.HasPrefix(s, stripped) {
			return false
		}

		if strings.ContainsRune(s, '.') {
			// If there is a dot, the part on the left of the dot will never change
			a := strings.Split(s, ".")
			b := strings.Split(stripped, ".")
			if a[0] != b[0] {
				return false
			}
		} else {
			// If there's no dot in the input, the output will always be the same as the input.
			if stripped != s {
				return false
			}
		}

		return true
	}, &quick.Config{
		MaxCount: 10000,
		Values: func(v []reflect.Value, r *rand.Rand) {
			rdigs := func(n int) string {
				digs := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
				var rv []rune
				for i := 0; i < n; i++ {
					rv = append(rv, digs[r.Intn(len(digs))])
				}
				return string(rv)
			}

			ls := r.Intn(20)
			rs := r.Intn(20)
			jc := "."
			if rs == 0 {
				jc = ""
			}
			s := rdigs(ls) + jc + rdigs(rs)
			digits := r.Intn(len(s) + 1)

			v[0] = reflect.ValueOf(s)
			v[1] = reflect.ValueOf(digits)
		},
	})

	if err != nil {
		t.Error(err)
	}
}
