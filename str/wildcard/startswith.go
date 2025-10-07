package wildcard

import (
	"unicode/utf8"

	"github.com/askasoft/pango/str"
)

// HasPrefixSimple tests whether the pattern string begins with prefix s.
// supports only '*' wildcard in the pattern.
// alias for StartsWithSimple.
func HasPrefixSimple(pattern, s string) bool {
	return StartsWithSimple(pattern, s)
}

// StartsWithSimple tests whether the pattern string begins with prefix s.
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

// HasPrefixSimpleFold tests whether the pattern string begins with prefix s.
// supports only '*' wildcard in the pattern. case insensitive.
// alias for StartsWithSimpleFold.
func HasPrefixSimpleFold(pattern, s string) bool {
	return StartsWithSimpleFold(pattern, s)
}

// StartsWithSimpleFold tests whether the pattern string begins with prefix s.
// supports only '*' wildcard in the pattern. case insensitive.
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

// HasPrefix tests whether the pattern string begins with prefix s.
// supports  '*' and '?' wildcards in the pattern string.
// alias for StartsWith.
func HasPrefix(pattern, s string) bool {
	return StartsWith(pattern, s)
}

// StartsWith tests whether the pattern string begins with prefix s.
// supports  '*' and '?' wildcards in the pattern string.
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

// HasPrefixFold tests whether the pattern string begins with prefix s.
// supports  '*' and '?' wildcards in the pattern string. case insensitive.
// alias for StartsWithFold.
func HasPrefixFold(pattern, s string) bool {
	return StartsWithFold(pattern, s)
}

// StartsWithFold tests whether the pattern string begins with prefix s.
// supports  '*' and '?' wildcards in the pattern string. case insensitive.
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
