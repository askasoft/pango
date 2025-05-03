package lut

import (
	"unicode"
	"unicode/utf8"

	"github.com/askasoft/pango/mag"
	"github.com/askasoft/pango/str"
)

// https://en.wikipedia.org/wiki/Halfwidth_and_Fullwidth_Forms_(Unicode_block)

var (
	// s2m single-byte to multi-byte rune rune
	s2m = map[rune]rune{
		// space
		'\u0020': '\u3000', // 　=>
		// Symbol
		'\u0021': '\uFF01', // ! => ！
		'\u0022': '\uFF02', // " => ＂
		'\u0023': '\uFF03', // # => ＃
		'\u0024': '\uFF04', // $ => ＄
		'\u0025': '\uFF05', // % => ％
		'\u0026': '\uFF06', // & => ＆
		'\u0027': '\uFF07', // ' => ＇
		'\u0028': '\uFF08', // ( => （
		'\u0029': '\uFF09', // ) => ）
		'\u002A': '\uFF0A', // * => ＊
		'\u002B': '\uFF0B', // + => ＋
		'\u002C': '\uFF0C', // , => ，
		'\u002D': '\uFF0D', // - => －
		'\u002E': '\uFF0E', // . => ．
		'\u002F': '\uFF0F', // / => ／
		// Number
		'\u0030': '\uFF10', // 0 => ０
		'\u0031': '\uFF11', // 1 => １
		'\u0032': '\uFF12', // 2 => ２
		'\u0033': '\uFF13', // 3 => ３
		'\u0034': '\uFF14', // 4 => ４
		'\u0035': '\uFF15', // 5 => ５
		'\u0036': '\uFF16', // 6 => ６
		'\u0037': '\uFF17', // 7 => ７
		'\u0038': '\uFF18', // 8 => ８
		'\u0039': '\uFF19', // 9 => ９
		// Symbol
		'\u003A': '\uFF1A', // : => ：
		'\u003B': '\uFF1B', // ; => ；
		'\u003C': '\uFF1C', // < => ＜
		'\u003D': '\uFF1D', // = => ＝
		'\u003E': '\uFF1E', // > => ＞
		'\u003F': '\uFF1F', // ? => ？
		'\u0040': '\uFF20', // @ => ＠
		// Letter
		'\u0041': '\uFF21', // A => Ａ
		'\u0042': '\uFF22', // B => Ｂ
		'\u0043': '\uFF23', // C => Ｃ
		'\u0044': '\uFF24', // D => Ｄ
		'\u0045': '\uFF25', // E => Ｅ
		'\u0046': '\uFF26', // F => Ｆ
		'\u0047': '\uFF27', // G => Ｇ
		'\u0048': '\uFF28', // H => Ｈ
		'\u0049': '\uFF29', // I => Ｉ
		'\u004A': '\uFF2A', // J => Ｊ
		'\u004B': '\uFF2B', // K => Ｋ
		'\u004C': '\uFF2C', // L => Ｌ
		'\u004D': '\uFF2D', // M => Ｍ
		'\u004E': '\uFF2E', // N => Ｎ
		'\u004F': '\uFF2F', // O => Ｏ
		'\u0050': '\uFF30', // P => Ｐ
		'\u0051': '\uFF31', // Q => Ｑ
		'\u0052': '\uFF32', // R => Ｒ
		'\u0053': '\uFF33', // S => Ｓ
		'\u0054': '\uFF34', // T => Ｔ
		'\u0055': '\uFF35', // U => Ｕ
		'\u0056': '\uFF36', // V => Ｖ
		'\u0057': '\uFF37', // W => Ｗ
		'\u0058': '\uFF38', // X => Ｘ
		'\u0059': '\uFF39', // Y => Ｙ
		'\u005A': '\uFF3A', // Z => Ｚ
		// Symbol
		'\u005B': '\uFF3B', // [ => ［
		'\u005C': '\uFF3C', // \ => ＼
		'\u005D': '\uFF3D', // ] => ］
		'\u005E': '\uFF3E', // ^ => ＾
		'\u005F': '\uFF3F', // _ => ＿
		'\u0060': '\uFF40', // ` => ｀
		// Letter
		'\u0061': '\uFF41', // a => ａ
		'\u0062': '\uFF42', // b => ｂ
		'\u0063': '\uFF43', // c => ｃ
		'\u0064': '\uFF44', // d => ｄ
		'\u0065': '\uFF45', // e => ｅ
		'\u0066': '\uFF46', // f => ｆ
		'\u0067': '\uFF47', // g => ｇ
		'\u0068': '\uFF48', // h => ｈ
		'\u0069': '\uFF49', // i => ｉ
		'\u006A': '\uFF4A', // j => ｊ
		'\u006B': '\uFF4B', // k => ｋ
		'\u006C': '\uFF4C', // l => ｌ
		'\u006D': '\uFF4D', // m => ｍ
		'\u006E': '\uFF4E', // n => ｎ
		'\u006F': '\uFF4F', // o => ｏ
		'\u0070': '\uFF50', // p => ｐ
		'\u0071': '\uFF51', // q => ｑ
		'\u0072': '\uFF52', // r => ｒ
		'\u0073': '\uFF53', // s => ｓ
		'\u0074': '\uFF54', // t => ｔ
		'\u0075': '\uFF55', // u => ｕ
		'\u0076': '\uFF56', // v => ｖ
		'\u0077': '\uFF57', // w => ｗ
		'\u0078': '\uFF58', // x => ｘ
		'\u0079': '\uFF59', // y => ｙ
		'\u007A': '\uFF5A', // z => ｚ
		// Symbol
		'\u007B': '\uFF5B', // { => ｛
		'\u007C': '\uFF5C', // | => ｜
		'\u007D': '\uFF5D', // } => ｝
		'\u007E': '\uFF5E', // ~ => ～
	}

	// m2s multi-byte rune to single-byte rune
	m2s = mag.Merge(mag.Reverse(s2m), map[rune]rune{
		// '\u00A0': '\u0020', //   =>  // non-breaking space
		// Symbol
		'\uFF61': '\u002E', // ｡ => .
		'\uFF62': '\u005B', // ｢ => [
		'\uFF63': '\u005D', // ｣ => ]
		'\uFF64': '\u002C', // ､ => ,
		'\u2010': '\u002D', // ‐ => -
		'\u2019': '\u0027', // ’ => '
		'\u201C': '\u0022', // “ => "
		'\u201D': '\u0022', // ” => "
		'\u2033': '\u0022', // ″ => "
		'\u2212': '\u002D', // − => -
		'\u301C': '\u007E', // 〜 => ~
		'\u30A0': '\u003D', // ゠ => =  // KATAKANA-HIRAGANA DOUBLE HYPHEN
		'\u30FC': '\u002D', // ー => -
	})
)

// ToASCIIRune convert multi byte rune c to single ascii rune
func ToASCIIRune(c rune) rune {
	if c < utf8.RuneSelf {
		return c
	}

	if s, ok := m2s[c]; ok {
		return s
	}
	return c
}

// ToFullRune convert single ascii rune c to multi byte rune
func ToFullRune(c rune) rune {
	if c < utf8.RuneSelf {
		if m, ok := s2m[c]; ok {
			return m
		}
	}
	return c
}

// ToASCII convert the string from multi byte to single byte ascii
func ToASCII(s string) string {
	return str.ReplaceFunc(s, ToASCIIRune)
}

// ToASCIIs convert all string from multi byte to single byte ascii in the string slice ss.
func ToASCIIs(ss []string) []string {
	for i, s := range ss {
		ss[i] = ToASCII(s)
	}
	return ss
}

// ToFullWidth convert the string from single ascii to multi ascii
func ToFullWidth(s string) string {
	return str.ReplaceFunc(s, ToFullRune)
}

// ToFullWidths convert all string from single ascii to multi ascii in the string slice ss.
func ToFullWidths(ss []string) []string {
	for i, s := range ss {
		ss[i] = ToFullWidth(s)
	}
	return ss
}

// IsHalfWidthRune checks if the rune c is unicode half-width char.
// https://www.asahi-net.or.jp/~ax2s-kmtn/ref/unicode/uff00.html
// \u0000-\u00FF\uFF61-\uFFDF\uFFE8-\uFFEE
func IsHalfWidthRune(c rune) bool {
	if c <= unicode.MaxASCII ||
		(c >= '\uFF61' && c <= '\uFFDF') ||
		(c >= '\uFFE8' && c <= '\uFFEF') {
		return true
	}
	return false
}

func IsFullWidthRune(c rune) bool {
	return !IsHalfWidthRune(c)
}

// HasHalfWidth checks if the string contains any unicode half-width chars.
func HasHalfWidth(s string) bool {
	return str.ContainsFunc(s, IsHalfWidthRune)
}

// IsHalfWidth checks if the string contains half-width chars only.
func IsHalfWidth(s string) bool {
	return str.ChecksFunc(s, IsHalfWidthRune)
}

// HasFullWidth checks if the string contains any full-width chars.
func HasFullWidth(s string) bool {
	return str.ContainsFunc(s, IsFullWidthRune)
}

// IsFullWidth checks if the string contains full-width chars only.
func IsFullWidth(s string) bool {
	return str.ChecksFunc(s, IsFullWidthRune)
}

// IsVariableWidth checks if the string contains a mixture of full and half-width chars.
func IsVariableWidth(s string) bool {
	if s == "" {
		return false
	}
	return HasHalfWidth(s) && HasFullWidth(s)
}

// CompareFold returns an integer comparing two strings case & width (full/single) insensitive.
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
func CompareFold(s, t string) int {
	// ASCII fast path
	i := 0
	for ; i < len(s) && i < len(t); i++ {
		sr := s[i]
		tr := t[i]
		if sr|tr >= utf8.RuneSelf {
			goto hasUnicode
		}

		if tr == sr {
			continue
		}

		// ASCII only, sr/tr must be upper/lower case
		if 'A' <= sr && sr <= 'Z' {
			sr += ('a' - 'A')
		}
		if 'A' <= tr && tr <= 'Z' {
			tr += ('a' - 'A')
		}

		switch {
		case sr < tr:
			return -1
		case sr > tr:
			return 1
		}
	}

	// Check if we've exhausted both strings.
	{
		r := len(s) - len(t)
		switch {
		case r < 0:
			return -1
		case r > 0:
			return 1
		default:
			return 0
		}
	}

hasUnicode:
	s = s[i:]
	t = t[i:]
	for _, sr := range s {
		// If t is exhausted the strings are not equal.
		if len(t) == 0 {
			return 1
		}

		// Extract first rune from second string.
		var tr rune
		if t[0] < utf8.RuneSelf {
			tr, t = rune(t[0]), t[1:]
		} else {
			r, size := utf8.DecodeRuneInString(t)
			tr, t = r, t[size:]
		}

		// If they match, keep going;
		if tr == sr {
			continue
		}

		// Fast check for ASCII.
		if sr < utf8.RuneSelf && tr < utf8.RuneSelf {
			// ASCII only, sr/tr must be upper/lower case
			if 'A' <= sr && sr <= 'Z' {
				sr += ('a' - 'A')
			}
			if 'A' <= tr && tr <= 'Z' {
				tr += ('a' - 'A')
			}

			switch {
			case sr < tr:
				return -1
			case sr > tr:
				return 1
			default:
				continue
			}
		}

		// check for full/single width and lower/upper case.
		sr = unicode.ToLower(ToASCIIRune(sr))
		tr = unicode.ToLower(ToASCIIRune(tr))
		switch {
		case sr < tr:
			return -1
		case sr > tr:
			return 1
		}
	}

	// First string is empty, so check if the second one is also empty.
	if len(t) == 0 {
		return 0
	}

	return -1
}
