package str

import (
	"testing"
)

func TestSubstrAfter(t *testing.T) {
	cs := []struct {
		w string
		s string
		b string
	}{
		{"", "", "c"},
		{"", "abc", "c"},
		{"", "abc", "d"},
		{"ot", "foot", "o"},
		{"bc", "abc", "a"},
		{"cba", "abcba", "b"},
		{"ba世界", "abc世界ba世界", "世界"},
	}

	for i, c := range cs {
		a := SubstrAfter(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] SubstrAfter(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestSubstrAfterByte(t *testing.T) {
	cs := []struct {
		w string
		s string
		b byte
	}{
		{"ot", "foot", 'o'},
		{"bc", "abc", 'a'},
		{"cba", "abcba", 'b'},
		{"", "abc", 'c'},
		{"", "abc", 'd'},
	}

	for i, c := range cs {
		a := SubstrAfterByte(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] SubstrAfterByte(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestSubstrAfterRune(t *testing.T) {
	cs := []struct {
		w string
		s string
		b rune
	}{
		{"", "abc", 'c'},
		{"", "abc", 'd'},
		{"ot", "foot", 'o'},
		{"bc", "abc", 'a'},
		{"cba", "abcba", 'b'},
		{"ba世界", "abc世界ba世界", '界'},
	}

	for i, c := range cs {
		a := SubstrAfterRune(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] SubstrAfterRune(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestSubstrAfterLast(t *testing.T) {
	cs := []struct {
		w string
		s string
		b string
	}{
		{"", "", ""},
		{"", "", "a"},
		{"", "foo", "b"},
		{"", "abc", "c"},
		{"", "", "d"},
		{"", "abc", ""},
		{"t", "foot", "o"},
		{"bc", "abc", "a"},
		{"a", "abcba", "b"},
		{"baz", "fooXXbarXXbaz", "XX"},
		{"ef", "ab世界cd世界ef", "世界"},
	}

	for i, c := range cs {
		a := SubstrAfterLast(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] SubstrAfterLast(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestSubstrAfterLastByte(t *testing.T) {
	cs := []struct {
		w string
		s string
		b byte
	}{
		{"", "", 'a'},
		{"", "foo", 'b'},
		{"t", "foot", 'o'},
		{"bc", "abc", 'a'},
		{"a", "abcba", 'b'},
		{"", "abc", 'c'},
		{"", "", 'd'},
	}

	for i, c := range cs {
		a := SubstrAfterLastByte(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] SubstrAfterLastByte(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestSubstrAfterLastRune(t *testing.T) {
	cs := []struct {
		w string
		s string
		b rune
	}{
		{"", "", 'a'},
		{"", "foo", 'b'},
		{"t", "foot", 'o'},
		{"bc", "abc", 'a'},
		{"a", "abcba", 'b'},
		{"", "abc", 'c'},
		{"", "", 'd'},
		{"bc", "a界世界bc", '界'},
	}

	for i, c := range cs {
		a := SubstrAfterLastRune(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] SubstrAfterLastRune(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestSubstrBefore(t *testing.T) {
	cs := []struct {
		w string
		s string
		b string
	}{
		{"", "", "a"},
		{"", "abc", "a"},
		{"f", "foot", "o"},
		{"a", "abcba", "b"},
		{"ab", "abc", "c"},
		{"abc", "abc", "d"},
		{"abc", "abc世界efg世界xyz", "世界"},
	}

	for i, c := range cs {
		a := SubstrBefore(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] SubstrBefore(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestSubstrBeforeByte(t *testing.T) {
	cs := []struct {
		w string
		s string
		b byte
	}{
		{"", "abc", 'a'},
		{"f", "foot", 'o'},
		{"a", "abcba", 'b'},
		{"ab", "abc", 'c'},
		{"abc", "abc", 'd'},
	}

	for i, c := range cs {
		a := SubstrBeforeByte(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] SubstrBeforeByte(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestSubstrBeforeRune(t *testing.T) {
	cs := []struct {
		w string
		s string
		b rune
	}{
		{"f", "foot", 'o'},
		{"", "abc", 'a'},
		{"a", "abcba", 'b'},
		{"ab", "abc", 'c'},
		{"abc", "abc", 'd'},
	}

	for i, c := range cs {
		a := SubstrBeforeRune(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] SubstrBeforeRune(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestSubstrBeforeLast(t *testing.T) {
	cs := []struct {
		w string
		s string
		b string
	}{
		{"", "", ""},
		{"", "", "a"},
		{"foo", "foo", "b"},
		{"fo", "foot", "o"},
		{"", "abc", "a"},
		{"abc", "abcba", "b"},
		{"ab", "abc", "c"},
		{"", "", "d"},
		{"abc", "abc", ""},
		{"abc", "abc", "d"},
		{"fooXXbar", "fooXXbarXXbaz", "XX"},
		{"abc世界efg", "abc世界efg世界xyz", "世界"},
	}

	for i, c := range cs {
		a := SubstrBeforeLast(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] SubstrBeforeLast(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestSubstrBeforeLastByte(t *testing.T) {
	cs := []struct {
		w string
		s string
		b byte
	}{
		{"", "", 'a'},
		{"foo", "foo", 'b'},
		{"fo", "foot", 'o'},
		{"", "abc", 'a'},
		{"abc", "abcba", 'b'},
		{"ab", "abc", 'c'},
		{"", "", 'd'},
		{"abc", "abc", 'd'},
	}

	for i, c := range cs {
		a := SubstrBeforeLastByte(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] SubstrBeforeLastByte(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestSubstrBeforeLastRune(t *testing.T) {
	cs := []struct {
		w string
		s string
		b rune
	}{
		{"", "", 'a'},
		{"foo", "foo", 'b'},
		{"fo", "foot", 'o'},
		{"", "abc", 'a'},
		{"abc", "abcba", 'b'},
		{"ab", "abc", 'c'},
		{"", "", 'd'},
		{"abc", "abc", 'd'},
		{"abc世界", "abc世界世ba", '世'},
	}

	for i, c := range cs {
		a := SubstrBeforeLastRune(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] SubstrBeforeLastRune(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}
