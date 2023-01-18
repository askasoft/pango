package fdk

import (
	"time"
)

const TimeFormat = "2006-01-02T15:04:05Z"

type Time struct {
	time.Time
}

func (t *Time) String() string {
	return t.Time.Format(TimeFormat)
}

const jsonTimeFormat = `"2006-01-02T15:04:05Z"`

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
