package squ

import (
	"testing"
)

func TestEscapeLike(t *testing.T) {
	cs := []struct {
		s, w string
	}{
		{"ab~c", "ab~~c"},
		{"ab%", "ab~%"},
		{"ab_", "ab~_"},
	}

	for i, c := range cs {
		a := EscapeLike(c.s)
		if a != c.w {
			t.Errorf("[%d] EscapeLike(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestEscapeString(t *testing.T) {
	cs := []struct {
		s, w string
	}{
		{"ab'c", "ab''c"},
	}

	for i, c := range cs {
		a := EscapeString(c.s)
		if a != c.w {
			t.Errorf("[%d] EscapeString(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestStringLike(t *testing.T) {
	cs := []struct {
		s, w string
	}{
		{"abc", "%abc%"},
		{"ab~c", "%ab~~c%"},
		{"ab%", "%ab~%%"},
		{"ab_", "%ab~_%"},
	}

	for i, c := range cs {
		a := StringLike(c.s)
		if a != c.w {
			t.Errorf("[%d] StringLike(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestStartsLike(t *testing.T) {
	cs := []struct {
		s, w string
	}{
		{"abc", "abc%"},
		{"ab~c", "ab~~c%"},
		{"ab%", "ab~%%"},
		{"ab_", "ab~_%"},
	}

	for i, c := range cs {
		a := StartsLike(c.s)
		if a != c.w {
			t.Errorf("[%d] StartsLike(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestEndsLike(t *testing.T) {
	cs := []struct {
		s, w string
	}{
		{"abc", "%abc"},
		{"ab~c", "%ab~~c"},
		{"ab%", "%ab~%"},
		{"ab_", "%ab~_"},
	}

	for i, c := range cs {
		a := EndsLike(c.s)
		if a != c.w {
			t.Errorf("[%d] EndsLike(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}
