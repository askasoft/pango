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
	return deepMatchSimple(pattern, s, false)
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
	return deepMatchSimple(pattern, s, true)
}

func deepMatchSimple(pattern, s string, fold bool) bool {
	for len(pattern) > 0 {
		pc, pz := utf8.DecodeRuneInString(pattern)
		switch pc {
		case '*':
			if deepMatchSimple(pattern[pz:], s, fold) {
				return true
			}
			if len(s) > 0 {
				_, sz := utf8.DecodeRuneInString(s)
				return deepMatchSimple(pattern, s[sz:], fold)
			}
			return false
		default:
			if len(s) == 0 {
				return false
			}
			sc, sz := utf8.DecodeRuneInString(s)
			if fold {
				if !str.RuneEqualFold(sc, pc) {
					return false
				}
			} else if sc != pc {
				return false
			}
			s = s[sz:]
			pattern = pattern[pz:]
		}
	}
	return len(s) == 0 && len(pattern) == 0
}

// Match -  finds whether the text matches/satisfies the pattern string.
// supports  '*' and '?' wildcards in the pattern string.
// case insensitive.
// unlike path.Match(), considers a path as a flat name space while matching the pattern.
func Match(pattern, s string) (matched bool) {
	if pattern == "" {
		return s == pattern
	}

	if pattern == "*" {
		return true
	}

	// Does extended wildcard '*' and '?' match.
	return deepMatchWild(pattern, s, false)
}

// MatchFold -  finds whether the text matches/satisfies the pattern string.
// supports  '*' and '?' wildcards in the pattern string.
// unlike path.Match(), considers a path as a flat name space while matching the pattern.
func MatchFold(pattern, s string) (matched bool) {
	if pattern == "" {
		return s == pattern
	}

	if pattern == "*" {
		return true
	}

	// Does extended wildcard '*' and '?' match.
	return deepMatchWild(pattern, s, true)
}

func deepMatchWild(pattern, s string, fold bool) bool {
	for len(pattern) > 0 {
		pc, pz := utf8.DecodeRuneInString(pattern)
		switch pc {
		case '*':
			if deepMatchWild(pattern[pz:], s, fold) {
				return true
			}
			if len(s) > 0 {
				_, sz := utf8.DecodeRuneInString(s)
				return deepMatchWild(pattern, s[sz:], fold)
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
			if fold {
				if !str.RuneEqualFold(sc, pc) {
					return false
				}
			} else if sc != pc {
				return false
			}
			s = s[sz:]
			pattern = pattern[pz:]
		}
	}
	return len(s) == 0 && len(pattern) == 0
}
