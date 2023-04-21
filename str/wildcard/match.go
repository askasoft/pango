package wildcard

import (
	"unicode/utf8"

	"github.com/pandafw/pango/str"
)

// MatchSimple - finds whether the text matches/satisfies the pattern string.
// supports only '*' wildcard in the pattern.
// considers a file system path as a flat name space.
func MatchSimple(pattern, s string) bool {
	if pattern == "" {
		return s == pattern
	}
	if pattern == "*" {
		return true
	}

	// Does only wildcard '*' match.
	return deepMatchSimple(pattern, s)
}

// MatchSimpleFold - finds whether the text matches/satisfies the pattern string.
// supports only '*' wildcard in the pattern.
// case insensitive.
// considers a file system path as a flat name space.
func MatchSimpleFold(pattern, s string) bool {
	if pattern == "" {
		return s == pattern
	}
	if pattern == "*" {
		return true
	}

	// Does only wildcard '*' match.
	return deepMatchSimpleFold(pattern, s)
}

// Match -  finds whether the text matches/satisfies the pattern string.
// supports  '*' and '?' wildcards in the pattern string.
// case insensitive.
// unlike path.Match(), considers a path as a flat name space while matching the pattern.
func Match(pattern, s string) bool {
	if pattern == "" {
		return s == pattern
	}

	if pattern == "*" {
		return true
	}

	// Does extended wildcard '*' and '?' match.
	return deepMatchWild(pattern, s)
}

// MatchFold -  finds whether the text matches/satisfies the pattern string.
// supports  '*' and '?' wildcards in the pattern string.
// unlike path.Match(), considers a path as a flat name space while matching the pattern.
func MatchFold(pattern, s string) bool {
	if pattern == "" {
		return s == pattern
	}

	if pattern == "*" {
		return true
	}

	// Does extended wildcard '*' and '?' match.
	return deepMatchWildFold(pattern, s)
}

func skipAsterisk(pattern string) string {
	for len(pattern) > 1 && pattern[1] == '*' {
		pattern = pattern[1:]
	}
	return pattern
}

func deepMatchSimple(pattern, s string) bool {
	for len(pattern) > 0 {
		pc, pz := utf8.DecodeRuneInString(pattern)
		switch pc {
		case '*':
			pattern = skipAsterisk(pattern)
			if deepMatchSimple(pattern[pz:], s) {
				return true
			}
			if len(s) > 0 {
				_, sz := utf8.DecodeRuneInString(s)
				return deepMatchSimple(pattern, s[sz:])
			}
			return false
		default:
			if len(s) == 0 {
				return false
			}
			sc, sz := utf8.DecodeRuneInString(s)
			if sc != pc {
				return false
			}
			s = s[sz:]
			pattern = pattern[pz:]
		}
	}
	return len(s) == 0 && len(pattern) == 0
}

func deepMatchSimpleFold(pattern, s string) bool {
	for len(pattern) > 0 {
		pc, pz := utf8.DecodeRuneInString(pattern)
		switch pc {
		case '*':
			pattern = skipAsterisk(pattern)
			if deepMatchSimpleFold(pattern[pz:], s) {
				return true
			}
			if len(s) > 0 {
				_, sz := utf8.DecodeRuneInString(s)
				return deepMatchSimpleFold(pattern, s[sz:])
			}
			return false
		default:
			if len(s) == 0 {
				return false
			}
			sc, sz := utf8.DecodeRuneInString(s)
			if !str.RuneEqualFold(sc, pc) {
				return false
			}
			s = s[sz:]
			pattern = pattern[pz:]
		}
	}
	return len(s) == 0 && len(pattern) == 0
}

func deepMatchWild(pattern, s string) bool {
	for len(pattern) > 0 {
		pc, pz := utf8.DecodeRuneInString(pattern)
		switch pc {
		case '*':
			pattern = skipAsterisk(pattern)
			if deepMatchWild(pattern[pz:], s) {
				return true
			}
			if len(s) > 0 {
				_, sz := utf8.DecodeRuneInString(s)
				return deepMatchWild(pattern, s[sz:])
			}
			return false
		case '?':
			if len(s) == 0 {
				return false
			}
			_, sz := utf8.DecodeRuneInString(s)
			s = s[sz:]
			pattern = pattern[pz:]
		default:
			if len(s) == 0 {
				return false
			}
			sc, sz := utf8.DecodeRuneInString(s)
			if sc != pc {
				return false
			}
			s = s[sz:]
			pattern = pattern[pz:]
		}
	}
	return len(s) == 0 && len(pattern) == 0
}

func deepMatchWildFold(pattern, s string) bool {
	for len(pattern) > 0 {
		pc, pz := utf8.DecodeRuneInString(pattern)
		switch pc {
		case '*':
			pattern = skipAsterisk(pattern)
			if deepMatchWildFold(pattern[pz:], s) {
				return true
			}
			if len(s) > 0 {
				_, sz := utf8.DecodeRuneInString(s)
				return deepMatchWildFold(pattern, s[sz:])
			}
			return false
		case '?':
			if len(s) == 0 {
				return false
			}
			_, sz := utf8.DecodeRuneInString(s)
			s = s[sz:]
			pattern = pattern[pz:]
		default:
			if len(s) == 0 {
				return false
			}
			sc, sz := utf8.DecodeRuneInString(s)
			if !str.RuneEqualFold(sc, pc) {
				return false
			}
			s = s[sz:]
			pattern = pattern[pz:]
		}
	}
	return len(s) == 0 && len(pattern) == 0
}
