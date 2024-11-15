package bye

import "testing"

func TestCountByte(t *testing.T) {
	cs := []struct {
		w int
		s string
		b byte
	}{
		{0, "", 'a'},
		{1, " a ", 'a'},
		{0, " b ", 'a'},
		{2, " a a ", 'a'},
	}

	for i, c := range cs {
		a := CountByte([]byte(c.s), c.b)
		if a != c.w {
			t.Errorf("[%d] CountByte(%q, %v) = %v, want %v", i, c.s, c.b, a, c.w)
		}
	}
}

const space = "\t\v\r\f\n\u0085\u00a0\u2000\u3000"

type bytesTest struct {
	in, out string
}

// Execute f on each test case.  funcName should be the name of f; it's used
// in failure reports.
func runBytesTests(t *testing.T, f func([]byte) []byte, funcName string, testCases []bytesTest) {
	for _, tc := range testCases {
		actual := string(f([]byte(tc.in)))
		if actual != tc.out {
			t.Errorf("%s(%q) = %q; want %q", funcName, tc.in, actual, tc.out)
		}
	}
}

func TestStrip(t *testing.T) {
	var tests = []bytesTest{
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

	runBytesTests(t, Strip, "Strip", tests)
}

func TestStripLeft(t *testing.T) {
	var tests = []bytesTest{
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

	runBytesTests(t, StripLeft, "StripLeft", tests)
}

func TestStripRight(t *testing.T) {
	var tests = []bytesTest{
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

	runBytesTests(t, StripRight, "StripRight", tests)
}
