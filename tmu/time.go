package tmu

import "time"

// Atod convert string to time.Duration.
// if not found or convert error, returns the defs[0] or zero.
func Atod(s string, defs ...time.Duration) time.Duration {
	if d, err := time.ParseDuration(s); err == nil {
		return d
	}
	if len(defs) > 0 {
		return defs[0]
	}
	return 0
}

// TruncateHours tuncate hours for time t.
// Returns the time (local: yyyy-MM-dd 00:00:00).
func TruncateHours(t time.Time) time.Time {
	return t.Truncate(time.Hour).Add(-time.Duration(t.Hour()) * time.Hour)
}

// TruncateMinutes tuncate minute for time t.
// Returns the time (yyyy-MM-dd hh:00:00).
func TruncateMinutes(t time.Time) time.Time {
	return t.Truncate(time.Hour)
}

// TruncateSeconds tuncate second for time t.
// Returns the time (yyyy-MM-dd hh:mm:00).
func TruncateSeconds(t time.Time) time.Time {
	return t.Truncate(time.Minute)
}

var GeneralLayouts = []string{time.RFC3339, "2006-1-2 15:04:05", "2006-1-2", "15:04:05"}

func ParseInLocation(value string, loc *time.Location, layouts ...string) (tt time.Time, err error) {
	if len(layouts) == 0 {
		layouts = GeneralLayouts
	}

	for _, f := range layouts {
		if tt, err = time.ParseInLocation(f, value, time.Local); err == nil {
			return //nolint: nilerr
		}
	}
	return
}

func Parse(value string, layouts ...string) (tt time.Time, err error) {
	if len(layouts) == 0 {
		layouts = GeneralLayouts
	}

	for _, f := range layouts {
		if tt, err = time.Parse(f, value); err == nil {
			return //nolint: nilerr
		}
	}
	return
}

func IsLeapYear(t time.Time) bool {
	return t.YearDay() > 365
}

// AddYear(2020-02-29, 1) = 2021-02-28
// AddYear(2020-02-29, 10) = 2030-02-28
// AddYear(2020-02-29, 12) = 2032-02-29
// AddYear(2020-02-29, -1) = 2019-02-28
// AddYear(2020-02-29, -10) = 2010-02-28
// AddYear(2020-02-29, -12) = 2008-02-29
func AddYear(t time.Time, y int) time.Time {
	n := t.AddDate(y, 0, 0)
	if t.Month() != n.Month() {
		n = n.AddDate(0, 0, -1)
	}
	return n
}

// AddMonth2019-12-31, 1) = 2020-01-31
// AddMonth2019-12-31, 2) = 2020-02-29
// AddMonth2019-12-31, 3) = 2020-03-31
// AddMonth2019-12-31, 12) = 2020-12-31
// AddMonth2019-12-31, 13) = 2021-01-31
// AddMonth2019-12-31, 14) = 2021-02-28
// AddMonth2020-03-31, -1) = 2020-02-29
// AddMonth2020-03-31, -2) = 2020-01-31
// AddMonth2020-03-31, -3) = 2019-12-31
// AddMonth2020-03-31, -4) = 2019-11-30
func AddMonth(t time.Time, m int) time.Time {
	n := t

	y := m / 12
	m = m % 12

	if m != 0 {
		n = t.AddDate(0, m, 0)
	}

	if y != 0 {
		n = AddYear(n, y)
	}

	if m != 0 && (int(t.Month())+m)%12 != int(n.Month()%12) {
		n = n.AddDate(0, 0, -n.Day())
	}

	return n
}
