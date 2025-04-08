package urlx

import (
	"testing"
)

func TestCleanURL(t *testing.T) {
	cs := []struct {
		s string
		w string
	}{
		{"http://user:pass@abc.com", "http://abc.com/"},
		{"http://user:pass@abc.com/?a=b#xyz", "http://abc.com/"},
	}

	for i, c := range cs {
		a, err := CleanURL(c.s)
		if err != nil {
			t.Errorf("#%d CleanURL(%q)\n     = %s\nWANT: %s", i, c.s, a, c.w)
		}
	}
}

func TestParentURL(t *testing.T) {
	cs := []struct {
		s string
		w string
	}{
		{"http://user:pass@abc.com/?a=b#xyz", "http://abc.com/"},
		{"http://user:pass@abc.com/abc?a=b#xyz", "http://abc.com/"},
		{"http://user:pass@abc.com/abc/?a=b#xyz", "http://abc.com/abc/"},
	}

	for i, c := range cs {
		a, err := ParentURL(c.s)
		if err != nil {
			t.Errorf("#%d ParentURL(%q)\n     = %s\nWANT: %s", i, c.s, a, c.w)
		}
	}
}
