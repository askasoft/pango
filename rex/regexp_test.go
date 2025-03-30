package rex

import (
	"html"
	"regexp"
	"testing"
)

func TestReplaceAllConvertString(t *testing.T) {
	cs := []struct {
		w string
		s string
		m string
		t string
	}{
		{
			`タイマー|たいまー|Timer|
|対処方法|たいしょほうほう|対処法|対応策`,
			`|タイマー|たいまー|Timer|
|対処方法|たいしょほうほう|対処法|対応策|`,
			`(^\s*\|+)|(\s*\|+$)`,
			``,
		},
		{
			`a|b|c`,
			"a+b+c",
			`\+`,
			`|`,
		},
		{
			`go &lt;to&gt; <a></a> ok?`,
			"go <to> [abc](http://a.b.com?a=1&b=2) ok?",
			`\[([^\]]*)\]\((https?:\/\/[\w~!@#\$%&\*\(\)_\-\+=\[\]\|:;,\.\?\/']+)\)`,
			`<a></a>`,
		},
		{
			`go &lt;to&gt; <a href="http://a.b.com?a=1&amp;b=2">abc</a> ok?`,
			"go <to> [abc](http://a.b.com?a=1&b=2) ok?",
			`\[([^\]]*)\]\((https?:\/\/[\w~!@#\$%&\*\(\)_\-\+=\[\]\|:;,\.\?\/']+)\)`,
			`<a href="$2">$1</a>`,
		},
	}

	cv := func(n int, name, value string) string {
		return html.EscapeString(value)
	}

	for i, c := range cs {
		re := regexp.MustCompile(c.m)

		a := ReplaceAllConvertString(c.s, re, c.t, cv)
		if a != c.w {
			t.Errorf("[%d] ReplaceAllConvertString(%q, %q, %q) = \n%q\nwant:\n%q", i, c.m, c.s, c.t, a, c.w)
		}
	}
}
