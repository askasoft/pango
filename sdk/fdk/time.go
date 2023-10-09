package fdk

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/askasoft/pango/bye"
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
	return []byte(d.Time.Format(jsonDateFormat)), nil
}

func (d *Date) UnmarshalJSON(data []byte) (err error) {
	// Ignore null, like in the main JSON package.
	js := bye.UnsafeString(data)
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
	return []byte(t.Time.UTC().Format(jsonTimeFormat)), nil
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	// Ignore null, like in the main JSON package.
	js := bye.UnsafeString(data)
	if js == "null" {
		return
	}

	t.Time, err = time.ParseInLocation(jsonTimeFormat, js, time.UTC)
	return
}

type TimeSpent int

func ParseTimeSpent(s string) (ts TimeSpent, err error) {
	min, sec := 0, 0

	s1, s2, ok := str.Cut(s, ":")
	if ok {
		min, err = strconv.Atoi(s1)
		if err == nil {
			sec, err = strconv.Atoi(s2)
			if err == nil {
				ts = TimeSpent(min*60 + sec)
				return //nolint: nilerr
			}
		}
	} else {
		sec, err = strconv.Atoi(s1)
		if err == nil {
			ts = TimeSpent(sec)
			return //nolint: nilerr
		}
	}

	return 0, fmt.Errorf(`ParseTimeSpent: "%s" is not a HH:MM string`, s)
}

func (ts TimeSpent) Seconds() int {
	return int(ts)
}

func (ts TimeSpent) String() string {
	min, sec := ts/60, ts%60
	return fmt.Sprintf("%02d:%02d", min, sec)
}

func (ts TimeSpent) MarshalJSON() ([]byte, error) {
	min, sec := ts/60, ts%60
	return []byte(fmt.Sprintf(`"%02d:%02d"`, min, sec)), nil
}

func (ts *TimeSpent) UnmarshalJSON(data []byte) (err error) {
	// Ignore null, like in the main JSON package.
	js := bye.UnsafeString(data)
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
