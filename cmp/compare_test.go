package cmp

import "testing"

func TestCompareInt8(t *testing.T) {
	cs := []struct {
		a int8
		b int8
		o int
	}{
		{1, 2, -1},
		{1, 1, 0},
		{2, 1, 1},
	}

	for i, c := range cs {
		o := CompareInt8(c.a, c.b)
		if o != c.o {
			t.Errorf("[%d] CompareInt8(%v, %v) = %v", i, c.a, c.b, o)
		}
	}
}

func TestCompareInt64(t *testing.T) {
	cs := []struct {
		a int64
		b int64
		o int
	}{
		{0, 9223372036854775807, -1},
		{1, 1, 0},
		{2, 1, 1},
	}

	for i, c := range cs {
		o := CompareInt64(c.a, c.b)
		if o != c.o {
			t.Errorf("[%d] CompareInt64(%v, %v) = %v", i, c.a, c.b, o)
		}
	}
}
