package funcs

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
		{1, 2, int64(3), nil},
		{1, uint(2), int64(3), nil},
		{1, 2.0, float64(3), nil},
		{uint(1), 2, int64(3), nil},
		{uint(1), uint(2), uint64(3), nil},
		{uint(1), 2.0, float64(3), nil},
		{1.0, 2, float64(3), nil},
		{1.0, uint(2), float64(3), nil},
		{1.0, 2.0, float64(3), nil},
		{"foo", 2, nil, errors.New("Add: unknown type for 'foo' (string)")},
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
		{2, 1, int64(1), nil},
		{uint(2), 1, int64(1), nil},
		{2.0, 1, float64(1), nil},
		{2, uint(1), int64(1), nil},
		{uint(2), uint(1), uint64(1), nil},
		{2.0, uint(1), float64(1), nil},
		{2, 1.0, float64(1), nil},
		{uint(2), 1.0, float64(1), nil},
		{2.0, 1.0, float64(1), nil},
		{"foo", 2, nil, errors.New("Subtract: unknown type for 'foo' (string)")},
	}

	for i, c := range cs {
		r, e := Subtract(c.a, c.b)
		if c.w != r || fmt.Sprint(c.e) != fmt.Sprint(e) {
			t.Errorf("[%d] Subtract(%v, %v) = (%T, %v, %v), want: (%T, %v, %v)", i, c.a, c.b, r, r, e, c.w, c.w, c.e)
		}
	}
}

func TestMultiply(t *testing.T) {
	cs := []struct {
		a, b, w any
		e       error
	}{
		{2, 3, int64(6), nil},
		{uint(3), 2, int64(6), nil},
		{3.0, 2, float64(6), nil},
		{3, uint(2), int64(6), nil},
		{uint(3), uint(2), uint64(6), nil},
		{3.0, uint(2), float64(6), nil},
		{3.0, 2, float64(6), nil},
		{uint(3), 2.0, float64(6), nil},
		{3.0, 2.0, float64(6), nil},
		{"foo", 2, nil, errors.New("Multiply: unknown type for 'foo' (string)")},
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
		{10, 2, int64(5), nil},
		{uint(10), 2, int64(5), nil},
		{10.0, 2, float64(5), nil},
		{10, uint(2), int64(5), nil},
		{uint(10), uint(2), uint64(5), nil},
		{10.0, uint(2), float64(5), nil},
		{10, 2.0, float64(5), nil},
		{uint(10), 2.0, float64(5), nil},
		{10.0, 2.0, float64(5), nil},
		{"foo", 2, nil, errors.New("Divide: unknown type for 'foo' (string)")},
	}

	for i, c := range cs {
		r, e := Divide(c.a, c.b)
		if c.w != r || fmt.Sprint(c.e) != fmt.Sprint(e) {
			t.Errorf("[%d] Divide(%v, %v) = (%T, %v, %v), want: (%T, %v, %v)", i, c.a, c.b, r, r, e, c.w, c.w, c.e)
		}
	}
}
