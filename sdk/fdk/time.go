package fdk

import (
	"time"

	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tmu"
)

const (
	DateFormat     = "2006-01-02"
	TimeFormat     = time.RFC3339 //"2006-01-02T15:04:05Z07:00"
	jsonDateFormat = `"` + DateFormat + `"`
	jsonTimeFormat = `"` + TimeFormat + `"`
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

type TimeSpent = tmu.HHMM

func ParseTimeSpent(s string) (TimeSpent, error) {
	return tmu.ParseHHMM(s)
}
