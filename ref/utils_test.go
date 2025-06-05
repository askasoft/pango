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

type Example struct {
	Value    int
	Multiply func(a, b int) int
}

func (e Example) Greet(name string) string {
	return "Hello, " + name
}

func (e *Example) Sum(a, b int) int {
	return a + b + e.Value
}

func TestInvokeMethod_MethodCall(t *testing.T) {
	obj := &Example{Value: 2}

	result, err := InvokeMethod(obj, "Sum", 3, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 || result[0] != 10 {
		t.Errorf("expected result 10, got %v", result)
	}
}

func TestInvokeMethod_FieldFuncCall(t *testing.T) {
	obj := Example{
		Multiply: func(a, b int) int {
			return a * b
		},
	}

	result, err := InvokeMethod(obj, "Multiply", 4, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 || result[0] != 20 {
		t.Errorf("expected 20, got %v", result)
	}
}

func TestInvokeMethod_EmptyName(t *testing.T) {
	_, err := InvokeMethod(Example{}, "")
	if err == nil || err.Error() != "ref: empty function name" {
		t.Errorf("expected error for empty method name, got %v", err)
	}
}

func TestInvokeMethod_InvalidMethodName(t *testing.T) {
	_, err := InvokeMethod(Example{}, "DoesNotExist")
	if err == nil { // just check error presence
		t.Errorf("expected error for missing method, got nil")
	}
}

func TestInvokeMethod_InvalidArgCount(t *testing.T) {
	obj := Example{
		Multiply: func(a, b int) int {
			return a * b
		},
	}

	_, err := InvokeMethod(obj, "Multiply", 3) // missing one argument
	if err == nil {
		t.Errorf("expected error for argument count mismatch, got %v", err)
	}
}

func TestInvokeMethod_InvalidArgType(t *testing.T) {
	obj := Example{
		Multiply: func(a, b int) int {
			return a * b
		},
	}

	_, err := InvokeMethod(obj, "Multiply", "three", "five") // wrong types
	if err == nil {
		t.Error("expected error for invalid argument types, got nil")
	}
}
