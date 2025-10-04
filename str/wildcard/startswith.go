package wildcard

import (
	"unicode/utf8"

	"github.com/askasoft/pango/str"
)

// StartsWithSimple - finds whether the pattern string contains the text.
// supports only '*' wildcard in the pattern.
func StartsWithSimple(pattern, s string) bool {
	switch pattern {
	case "":
		return pattern == s
	case "*":
		return true
	default:
		return deepStartsWithSimple(pattern, s, runeEqual)
	}
}

// StartsWithSimpleFold - finds whether the pattern string contains the text.
// supports only '*' wildcard in the pattern.
// case insensitive.
func StartsWithSimpleFold(pattern, s string) bool {
	switch pattern {
	case "":
		return pattern == s
	case "*":
		return true
	default:
		return deepStartsWithSimple(pattern, s, str.RuneEqualFold)
	}
}

// StartsWith - finds whether the pattern string contains the text.
// supports  '*' and '?' wildcards in the pattern string.
// case insensitive.
func StartsWith(pattern, s string) bool {
	switch pattern {
	case "":
		return pattern == s
	case "*":
		return true
	default:
		return deepStartsWithWild(pattern, s, runeEqual)
	}
}

// StartsWithFold - finds whether the pattern string contains the text.
// supports  '*' and '?' wildcards in the pattern string.
func StartsWithFold(pattern, s string) bool {
	switch pattern {
	case "":
		return pattern == s
	case "*":
		return true
	default:
		return deepStartsWithWild(pattern, s, str.RuneEqualFold)
	}
}

func deepStartsWithSimple(p, s string, eq equal) bool {
	for len(p) > 0 {
		pc, pz := utf8.DecodeRuneInString(p)
		switch pc {
		case '*':
			return true
		default:
			if len(s) == 0 {
				return true
			}
			sc, sz := utf8.DecodeRuneInString(s)
			if !eq(sc, pc) {
				return false
			}
			p, s = p[pz:], s[sz:]
		}
	}
	return len(s) == 0
}

func deepStartsWithWild(p, s string, eq equal) bool {
	for len(p) > 0 {
		pc, pz := utf8.DecodeRuneInString(p)
		switch pc {
		case '*':
			return true
		case '?':
			if len(s) == 0 {
				return true
			}
			_, sz := utf8.DecodeRuneInString(s)
			p, s = p[pz:], s[sz:]
		default:
			if len(s) == 0 {
				return true
			}
			sc, sz := utf8.DecodeRuneInString(s)
			if !eq(sc, pc) {
				return false
			}
			p, s = p[pz:], s[sz:]
		}
	}
	return len(s) == 0
}
