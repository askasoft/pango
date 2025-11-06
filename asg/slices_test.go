package asg

import (
	"cmp"
	"fmt"
	"math"
	"strings"
	"testing"
)

var equalIntTests = []struct {
	s1, s2 []int
	want   bool
}{
	{
		[]int{1},
		nil,
		false,
	},
	{
		[]int{},
		nil,
		true,
	},
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3},
		true,
	},
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3, 4},
		false,
	},
}

var equalFloatTests = []struct {
	s1, s2       []float64
	wantEqual    bool
	wantEqualNaN bool
}{
	{
		[]float64{1, 2},
		[]float64{1, 2},
		true,
		true,
	},
	{
		[]float64{1, 2, math.NaN()},
		[]float64{1, 2, math.NaN()},
		false,
		true,
	},
}

func TestEqual(t *testing.T) {
	for _, test := range equalIntTests {
		if got := Equal(test.s1, test.s2); got != test.want {
			t.Errorf("Equal(%v, %v) = %t, want %t", test.s1, test.s2, got, test.want)
		}
	}
	for _, test := range equalFloatTests {
		if got := Equal(test.s1, test.s2); got != test.wantEqual {
			t.Errorf("Equal(%v, %v) = %t, want %t", test.s1, test.s2, got, test.wantEqual)
		}
	}
}

// equal is simply ==.
func equal[T comparable](v1, v2 T) bool {
	return v1 == v2
}

// equalNaN is like == except that all NaNs are equal.
func equalNaN[T comparable](v1, v2 T) bool {
	isNaN := func(f T) bool { return f != f }
	return v1 == v2 || (isNaN(v1) && isNaN(v2))
}

// offByOne returns true if integers v1 and v2 differ by 1.
func offByOne(v1, v2 int) bool {
	return v1 == v2+1 || v1 == v2-1
}

func TestEqualFunc(t *testing.T) {
	for _, test := range equalIntTests {
		if got := EqualFunc(test.s1, test.s2, equal[int]); got != test.want {
			t.Errorf("EqualFunc(%v, %v, equal[int]) = %t, want %t", test.s1, test.s2, got, test.want)
		}
	}
	for _, test := range equalFloatTests {
		if got := EqualFunc(test.s1, test.s2, equal[float64]); got != test.wantEqual {
			t.Errorf("Equal(%v, %v, equal[float64]) = %t, want %t", test.s1, test.s2, got, test.wantEqual)
		}
		if got := EqualFunc(test.s1, test.s2, equalNaN[float64]); got != test.wantEqualNaN {
			t.Errorf("Equal(%v, %v, equalNaN[float64]) = %t, want %t", test.s1, test.s2, got, test.wantEqualNaN)
		}
	}

	s1 := []int{1, 2, 3}
	s2 := []int{2, 3, 4}
	if EqualFunc(s1, s1, offByOne) {
		t.Errorf("EqualFunc(%v, %v, offByOne) = true, want false", s1, s1)
	}
	if !EqualFunc(s1, s2, offByOne) {
		t.Errorf("EqualFunc(%v, %v, offByOne) = false, want true", s1, s2)
	}

	s3 := []string{"a", "b", "c"}
	s4 := []string{"A", "B", "C"}
	if !EqualFunc(s3, s4, strings.EqualFold) {
		t.Errorf("EqualFunc(%v, %v, strings.EqualFold) = false, want true", s3, s4)
	}

	cmpIntString := func(v1 int, v2 string) bool {
		return string(rune(v1)-1+'a') == v2
	}
	if !EqualFunc(s1, s3, cmpIntString) {
		t.Errorf("EqualFunc(%v, %v, cmpIntString) = false, want true", s1, s3)
	}
}

func BenchmarkEqualFunc_Large(b *testing.B) {
	type Large [4 * 1024]byte

	xs := make([]Large, 1024)
	ys := make([]Large, 1024)
	for i := 0; i < b.N; i++ {
		_ = EqualFunc(xs, ys, func(x, y Large) bool { return x == y })
	}
}

var indexTests = []struct {
	s    []int
	v    int
	want int
}{
	{
		nil,
		0,
		-1,
	},
	{
		[]int{},
		0,
		-1,
	},
	{
		[]int{1, 2, 3},
		2,
		1,
	},
	{
		[]int{1, 2, 2, 3},
		2,
		1,
	},
	{
		[]int{1, 2, 3, 2},
		2,
		1,
	},
}

func TestIndex(t *testing.T) {
	for _, test := range indexTests {
		if got := Index(test.s, test.v); got != test.want {
			t.Errorf("Index(%v, %v) = %d, want %d", test.s, test.v, got, test.want)
		}
	}
}

func equalToIndex[T any](f func(T, T) bool, v1 T) func(T) bool {
	return func(v2 T) bool {
		return f(v1, v2)
	}
}

func BenchmarkIndex_Large(b *testing.B) {
	type Large [4 * 1024]byte

	ss := make([]Large, 1024)
	for i := 0; i < b.N; i++ {
		_ = Index(ss, Large{1})
	}
}

func TestIndexFunc(t *testing.T) {
	for _, test := range indexTests {
		if got := IndexFunc(test.s, equalToIndex(equal[int], test.v)); got != test.want {
			t.Errorf("IndexFunc(%v, equalToIndex(equal[int], %v)) = %d, want %d", test.s, test.v, got, test.want)
		}
	}

	s1 := []string{"hi", "HI"}
	if got := IndexFunc(s1, equalToIndex(equal[string], "HI")); got != 1 {
		t.Errorf("IndexFunc(%v, equalToIndex(equal[string], %q)) = %d, want %d", s1, "HI", got, 1)
	}
	if got := IndexFunc(s1, equalToIndex(strings.EqualFold, "HI")); got != 0 {
		t.Errorf("IndexFunc(%v, equalToIndex(strings.EqualFold, %q)) = %d, want %d", s1, "HI", got, 0)
	}
}

func BenchmarkIndexFunc_Large(b *testing.B) {
	type Large [4 * 1024]byte

	ss := make([]Large, 1024)
	for i := 0; i < b.N; i++ {
		_ = IndexFunc(ss, func(e Large) bool {
			return e == Large{1}
		})
	}
}

func TestFindFunc(t *testing.T) {
	for _, test := range indexTests {
		if got, ok := FindFunc(test.s, equalToIndex(equal[int], test.v)); ok != (test.want != -1) || (test.want == -1 && got != 0) || (test.want != -1 && got != test.s[test.want]) {
			t.Errorf("IndexFunc(%v, equalToIndex(equal[int], %v)) = (%d, %v), want %d", test.s, test.v, got, ok, test.want)
		}
	}

	s1 := []string{"hi", "HI"}
	if got, ok := FindFunc(s1, equalToIndex(equal[string], "HI")); !ok || got != "HI" {
		t.Errorf("FindFunc(%v, equalToIndex(equal[string], %q)) = (%q, %v), want (%q, %v)", s1, "HI", got, ok, "HI", true)
	}
	if got, ok := FindFunc(s1, equalToIndex(strings.EqualFold, "HI")); !ok || got != "hi" {
		t.Errorf("FindFunc(%v, equalToIndex(strings.EqualFold, %q)) = (%q, %v), want (%q, %v)", s1, "HI", got, ok, "hi", true)
	}
}

func TestContains(t *testing.T) {
	for _, test := range indexTests {
		if got := Contains(test.s, test.v); got != (test.want != -1) {
			t.Errorf("Contains(%v, %v) = %t, want %t", test.s, test.v, got, test.want != -1)
		}
	}
}

func TestContainsFunc(t *testing.T) {
	for _, test := range indexTests {
		if got := ContainsFunc(test.s, equalToIndex(equal[int], test.v)); got != (test.want != -1) {
			t.Errorf("ContainsFunc(%v, equalToIndex(equal[int], %v)) = %t, want %t", test.s, test.v, got, test.want != -1)
		}
	}

	s1 := []string{"hi", "HI"}
	if got := ContainsFunc(s1, equalToIndex(equal[string], "HI")); got != true {
		t.Errorf("ContainsFunc(%v, equalToContains(equal[string], %q)) = %t, want %t", s1, "HI", got, true)
	}
	if got := ContainsFunc(s1, equalToIndex(equal[string], "hI")); got != false {
		t.Errorf("ContainsFunc(%v, equalToContains(strings.EqualFold, %q)) = %t, want %t", s1, "hI", got, false)
	}
	if got := ContainsFunc(s1, equalToIndex(strings.EqualFold, "hI")); got != true {
		t.Errorf("ContainsFunc(%v, equalToContains(strings.EqualFold, %q)) = %t, want %t", s1, "hI", got, true)
	}
}

var deleteTests = []struct {
	s    []int
	i, j int
	want []int
}{
	{
		[]int{1, 2, 3},
		0,
		0,
		[]int{1, 2, 3},
	},
	{
		[]int{1, 2, 3},
		0,
		1,
		[]int{2, 3},
	},
	{
		[]int{1, 2, 3},
		3,
		3,
		[]int{1, 2, 3},
	},
	{
		[]int{1, 2, 3},
		0,
		2,
		[]int{3},
	},
	{
		[]int{1, 2, 3},
		0,
		3,
		[]int{},
	},
}

func TestDelete(t *testing.T) {
	for _, test := range deleteTests {
		copy := Clone(test.s)
		if got := Delete(copy, test.i, test.j); !Equal(got, test.want) {
			t.Errorf("Delete(%v, %d, %d) = %v, want %v", test.s, test.i, test.j, got, test.want)
		}
	}
}

var deleteFuncTests = []struct {
	s    []int
	fn   func(int) bool
	want []int
}{
	{
		nil,
		func(int) bool { return true },
		nil,
	},
	{
		[]int{1, 2, 3},
		func(int) bool { return true },
		nil,
	},
	{
		[]int{1, 2, 3},
		func(int) bool { return false },
		[]int{1, 2, 3},
	},
	{
		[]int{1, 2, 3},
		func(i int) bool { return i > 2 },
		[]int{1, 2},
	},
	{
		[]int{1, 2, 3},
		func(i int) bool { return i < 2 },
		[]int{2, 3},
	},
	{
		[]int{10, 2, 30},
		func(i int) bool { return i >= 10 },
		[]int{2},
	},
}

func TestDeleteFunc(t *testing.T) {
	for i, test := range deleteFuncTests {
		copy := Clone(test.s)
		if got := DeleteFunc(copy, test.fn); !Equal(got, test.want) {
			t.Errorf("DeleteFunc case %d: got %v, want %v", i, got, test.want)
		}
	}
}

var deleteEqualTests = []struct {
	s    []int
	v    int
	want []int
}{
	{
		nil,
		0,
		nil,
	},
	{
		[]int{1, 2, 3},
		4,
		[]int{1, 2, 3},
	},
	{
		[]int{1, 2, 3},
		3,
		[]int{1, 2},
	},
	{
		[]int{1, 2, 3},
		2,
		[]int{1, 3},
	},
	{
		[]int{1, 2, 3},
		1,
		[]int{2, 3},
	},
	{
		[]int{1, 1, 2, 3},
		1,
		[]int{2, 3},
	},
}

func TestDeleteEqual(t *testing.T) {
	for i, test := range deleteEqualTests {
		copy := Clone(test.s)
		if got := DeleteEqual(copy, test.v); !Equal(got, test.want) {
			t.Errorf("DeleteEqual case %d: got %v, want %v", i, got, test.want)
		}
	}
}

func panics(f func()) (b bool) {
	defer func() {
		if x := recover(); x != nil {
			b = true
		}
	}()
	f()
	return false
}

func TestDeletePanics(t *testing.T) {
	for _, test := range []struct {
		name string
		s    []int
		i, j int
	}{
		{"with negative first index", []int{42}, -2, 1},
		{"with negative second index", []int{42}, 1, -1},
		{"with out-of-bounds first index", []int{42}, 2, 3},
		{"with out-of-bounds second index", []int{42}, 0, 2},
		{"with invalid i>j", []int{42}, 1, 0},
	} {
		if !panics(func() { Delete(test.s, test.i, test.j) }) {
			t.Errorf("Delete %s: got no panic, want panic", test.name)
		}
	}
}

func TestClone(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := Clone(s1)
	if !Equal(s1, s2) {
		t.Errorf("Clone(%v) = %v, want %v", s1, s2, s1)
	}
	s1[0] = 4
	want := []int{1, 2, 3}
	if !Equal(s2, want) {
		t.Errorf("Clone(%v) changed unexpectedly to %v", want, s2)
	}
	if got := Clone([]int(nil)); got != nil {
		t.Errorf("Clone(nil) = %#v, want nil", got)
	}
	if got := Clone(s1[:0]); got == nil || len(got) != 0 {
		t.Errorf("Clone(%v) = %#v, want %#v", s1[:0], got, s1[:0])
	}
}

var compactTests = []struct {
	name string
	s    []int
	want []int
}{
	{
		"nil",
		nil,
		nil,
	},
	{
		"one",
		[]int{1},
		[]int{1},
	},
	{
		"sorted",
		[]int{1, 2, 3},
		[]int{1, 2, 3},
	},
	{
		"1 item",
		[]int{1, 1, 2},
		[]int{1, 2},
	},
	{
		"unsorted",
		[]int{1, 2, 1},
		[]int{1, 2, 1},
	},
	{
		"many",
		[]int{1, 2, 2, 3, 3, 4},
		[]int{1, 2, 3, 4},
	},
}

func TestCompact(t *testing.T) {
	for _, test := range compactTests {
		copy := Clone(test.s)
		if got := Compact(copy); !Equal(got, test.want) {
			t.Errorf("Compact(%v) = %v, want %v", test.s, got, test.want)
		}
	}
}

func BenchmarkCompact(b *testing.B) {
	for _, c := range compactTests {
		b.Run(c.name, func(b *testing.B) {
			ss := make([]int, 0, 64)
			for k := 0; k < b.N; k++ {
				ss = ss[:0]
				ss = append(ss, c.s...)
				_ = Compact(ss)
			}
		})
	}
}

func BenchmarkCompact_Large(b *testing.B) {
	type Large [4 * 1024]byte

	ss := make([]Large, 1024)
	for i := 0; i < b.N; i++ {
		_ = Compact(ss)
	}
}

func TestCompactFunc(t *testing.T) {
	for _, test := range compactTests {
		copy := Clone(test.s)
		if got := CompactFunc(copy, equal[int]); !Equal(got, test.want) {
			t.Errorf("CompactFunc(%v, equal[int]) = %v, want %v", test.s, got, test.want)
		}
	}

	s1 := []string{"a", "a", "A", "B", "b"}
	copy := Clone(s1)
	want := []string{"a", "B"}
	if got := CompactFunc(copy, strings.EqualFold); !Equal(got, want) {
		t.Errorf("CompactFunc(%v, strings.EqualFold) = %v, want %v", s1, got, want)
	}
}

func BenchmarkCompactFunc_Large(b *testing.B) {
	type Large [4 * 1024]byte

	ss := make([]Large, 1024)
	for i := 0; i < b.N; i++ {
		_ = CompactFunc(ss, func(a, b Large) bool { return a == b })
	}
}

func TestGrow(t *testing.T) {
	s1 := []int{1, 2, 3}

	copy := Clone(s1)
	s2 := Grow(copy, 1000)
	if !Equal(s1, s2) {
		t.Errorf("Grow(%v) = %v, want %v", s1, s2, s1)
	}
	if cap(s2) < 1000+len(s1) {
		t.Errorf("after Grow(%v) cap = %d, want >= %d", s1, cap(s2), 1000+len(s1))
	}

	// Test mutation of elements between length and capacity.
	copy = Clone(s1)
	s3 := Grow(copy[:1], 2)[:3]
	if !Equal(s1, s3) {
		t.Errorf("Grow should not mutate elements between length and capacity")
	}
	s3 = Grow(copy[:1], 1000)[:3]
	if !Equal(s1, s3) {
		t.Errorf("Grow should not mutate elements between length and capacity")
	}

	// Test number of allocations.
	if n := testing.AllocsPerRun(100, func() { Grow(s2, cap(s2)-len(s2)) }); n != 0 {
		t.Errorf("Grow should not allocate when given sufficient capacity; allocated %v times", n)
	}
	if n := testing.AllocsPerRun(100, func() { Grow(s2, cap(s2)-len(s2)+1) }); n != 1 {
		errorf := t.Errorf
		errorf("Grow should allocate once when given insufficient capacity; allocated %v times", n)
	}

	// Test for negative growth sizes.
	var gotPanic bool
	func() {
		defer func() { gotPanic = recover() != nil }()
		Grow(s1, -1)
	}()
	if !gotPanic {
		t.Errorf("Grow(-1) did not panic; expected a panic")
	}
}

func TestClip(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5, 6}[:3]
	orig := Clone(s1)
	if len(s1) != 3 {
		t.Errorf("len(%v) = %d, want 3", s1, len(s1))
	}
	if cap(s1) < 6 {
		t.Errorf("cap(%v[:3]) = %d, want >= 6", orig, cap(s1))
	}
	s2 := Clip(s1)
	if !Equal(s1, s2) {
		t.Errorf("Clip(%v) = %v, want %v", s1, s2, s1)
	}
	if cap(s2) != 3 {
		t.Errorf("cap(Clip(%v)) = %d, want 3", orig, cap(s2))
	}
}

func TestReverse(t *testing.T) {
	even := []int{3, 1, 4, 1, 5, 9} // len = 6
	Reverse(even)
	if want := []int{9, 5, 1, 4, 1, 3}; !Equal(even, want) {
		t.Errorf("Reverse(even) = %v, want %v", even, want)
	}

	odd := []int{3, 1, 4, 1, 5, 9, 2} // len = 7
	Reverse(odd)
	if want := []int{2, 9, 5, 1, 4, 1, 3}; !Equal(odd, want) {
		t.Errorf("Reverse(odd) = %v, want %v", odd, want)
	}

	words := strings.Fields("one two three")
	Reverse(words)
	if want := strings.Fields("three two one"); !Equal(words, want) {
		t.Errorf("Reverse(words) = %v, want %v", words, want)
	}

	singleton := []string{"one"}
	Reverse(singleton)
	if want := []string{"one"}; !Equal(singleton, want) {
		t.Errorf("Reverse(singeleton) = %v, want %v", singleton, want)
	}

	Reverse[[]string](nil)
}

var minMaxTests = []struct {
	name string
	s    []int
	min  int
	max  int
}{
	{"one", []int{1}, 1, 1},
	{"two", []int{1, 2}, 1, 2},
}

func TestMin(t *testing.T) {
	for _, c := range minMaxTests {
		a := Min(c.s)
		if c.min != a {
			t.Errorf("[%s] Min(%v) = %v, WANT %v", c.name, c.s, a, c.min)
		}
	}
}

func TestMinFunc(t *testing.T) {
	for _, c := range minMaxTests {
		a := MinFunc(c.s, cmp.Compare)
		if c.min != a {
			t.Errorf("[%s] MinFunc(%v) = %v, WANT %v", c.name, c.s, a, c.min)
		}
	}
}

func TestMax(t *testing.T) {
	for _, c := range minMaxTests {
		a := Max(c.s)
		if c.max != a {
			t.Errorf("[%s] Max(%v) = %v, WANT %v", c.name, c.s, a, c.max)
		}
	}
}

func TestMaxFunc(t *testing.T) {
	for _, c := range minMaxTests {
		a := MaxFunc(c.s, cmp.Compare)
		if c.max != a {
			t.Errorf("[%s] MaxFunc(%v) = %v, WANT %v", c.name, c.s, a, c.max)
		}
	}
}

func TestJoin(t *testing.T) {
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
		a := Join(c.s, " ", c.f...)
		if a != c.w {
			t.Errorf("[%d] Join(%v) = %v, want %v", i, c.s, a, c.w)
		}
	}
}
