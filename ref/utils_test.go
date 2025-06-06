package ref

import (
	"strings"
	"testing"
	"time"
)

func somefunction() {
	// this empty function is used by TestFunctionName()
}

func TestNameOfFunc(t *testing.T) {
	a := NameOfFunc(somefunction)
	if !strings.HasSuffix(a, "github.com/askasoft/pango/ref.somefunction") {
		t.Errorf("NameOfFunc(somefunction) = %v", a)
	}
}

func TestIsZero(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  bool
	}{
		{"nil", nil, true},
		{"true", true, false},
		{"false", false, true},
		{"zero int", 0, true},
		{"non-zero int", 42, false},
		{"zero string", "", true},
		{"non-zero string", "hello", false},
		{"zero struct", struct{}{}, true},
		{"zero time.Time", time.Time{}, true},
		{"non-zero time.Time", time.Now(), false},
		{"zero slice", []int(nil), true},
		{"non-zero slice", []int{1}, false},
		{"zero pointer", (*int)(nil), true},
		{"non-zero pointer", new(int), false},
		{"zero map", map[string]int(nil), true},
		{"non-zero map", map[string]int{"a": 1}, false},
		{"zero interface", any(nil), true},
		{"non-zero interface", any(42), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsZero(tt.input)
			if got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}
