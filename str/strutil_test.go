package str

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCompare(t *testing.T) {
	cs := []struct {
		w int
		s string
		t string
	}{
		{0, " a ", " a "},
		{-1, " a ", " b "},
		{1, " b ", " a "},
	}

	for i, c := range cs {
		a := Compare(c.s, c.t)
		if a != c.w {
			t.Errorf("[%d] Compare(%q, %q) = %v, want %v", i, c.s, c.t, a, c.w)
		}
	}
}

const space = "\t\v\r\f\n\u0085\u00a0\u2000\u3000"

type StringTest struct {
	in, out string
}

// Execute f on each test case.  funcName should be the name of f; it's used
// in failure reports.
func runStringTests(t *testing.T, f func(string) string, funcName string, testCases []StringTest) {
	for _, tc := range testCases {
		actual := f(tc.in)
		if actual != tc.out {
			t.Errorf("%s(%q) = %q; want %q", funcName, tc.in, actual, tc.out)
		}
	}
}

func TestCapitalize(t *testing.T) {
	var tests = []StringTest{
		{"", ""},
		{"abc", "Abc"},
		{"abc abc", "Abc abc"},
		{"хлеб", "Хлеб"},
	}

	runStringTests(t, Capitalize, "Capitalize", tests)
}

func TestCamelCase(t *testing.T) {
	var tests = []StringTest{
		{"", ""},
		{"abc_abc", "abcAbc"},
		{"ABC", "abc"},
		{"Хлеб", "хлеб"},
		{"Хлеб-Хлеб", "хлебХлеб"},
	}

	runStringTests(t, func(s string) string { return CamelCase(s) }, "CamelCase", tests)
}

func TestPascalCase(t *testing.T) {
	var tests = []StringTest{
		{"", ""},
		{"abc_abc", "AbcAbc"},
		{"ABC", "Abc"},
		{"хлеб", "Хлеб"},
		{"хлеб-Хлеб", "ХлебХлеб"},
	}

	runStringTests(t, func(s string) string { return PascalCase(s) }, "PascalCase", tests)
}

func TestSnakeCase(t *testing.T) {
	var tests = []StringTest{
		{"", ""},
		{"abcAbc", "abc_abc"},
		{"abc", "abc"},
		{"ABC", "a_b_c"},
		{"Хлеб", "хлеб"},
		{"ХлебХлеб", "хлеб_хлеб"},
	}

	runStringTests(t, func(s string) string { return SnakeCase(s) }, "SnakeCase", tests)
}

func TestStrip(t *testing.T) {
	var tests = []StringTest{
		{"", ""},
		{"abc", "abc"},
		{space + "abc" + space, "abc"},
		{" ", ""},
		{" \t\r\n \t\t\r\r\n\n ", ""},
		{" \t\r\n x\t\t\r\r\n\n ", "x"},
		{" \u2000\t\r\n x\t\t\r\r\ny\n \u3000", "x\t\t\r\r\ny"},
		{"1 \t\r\n2", "1 \t\r\n2"},
		{" x\x80", "x\x80"},
		{" x\xc0", "x\xc0"},
		{"x \xc0\xc0 ", "x \xc0\xc0"},
		{"x \xc0", "x \xc0"},
		{"x \xc0 ", "x \xc0"},
		{"x \xc0\xc0 ", "x \xc0\xc0"},
		{"x ☺\xc0\xc0 ", "x ☺\xc0\xc0"},
		{"x ☺ ", "x ☺"},
	}

	runStringTests(t, Strip, "Strip", tests)
}

func TestStripLeft(t *testing.T) {
	var tests = []StringTest{
		{"", ""},
		{"abc", "abc"},
		{space + "abc" + space, "abc" + space},
		{" ", ""},
		{" \t\r\n \t\t\r\r\n\n ", ""},
		{" \t\r\n x\t\t\r\r\n\n ", "x\t\t\r\r\n\n "},
		{" \u2000\t\r\n x\t\t\r\r\ny\n \u3000", "x\t\t\r\r\ny\n \u3000"},
		{"1 \t\r\n2", "1 \t\r\n2"},
		{" x\x80", "x\x80"},
		{" x\xc0", "x\xc0"},
		{"x \xc0\xc0 ", "x \xc0\xc0 "},
		{"x \xc0", "x \xc0"},
		{"x \xc0 ", "x \xc0 "},
		{"x \xc0\xc0 ", "x \xc0\xc0 "},
		{"x ☺\xc0\xc0 ", "x ☺\xc0\xc0 "},
		{"x ☺ ", "x ☺ "},
	}

	runStringTests(t, StripLeft, "StripLeft", tests)
}

func TestStripRight(t *testing.T) {
	var tests = []StringTest{
		{"", ""},
		{"abc", "abc"},
		{space + "abc" + space, space + "abc"},
		{" ", ""},
		{" \t\r\n \t\t\r\r\n\n ", ""},
		{" \t\r\n x\t\t\r\r\n\n ", " \t\r\n x"},
		{" \u2000\t\r\n x\t\t\r\r\ny\n \u3000", " \u2000\t\r\n x\t\t\r\r\ny"},
		{"1 \t\r\n2", "1 \t\r\n2"},
		{" x\x80", " x\x80"},
		{" x\xc0", " x\xc0"},
		{"x \xc0\xc0 ", "x \xc0\xc0"},
		{"x \xc0", "x \xc0"},
		{" x \xc0 ", " x \xc0"},
		{"x \xc0\xc0 ", "x \xc0\xc0"},
		{" x ☺\xc0\xc0 ", " x ☺\xc0\xc0"},
		{" x ☺ ", " x ☺"},
	}

	runStringTests(t, StripRight, "StripRight", tests)
}

func TestRuneCount(t *testing.T) {
	cs := []struct {
		w int
		s string
	}{
		{0, ""},
		{4, "qeed"},
		{1, "あ"},
	}

	for i, c := range cs {
		a := RuneCount(c.s)
		if a != c.w {
			t.Errorf("[%d] RuneCount(%q) = %q, want %q", i, c.s, a, c.w)
		}
	}
}

func TestRuneEqualFold(t *testing.T) {
	cs := []struct {
		w bool
		s rune
		t rune
	}{
		{true, 'a', 'A'},
		{true, 'k', '\u212A'},
		{false, 'a', 'B'},
	}

	for i, c := range cs {
		a := RuneEqualFold(c.s, c.t)
		if a != c.w {
			t.Errorf("[%d] RuneEqualFold(%q, %q) = %v, want %v", i, c.s, c.t, a, c.w)
		}
	}
}

func TestRemoveByte(t *testing.T) {
	cs := []struct {
		w string
		s string
		b byte
	}{
		{"", "", 'a'},
		{"qeed", "queued", 'u'},
		{"queued", "queued", 'z'},
	}

	for i, c := range cs {
		a := RemoveByte(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] RemoveByte(%q, %q) = %q, want %q", i, c.s, c.b, a, c.w)
		}
	}
}

func TestRemoveRune(t *testing.T) {
	cs := []struct {
		w string
		s string
		r rune
	}{
		{"", "", 'a'},
		{"ueued", "queued", 'q'},
		{"queue", "queued", 'd'},
		{"qeed", "queued", 'u'},
		{"queued", "queued", 'z'},
		{"ありとうございます。", "ありがとうございます。", 'が'},
	}

	for i, c := range cs {
		a := RemoveRune(c.s, c.r)
		if a != c.w {
			t.Errorf("[%d] RemoveRune(%q, %v) = %q, want %q", i, c.s, c.r, a, c.w)
		}
	}
}

func TestRemoveAny(t *testing.T) {
	cs := []struct {
		w string
		s string
		b string
	}{
		{"", "", "ab"},
		{"qee", "queued", "ud"},
		{"queued", "queued", "z"},
		{"ありとういます。", "ありがとうございます。", "がござ"},
	}

	for i, c := range cs {
		a := RemoveAny(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] RemoveAny(%q, %q) = %q, want %q", i, c.s, c.b, a, c.w)
		}
	}
}

func TestRemoveFunc(t *testing.T) {
	cs := []struct {
		w string
		s string
		f func(r rune) bool
	}{
		{"", "", func(r rune) bool { return r == 'a' || r == 'b' }},
		{"qee", "queued", func(r rune) bool { return r == 'u' || r == 'd' }},
		{"queued", "queued", func(r rune) bool { return r == 'z' }},
		{"ありとういます。", "ありがとうございます。", func(r rune) bool { return r == 'が' || r == 'ご' || r == 'ざ' }},
	}

	for i, c := range cs {
		a := RemoveFunc(c.s, c.f)
		if a != c.w {
			t.Errorf("[%d] RemoveFunc(%q) = %q, want %q", i, c.s, a, c.w)
		}
	}
}

func TestRemove(t *testing.T) {
	cs := []struct {
		w string
		s string
		b string
	}{
		{"", "", "ab"},
		{"qd", "queued", "ue"},
		{"queued", "queued", "z"},
		{"ありございます。", "ありがとうございます。", "がとう"},
	}

	for i, c := range cs {
		a := Remove(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] Remove(%q, %q) = %q, want %q", i, c.s, c.b, a, c.w)
		}
	}
}

func TestTrimSpaces(t *testing.T) {
	cs := []struct {
		s []string
		w []string
	}{
		{[]string{" a "}, []string{"a"}},
	}

	for i, c := range cs {
		a := TrimSpaces(c.s)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] TrimSpaces(%v) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestRemoveEmptys(t *testing.T) {
	cs := []struct {
		s []string
		w []string
	}{
		{[]string{"a", ""}, []string{"a"}},
		{[]string{"a", "", "b", "", "c"}, []string{"a", "b", "c"}},
	}

	for i, c := range cs {
		a := RemoveEmptys(c.s)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] RemoveEmptys(%v) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestJoinInts(t *testing.T) {
	cs := []struct {
		s []int
		f []func(int) string
		w string
	}{
		{[]int{}, nil, ""},
		{[]int{1}, nil, "1"},
		{[]int{1, 2}, nil, "1 2"},
		{[]int{1, 10}, []func(int) string{func(i int) string { return fmt.Sprintf("0x%x", i) }}, "0x1 0xa"},
	}

	for i, c := range cs {
		a := JoinInts(c.s, " ", c.f...)
		if a != c.w {
			t.Errorf("[%d] JoinInts(%v) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestJoinInt64s(t *testing.T) {
	cs := []struct {
		s []int64
		f []func(int64) string
		w string
	}{
		{[]int64{}, nil, ""},
		{[]int64{1}, nil, "1"},
		{[]int64{1, 2}, nil, "1 2"},
		{[]int64{1, 10}, []func(int64) string{func(i int64) string { return fmt.Sprintf("0x%x", i) }}, "0x1 0xa"},
	}

	for i, c := range cs {
		a := JoinInt64s(c.s, " ", c.f...)
		if a != c.w {
			t.Errorf("[%d] JoinInt64s(%v) = %v, want %v", i, c.s, a, c.w)
		}
	}
}
