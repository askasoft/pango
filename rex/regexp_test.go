package rex

import (
	"html"
	"regexp"
	"testing"
)

func TestSubmatchReplace(t *testing.T) {
	cs := []struct {
		w string
		s string
		m string
		t string
	}{
		{"〒(123) 456-7890 〒", "〒123-456-7890 〒", `(\d{3})-(\d{3})-(\d{4})`, "($1) $2-$3"},
		{"Border: Red, Line: RED", "Border Color is Red, Line Color is RED", ` Color is (Red|RED)`, ": $1"},
		{"(123) 456-7890\n12-345-6789", "123-456-7890\n12-345-6789", `(\d{3})-(\d{3})-(\d{4})`, "($1) $2-$3"},
	}

	for i, c := range cs {
		re := regexp.MustCompile(c.m)
		a := SubmatchReplace(re, []byte(c.s), []byte(c.t))
		if string(a) != c.w {
			t.Errorf("[%d] SubmatchReplace(%q, %q, %q) = %q, want %q", i, c.m, c.s, c.t, string(a), c.w)
		}
	}
}

func TestSubmatchReplaceString(t *testing.T) {
	cs := []struct {
		w string
		s string
		m string
		t string
	}{
		{"〒(123) 456-7890 〒", "〒123-456-7890 〒", `(\d{3})-(\d{3})-(\d{4})`, "($1) $2-$3"},
		{"Border: Red, Line: RED", "Border Color is Red, Line Color is RED", ` Color is (Red|RED)`, ": $1"},
		{"(123) 456-7890\n12-345-6789", "123-456-7890\n12-345-6789", `(\d{3})-(\d{3})-(\d{4})`, "($1) $2-$3"},
	}

	for i, c := range cs {
		re := regexp.MustCompile(c.m)
		a := SubmatchReplaceString(re, c.s, c.t)
		if a != c.w {
			t.Errorf("[%d] SubmatchReplaceString(%q, %q, %q) = %q, want %q", i, c.m, c.s, c.t, a, c.w)
		}
	}
}

func TestSubmatchReplacer(t *testing.T) {
	cs := []struct {
		w string
		s string
		m string
		t string
	}{
		{
			`go &lt;to&gt; <a href="http://a.b.com?a=1&amp;b=2">abc</a> ok?`,
			"go <to> [abc](http://a.b.com?a=1&b=2) ok?",
			`\[([^\]]*)\]\((https?:\/\/[\w~!@#\$%&\*\(\)_\-\+=\[\]\|:;,\.\?\/']+)\)`,
			`<a href="$2">$1</a>`,
		},
	}

	sr := SubmatchReplacer{
		Converter: func(n int, name, value string) string {
			return html.EscapeString(value)
		},
	}

	for i, c := range cs {
		sr.Pattern = regexp.MustCompile(c.m)
		sr.Template = c.t

		a := sr.Replace(c.s)
		if a != c.w {
			t.Errorf("[%d] SubmatchReplacer.Replace(%q, %q, %q) = %q, want %q", i, c.m, c.s, c.t, a, c.w)
		}
	}
}
