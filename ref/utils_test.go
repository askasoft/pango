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

func (e *Example) Sums(a int, bs ...int) (r int) {
	r = a
	for _, b := range bs {
		r += b
	}
	return r + e.Value
}

func TestCallMethod_MethodCall(t *testing.T) {
	cs := []struct {
		v int
		m string
		a []any
		w int
		e bool
	}{
		{2, "Sum", []any{3, 5}, 10, false},
		{2, "Sums", []any{}, 0, true},
		{2, "Sums", []any{3}, 5, false},
		{2, "Sums", []any{3, 5}, 10, false},
		{2, "Sums", []any{3, 5, 6}, 16, false},
	}

	for i, c := range cs {
		obj := &Example{Value: c.v}

		r, err := CallMethod(obj, c.m, c.a...)
		if c.e {
			if err == nil {
				t.Fatalf("#%d %s() want error", i, c.m)
			}
			continue
		}

		if err != nil {
			t.Fatalf("#%d %s() unexpected error: %v", i, c.m, err)
		}

		if len(r) != 1 || r[0] != c.w {
			t.Errorf("#%d %s(%v) = %v, want %v", i, c.m, c.a, r, c.w)
		}
	}
}

func TestCallMethod_FieldFuncCall(t *testing.T) {
	obj := Example{
		Multiply: func(a, b int) int {
			return a * b
		},
	}

	result, err := CallMethod(obj, "Multiply", 4, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 || result[0] != 20 {
		t.Errorf("expected 20, got %v", result)
	}
}

func TestCallMethod_EmptyName(t *testing.T) {
	_, err := CallMethod(Example{}, "")
	if err == nil || err.Error() != "ref: empty function name" {
		t.Errorf("expected error for empty method name, got %v", err)
	}
}

func TestCallMethod_InvalidMethodName(t *testing.T) {
	_, err := CallMethod(Example{}, "DoesNotExist")
	if err == nil { // just check error presence
		t.Errorf("expected error for missing method, got nil")
	}
}

func TestCallMethod_InvalidArgCount(t *testing.T) {
	obj := Example{
		Multiply: func(a, b int) int {
			return a * b
		},
	}

	_, err := CallMethod(obj, "Multiply", 3) // missing one argument
	if err == nil {
		t.Errorf("expected error for argument count mismatch, got %v", err)
	}
}

func TestCallMethod_InvalidArgType(t *testing.T) {
	obj := Example{
		Multiply: func(a, b int) int {
			return a * b
		},
	}

	_, err := CallMethod(obj, "Multiply", "three", "five") // wrong types
	if err == nil {
		t.Error("expected error for invalid argument types, got nil")
	}
}

// Sample struct for testing
type TestStruct struct {
	Name   string
	Age    int
	hidden bool // unexported, should not appear
}

func TestStructFieldsToMap_StructValue(t *testing.T) {
	obj := TestStruct{Name: "Alice", Age: 30, hidden: true}

	m, err := StructFieldsToMap(obj)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(m) != 2 {
		t.Fatalf("expected 2 exported fields, got %d", len(m))
	}

	if m["Name"] != "Alice" {
		t.Errorf("expected Name to be 'Alice', got %v", m["Name"])
	}
	if m["Age"] != 30 {
		t.Errorf("expected Age to be 30, got %v", m["Age"])
	}
	if _, ok := m["hidden"]; ok {
		t.Error("unexported field 'hidden' should not be present")
	}
}

func TestStructFieldsToMap_StructPointer(t *testing.T) {
	obj := &TestStruct{Name: "Bob", Age: 25}

	m, err := StructFieldsToMap(obj)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(m) != 2 || m["Name"] != "Bob" || m["Age"] != 25 {
		t.Errorf("unexpected map output: %v", m)
	}
}

func TestStructFieldsToMap_InvalidInput(t *testing.T) {
	cases := []struct {
		input any
	}{
		{input: "not a struct"},
		{input: 123},
		{input: []string{"a", "b"}},
		{input: nil},
	}

	for _, c := range cases {
		_, err := StructFieldsToMap(c.input)
		if err == nil {
			t.Errorf("expected error for input %T, got nil", c.input)
		}
	}
}

func TestStructFieldsToMap_EmptyStruct(t *testing.T) {
	type Empty struct{}
	obj := Empty{}

	m, err := StructFieldsToMap(obj)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}
