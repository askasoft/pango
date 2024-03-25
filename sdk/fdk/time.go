package fdk

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/askasoft/pango/str"
)

const (
	DateFormat     = "2006-01-02"
	TimeFormat     = "2006-01-02T15:04:05Z"
	jsonDateFormat = `"2006-01-02"`
	jsonTimeFormat = `"2006-01-02T15:04:05Z"`
)

type Date struct {
	time.Time
}

func ParseDate(s string) (*Date, error) {
	t, err := time.ParseInLocation(DateFormat, s, time.UTC)
	if err != nil {
		return nil, err
	}
	return &Date{t}, nil
}

func (d *Date) String() string {
	return d.Time.Format(DateFormat)
}

func (d *Date) MarshalJSON() ([]byte, error) {
	bs := make([]byte, 0, len(jsonDateFormat))
	bs = d.Time.AppendFormat(bs, jsonDateFormat)
	return bs, nil
}

func (d *Date) UnmarshalJSON(data []byte) (err error) {
	// Ignore null, like in the main JSON package.
	js := str.UnsafeString(data)
	if js == "null" {
		return
	}

	d.Time, err = time.Parse(jsonDateFormat, js)
	return
}

type Time struct {
	time.Time
}

func ParseTime(s string) (*Time, error) {
	t, err := time.ParseInLocation(TimeFormat, s, time.UTC)
	if err != nil {
		return nil, err
	}
	return &Time{t}, nil
}

func (t *Time) String() string {
	return t.Time.UTC().Format(TimeFormat)
}

func (t *Time) MarshalJSON() ([]byte, error) {
	bs := make([]byte, 0, len(jsonTimeFormat))
	bs = t.Time.UTC().AppendFormat(bs, jsonTimeFormat)
	return bs, nil
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	// Ignore null, like in the main JSON package.
	js := str.UnsafeString(data)
	if js == "null" {
		return
	}

	t.Time, err = time.ParseInLocation(jsonTimeFormat, js, time.UTC)
	return
}

// TimeSpent HH:MM (Minutes)
type TimeSpent int

func ParseTimeSpent(s string) (TimeSpent, error) {
	s1, s2, ok := str.Cut(s, ":")
	if ok {
		hour, err := strconv.Atoi(s1)
		if err == nil {
			min, err := strconv.Atoi(s2)
			if err == nil {
				return TimeSpent(hour*60 + min), nil
			}
		}
		return 0, fmt.Errorf(`ParseTimeSpent: "%s" is not a HH:MM string`, s)
	}

	min, err := strconv.Atoi(s1)
	if err == nil {
		return TimeSpent(min), nil
	}
	return 0, fmt.Errorf(`ParseTimeSpent: "%s" is not a numeric string`, s)
}

func (ts TimeSpent) Minutes() int {
	return int(ts)
}

func (ts TimeSpent) String() string {
	hour, min := ts/60, ts%60
	return fmt.Sprintf("%02d:%02d", hour, min)
}

func (ts TimeSpent) MarshalJSON() ([]byte, error) {
	hour, min := ts/60, ts%60
	return []byte(fmt.Sprintf(`"%02d:%02d"`, hour, min)), nil
}

func (ts *TimeSpent) UnmarshalJSON(data []byte) (err error) {
	// Ignore null, like in the main JSON package.
	js := str.UnsafeString(data)
	if js == "null" {
		return
	}

	if len(js) < 2 || js[0] != '"' || js[len(js)-1] != '"' {
		return errors.New("TimeSpent.UnmarshalJSON: input is not a JSON string")
	}
	js = js[len(`"`) : len(js)-len(`"`)]

	*ts, err = ParseTimeSpent(js)
	return
}
