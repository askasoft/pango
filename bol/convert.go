package bol

import (
	"strconv"
)

// Atob use strconv.ParseBool(s, 0, strconv.IntSize) to parse string 's' to int, return n[0] if error.
func Atob(s string, b ...bool) bool {
	if s == "" {
		if len(b) > 0 {
			return b[0]
		}
		return false
	}

	r, err := strconv.ParseBool(s)
	if err != nil && len(b) > 0 {
		return b[0]
	}
	return r
}

func Btoa(b bool) string {
	return strconv.FormatBool(b)
}
