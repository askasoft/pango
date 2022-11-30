package funcs

import (
	"errors"
	"fmt"
	"html/template"
	"testing"
)

func TestHTML(t *testing.T) {
	cs := []struct {
		a any
		w template.HTML
		e error
	}{
		{"1", "1", nil},
		{[]byte{'2'}, "2", nil},
		{'3', "3", nil},
		{1.1, "", errors.New("HTML: unknown type for '1.1' (float64)")},
	}

	for i, c := range cs {
		r, e := HTML(c.a)
		if c.w != r || fmt.Sprint(c.e) != fmt.Sprint(e) {
			t.Errorf("[%d] HTML(%v) = (%T, %v, %v), want: (%T, %v, %v)", i, c.a, r, r, e, c.w, c.w, c.e)
		}
	}
}

func TestCSS(t *testing.T) {
	cs := []struct {
		a any
		w template.CSS
		e error
	}{
		{"1", "1", nil},
		{[]byte{'2'}, "2", nil},
		{'3', "3", nil},
		{1.1, "", errors.New("CSS: unknown type for '1.1' (float64)")},
	}

	for i, c := range cs {
		r, e := CSS(c.a)
		if c.w != r || fmt.Sprint(c.e) != fmt.Sprint(e) {
			t.Errorf("[%d] CSS(%v) = (%T, %v, %v), want: (%T, %v, %v)", i, c.a, r, r, e, c.w, c.w, c.e)
		}
	}
}

func TestJS(t *testing.T) {
	cs := []struct {
		a any
		w template.JS
		e error
	}{
		{"1", "1", nil},
		{[]byte{'2'}, "2", nil},
		{'3', "3", nil},
		{1.1, "", errors.New("JS: unknown type for '1.1' (float64)")},
	}

	for i, c := range cs {
		r, e := JS(c.a)
		if c.w != r || fmt.Sprint(c.e) != fmt.Sprint(e) {
			t.Errorf("[%d] JS(%v) = (%T, %v, %v), want: (%T, %v, %v)", i, c.a, r, r, e, c.w, c.w, c.e)
		}
	}
}

func TestJSStr(t *testing.T) {
	cs := []struct {
		a any
		w template.JSStr
		e error
	}{
		{"1", "1", nil},
		{[]byte{'2'}, "2", nil},
		{'3', "3", nil},
		{1.1, "", errors.New("JSStr: unknown type for '1.1' (float64)")},
	}

	for i, c := range cs {
		r, e := JSStr(c.a)
		if c.w != r || fmt.Sprint(c.e) != fmt.Sprint(e) {
			t.Errorf("[%d] JSStr(%v) = (%T, %v, %v), want: (%T, %v, %v)", i, c.a, r, r, e, c.w, c.w, c.e)
		}
	}
}

func TestURL(t *testing.T) {
	cs := []struct {
		a any
		w template.URL
		e error
	}{
		{"1", "1", nil},
		{[]byte{'2'}, "2", nil},
		{'3', "3", nil},
		{1.1, "", errors.New("URL: unknown type for '1.1' (float64)")},
	}

	for i, c := range cs {
		r, e := URL(c.a)
		if c.w != r || fmt.Sprint(c.e) != fmt.Sprint(e) {
			t.Errorf("[%d] URL(%v) = (%T, %v, %v), want: (%T, %v, %v)", i, c.a, r, r, e, c.w, c.w, c.e)
		}
	}
}

func TestSrcset(t *testing.T) {
	cs := []struct {
		a any
		w template.Srcset
		e error
	}{
		{"1", "1", nil},
		{[]byte{'2'}, "2", nil},
		{'3', "3", nil},
		{1.1, "", errors.New("Srcset: unknown type for '1.1' (float64)")},
	}

	for i, c := range cs {
		r, e := Srcset(c.a)
		if c.w != r || fmt.Sprint(c.e) != fmt.Sprint(e) {
			t.Errorf("[%d] Srcset(%v) = (%T, %v, %v), want: (%T, %v, %v)", i, c.a, r, r, e, c.w, c.w, c.e)
		}
	}
}
