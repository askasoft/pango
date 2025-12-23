package str

import (
	"strings"

	"github.com/askasoft/pango/asg"
)

// ToValidUTF8s returns a copy of the string s with each run of invalid UTF-8 byte sequences
// replaced by the replacement string, which may be empty.
func ToValidUTF8s(ss []string, replacement string) []string {
	for i, s := range ss {
		ss[i] = ToValidUTF8(s, replacement)
	}
	return ss
}

// ToLowers lowercase all string in the string array ss.
func ToLowers(ss []string) []string {
	for i, s := range ss {
		ss[i] = strings.ToLower(s)
	}
	return ss
}

// ToUppers uppercase all string in the string array ss.
func ToUppers(ss []string) []string {
	for i, s := range ss {
		ss[i] = strings.ToUpper(s)
	}
	return ss
}

// Strips strip all string in the string array ss.
func Strips(ss []string) []string {
	return TrimSpaces(ss)
}

// StripLefts left strip all string in the string array ss.
func StripLefts(ss []string) []string {
	for i, s := range ss {
		ss[i] = StripLeft(s)
	}
	return ss
}

// StripRights right strip all string in the string array ss.
func StripRights(ss []string) []string {
	for i, s := range ss {
		ss[i] = StripRight(s)
	}
	return ss
}

// TrimSpaces trim every string in the string array.
func TrimSpaces(ss []string) []string {
	for i, s := range ss {
		ss[i] = strings.TrimSpace(s)
	}
	return ss
}

// RemoveEmpties remove empty string in the string array 'ss', and returns the string array 'ss'
func RemoveEmpties(ss []string) []string {
	return Removes(ss, "")
}

// Removes remove string 'v' in the string array 'ss', and returns the string array 'ss'
func Removes(ss []string, v string) []string {
	return asg.DeleteEqual(ss, v)
}
