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

func TestSubstrAfterAny(t *testing.T) {
	cs := []struct {
		w string
		s string
		b string
	}{
		{"", "", "c"},
		{"", "abc", "zc"},
		{"", "abc", "zd"},
		{"ot", "foot", "zo"},
		{"bc", "abc", "za"},
		{"cba", "abcba", "zb"},
		{"界ba世界", "abc世界ba世界", "世界"},
	}

	for i, c := range cs {
		a := SubstrAfterAny(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] SubstrAfterAny(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
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

func TestSubstrAfterLastAny(t *testing.T) {
	cs := []struct {
		w string
		s string
		b string
	}{
		{"", "", ""},
		{"", "", "a"},
		{"", "foo", "b"},
		{"", "abc", "zc"},
		{"", "abc", ""},
		{"t", "foot", "zo"},
		{"bc", "abc", "za"},
		{"a", "abcba", "zb"},
		{"baz", "fooXXbarXXbaz", "XX"},
		{"ef", "ab世界cd世界ef", "世界"},
	}

	for i, c := range cs {
		a := SubstrAfterLastAny(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] SubstrAfterLastAny(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
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

func TestSubstrBeforeAny(t *testing.T) {
	cs := []struct {
		w string
		s string
		b string
	}{
		{"", "", "a"},
		{"", "abc", "az"},
		{"f", "foot", "zo"},
		{"a", "abcba", "zb"},
		{"ab", "abc", "zc"},
		{"abc", "abc", "zd"},
		{"abc", "abc世界efg世界xyz", "世界"},
	}

	for i, c := range cs {
		a := SubstrBeforeAny(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] SubstrBeforeAny(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
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

func TestSubstrBeforeLastAny(t *testing.T) {
	cs := []struct {
		w string
		s string
		b string
	}{
		{"", "", ""},
		{"", "", "a"},
		{"foo", "foo", "zb"},
		{"fo", "foot", "zo"},
		{"", "abc", "az"},
		{"abc", "abcba", "zb"},
		{"ab", "abc", "zc"},
		{"abc", "abc", ""},
		{"abc", "abc", "zd"},
		{"fooXXbarX", "fooXXbarXXbaz", "XX"},
		{"abc世界efg世", "abc世界efg世界xyz", "世界"},
	}

	for i, c := range cs {
		a := SubstrBeforeLastAny(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] SubstrBeforeLastAny(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
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

func TestCutPrefix(t *testing.T) {
	var cutPrefixTests = []struct {
		s, sep string
		after  string
		found  bool
	}{
		{"abc", "a", "bc", true},
		{"abc", "abc", "", true},
		{"abc", "", "abc", true},
		{"abc", "d", "abc", false},
		{"", "d", "", false},
		{"", "", "", true},
	}

	for _, tt := range cutPrefixTests {
		if after, found := CutPrefix(tt.s, tt.sep); after != tt.after || found != tt.found {
			t.Errorf("CutPrefix(%q, %q) = %q, %v, want %q, %v", tt.s, tt.sep, after, found, tt.after, tt.found)
		}
	}
}

func TestCutSuffix(t *testing.T) {
	var cutSuffixTests = []struct {
		s, sep string
		after  string
		found  bool
	}{
		{"abc", "bc", "a", true},
		{"abc", "abc", "", true},
		{"abc", "", "abc", true},
		{"abc", "d", "abc", false},
		{"", "d", "", false},
		{"", "", "", true},
	}

	for _, tt := range cutSuffixTests {
		if after, found := CutSuffix(tt.s, tt.sep); after != tt.after || found != tt.found {
			t.Errorf("CutSuffix(%q, %q) = %q, %v, want %q, %v", tt.s, tt.sep, after, found, tt.after, tt.found)
		}
	}
}

func TestCut(t *testing.T) {
	cutTests := []struct {
		s, sep        string
		before, after string
		found         bool
	}{
		{"abc", "b", "a", "c", true},
		{"abc", "a", "", "bc", true},
		{"abc", "c", "ab", "", true},
		{"abc", "abc", "", "", true},
		{"abc", "", "", "abc", true},
		{"abc", "d", "abc", "", false},
		{"", "d", "", "", false},
		{"", "", "", "", true},
	}

	for _, tt := range cutTests {
		if before, after, found := Cut(tt.s, tt.sep); before != tt.before || after != tt.after || found != tt.found {
			t.Errorf("Cut(%q, %q) = %q, %q, %v, want %q, %q, %v", tt.s, tt.sep, before, after, found, tt.before, tt.after, tt.found)
		}
	}
}

func TestCutByte(t *testing.T) {
	cutTests := []struct {
		s             string
		sep           byte
		before, after string
		found         bool
	}{
		{"abc", 'b', "a", "c", true},
		{"abc", 'a', "", "bc", true},
		{"abc", 'c', "ab", "", true},
		{"abc", 'd', "abc", "", false},
		{"", 'd', "", "", false},
	}

	for _, tt := range cutTests {
		if before, after, found := CutByte(tt.s, tt.sep); before != tt.before || after != tt.after || found != tt.found {
			t.Errorf("CutByte(%q, %q) = %q, %q, %v, want %q, %q, %v", tt.s, tt.sep, before, after, found, tt.before, tt.after, tt.found)
		}
	}
}

func TestCutRune(t *testing.T) {
	cutTests := []struct {
		s             string
		sep           rune
		before, after string
		found         bool
	}{
		{"abc", 'b', "a", "c", true},
		{"abc", 'a', "", "bc", true},
		{"abc", 'c', "ab", "", true},
		{"abc", 'd', "abc", "", false},
		{"", 'd', "", "", false},
	}

	for _, tt := range cutTests {
		if before, after, found := CutRune(tt.s, tt.sep); before != tt.before || after != tt.after || found != tt.found {
			t.Errorf("CutRune(%q, %q) = %q, %q, %v, want %q, %q, %v", tt.s, tt.sep, before, after, found, tt.before, tt.after, tt.found)
		}
	}
}

func TestLastCut(t *testing.T) {
	cutTests := []struct {
		s, sep        string
		before, after string
		found         bool
	}{
		{"abcabc", "b", "abca", "c", true},
		{"abcabc", "a", "abc", "bc", true},
		{"abcabc", "c", "abcab", "", true},
		{"abcabc", "abc", "abc", "", true},
		{"abcabc", "", "abcabc", "", true},
		{"abcabc", "d", "abcabc", "", false},
		{"", "d", "", "", false},
		{"", "", "", "", true},
	}

	for _, tt := range cutTests {
		if before, after, found := LastCut(tt.s, tt.sep); before != tt.before || after != tt.after || found != tt.found {
			t.Errorf("LastCut(%q, %q) = %q, %q, %v, want %q, %q, %v", tt.s, tt.sep, before, after, found, tt.before, tt.after, tt.found)
		}
	}
}

func TestLastCutByte(t *testing.T) {
	cutTests := []struct {
		s             string
		sep           byte
		before, after string
		found         bool
	}{
		{"abcabc", 'b', "abca", "c", true},
		{"abcabc", 'a', "abc", "bc", true},
		{"abcabc", 'c', "abcab", "", true},
		{"abcabc", 'd', "abcabc", "", false},
		{"", 'd', "", "", false},
	}

	for _, tt := range cutTests {
		if before, after, found := LastCutByte(tt.s, tt.sep); before != tt.before || after != tt.after || found != tt.found {
			t.Errorf("LastCutByte(%q, %q) = %q, %q, %v, want %q, %q, %v", tt.s, tt.sep, before, after, found, tt.before, tt.after, tt.found)
		}
	}
}

func TestLastCutRune(t *testing.T) {
	cutTests := []struct {
		s             string
		sep           rune
		before, after string
		found         bool
	}{
		{"abcabc", 'b', "abca", "c", true},
		{"abcabc", 'a', "abc", "bc", true},
		{"abcabc", 'c', "abcab", "", true},
		{"abcabc", 'd', "abcabc", "", false},
		{"", 'd', "", "", false},
	}

	for _, tt := range cutTests {
		if before, after, found := LastCutRune(tt.s, tt.sep); before != tt.before || after != tt.after || found != tt.found {
			t.Errorf("LastCutRune(%q, %q) = %q, %q, %v, want %q, %q, %v", tt.s, tt.sep, before, after, found, tt.before, tt.after, tt.found)
		}
	}
}
