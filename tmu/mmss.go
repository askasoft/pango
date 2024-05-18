package tmu

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/askasoft/pango/str"
)

// MMSS MM:SS (Seconds)
type MMSS int

func ParseMMSS(s string) (MMSS, error) {
	s1, s2, ok := str.Cut(s, ":")
	if ok {
		min, err := strconv.Atoi(s1)
		if err == nil {
			sec, err := strconv.Atoi(s2)
			if err == nil {
				return MMSS(min*60 + sec), nil
			}
		}
		return 0, fmt.Errorf(`ParseMMSS: "%s" is not a MM:SS string`, s)
	}

	sec, err := strconv.Atoi(s1)
	if err == nil {
		return MMSS(sec), nil
	}
	return 0, fmt.Errorf(`ParseMMSS: "%s" is not a numeric string`, s)
}

func (ms MMSS) Seconds() int {
	return int(ms)
}

func (ms MMSS) String() string {
	min, sec := ms/60, ms%60
	return fmt.Sprintf("%02d:%02d", min, sec)
}

func (ms MMSS) MarshalJSON() ([]byte, error) {
	min, sec := ms/60, ms%60
	return []byte(fmt.Sprintf(`"%02d:%02d"`, min, sec)), nil
}

func (ms *MMSS) UnmarshalJSON(data []byte) (err error) {
	// Ignore null, like in the main JSON package.
	js := str.UnsafeString(data)
	if js == "null" {
		return
	}

	if len(js) < 2 || js[0] != '"' || js[len(js)-1] != '"' {
		return errors.New("MMSS.UnmarshalJSON: input is not a JSON string")
	}
	js = js[len(`"`) : len(js)-len(`"`)]

	*ms, err = ParseMMSS(js)
	return
}
