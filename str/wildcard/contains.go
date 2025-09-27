package wildcard

import (
	"unicode/utf8"

	"github.com/askasoft/pango/str"
)

// ContainsSimple - finds whether the pattern string contains the text.
// supports only '*' wildcard in the pattern.
func ContainsSimple(pattern, s string) bool {
	switch pattern {
	case "":
		return pattern == s
	case "*":
		return true
	default:
		return deepContainsSimple(pattern, s, runeEqual)
	}
}

// ContainsSimpleFold - finds whether the pattern string contains the text.
// supports only '*' wildcard in the pattern.
// case insensitive.
func ContainsSimpleFold(pattern, s string) bool {
	switch pattern {
	case "":
		return pattern == s
	case "*":
		return true
	default:
		return deepContainsSimple(pattern, s, str.RuneEqualFold)
	}
}

// Contains - finds whether the pattern string contains the text.
// supports  '*' and '?' wildcards in the pattern string.
// case insensitive.
func Contains(pattern, s string) bool {
	switch pattern {
	case "":
		return pattern == s
	case "*":
		return true
	default:
		return deepContainsWild(pattern, s, runeEqual)
	}
}

// ContainsFold - finds whether the pattern string contains the text.
// supports  '*' and '?' wildcards in the pattern string.
func ContainsFold(pattern, s string) bool {
	switch pattern {
	case "":
		return pattern == s
	case "*":
		return true
	default:
		return deepContainsWild(pattern, s, str.RuneEqualFold)
	}
}

func deepContainsSimple(p, s string, eq equal) bool {
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

func deepContainsWild(p, s string, eq equal) bool {
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
