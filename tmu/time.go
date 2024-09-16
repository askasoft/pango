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
