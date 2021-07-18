package str

import (
	"testing"
)

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

func TestRemoveAnyByte(t *testing.T) {
	cs := []struct {
		w string
		s string
		b string
	}{
		{"", "", "ab"},
		{"qee", "queued", "ud"},
		{"queued", "queued", "z"},
	}

	for i, c := range cs {
		a := RemoveAnyByte(c.s, c.b)
		if a != c.w {
			t.Errorf("[%d] RemoveAnyByte(%q, %q) = %q, want %q", i, c.s, c.b, a, c.w)
		}
	}
}
