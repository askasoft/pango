package zho

import (
	"unicode"

	"github.com/askasoft/pango/str"
)

type Kind int

const (
	Unknown Kind = iota
	Hans
	Hant
)

func (k Kind) String() string {
	switch k {
	case Hans:
		return "zh-hans"
	case Hant:
		return "zh-hant"
	default:
		return ""
	}
}

func IsChineseSimplifiedRune(c rune) bool {
	if unicode.Is(unicode.Han, c) {
		if k, ok := chinese_variants[c]; ok && k == Hans {
			return true
		}
	}
	return false
}

func IsChineseTraditionalRune(c rune) bool {
	if unicode.Is(unicode.Han, c) {
		if k, ok := chinese_variants[c]; ok && k == Hant {
			return true
		}
	}
	return false
}

func IsChineseSimplified(s string) bool {
	if s == "" {
		return false
	}

	return !str.ContainsFunc(s, IsChineseTraditionalRune)
}

func IsChineseTraditional(s string) bool {
	if s == "" {
		return false
	}

	return !str.ContainsFunc(s, IsChineseSimplifiedRune)
}

func DetectChinese(s string) Kind {
	if s == "" {
		return Unknown
	}

	hans, hant := 0, 0
	for _, c := range s {
		if IsChineseSimplifiedRune(c) {
			hans++
		} else if IsChineseTraditionalRune(c) {
			hant++
		}
	}

	if hant > hans {
		return Hant
	}
	return Hans
}
