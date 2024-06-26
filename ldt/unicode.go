package ldt

import "unicode"

// isStopChar returns true if r is space, punctuation or digit.
func isStopChar(r rune) bool {
	if unicode.IsSymbol(r) || unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsDigit(r) {
		return true
	}
	return false
}

func isLatin(r rune) bool {
	return unicode.Is(unicode.Latin, r)
}

func isCyrillic(r rune) bool {
	return unicode.Is(unicode.Cyrillic, r)
}

func isArabic(r rune) bool {
	return unicode.Is(unicode.Arabic, r)
}

func isDevanagari(r rune) bool {
	return unicode.Is(unicode.Devanagari, r)
}

func isEthiopic(r rune) bool {
	return unicode.Is(unicode.Ethiopic, r)
}

func isHebrew(r rune) bool {
	return unicode.Is(unicode.Hebrew, r)
}

func isBengali(r rune) bool {
	return unicode.Is(unicode.Bengali, r)
}

func isGeorgian(r rune) bool {
	return unicode.Is(unicode.Georgian, r)
}

func isGreek(r rune) bool {
	return unicode.Is(unicode.Greek, r)
}

func isKannada(r rune) bool {
	return unicode.Is(unicode.Kannada, r)
}

func isTamil(r rune) bool {
	return unicode.Is(unicode.Tamil, r)
}

func isThai(r rune) bool {
	return unicode.Is(unicode.Thai, r)
}

func isGujarati(r rune) bool {
	return unicode.Is(unicode.Gujarati, r)
}

func isGurmukhi(r rune) bool {
	return unicode.Is(unicode.Gurmukhi, r)
}

func isTelugu(r rune) bool {
	return unicode.Is(unicode.Telugu, r)
}

func isMalayalam(r rune) bool {
	return unicode.Is(unicode.Malayalam, r)
}

func isOriya(r rune) bool {
	return unicode.Is(unicode.Oriya, r)
}

func isMyanmar(r rune) bool {
	return unicode.Is(unicode.Myanmar, r)
}

func isSinhala(r rune) bool {
	return unicode.Is(unicode.Sinhala, r)
}

func isKhmer(r rune) bool {
	return unicode.Is(unicode.Khmer, r)
}

func isHan(r rune) bool {
	return unicode.Is(unicode.Han, r)
}

func isKana(r rune) bool {
	return unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r)
}

func isHangul(r rune) bool {
	return unicode.Is(unicode.Hangul, r)
}
