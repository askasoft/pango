package num

import "testing"

func TestIfZero(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"a non-zero", 5, 10, 5},
		{"a zero", 0, 10, 10},
		{"both zero", 0, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IfZero(tt.a, tt.b); got != tt.want {
				t.Errorf("IfZero(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestAbsInt(t *testing.T) {
	cases := []struct {
		in, want int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
	}
	for _, c := range cases {
		if got := Abs(c.in); got != c.want {
			t.Errorf("AbsInt(%d) = %d; want %d", c.in, got, c.want)
		}
	}
}

func TestAbsInt16(t *testing.T) {
	cases := []struct {
		in, want int16
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
	}
	for _, c := range cases {
		if got := Abs(c.in); got != c.want {
			t.Errorf("AbsInt16(%d) = %d; want %d", c.in, got, c.want)
		}
	}
}

func TestAbsInt32(t *testing.T) {
	cases := []struct {
		in, want int32
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
	}
	for _, c := range cases {
		if got := Abs(c.in); got != c.want {
			t.Errorf("AbsInt32(%d) = %d; want %d", c.in, got, c.want)
		}
	}
}

func TestAbsInt64(t *testing.T) {
	cases := []struct {
		in, want int64
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
	}
	for _, c := range cases {
		if got := Abs(c.in); got != c.want {
			t.Errorf("AbsInt64(%d) = %d; want %d", c.in, got, c.want)
		}
	}
}

func TestAbsFloat32(t *testing.T) {
	cases := []struct {
		in, want float32
	}{
		{5.5, 5.5},
		{-5.5, 5.5},
		{0, 0},
	}
	for _, c := range cases {
		if got := Abs(c.in); got != c.want {
			t.Errorf("AbsFloat32(%f) = %f; want %f", c.in, got, c.want)
		}
	}
}

func TestAbsFloat64(t *testing.T) {
	cases := []struct {
		in, want float64
	}{
		{5.5, 5.5},
		{-5.5, 5.5},
		{0, 0},
	}
	for _, c := range cases {
		if got := Abs(c.in); got != c.want {
			t.Errorf("AbsFloat64(%f) = %f; want %f", c.in, got, c.want)
		}
	}
}
