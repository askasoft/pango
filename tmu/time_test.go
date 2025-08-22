package tmu

import (
	"testing"
	"time"
)

func TestLocalFormat(t *testing.T) {
	layout := "2006-01-02 15:04:05"

	// A fixed reference time (UTC)
	refTime := time.Date(2023, 5, 10, 14, 30, 0, 0, time.UTC)
	refLocal := refTime.Local().Format(layout)

	tests := []struct {
		name   string
		input  any
		format string
		want   string
	}{
		{
			name:   "nil input",
			input:  nil,
			format: layout,
			want:   "",
		},
		{
			name:   "zero time",
			input:  time.Time{},
			format: layout,
			want:   "",
		},
		{
			name:   "non-zero time",
			input:  refTime,
			format: layout,
			want:   refLocal,
		},
		{
			name:   "pointer to non-zero time",
			input:  &refTime,
			format: layout,
			want:   refLocal,
		},
		{
			name:   "pointer to zero time",
			input:  &time.Time{},
			format: layout,
			want:   "",
		},
		{
			name:   "integer input",
			input:  123,
			format: layout,
			want:   "123",
		},
		{
			name:   "string input",
			input:  "hello",
			format: layout,
			want:   "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LocalFormat(tt.input, tt.format)
			if got != tt.want {
				t.Errorf("LocalFormat(%v, %q) = %q, want %q",
					tt.input, tt.format, got, tt.want)
			}
		})
	}
}

func TestLocalFormatHelpers(t *testing.T) {
	refTime := time.Date(2024, 12, 25, 8, 45, 30, 0, time.UTC)

	tests := []struct {
		name string
		f    func(any) string
		arg  any
		want string
	}{
		{
			name: "LocalFormatDateTime with time.Time",
			f:    LocalFormatDateTime,
			arg:  refTime,
			want: refTime.Local().Format(time.DateTime),
		},
		{
			name: "LocalFormatDate with *time.Time",
			f:    LocalFormatDate,
			arg:  &refTime,
			want: refTime.Local().Format(time.DateOnly),
		},
		{
			name: "LocalFormatTime with time.Time",
			f:    LocalFormatTime,
			arg:  refTime,
			want: refTime.Local().Format(time.TimeOnly),
		},
		{
			name: "LocalFormatDateTime with string input",
			f:    LocalFormatDateTime,
			arg:  "not-a-time",
			want: "not-a-time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.f(tt.arg)
			if got != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	cs := []struct {
		s string
		w string
	}{
		{"2020-2-3 06:07:08", "2020-02-03T06:07:08"},
		{"2020-02-03 06:07:08", "2020-02-03T06:07:08"},
		{"2020-12-13 06:07:08", "2020-12-13T06:07:08"},
		{"2020-2-3", "2020-02-03T00:00:00"},
		{"2020-02-03", "2020-02-03T00:00:00"},
		{"2020-12-13", "2020-12-13T00:00:00"},
		{"06:07:08", "0000-01-01T06:07:08"},
	}

	for i, c := range cs {
		d, err := Parse(c.s)
		if err != nil {
			t.Fatal(err)
		}

		a := d.Format("2006-01-02T15:04:05")
		if err != nil || a != c.w {
			t.Errorf("[%d] Parse(%q) = %q, want %q", i, c.s, a, c.w)
		}
	}
}

func TestAddYear(t *testing.T) {
	cs := []struct {
		s string
		y int
		w string
	}{
		{"2020-02-29", 1, "2021-02-28"},
		{"2020-02-29", 10, "2030-02-28"},
		{"2020-02-29", 12, "2032-02-29"},
		{"2020-02-29", -1, "2019-02-28"},
		{"2020-02-29", -10, "2010-02-28"},
		{"2020-02-29", -12, "2008-02-29"},
	}

	for i, c := range cs {
		d, err := time.Parse("2006-01-02", c.s)
		if err != nil {
			t.Fatal(err)
		}

		a := AddYear(d, c.y).Format("2006-01-02")
		if err != nil || a != c.w {
			t.Errorf("[%d] AddYear(%q, %v) = %q, want %q", i, c.s, c.y, a, c.w)
		}
	}
}

func TestAddMonth(t *testing.T) {
	cs := []struct {
		s string
		m int
		w string
	}{
		{"2019-11-30", 1, "2019-12-30"},
		{"2019-12-31", 0, "2019-12-31"},
		{"2019-12-31", 1, "2020-01-31"},
		{"2019-12-31", 2, "2020-02-29"},
		{"2019-12-31", 3, "2020-03-31"},
		{"2019-12-31", 12, "2020-12-31"},
		{"2019-12-31", 13, "2021-01-31"},
		{"2019-12-31", 14, "2021-02-28"},
		{"2019-12-31", 15, "2021-03-31"},
		{"2020-01-31", 1, "2020-02-29"},
		{"2020-01-31", 2, "2020-03-31"},
		{"2020-01-31", 3, "2020-04-30"},
		{"2019-01-31", 13, "2020-02-29"},
		{"2019-01-31", 14, "2020-03-31"},
		{"2019-01-31", 15, "2020-04-30"},

		{"2019-12-31", -1, "2019-11-30"},
		{"2019-12-31", -2, "2019-10-31"},
		{"2019-12-31", -3, "2019-09-30"},
		{"2019-12-31", -12, "2018-12-31"},
		{"2019-12-31", -13, "2018-11-30"},
		{"2019-12-31", -14, "2018-10-31"},
		{"2019-12-31", -15, "2018-09-30"},
		{"2020-03-31", -1, "2020-02-29"},
		{"2020-03-31", -2, "2020-01-31"},
		{"2020-03-31", -3, "2019-12-31"},
		{"2020-03-31", -4, "2019-11-30"},
	}

	for i, c := range cs {
		d, err := time.Parse("2006-01-02", c.s)
		if err != nil {
			t.Fatal(err)
		}

		a := AddMonth(d, c.m).Format("2006-01-02")
		if err != nil || a != c.w {
			t.Errorf("[%d] AddMonth(%q, %v) = %q, want %q", i, c.s, c.m, a, c.w)
		}
	}
}

func TestTruncateHours(t *testing.T) {
	cs := []struct {
		s string
		w string
	}{
		{"2020-01-02T01:02:03Z", "2020-01-02T00:00:00Z"},
		{"2020-01-02T01:02:03+09:00", "2020-01-02T00:00:00+09:00"},
	}

	for i, c := range cs {
		d, err := Parse(c.s)
		if err != nil {
			t.Fatal(err)
		}

		a := TruncateHours(d).Format(time.RFC3339)
		if err != nil || a != c.w {
			t.Errorf("[%d] TruncateHours(%q) = %q, want %q", i, c.s, a, c.w)
		}
	}
}
