package num

import (
	"testing"
)

func TestAtoi(t *testing.T) {
	cs := []struct {
		w int
		s string
		n []int
	}{
		{1, "1", nil},
		{2, "02", nil},
		{0660, "0660", nil},
		{0xf, "0xf", nil},
		{-1, "a", []int{0, -1}},
	}

	for i, c := range cs {
		a := Atoi(c.s, c.n...)
		if a != c.w {
			t.Errorf("[%d] Atoi(%s, %v) = %v, want %v", i, c.s, c.n, a, c.w)
		}
	}
}

func TestAtol(t *testing.T) {
	cs := []struct {
		w int64
		s string
		n []int64
	}{
		{1, "1", nil},
		{2, "02", nil},
		{0660, "0660", nil},
		{0xf, "0xf", nil},
		{-1, "a", []int64{0, -1}},
	}

	for i, c := range cs {
		a := Atol(c.s, c.n...)
		if a != c.w {
			t.Errorf("[%d] Atol(%s, %v) = %v, want %v", i, c.s, c.n, a, c.w)
		}
	}
}

func TestAtof(t *testing.T) {
	cs := []struct {
		w float64
		s string
		n []float64
	}{
		{1, "1", nil},
		{-1, "a", []float64{0, -1}},
	}

	for i, c := range cs {
		a := Atof(c.s, c.n...)
		if a != c.w {
			t.Errorf("[%d] Atof(%s, %v) = %v, want %v", i, c.s, c.n, a, c.w)
		}
	}
}

func TestFtoa(t *testing.T) {
	cs := []struct {
		w string
		n float64
	}{
		{"200", 200},
		{"2", 2},
		{"2.2", 2.2},
		{"2.02", 2.02},
		{"200.02", 200.02},
	}

	for i, c := range cs {
		a := Ftoa(c.n)
		if a != c.w {
			t.Errorf("[%d] Ftoa(%f) = %v, want %v", i, c.n, a, c.w)
		}
	}
}

func TestFtoaWithDigits(t *testing.T) {
	cs := []struct {
		w string
		n float64
		d int
	}{
		{"1", 1.23, 0},
		{"1.2", 1.23, 1},
		{"1.3", 1.26, 1},
		{"1.23", 1.23, 2},
		{"1.23", 1.23, 3},
		{"1.234", 1.234, 3},
		{"1.235", 1.2346, 3},
	}

	for i, c := range cs {
		a := FtoaWithDigits(c.n, c.d)
		if a != c.w {
			t.Errorf("[%d] FtoaWithDigits(%f, %d) = %v, want %v", i, c.n, c.d, a, c.w)
		}
	}
}
