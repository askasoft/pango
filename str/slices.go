package str

import (
	"strconv"
	"strings"

	"github.com/askasoft/pango/num"
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
	for i, s := range ss {
		if s == v {
			for j := i + 1; j < len(ss); j++ {
				if s := ss[j]; s != v {
					ss[i] = s
					i++
				}
			}
			return ss[:i]
		}
	}
	return ss
}

// JoinInts concatenates the elements of its first argument to create a single string. The separator
// string sep is placed between elements in the resulting string.
func JoinInts(elems []int, sep string, fmt ...func(int) string) string {
	itoa := strconv.Itoa
	if len(fmt) > 0 {
		itoa = fmt[0]
	}

	switch len(elems) {
	case 0:
		return ""
	case 1:
		return itoa(elems[0])
	}

	var b Builder
	b.WriteString(itoa(elems[0]))
	for _, n := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(itoa(n))
	}
	return b.String()
}

// JoinInt64s concatenates the elements of its first argument to create a single string. The separator
// string sep is placed between elements in the resulting string.
func JoinInt64s(elems []int64, sep string, fmt ...func(int64) string) string {
	ltoa := num.Ltoa
	if len(fmt) > 0 {
		ltoa = fmt[0]
	}

	switch len(elems) {
	case 0:
		return ""
	case 1:
		return ltoa(elems[0])
	}

	var b Builder
	b.WriteString(ltoa(elems[0]))
	for _, n := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(ltoa(n))
	}
	return b.String()
}
