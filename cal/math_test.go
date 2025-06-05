package cal

import (
	"errors"
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	cs := []struct {
		a, b, w any
		e       error
	}{
		{1, 2, 3, nil},
		{1, uint(2), 3, nil},
		{1, 2.0, 3.0, nil},
		{uint(1), 2, 3, nil},
		{uint(1), uint(2), uint(3), nil},
		{uint(1), 2.0, 3.0, nil},
		{1.0, 2, 3.0, nil},
		{1.0, uint(2), 3.0, nil},
		{1.0, 2.0, 3.0, nil},
		{"foo", 2, "foo2", nil},
		{2, "foo", "2foo", nil},
	}

	for i, c := range cs {
		r, e := Add(c.a, c.b)
		if c.w != r || fmt.Sprint(c.e) != fmt.Sprint(e) {
			t.Errorf("[%d] Add(%v, %v) = (%T, %v, %v), want: (%T, %v, %v)", i, c.a, c.b, r, r, e, c.w, c.w, c.e)
		}
	}
}

func TestSubtract(t *testing.T) {
	cs := []struct {
		a, b, w any
		e       error
	}{
		{2, 1, 1, nil},
		{uint(2), 1, 1, nil},
		{2.0, 1, 1.0, nil},
		{2, uint(1), 1, nil},
		{uint(2), uint(1), uint(1), nil},
		{2.0, uint(1), 1.0, nil},
		{2, 1.0, 1.0, nil},
		{uint(2), 1.0, 1.0, nil},
		{2.0, 1.0, 1.0, nil},
		{"foo", 2, "foo", errors.New("Subtract: unknown type for 'string'")},
	}

	for i, c := range cs {
		r, e := Subtract(c.a, c.b)
		if c.w != r || fmt.Sprint(c.e) != fmt.Sprint(e) {
			t.Errorf("[%d] Subtract(%v, %v) = (%T, %v, %v), want: (%T, %v, %v)", i, c.a, c.b, r, r, e, c.w, c.w, c.e)
		}
	}
}

func TestNegate(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		want      any
		expectErr bool
	}{
		// Signed integers
		{"int", int(10), int(-10), false},
		{"int8", int8(10), int8(-10), false},
		{"int16", int16(10), int16(-10), false},
		{"int32", int32(10), int32(-10), false},
		{"int64", int64(10), int64(-10), false},

		// Unsigned integers
		{"uint", uint(10), int(-10), false},
		{"uint8", uint8(10), int8(-10), false},
		{"uint16", uint16(10), int16(-10), false},
		{"uint32", uint32(10), int32(-10), false},
		{"uint64", uint64(10), int64(-10), false},

		// Floating points
		{"float32", float32(3.14), float32(-3.14), false},
		{"float64", float64(2.718), float64(-2.718), false},

		// Unknown type
		{"string", "hello", "hello", true},
		{"nil", nil, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Negate(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("Negate(%v) = %v, want error", tt.input, got)
				}
			} else {
				if err != nil {
					t.Errorf("Negate(%v) = %v", tt.input, err)
				}
				if tt.want != got {
					t.Errorf("Negate(%v) = %v, want %v", tt.input, got, tt.want)
				}
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	cs := []struct {
		a, b, w any
		e       error
	}{
		{2, 3, 6, nil},
		{uint(3), 2, 6, nil},
		{3.0, 2, 6.0, nil},
		{3, uint(2), 6, nil},
		{uint(3), uint(2), uint(6), nil},
		{3.0, uint(2), 6.0, nil},
		{3.0, 2, 6.0, nil},
		{uint(3), 2.0, 6.0, nil},
		{3.0, 2.0, 6.0, nil},
		{"foo", 2, "foo", errors.New("Multiply: unknown type for 'string'")},
	}

	for i, c := range cs {
		r, e := Multiply(c.a, c.b)
		if c.w != r || fmt.Sprint(c.e) != fmt.Sprint(e) {
			t.Errorf("[%d] Multiply(%v, %v) = (%T, %v, %v), want: (%T, %v, %v)", i, c.a, c.b, r, r, e, c.w, c.w, c.e)
		}
	}
}

func TestDivide(t *testing.T) {
	cs := []struct {
		a, b, w any
		e       error
	}{
		{10, 2, 5, nil},
		{uint(10), 2, 5, nil},
		{10.0, 2, 5.0, nil},
		{10, uint(2), 5, nil},
		{uint(10), uint(2), uint(5), nil},
		{10.0, uint(2), 5.0, nil},
		{10, 2.0, 5.0, nil},
		{uint(10), 2.0, 5.0, nil},
		{10.0, 2.0, 5.0, nil},
		{"foo", 2, "foo", errors.New("Divide: unknown type for 'string'")},
	}

	for i, c := range cs {
		r, e := Divide(c.a, c.b)
		if c.w != r || fmt.Sprint(c.e) != fmt.Sprint(e) {
			t.Errorf("[%d] Divide(%v, %v) = (%T, %v, %v), want: (%T, %v, %v)", i, c.a, c.b, r, r, e, c.w, c.w, c.e)
		}
	}
}
