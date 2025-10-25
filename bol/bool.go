package bol

import (
	"strconv"

	"github.com/askasoft/pango/asg"
)

// Atob use strconv.ParseBool(s) to parse string 's' to int,
// returns first value of defs if error.
func Atob(s string, defs ...bool) bool {
	if s == "" {
		return asg.First(defs)
	}

	if r, err := strconv.ParseBool(s); err == nil {
		return r
	}
	return asg.First(defs)
}

// Btoa use strconv.FormatBool(b) to convert a boolean value to its string representation.
func Btoa(b bool) string {
	return strconv.FormatBool(b)
}
