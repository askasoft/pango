package str

import (
	"regexp"
	"testing"
)

func TestRegexpSubmatchReplace(t *testing.T) {
	cs := []struct {
		w string
		r string
		s string
		p string
	}{
		{"〒(123) 456-7890", `(\d{3})-(\d{3})-(\d{4})`, "〒123-456-7890", "($1) $2-$3"},
		{"Border: Red, Line: RED", ` Color is (Red|RED)`, "Border Color is Red, Line Color is RED", ": $1"},
	}

	for i, c := range cs {
		re := regexp.MustCompile(c.r)
		a := RegexpSubmatchReplace(re, c.s, c.p)
		if a != c.w {
			t.Errorf("[%d] RegexpSubmatchReplace(%q, %q, %q) = %q, want %q", i, c.r, c.s, c.p, a, c.w)
		}
	}
}

func TestRegexpSubmatchReplacer(t *testing.T) {
	cs := []struct {
		w   string
		s   string
		rps []string
	}{
		{"123-456-7890\n12-345-6789", "123-456-7890\n12-345-6789", []string{}},
		{"(123) 456-7890\n[12] 345-6789", "123-456-7890\n12-345-6789", []string{`(\d{3})-(\d{3})-(\d{4})`, "($1) $2-$3", `(\d{2})-(\d{3})-(\d{4})`, "[$1] $2-$3"}},
	}

	for i, c := range cs {
		rsr, err := NewRegexpSubmatchReplacer(c.rps...)
		if err != nil {
			t.Fatalf("[%d] NewRegexpSubmatchReplacer(%v) = %v", i, c.rps, err)
		}

		a := rsr.Replace(c.s)
		if a != c.w {
			t.Errorf("[%d] RegexpSubmatchReplacer.Replace(%q) = %q, want %q", i, c.s, a, c.w)
		}
	}
}
