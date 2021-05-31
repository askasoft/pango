package mbs

import (
	"strings"
	"unicode/utf8"
)

var (
	// m2s multi-byte rune to single-byte rune
	m2s = map[rune]rune{
		// Number
		'\uFF10': '\u0030', // ０ => 0
		'\uFF11': '\u0031', // １ => 1
		'\uFF12': '\u0032', // ２ => 2
		'\uFF13': '\u0033', // ３ => 3
		'\uFF14': '\u0034', // ４ => 4
		'\uFF15': '\u0035', // ５ => 5
		'\uFF16': '\u0036', // ６ => 6
		'\uFF17': '\u0037', // ７ => 7
		'\uFF18': '\u0038', // ８ => 8
		'\uFF19': '\u0039', // ９ => 9
		// Letter
		'\uFF21': '\u0041', // Ａ => A
		'\uFF22': '\u0042', // Ｂ => B
		'\uFF23': '\u0043', // Ｃ => C
		'\uFF24': '\u0044', // Ｄ => D
		'\uFF25': '\u0045', // Ｅ => E
		'\uFF26': '\u0046', // Ｆ => F
		'\uFF27': '\u0047', // Ｇ => G
		'\uFF28': '\u0048', // Ｈ => H
		'\uFF29': '\u0049', // Ｉ => I
		'\uFF2A': '\u004A', // Ｊ => J
		'\uFF2B': '\u004B', // Ｋ => K
		'\uFF2C': '\u004C', // Ｌ => L
		'\uFF2D': '\u004D', // Ｍ => M
		'\uFF2E': '\u004E', // Ｎ => N
		'\uFF2F': '\u004F', // Ｏ => O
		'\uFF30': '\u0050', // Ｐ => P
		'\uFF31': '\u0051', // Ｑ => Q
		'\uFF32': '\u0052', // Ｒ => R
		'\uFF33': '\u0053', // Ｓ => S
		'\uFF34': '\u0054', // Ｔ => T
		'\uFF35': '\u0055', // Ｕ => U
		'\uFF36': '\u0056', // Ｖ => V
		'\uFF37': '\u0057', // Ｗ => W
		'\uFF38': '\u0058', // Ｘ => X
		'\uFF39': '\u0059', // Ｙ => Y
		'\uFF3A': '\u005A', // Ｚ => Z
		'\uFF41': '\u0061', // ａ => a
		'\uFF42': '\u0062', // ｂ => b
		'\uFF43': '\u0063', // ｃ => c
		'\uFF44': '\u0064', // ｄ => d
		'\uFF45': '\u0065', // ｅ => e
		'\uFF46': '\u0066', // ｆ => f
		'\uFF47': '\u0067', // ｇ => g
		'\uFF48': '\u0068', // ｈ => h
		'\uFF49': '\u0069', // ｉ => i
		'\uFF4A': '\u006A', // ｊ => j
		'\uFF4B': '\u006B', // ｋ => k
		'\uFF4C': '\u006C', // ｌ => l
		'\uFF4D': '\u006D', // ｍ => m
		'\uFF4E': '\u006E', // ｎ => n
		'\uFF4F': '\u006F', // ｏ => o
		'\uFF50': '\u0070', // ｐ => p
		'\uFF51': '\u0071', // ｑ => q
		'\uFF52': '\u0072', // ｒ => r
		'\uFF53': '\u0073', // ｓ => s
		'\uFF54': '\u0074', // ｔ => t
		'\uFF55': '\u0075', // ｕ => u
		'\uFF56': '\u0076', // ｖ => v
		'\uFF57': '\u0077', // ｗ => w
		'\uFF58': '\u0078', // ｘ => x
		'\uFF59': '\u0079', // ｙ => y
		'\uFF5A': '\u007A', // ｚ => z
		// Symbol
		'\u3000': '\u0020', // 　 =>
		'\uFF01': '\u0021', // ！ => !
		'\u2033': '\u0022', // " => "
		'\u201C': '\u0022', // “ => "
		'\u201D': '\u0022', // ” => "
		'\uFF03': '\u0023', // ＃ => #
		'\uFF04': '\u0024', // ＄ => $
		'\uFF05': '\u0025', // ％ => %
		'\uFF06': '\u0026', // ＆ => &
		'\u2019': '\u0027', // ’ => '
		'\uFF08': '\u0028', // （ => (
		'\uFF09': '\u0029', // ） => )
		'\uFF0A': '\u002A', // ＊ => *
		'\uFF0B': '\u002B', // ＋ => +
		'\uFF0C': '\u002C', // ， => ,
		'\uFF0D': '\u002D', // － => -
		'\u2010': '\u002D', // ‐ => -
		'\u2212': '\u002D', // − => -
		'\u30FC': '\u002D', // ー => -
		'\uFF0E': '\u002E', // ． => .
		'\uFF0F': '\u002F', // ／ => /
		'\uFF1A': '\u003A', // ： => :
		'\uFF1B': '\u003B', // ； => ;
		'\uFF1C': '\u003C', // ＜ => <
		'\uFF1D': '\u003D', // ＝ => =
		'\uFF1E': '\u003E', // ＞ => >
		'\uFF1F': '\u003F', // ？ => ?
		'\uFF20': '\u0040', // ＠ => @
		'\uFF3B': '\u005B', // ［ => [
		'\uFFE5': '\u005C', // ￥ => \
		'\uFF3D': '\u005D', // ］ => ]
		'\uFF3E': '\u005E', // ＾ => ^
		'\uFF3F': '\u005F', // ＿ => _
		'\uFF40': '\u0060', // ｀ => `
		'\uFF5B': '\u007B', // ｛ => {
		'\uFF5C': '\u007C', // ｜ => |
		'\uFF5D': '\u007D', // ｝ => }
		'\uFF5E': '\u007E', // ～ => ~
		'\u301C': '\u007E', // 〜 => ~
	}

	// m2s single-byte rune to multi-byte rune
	s2m = initS2M()
)

func merge(ms ...map[rune]rune) map[rune]rune {
	sz := 0
	for _, m := range ms {
		sz += len(m)
	}

	mm := make(map[rune]rune, sz)
	for _, m := range ms {
		for k, v := range m {
			mm[k] = v
		}
	}
	return mm
}

func reverse(m map[rune]rune) map[rune]rune {
	rm := make(map[rune]rune, len(m))
	for k, v := range m {
		rm[v] = k
	}
	return rm
}

func initS2M() map[rune]rune {
	s2m := reverse(m2s)
	s2m['\u002D'] = '\uFF0D' // - => －
	s2m['\u0022'] = '\u2033' // ″ => "
	s2m['\u007E'] = '\uFF5E' // ~ => ～
	return s2m
}

// ToHalfRune convert multi ascii rune c to single ascii rune
func ToHalfRune(c rune) rune {
	if c < utf8.RuneSelf {
		return c
	}

	if s, ok := m2s[c]; ok {
		return s
	}
	return c
}

// ToFullRune convert single ascii rune c to multi ascii rune
func ToFullRune(c rune) rune {
	if c < utf8.RuneSelf {
		if m, ok := s2m[c]; ok {
			return m
		}
	}
	return c
}

// ToHalfWidth convert the string from multi ascii to single ascii
func ToHalfWidth(s string) string {
	if s == "" {
		return s
	}

	sb := &strings.Builder{}
	for i, c := range s {
		r := ToHalfRune(c)
		if r != c {
			if sb.Len() == 0 {
				sb.Grow(len(s))
				sb.WriteString(s[:i])
			}
			sb.WriteRune(r)
			continue
		}

		if sb.Len() > 0 {
			sb.WriteRune(c)
		}
	}

	if sb.Len() > 0 {
		return sb.String()
	}

	return s
}

// ToFullWidth convert the string from single ascii to multi ascii
func ToFullWidth(s string) string {
	if s == "" {
		return s
	}

	sb := &strings.Builder{}
	for i, c := range s {
		r := ToFullRune(c)
		if r != c {
			if sb.Len() == 0 {
				sb.Grow(len(s))
				sb.WriteString(s[:i])
			}
			sb.WriteRune(r)
			continue
		}

		if sb.Len() > 0 {
			sb.WriteRune(c)
		}
	}

	if sb.Len() > 0 {
		return sb.String()
	}

	return s
}
