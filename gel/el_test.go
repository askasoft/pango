package gel

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/askasoft/pango/str"
)

type testcase1 struct {
	w any
	s string
}

func testCalculate1(t *testing.T, cs []testcase1) {
	for i, c := range cs {
		el, err := Compile(c.s)
		if err != nil {
			t.Fatalf("[%d] Compile(%q) = %v", i, c.s, err)
		}
		if a := el.String(); a != c.s {
			t.Fatalf("[%d] String(%q) = %v, want %v", i, c.s, a, c.s)
		}

		a, err := el.Calculate(nil)
		if err != nil {
			t.Fatalf("[%d] Calculate(%q) = %v", i, c.s, err)
		}
		if !reflect.DeepEqual(a, c.w) {
			t.Fatalf("[%d] Calculate(%q) = (%v, %T), want (%v, %T)", i, c.s, a, a, c.w, c.w)
		}
		fmt.Printf("[%d] Calculate(%q) = (%v, %T)\n", i, c.s, a, a)
	}
}

type testcase2 struct {
	w any
	s string
	d any
}

func testCalculate2(t *testing.T, cs []testcase2) {
	for i, c := range cs {
		a, err := Calculate(c.s, c.d)
		if err != nil {
			t.Fatalf("[%d] Calculate(%q, %v) = %v", i, c.s, c.d, err)
		}
		if !reflect.DeepEqual(a, c.w) {
			t.Fatalf("[%d] Calculate(%q, %v) = (%v, %T), want (%v, %T)", i, c.s, c.d, a, a, c.w, c.w)
		}
		fmt.Printf("[%d] Calculate(%q, %v) = (%v, %T)\n", i, c.s, c.d, a, a)
	}
}

func testCalculate2s(t *testing.T, cs []testcase2) {
	for i, c := range cs {
		a, err := CalculateStrict(c.s, c.d)
		if wer, ok := c.w.(error); ok {
			if wer.Error() != fmt.Sprint(err) {
				t.Fatalf("[%d] CalculateStrict(%q, %v) = (%v, %T), want (%v, %T)", i, c.s, c.d, err, err, c.w, c.w)
			}
			continue
		}
		if err != nil {
			t.Fatalf("[%d] CalculateStrict(%q, %v) = %v", i, c.s, c.d, err)
		}
		if !reflect.DeepEqual(a, c.w) {
			t.Fatalf("[%d] CalculateStrict(%q, %v) = (%v, %T), want (%v, %T)", i, c.s, c.d, a, a, c.w, c.w)
		}
		fmt.Printf("[%d] CalculateStrict(%q, %v) = (%v, %T)\n", i, c.s, c.d, a, a)
	}
}

func TestOneValue(t *testing.T) {
	cs := []testcase1{
		{1, "1"},
		{float64(0.1), ".1"},
		{float64(0.1), "0.1"},
		{float32(0.1), "0.1f"},
		{float64(0.1), "0.1d"},
		{true, "true"},
		{false, "false"},
		{"jk", "'jk'"},
	}
	testCalculate1(t, cs)
}

func TestBit(t *testing.T) {
	cs := []testcase1{
		{-5 << 3, "-5<<3"},
		{-5 >> 3, "-5>>3"},
		{5 & 3, "5&3"},
		{5 | 3, "5|3"},
		{^5, "~5"},
		{^5, "^5"},
		{5 ^ 3, "5^3"},
		{6 + ^5, "6 + ~5"},
		{6 + ^5, "6 + ^5"},
		{1 + 1 + ^11, "1 + 1 + ^11"},
	}
	testCalculate1(t, cs)
}

func TestMathSingle(t *testing.T) {
	cs := []testcase1{
		{2, "1+1"},
		{2.2, "1.1+1.1"},
		{1, "2-1"},
		{9, "3*3"},
		{0, "3*0"},
		{3, "9/3"},
		{2.2, "4.4/2"},
		{9.9 / float64(3.0), "9.9/3"},
		{1, "5%2"},
	}
	testCalculate1(t, cs)
}

func TestMathMulti(t *testing.T) {
	cs := []testcase1{
		{3, "1 + 1 + 1"},
		{1, "  1+1-1  "},
		{-1, "1-1-1"},
		{1, "1-(1-1)"},
		{7, "1+2*3"},
		{2*4 + 2*3 + 4*5, "2*4+2*3+4*5"},
		{9 + 8*7 + (6+5)*(4-1*2+3), "9+8*7+(6+5)*((4-1*2+3))"},
		{.3 + .2*.5, ".3+.2*.5"},
		{(.5 + 0.1) * .9, "(.5 + 0.1)*.9"},
		{1/1 + 10*(1400-1400)/400, "1/1+10*(1400-1400)/400"},
		{0.1354 * ((70 - 8) % 70) * 100, "0.1354 * ((70 - 8) % 70) * 100"},
		{0.5006 * (70 / 600 * 100), "0.5006 * (70 / 600 * 100)"},
		{2 + (-3), "2+(-3)"},
		{2 + -3, "2+-3"},
		{2 * -3, "2*-3"},
		{-2 * -3, "-2*-3"},
		{2 / -3, "2/-3"},
		{2 % -3, "2%-3"},
		{1000 + 100.0*99 - (600-3*15)%(((68-9)-3)*2-100) + 10000%7*71, "1000+100.0*99-(600-3*15)%(((68-9)-3)*2-100)+10000%7*71"},
		{1, "6.7-100>39.6 ? 5==5? 4+5:6-1 : !(100%3-39.0<27) ? 8*2-199: 100%3"},
	}
	testCalculate1(t, cs)
}

func TestNil(t *testing.T) {
	cs := []testcase1{
		{nil, "nil"},
	}
	testCalculate1(t, cs)
}

func TestLogical(t *testing.T) {
	cs := []testcase1{
		{true, "2 > 1"},
		{false, "2 < 1"},
		{true, "2 >= 2"},
		{true, "2 <= 2"},
		{true, "2 == 2 "},
		{1 != 2, "1 != 2"},
		{true, "!(1 == 2)"},
		{false, "!(1 == 1)"},
		{!false == false, "!false == false"},
		{!false, "!false"},
		{true || false, "true || false"},
		{true && false, "true && false"},
		{false || true && false, "false || true && false"},
		{true, `"a" == "a"`},
		{true, `"abc" ~= "^a.*$"`},
		{true, `"abc" ~= "b"`},
		{false, `"abc" ~= "abz"`},
		{false, `"a" == !!$`},
		{true, `nil == !!$`},
	}
	testCalculate1(t, cs)
}

func TestTernary(t *testing.T) {
	cs := []testcase1{
		{2, "1>0?2:3"},
		{2, "1>0&&1<2?2:3"},
	}
	testCalculate1(t, cs)
}

func TestString(t *testing.T) {
	cs := []testcase1{
		{"jk", "'jk'"},
		{"j\r\n\t '\"　k", "\"j\\r\\n\\t\\x20\\'\\\"\\u3000k\""},
		{"jk", "'j' + 'k'"},
		{"j0", "'j' + 0"},
	}
	testCalculate1(t, cs)
}

func TestNegative(t *testing.T) {
	cs := []testcase1{
		{-1, "-1"},
		{0, "-1+1"},
		{-1 - -1, "-1 - -1"},
		{9 + 8*7 + (6+5)*(-(4 - 1*2 + 3)), "9+8*7+(6+5)*(-(4-1*2+3))"},
	}
	testCalculate1(t, cs)
}

type teststr string

func (ts teststr) Len() int {
	return len(ts)
}

func (ts teststr) Left(i int) string {
	return str.Left(string(ts), i)
}

func (ts teststr) Substring(i, n int) string {
	return string(ts)[i:n]
}

func (ts teststr) IndexOf(s string) int {
	return str.Index(string(ts), s)
}

func (ts teststr) Contains(s string) bool {
	return str.Contains(string(ts), s)
}

func (ts teststr) Strip() string {
	return str.Strip(string(ts))
}

type pet struct {
	name string
	Age  int
	Fget func() string
	Fset func(string)
}

func (p *pet) SetName(name string) {
	p.name = name
}

func (p *pet) GetName() string {
	return p.name
}

func (p *pet) Display() string {
	return p.name
}

func TestObject(t *testing.T) {
	pet := &pet{
		name: "XiaoBai",
		Age:  10,
	}

	cs := []testcase2{
		{10, ".age", pet},
		{10, ".Age", pet},
		{"XiaoBai", ".name", pet},
		{"XiaoBai", ".Name", pet},
		{"XiaoBai", ".display()", pet},
		{"XiaoBai", ".Display()", pet},
	}
	testCalculate2(t, cs)
}

func TestCallFunc(t *testing.T) {
	pet := &pet{}
	pet.Fget = pet.GetName
	pet.Fset = pet.SetName

	m := map[string]any{
		"get": pet.Fget,
		"set": pet.Fset,
	}

	cs := []testcase2{
		{"ab", "Left(2)", teststr("abcde")},
		{"b", ".Substring(1,2)", teststr("abcde")},
		{true, ".Contains('cd')", teststr("abcde")},
		{"abab", ".Strip()", teststr("  abab  ")},
		{5, ".Len()", teststr("abcde")},
		{nil, "SetName('XiaoBai')", pet},
		{"XiaoBai", ".GetName()", pet},
		{nil, ".Fset('XiaoHei')", pet},
		{"XiaoHei", "Fget()", pet},
		{nil, ".set('XiaoHui')", m},
		{"XiaoHui", "get()", m},
	}
	testCalculate2(t, cs)
}

func TestArray(t *testing.T) {
	m := map[string]any{
		"a": []string{"a", "b", "c"},
		"b": [][]string{{"a", "b"}, {"c", "d"}},
	}
	cs := []testcase2{
		{[]any{}, "{}", m},
		{[]any{1}, "{1}", m},
		{"b", "a[1]", m},
		{"b", "a[2-1]", m},
		{"d", "b[1][1]", m},
	}
	testCalculate2(t, cs)
}

func TestMap(t *testing.T) {
	m := map[string]any{
		"a": map[string]any{"x": 10, "y": 50, "txt": "Hello"},
		"b": map[string]any{"c": map[string]any{"x": 10, "y": 50, "txt": "Hello"}},
	}

	m2 := map[string]any{
		"i":    100,
		"pi":   3.14,
		"d":    -3.9,
		"b":    uint8(4),
		"bool": false,
		"t":    "",
	}

	cs := []testcase2{
		{100, "a.x*10", m},
		{100, "b.c.x*10", m},
		{100, "$.b.c.x*10", m},
		{100, "a['x']*10", m},
		{100, "b.c['x']*10", m},
		{100, "$.b.c['x']*10", m},
		{50, "a.x > a.y ? a.x : a.y", m},
		{50, "b.c.x > b.c.y ? b.c.x : b.c.y", m},
		{"Hello-40", "a['txt']+(a.x-a.y)", m},
		{"Hello-40", "b.c['txt']+(b.c.x-b.c.y)", m},
		{true, "i * pi + (d * b - 199) / (1 - d * pi) - (2 + 100 - i / pi) % 99 ==i * pi + (d * b - 199) / (1 - d * pi) - (2 + 100 - i / pi) % 99", m2},
		{true, "'A' == 'A' || 'B' == 'B' && 'ABCD' == t &&  'A' == 'A'", m2},
		{">= 1", "(min != nil && max != nil) ? (min + '~' + max) : (min != nil ? ('>= ' + min) : (max != nil ? ('<= ' + max) : ''))", map[string]any{"min": 1}},
		{"<= 2", "(min != nil && max != nil) ? (min + '~' + max) : (min != nil ? ('>= ' + min) : (max != nil ? ('<= ' + max) : ''))", map[string]any{"max": 2}},
		{"1~2", "(min != nil && max != nil) ? (min + '~' + max) : (min != nil ? ('>= ' + min) : (max != nil ? ('<= ' + max) : ''))", map[string]any{"min": 1, "max": 2}},
		{nil, "a['z']", m},
		{nil, "a['\\'z']", m},
	}
	testCalculate2(t, cs)
}

func TestOrable(t *testing.T) {
	m := map[string]any{}
	m["obj"] = map[string]any{"pet": nil}
	m["girls"] = []string{}

	cs1 := []testcase2{
		{"cat", "obj.pet.name ||| 'cat'", m},
		{"dog", "obj.girls ||| 'dog'", m},
	}
	testCalculate2(t, cs1)

	cs2 := []testcase2{
		{"cat", "!!(obj.pet.name) ||| 'cat'", m},
		{"dog", "!!(obj.girls) ||| 'dog'", m},
		{"cat", "!!obj.pet.name ||| 'cat'", m},
		{"dog", "!!obj.girls ||| 'dog'", m},
	}
	testCalculate2s(t, cs2)
}

func TestStrict(t *testing.T) {
	m := map[string]any{
		"obj": map[string]any{"pet": nil},
	}

	cs := []testcase2{
		{errors.New("gel: can't get nil.name"), "(obj.pet.name) == nil", m},
	}
	testCalculate2s(t, cs)
}
