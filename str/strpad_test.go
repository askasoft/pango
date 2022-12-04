package str

import (
	"testing"
)

func TestRepeatByte(t *testing.T) {
	cs := []struct {
		in    byte
		out   string
		count int
	}{
		{'-', "", 0},
		{'-', "-", 1},
		{'-', "----------", 10},
	}

	for i, tt := range cs {
		a := RepeatByte(tt.in, tt.count)
		if a != tt.out {
			t.Errorf("[%d] RepeatByte(%q, %d) = %q; want %q", i, tt.in, tt.count, a, tt.out)
			continue
		}
	}
}

func TestRepeatRune(t *testing.T) {
	cs := []struct {
		in    rune
		out   string
		count int
	}{
		{'-', "", 0},
		{'-', "-", 1},
		{'-', "----------", 10},
	}

	for i, tt := range cs {
		a := RepeatRune(tt.in, tt.count)
		if a != tt.out {
			t.Errorf("[%d] RepeatRune(%q, %d) = %q; want %q", i, tt.in, tt.count, a, tt.out)
			continue
		}
	}
}

func TestPadCenterByte(t *testing.T) {
	cs := []struct {
		w string
		s string
		n int
		b byte
	}{
		{"    ", "", 4, ' '},
		{"ab", "ab", 0, ' '},
		{"ab", "ab", -1, ' '},
		{"ab", "ab", 1, ' '},
		{"    ", "", 4, ' '},
		{" ab ", "ab", 4, ' '},
		{"abcd", "abcd", 2, ' '},
		{" a  ", "a", 4, ' '},
		{"  a  ", "a", 5, ' '},
		{"xxaxx", "a", 5, 'x'},
	}

	for i, c := range cs {
		a := PadCenterByte(c.s, c.n, c.b)
		if a != c.w {
			t.Errorf("[%d] PadCenterByte(%q, %d, %q) = %v, want %v", i, c.s, c.n, c.b, a, c.w)
		}
	}
}

func TestPadCenterRune(t *testing.T) {
	cs := []struct {
		w string
		s string
		n int
		r rune
	}{
		{"    ", "", 4, ' '},
		{"ab", "ab", 0, ' '},
		{"ab", "ab", -1, ' '},
		{"ab", "ab", 1, ' '},
		{"    ", "", 4, ' '},
		{" ab ", "ab", 4, ' '},
		{"abcd", "abcd", 2, ' '},
		{" a  ", "a", 4, ' '},
		{"  a  ", "a", 5, ' '},
		{"xxaxx", "a", 5, 'x'},
	}

	for i, c := range cs {
		a := PadCenterRune(c.s, c.n, c.r)
		if a != c.w {
			t.Errorf("[%d] PadCenterRune(%q, %d, %q) = %v, want %v", i, c.s, c.n, c.r, a, c.w)
		}
	}
}

func TestPadCenter(t *testing.T) {
	cs := []struct {
		w string
		s string
		n int
		p string
	}{
		{"    ", "", 4, " "},
		{"ab", "ab", 0, " "},
		{"ab", "ab", -1, " "},
		{"ab", "ab", 1, " "},
		{"    ", "", 4, " "},
		{" ab ", "ab", 4, " "},
		{"abcd", "abcd", 2, " "},
		{" a  ", "a", 4, " "},
		{"yayz", "a", 4, "yz"},
		{"yzyayzy", "a", 7, "yz"},
		{"abc", "abc", 7, ""},
	}

	for i, c := range cs {
		a := PadCenter(c.s, c.n, c.p)
		if a != c.w {
			t.Errorf("[%d] PadCenter(%q, %d, %q) = %v, want %v", i, c.s, c.n, c.p, a, c.w)
		}
	}
}

func TestPadLeftByte(t *testing.T) {
	cs := []struct {
		w string
		s string
		n int
		b byte
	}{
		{"     ", "", 5, ' '},
		{"  abc", "abc", 5, ' '},
		{"xxabc", "abc", 5, 'x'},
		{"\xff\xffabc", "abc", 5, '\xff'},
		{"abc", "abc", 2, ' '},
	}

	for i, c := range cs {
		a := PadLeftByte(c.s, c.n, c.b)
		if a != c.w {
			t.Errorf("[%d] PadLeftByte(%q, %d, %q) = %v, want %v", i, c.s, c.n, c.b, a, c.w)
		}
	}
}

func TestPadLeftRune(t *testing.T) {
	cs := []struct {
		w string
		s string
		n int
		r rune
	}{
		{"     ", "", 5, ' '},
		{"  abc", "abc", 5, ' '},
		{"xxabc", "abc", 5, 'x'},
		{"\uffff\uffffabc", "abc", 5, '\uffff'},
		{"abc", "abc", 2, ' '},
	}

	for i, c := range cs {
		a := PadLeftRune(c.s, c.n, c.r)
		if a != c.w {
			t.Errorf("[%d] PadLeftRune(%q, %d, %q) = %v, want %v", i, c.s, c.n, c.r, a, c.w)
		}
	}
}

func TestPadLeft(t *testing.T) {
	cs := []struct {
		w string
		s string
		n int
		p string
	}{
		{"     ", "", 5, " "},
		{"-+-+abc", "abc", 7, "-+"},
		{"-+~abc", "abc", 6, "-+~"},
		{"-+abc", "abc", 5, "-+~"},
		{"abc", "abc", 2, " "},
		{"abc", "abc", -1, " "},
		{"abc", "abc", 5, ""},
		{"a一a一aabc", "abc", 8, "a一"},
	}

	for i, c := range cs {
		a := PadLeft(c.s, c.n, c.p)
		if a != c.w {
			t.Errorf("[%d] PadLeft(%q, %d, %q) = %v, want %v", i, c.s, c.n, c.p, a, c.w)
		}
	}
}

func TestPadRightByte(t *testing.T) {
	cs := []struct {
		w string
		s string
		n int
		b byte
	}{
		{"     ", "", 5, ' '},
		{"abc  ", "abc", 5, ' '},
		{"abc", "abc", 2, ' '},
		{"abc", "abc", -1, ' '},
		{"abcxx", "abc", 5, 'x'},
	}

	for i, c := range cs {
		a := PadRightByte(c.s, c.n, c.b)
		if a != c.w {
			t.Errorf("[%d] PadRightByte(%q, %d, %q) = %v, want %v", i, c.s, c.n, c.b, a, c.w)
		}
	}
}

func TestPadRightRune(t *testing.T) {
	cs := []struct {
		w string
		s string
		n int
		r rune
	}{
		{"     ", "", 5, ' '},
		{"abc  ", "abc", 5, ' '},
		{"abc", "abc", 2, ' '},
		{"abc", "abc", -1, ' '},
		{"abcxx", "abc", 5, 'x'},
	}

	for i, c := range cs {
		a := PadRightRune(c.s, c.n, c.r)
		if a != c.w {
			t.Errorf("[%d] PadRightRune(%q, %d, %q) = %v, want %v", i, c.s, c.n, c.r, a, c.w)
		}
	}
}

func TestPadRight(t *testing.T) {
	cs := []struct {
		w string
		s string
		n int
		p string
	}{
		{"     ", "", 5, " "},
		{"abc-+-+", "abc", 7, "-+"},
		{"abc-+~", "abc", 6, "-+~"},
		{"abc-+", "abc", 5, "-+~"},
		{"abc", "abc", 2, " "},
		{"abc", "abc", -1, " "},
		{"abc", "abc", 5, ""},
		{"abca一a一a", "abc", 8, "a一"},
	}

	for i, c := range cs {
		a := PadRight(c.s, c.n, c.p)
		if a != c.w {
			t.Errorf("[%d] PadRight(%q, %d, %q) = %v, want %v", i, c.s, c.n, c.p, a, c.w)
		}
	}
}
