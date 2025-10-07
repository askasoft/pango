package wildcard

import (
	"unicode/utf8"

	"github.com/askasoft/pango/str"
)

type equal func(a, b rune) bool

func runeEqual(a, b rune) bool {
	return a == b
}

// MatchSimple tests whether the string s matches/satisfies the pattern string.
// supports only '*' wildcard in the pattern.
// considers a file system path as a flat name space.
func MatchSimple(pattern, s string) bool {
	switch pattern {
	case "":
		return pattern == s
	case "*":
		return true
	default:
		return deepMatchSimple(pattern, s, runeEqual)
	}
}

// MatchSimpleFold tests whether the string s matches/satisfies the pattern string.
// supports only '*' wildcard in the pattern.
// case insensitive.
// considers a file system path as a flat name space.
func MatchSimpleFold(pattern, s string) bool {
	switch pattern {
	case "":
		return pattern == s
	case "*":
		return true
	default:
		return deepMatchSimple(pattern, s, str.RuneEqualFold)
	}
}

// Match tests whether the string s matches/satisfies the pattern string.
// supports  '*' and '?' wildcards in the pattern string.
// case insensitive.
// unlike path.Match(), considers a path as a flat name space while matching the pattern.
func Match(pattern, s string) bool {
	switch pattern {
	case "":
		return pattern == s
	case "*":
		return true
	default:
		return deepMatchWild(pattern, s, runeEqual)
	}
}

// MatchFold tests whether the string s matches/satisfies the pattern string.
// supports  '*' and '?' wildcards in the pattern string.
// unlike path.Match(), considers a path as a flat name space while matching the pattern.
func MatchFold(pattern, s string) bool {
	switch pattern {
	case "":
		return pattern == s
	case "*":
		return true
	default:
		return deepMatchWild(pattern, s, str.RuneEqualFold)
	}
}

func skipAsterisk(p string) string {
	for len(p) > 1 && p[1] == '*' {
		p = p[1:]
	}
	return p
}

func deepMatchSimple(p, s string, eq equal) bool {
	for len(p) > 0 {
		pc, pz := utf8.DecodeRuneInString(p)
		switch pc {
		case '*':
			p = skipAsterisk(p)
			if deepMatchSimple(p[pz:], s, eq) {
				return true
			}
			if len(s) > 0 {
				_, sz := utf8.DecodeRuneInString(s)
				return deepMatchSimple(p, s[sz:], eq)
			}
			return false
		default:
			if len(s) == 0 {
				return false
			}
			sc, sz := utf8.DecodeRuneInString(s)
			if !eq(sc, pc) {
				return false
			}
			p, s = p[pz:], s[sz:]
		}
	}
	return len(s) == 0 && len(p) == 0
}

func deepMatchWild(p, s string, eq equal) bool {
	for len(p) > 0 {
		pc, pz := utf8.DecodeRuneInString(p)
		switch pc {
		case '*':
			p = skipAsterisk(p)
			if deepMatchWild(p[pz:], s, eq) {
				return true
			}
			if len(s) > 0 {
				_, sz := utf8.DecodeRuneInString(s)
				return deepMatchWild(p, s[sz:], eq)
			}
			return false
		case '?':
			if len(s) == 0 {
				return false
			}
			_, sz := utf8.DecodeRuneInString(s)
			p, s = p[pz:], s[sz:]
		default:
			if len(s) == 0 {
				return false
			}
			sc, sz := utf8.DecodeRuneInString(s)
			if !eq(sc, pc) {
				return false
			}
			p, s = p[pz:], s[sz:]
		}
	}
	return len(s) == 0 && len(p) == 0
}
