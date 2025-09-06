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

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		// ✅ Valid durations
		{"zero", "0", 0, false},
		{"seconds", "10s", 10 * time.Second, false},
		{"minutes", "2m", 2 * time.Minute, false},
		{"hours", "3h", 3 * time.Hour, false},
		{"days", "1d", 24 * time.Hour, false},
		{"milliseconds", "500ms", 500 * time.Millisecond, false},
		{"microseconds (us)", "250us", 250 * time.Microsecond, false},
		{"nanoseconds", "100ns", 100 * time.Nanosecond, false},

		// ✅ With fractions
		{"fractional seconds", "1.5s", 1500 * time.Millisecond, false},
		{"fractional minutes", "0.5m", 30 * time.Second, false},
		{"fractional hours", "1.25h", 75 * time.Minute, false},

		// ✅ Mixed units
		{"hours and minutes", "1h30m", 90 * time.Minute, false},
		{"complex", "2h45m30s", 2*time.Hour + 45*time.Minute + 30*time.Second, false},

		// ✅ Negative durations
		{"negative duration", "-2h", -2 * time.Hour, false},
		{"negative mixed", "-1h30m", -(time.Hour + 30*time.Minute), false},

		// ❌ Errors
		{"empty string", "", 0, true},
		{"missing unit", "123", 0, true},
		{"unknown unit", "10x", 0, true},
		{"just dot", ".", 0, true},
		{"dot unit", ".s", 0, true},
		{"overflow value", "999999999999999999999999h", 0, true},
		{"invalid hour", "25h", 25 * time.Hour, false}, // allowed but large
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseDuration(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseDuration(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
