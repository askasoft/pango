package tmu

import (
	"fmt"
	"time"

	"github.com/askasoft/pango/str"
)

const (
	RFC1123   = time.RFC1123  // "Mon, 02 Jan 2006 15:04:05 MST"
	RFC3339   = time.RFC3339  // "2006-01-02T15:04:05Z07:00"
	DateTime  = time.DateTime // "2006-01-02 15:04:05"
	DateOnly  = time.DateOnly // "2006-01-02"
	DateMonth = "2006-01"     // "2006-01"
	TimeOnly  = time.TimeOnly // "15:04:05"
	TimeHHMM  = "15:04"
)

var GeneralLayouts = []string{
	time.RFC3339,
	"2006-1-2 15:04:05",
	"2006-1-2",
	"2006-1",
	TimeOnly,
}

func Format(a any, f string) string {
	if a != nil {
		switch t := a.(type) {
		case time.Time:
			if !t.IsZero() {
				return t.Format(f)
			}
		case *time.Time:
			if t != nil && !t.IsZero() {
				return t.Format(f)
			}
		default:
			return fmt.Sprint(a)
		}
	}
	return ""
}

func FormatDateMonth(a any) string {
	return Format(a, DateMonth)
}

func FormatDateTime(a any) string {
	return Format(a, DateTime)
}

func FormatDate(a any) string {
	return Format(a, DateOnly)
}

func FormatTime(a any) string {
	return Format(a, TimeOnly)
}

func LocalFormat(a any, f string) string {
	if a != nil {
		switch t := a.(type) {
		case time.Time:
			if !t.IsZero() {
				return t.Local().Format(f)
			}
		case *time.Time:
			if t != nil && !t.IsZero() {
				return t.Local().Format(f)
			}
		default:
			return fmt.Sprint(a)
		}
	}
	return ""
}

func LocalFormatDateMonth(a any) string {
	return LocalFormat(a, DateMonth)
}

func LocalFormatDateTime(a any) string {
	return LocalFormat(a, DateTime)
}

func LocalFormatDate(a any) string {
	return LocalFormat(a, DateOnly)
}

func LocalFormatTime(a any) string {
	return LocalFormat(a, TimeOnly)
}

func ParseInLocation(value string, loc *time.Location, layouts ...string) (tt time.Time, err error) {
	if len(layouts) == 0 {
		layouts = GeneralLayouts
	}

	for _, f := range layouts {
		if tt, err = time.ParseInLocation(f, value, time.Local); err == nil {
			return //nolint: nilerr
		}
	}

	err = fmt.Errorf("tmu: cannot parse %q as [%s]", value, str.Join(layouts, ", "))
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

	err = fmt.Errorf("tmu: cannot parse %q as [%s]", value, str.Join(layouts, ", "))
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

// TruncateHours tuncate hours for time t.
// Returns the time (local: yyyy-MM-dd 00:00:00).
func TruncateHours(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	return t.Truncate(time.Hour).Add(-time.Duration(t.Hour()) * time.Hour)
}

// TruncateMinutes tuncate minutes for time t.
// Returns the time (yyyy-MM-dd hh:00:00).
func TruncateMinutes(t time.Time) time.Time {
	return t.Truncate(time.Hour)
}

// TruncateSeconds tuncate seconds for time t.
// Returns the time (yyyy-MM-dd hh:mm:00).
func TruncateSeconds(t time.Time) time.Time {
	return t.Truncate(time.Minute)
}
