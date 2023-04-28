package fdk

import (
	"time"
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

func (d *Date) String() string {
	return d.Time.Format(DateFormat)
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return []byte(d.Time.Format(jsonDateFormat)), nil
}

func (d *Date) UnmarshalJSON(data []byte) (err error) {
	// Ignore null, like in the main JSON package.
	str := string(data)
	if str == "null" {
		return
	}

	d.Time, err = time.Parse(jsonDateFormat, str)
	return
}

type Time struct {
	time.Time
}

func (t *Time) String() string {
	return t.Time.Format(TimeFormat)
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(t.Time.Format(jsonTimeFormat)), nil
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	// Ignore null, like in the main JSON package.
	str := string(data)
	if str == "null" {
		return
	}

	t.Time, err = time.Parse(jsonTimeFormat, str)
	return
}
