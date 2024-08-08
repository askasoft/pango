package funcs

import (
	"errors"
	"fmt"
	"testing"
)

func TestComma(t *testing.T) {
	cs := []struct {
		n any
		d int
		w string
		e error
	}{
		{1234, 0, "1,234", nil},
		{2345.0, 3, "2,345", nil},
		{2345.1, 3, "2,345.1", nil},
		{2345.12, 3, "2,345.12", nil},
		{2345.1234, 3, "2,345.123", nil},
		{2345.1235, 3, "2,345.124", nil},
		{"1.1", 0, "", errors.New("Comma: unknown type for '1.1' (string)")},
	}

	for i, c := range cs {
		r, e := Comma(c.n, c.d)
		if c.w != r || fmt.Sprint(c.e) != fmt.Sprint(e) {
			t.Errorf("[%d] Comma(%v, %d) = (%T, %v, %v), want: (%T, %v, %v)", i, c.n, c.d, r, r, e, c.w, c.w, c.e)
		}
	}
}

func TestHumanSize(t *testing.T) {
	cs := []struct {
		a any
		w string
		e error
	}{
		{1234, "1.21 KB", nil},
		{2345.0, "2.29 KB", nil},
		{"1.1", "", errors.New("HumanSize: unknown type for '1.1' (string)")},
	}

	for i, c := range cs {
		r, e := HumanSize(c.a)
		if c.w != r || fmt.Sprint(c.e) != fmt.Sprint(e) {
			t.Errorf("[%d] HumanSize(%v) = (%T, %v, %v), want: (%T, %v, %v)", i, c.a, r, r, e, c.w, c.w, c.e)
		}
	}
}
