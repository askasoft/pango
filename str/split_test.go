package str

import (
	"reflect"
	"testing"
)

func TestSplitLength(t *testing.T) {
	cs := []struct {
		s string
		n int
		w []string
	}{
		{"", 0, []string{""}},
		{"a", 0, []string{"a"}},
		{"ab", 0, []string{"ab"}},
		{"", 1, []string{""}},
		{"a", 1, []string{"a"}},
		{"abc", 1, []string{"a", "b", "c"}},
		{"あ", 1, []string{"あ"}},
		{"あaい", 1, []string{"あ", "a", "い"}},
		{"あい", 1, []string{"あ", "い"}},
		{"あい", 2, []string{"あ", "い"}},
		{"あいうえお", 10, []string{"あいう", "えお"}},
		{"あいうえおか", 9, []string{"あいう", "えおか"}},
	}

	for i, c := range cs {
		a := SplitLength(c.s, c.n)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] SplitLength(%q, %d) = %v, want %v", i, c.s, c.n, a, c.w)
		}
	}
}

func TestSplitCount(t *testing.T) {
	cs := []struct {
		s string
		n int
		w []string
	}{
		{"", 0, []string{""}},
		{"a", 0, []string{"a"}},
		{"ab", 0, []string{"ab"}},
		{"", 1, []string{""}},
		{"a", 1, []string{"a"}},
		{"abc", 1, []string{"a", "b", "c"}},
		{"あ", 1, []string{"あ"}},
		{"あaい", 1, []string{"あ", "a", "い"}},
		{"あい", 1, []string{"あ", "い"}},
		{"あい", 2, []string{"あい"}},
		{"あいうえお", 3, []string{"あいう", "えお"}},
		{"あいうえおか", 3, []string{"あいう", "えおか"}},
	}

	for i, c := range cs {
		a := SplitCount(c.s, c.n)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] SplitCount(%q, %d) = %v, want %v", i, c.s, c.n, a, c.w)
		}
	}
}

func TestSplitAny(t *testing.T) {
	cs := []struct {
		w []string
		s string
		b string
	}{
		{[]string{""}, "", "c"},
		{[]string{""}, "", ".c"},
		{[]string{"http://a", "b-", ""}, "http://a.b-c", ".c"},
		{[]string{"http", "", "", "a", "b", "c"}, "http://a.b.c", ":/."},
		{[]string{"http", "", "", "一", "二", "三"}, "http://一.二.三", ":/."},
		{[]string{"http", "", "", "一", "二", "三"}, "http://一。二。三", ":/。."},
	}

	for i, c := range cs {
		a := SplitAny(c.s, c.b)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] SplitAny(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestSplitByte(t *testing.T) {
	cs := []struct {
		w []string
		s string
		b byte
	}{
		{[]string{}, "", 'c'},
		{[]string{"http://a", "b", "c", ""}, "http://a.b.c.", '.'},
		{[]string{"http:", "", "a.b.c"}, "http://a.b.c", '/'},
	}

	for i, c := range cs {
		a := SplitByte(c.s, c.b)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] SplitByte(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestSplitRune(t *testing.T) {
	cs := []struct {
		w []string
		s string
		b rune
	}{
		{[]string{}, "", 'c'},
		{[]string{"http://a", "b", "c", ""}, "http://a.b.c.", '.'},
		{[]string{"http:", "", "a.b.c"}, "http://a.b.c", '/'},
		{[]string{"http://一", "二", "三"}, "http://一.二.三", '.'},
		{[]string{"http://一", "二", "三"}, "http://一。二。三", '。'},
	}

	for i, c := range cs {
		a := SplitRune(c.s, c.b)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] SplitRune(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestSplitFuncIter(t *testing.T) {
	cs := []struct {
		w []string
		s string
		b rune
	}{
		{[]string{}, "", 'c'},
		{[]string{"http://a", "b", "c", ""}, "http://a.b.c.", '.'},
		{[]string{"http:", "", "a.b.c", ""}, "http://a.b.c/", '/'},
		{[]string{"http://一", "二", "三"}, "http://一.二.三", '.'},
		{[]string{"http://一", "二", "三"}, "http://一。二。三", '。'},
	}

	for i, c := range cs {
		a := []string{}
		SplitFuncIter(c.s, func(r rune) bool { return c.b == r }, func(s string) bool {
			a = append(a, s)
			return true
		})

		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] SplitFuncIter(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestFieldsAny(t *testing.T) {
	cs := []struct {
		w []string
		s string
		b string
	}{
		{[]string{}, "", "c"},
		{[]string{}, "", ".c"},
		{[]string{"http://a", "b"}, "http://a.b.c", ".c"},
		{[]string{"http", "a", "b", "c"}, "http://a.b.c", ":/."},
		{[]string{"http", "一", "二", "三"}, "http://一.二.三", ":/."},
		{[]string{"http", "一", "二", "三"}, "http://一。二。三", ":/。."},
	}

	for i, c := range cs {
		a := FieldsAny(c.s, c.b)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] FieldsAny(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestFieldsByte(t *testing.T) {
	cs := []struct {
		w []string
		s string
		b byte
	}{
		{[]string{}, "", 'c'},
		{[]string{"http://a", "b", "c"}, "http://a.b.c", '.'},
		{[]string{"http:", "a.b.c"}, "http://a.b.c", '/'},
	}

	for i, c := range cs {
		a := FieldsByte(c.s, c.b)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] FieldsByte(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestFieldsRune(t *testing.T) {
	cs := []struct {
		w []string
		s string
		b rune
	}{
		{[]string{}, "", 'c'},
		{[]string{"http://a", "b", "c"}, "http://a.b.c", '.'},
		{[]string{"http:", "a.b.c"}, "http://a.b.c", '/'},
		{[]string{"http://一", "二", "三"}, "http://一.二.三", '.'},
		{[]string{"http://一", "二", "三"}, "http://一。二。三", '。'},
	}

	for i, c := range cs {
		a := FieldsRune(c.s, c.b)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] FieldsRune(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

func TestFieldsFuncIter(t *testing.T) {
	cs := []struct {
		w []string
		s string
		b rune
	}{
		{[]string{}, "", 'c'},
		{[]string{"http://a", "b", "c"}, "http://a.b.c", '.'},
		{[]string{"http:", "a.b.c"}, "http://a.b.c/", '/'},
		{[]string{"http://一", "二", "三"}, "http://一.二.三", '.'},
		{[]string{"http://一", "二", "三"}, "http://一。二。三", '。'},
	}

	for i, c := range cs {
		a := []string{}
		FieldsFuncIter(c.s, func(r rune) bool { return c.b == r }, func(s string) bool {
			a = append(a, s)
			return true
		})

		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] FieldsFuncIter(%q, %q) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}
