package str

import (
	"testing"
)

func TestCountRune(t *testing.T) {
	cs := []struct {
		w int
		s string
		b rune
	}{
		{0, "123", '0'},
		{1, "123", '2'},
		{2, "12一一3", '一'},
	}

	for i, c := range cs {
		a := CountRune(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] CountRune(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestCountAny(t *testing.T) {
	cs := []struct {
		w int
		s string
		b string
	}{
		{0, "123", "04"},
		{1, "123", "2"},
		{4, "12一一3うう", "一う"},
	}

	for i, c := range cs {
		a := CountAny(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] CountAny(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestContainsFold(t *testing.T) {
	cs := []struct {
		w bool
		s string
		b string
	}{
		{true, "ABCDEF", "abc"},
		{false, "ABCDEF", "Z"},
	}

	for i, c := range cs {
		a := ContainsFold(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] ContainsFold(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestIndexFold(t *testing.T) {
	cs := []struct {
		w int
		s string
		b string
	}{
		{0, "", ""},
		{-1, "", "a"},
		{0, "ABCDEF", ""},
		{0, "ABCDEF", "abc"},
		{1, "ABCDEF", "bc"},
		{4, "ABCDEF", "ef"},
		{6, "一BCDEF", "ef"},
		{4, "ABCD四F", "四f"},
	}

	for i, c := range cs {
		a := IndexFold(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] IndexFold(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestStartsWith(t *testing.T) {
	cs := []struct {
		w bool
		s string
		b string
	}{
		{true, "", ""},
		{true, "foobar", ""},
		{false, "", "f"},

		{true, "foobar", "foo"},
		{true, "一二三四五", "一"},

		{false, "f", "oo"},
		{false, "一", "一二三四五"},
		{false, "foobar", "oo"},
		{false, "一二三四五", "二"},
	}

	for i, c := range cs {
		a := StartsWith(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] StartsWith(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestEndsWith(t *testing.T) {
	cs := []struct {
		w bool
		s string
		b string
	}{
		{true, "", ""},
		{true, "foobar", ""},
		{false, "", "f"},

		{true, "foobar", "bar"},
		{true, "一二三四五", "四五"},

		{false, "f", "oo"},
		{false, "一", "一二三四五"},
		{false, "foobar", "oo"},
		{false, "一二三四五", "二"},
	}

	for i, c := range cs {
		a := EndsWith(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] EndsWith(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestLastIndexRune(t *testing.T) {
	cs := []struct {
		w int
		s string
		b rune
	}{
		{3, "aabbcc", 'b'},
		{9, "一一二二うう", '二'},
	}

	for i, c := range cs {
		a := LastIndexRune(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] LastIndexRune(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestLastIndexFold(t *testing.T) {
	cs := []struct {
		w int
		s string
		b string
	}{
		{0, "", ""},
		{-1, "", "a"},
		{-1, "", "foo"},
		{-1, "fo", "foo"},
		{0, "foo", "foo"},
		{0, "foo", "f"},
		{7, "oofofoofooo", "f"},
		{7, "oofofoofooo", "foo"},
		{9, "barfoobarfoo", "foo"},
		{3, "foo", ""},
		{2, "foo", "o"},
		{6, "abcABCabc", "A"},
		{6, "abcABCabc", "a"},
	}

	for i, c := range cs {
		a := LastIndexFold(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] LastIndexFold(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}
