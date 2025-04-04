package tmu

import (
	"testing"
	"time"
)

func TestAtod(t *testing.T) {
	cs := []struct {
		w time.Duration
		s string
		n []time.Duration
	}{
		{time.Second, "1s", nil},
		{time.Minute, "x", []time.Duration{0, time.Minute}},
	}

	for i, c := range cs {
		a := Atod(c.s, c.n...)
		if a != c.w {
			t.Errorf("[%d] Atod(%s, %v) = %v, want %v", i, c.s, c.n, a, c.w)
		}
	}
}

func TestHumanDuration(t *testing.T) {
	cs := []struct {
		d time.Duration
		w string
	}{
		{Day + Hour + Minute + Second + Millisecond + Microsecond + Nanosecond, "1d1h1m1s"},
		{-(Day + Hour + Minute + Second + Millisecond + Microsecond + Nanosecond), "-1d1h1m1s"},
		{Day + Minute + Second + Millisecond + Microsecond + Nanosecond, "1d1m1s"},
		{-(Day + Minute + Second + Millisecond + Microsecond + Nanosecond), "-1d1m1s"},
		{Hour + Minute + Second + Millisecond + Microsecond + Nanosecond, "1h1m1s"},
		{-(Hour + Minute + Second + Millisecond + Microsecond + Nanosecond), "-1h1m1s"},
		{Hour + Second + Millisecond + Microsecond + Nanosecond, "1h1s"},
		{-(Hour + Second + Millisecond + Microsecond + Nanosecond), "-1h1s"},
		{Minute + Second + Millisecond + Microsecond + Nanosecond, "1m1s"},
		{-(Minute + Second + Millisecond + Microsecond + Nanosecond), "-1m1s"},
		{Second + Millisecond + Microsecond + Nanosecond, "1.001s"},
		{-(Second + Millisecond + Microsecond + Nanosecond), "-1.001s"},
		{Second + Microsecond + Nanosecond, "1s"},
		{-(Second + Microsecond + Nanosecond), "-1s"},
		{Millisecond + Microsecond + Nanosecond, "1.001ms"},
		{-(Millisecond + Microsecond + Nanosecond), "-1.001ms"},
		{Microsecond + Nanosecond, "1.001µs"},
		{-(Microsecond + Nanosecond), "-1.001µs"},
		{Nanosecond, "1ns"},
		{-Nanosecond, "-1ns"},
	}

	for i, c := range cs {
		a := HumanDuration(c.d)
		if a != c.w {
			t.Errorf("[%d] HumanDuration(%v) = (%v), want %v", i, c.d, a, c.w)
		}
	}
}
