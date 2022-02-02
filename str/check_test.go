package str

import "testing"

func TestIsASCII(t *testing.T) {
	t.Parallel()

	cs := []struct {
		w bool
		s string
	}{
		{false, ""},
		{false, "ｆｏｏbar"},
		{false, "ｘｙｚ０９８"},
		{false, "１２３456"},
		{false, "ｶﾀｶﾅ"},
		{true, "foobar"},
		{true, "0987654321"},
		{true, "test@example.com"},
		{true, "1234abcDEF"},
		{false, ""},
	}

	for i, c := range cs {
		a := IsASCII(c.s)
		if a != c.w {
			t.Errorf("[%d] IsASCII(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsPrintableASCII(t *testing.T) {
	t.Parallel()

	cs := []struct {
		w bool
		s string
	}{
		{false, ""},
		{false, "ｆｏｏbar"},
		{false, "ｘｙｚ０９８"},
		{false, "１２３456"},
		{false, "ｶﾀｶﾅ"},
		{true, "foobar"},
		{true, "0987654321"},
		{true, "test@example.com"},
		{true, "1234abcDEF"},
		{false, "newline\n"},
		{false, "\x19test\x7F"},
	}

	for i, c := range cs {
		a := IsPrintableASCII(c.s)
		if a != c.w {
			t.Errorf("[%d] IsPrintableASCII(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}
