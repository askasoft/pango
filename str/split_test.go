package str

import (
	"reflect"
	"testing"
)

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
