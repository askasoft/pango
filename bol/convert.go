package bol

import (
	"strconv"
)

// NonFalse returns first non-first value of defs if error.
func NonFalse(bs ...bool) bool {
	for _, b := range bs {
		if b {
			return b
		}
	}
	return false
}

// Atob use strconv.ParseBool(s, 0, strconv.IntSize) to parse string 's' to int,
// returns first non-first value of defs if error.
func Atob(s string, defs ...bool) bool {
	if s == "" {
		return NonFalse(defs...)
	}

	if r, err := strconv.ParseBool(s); err == nil {
		return r
	}
	return NonFalse(defs...)
}

func Btoa(b bool) string {
	return strconv.FormatBool(b)
}
