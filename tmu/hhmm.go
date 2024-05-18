package tmu

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/askasoft/pango/str"
)

// HHMM HH:MM (Minutes)
type HHMM int

func ParseHHMM(s string) (HHMM, error) {
	s1, s2, ok := str.Cut(s, ":")
	if ok {
		hour, err := strconv.Atoi(s1)
		if err == nil {
			min, err := strconv.Atoi(s2)
			if err == nil {
				return HHMM(hour*60 + min), nil
			}
		}
		return 0, fmt.Errorf(`ParseHHMM: "%s" is not a HH:MM string`, s)
	}

	min, err := strconv.Atoi(s1)
	if err == nil {
		return HHMM(min), nil
	}
	return 0, fmt.Errorf(`ParseHHMM: "%s" is not a numeric string`, s)
}

func (hm HHMM) Minutes() int {
	return int(hm)
}

func (hm HHMM) String() string {
	hour, min := hm/60, hm%60
	return fmt.Sprintf("%02d:%02d", hour, min)
}

func (hm HHMM) MarshalJSON() ([]byte, error) {
	hour, min := hm/60, hm%60
	return []byte(fmt.Sprintf(`"%02d:%02d"`, hour, min)), nil
}

func (hm *HHMM) UnmarshalJSON(data []byte) (err error) {
	// Ignore null, like in the main JSON package.
	js := str.UnsafeString(data)
	if js == "null" {
		return
	}

	if len(js) < 2 || js[0] != '"' || js[len(js)-1] != '"' {
		return errors.New("HHMM.UnmarshalJSON: input is not a JSON string")
	}
	js = js[len(`"`) : len(js)-len(`"`)]

	*hm, err = ParseHHMM(js)
	return
}
