package ldt

import "testing"

func TestAbs(t *testing.T) {
	tests := map[int]int{
		1:      1,
		-0:     0,
		69:     69,
		-65535: 65535,
		65535:  65535,
	}

	for x, want := range tests {
		got := abs(x)
		if got != want {
			t.Fatalf("want %d got %d", want, got)
		}
	}
}
