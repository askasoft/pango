package num

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"testing/quick"
)

func TestAtoi(t *testing.T) {
	cs := []struct {
		w int
		s string
		n []int
	}{
		{1, "1", []int{}},
		{-1, "a", []int{-1}},
	}

	for i, c := range cs {
		a := Atoi(c.s, c.n...)
		if a != c.w {
			t.Errorf("[%d] Atoi(%s, %v) = %v, want %v", i, c.s, c.n, a, c.w)
		}
	}
}

func TestAtol(t *testing.T) {
	cs := []struct {
		w int64
		s string
		n []int64
	}{
		{1, "1", []int64{}},
		{-1, "a", []int64{-1}},
	}

	for i, c := range cs {
		a := Atol(c.s, c.n...)
		if a != c.w {
			t.Errorf("[%d] Atol(%s, %v) = %v, want %v", i, c.s, c.n, a, c.w)
		}
	}
}

func TestAtof(t *testing.T) {
	cs := []struct {
		w float64
		s string
		n []float64
	}{
		{1, "1", []float64{}},
		{-1, "a", []float64{-1}},
	}

	for i, c := range cs {
		a := Atof(c.s, c.n...)
		if a != c.w {
			t.Errorf("[%d] Atof(%s, %v) = %v, want %v", i, c.s, c.n, a, c.w)
		}
	}
}

func TestFtoa(t *testing.T) {
	cs := []struct {
		w string
		n float64
	}{
		{"200", 200},
		{"2", 2},
		{"2.2", 2.2},
		{"2.02", 2.02},
		{"200.02", 200.02},
	}

	for i, c := range cs {
		a := Ftoa(c.n)
		if a != c.w {
			t.Errorf("[%d] Ftoa(%f) = %v, want %v", i, c.n, a, c.w)
		}
	}
}

func TestFtoaWithDigits(t *testing.T) {
	cs := []struct {
		w string
		n float64
		d int
	}{
		{"1", 1.23, 0},
		{"1.2", 1.23, 1},
		{"1.23", 1.23, 2},
		{"1.23", 1.23, 3},
	}

	for i, c := range cs {
		a := FtoaWithDigits(c.n, c.d)
		if a != c.w {
			t.Errorf("[%d] FtoaWithDigits(%f, %d) = %v, want %v", i, c.n, c.d, a, c.w)
		}
	}
}

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
