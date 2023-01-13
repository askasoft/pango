package mbs

import (
	"strings"
	"unicode/utf8"
)

// http://www.asahi-net.or.jp/~ax2s-kmtn/ref/unicode/uff00.html
var (
	// z2hMark 全角: 。「」、・゛゜
	z2hMark = map[rune]rune{
		'\u3002': '\uFF61', // 。 => ｡
		'\u300C': '\uFF62', // 「 => ｢
		'\u300D': '\uFF63', // 」 => ｣
		'\u3001': '\uFF64', // 、 => ､
		'\u30FB': '\uFF65', // ・ => ･
		'\u309B': '\uFF9E', // ゛ => ﾞ
		'\u309C': '\uFF9F', // ゜ => ﾟ
	}

	// h2zMark 半角: ｡｢｣､･ﾞﾟ
	h2zMark = reverse(z2hMark)

	// z2hAyatu 全角: ァィゥェォャュョッー
	z2hAyatu = map[rune]rune{
		'\u30A1': '\uFF67', // ァ => ｧ
		'\u30A3': '\uFF68', // ィ => ｨ
		'\u30A5': '\uFF69', // ゥ => ｩ
		'\u30A7': '\uFF6A', // ェ => ｪ
		'\u30A9': '\uFF6B', // ォ => ｫ
		'\u30E3': '\uFF6C', // ャ => ｬ
		'\u30E5': '\uFF6D', // ュ => ｭ
		'\u30E7': '\uFF6E', // ョ => ｮ
		'\u30C3': '\uFF6F', // ッ => ｯ
		'\u30FC': '\uFF70', // ー => ｰ
	}

	// h2zAyatu 半角: ｧｨｩｪｫｬｭｮｯｰ
	h2zAyatu = reverse(z2hAyatu)

	// z2hAnamayara 全角: アイエオナニヌネノマミムメモヤユヨラリルレロン
	z2hAnamayara = map[rune]rune{
		'\u30A2': '\uFF71', // ア => ｱ
		'\u30A4': '\uFF72', // イ => ｲ
		'\u30A8': '\uFF74', // エ => ｴ
		'\u30AA': '\uFF75', // オ => ｵ
		'\u30CA': '\uFF85', // ナ => ﾅ
		'\u30CB': '\uFF86', // ニ => ﾆ
		'\u30CC': '\uFF87', // ヌ => ﾇ
		'\u30CD': '\uFF88', // ネ => ﾈ
		'\u30CE': '\uFF89', // ノ => ﾉ
		'\u30DE': '\uFF8F', // マ => ﾏ
		'\u30DF': '\uFF90', // ミ => ﾐ
		'\u30E0': '\uFF91', // ム => ﾑ
		'\u30E1': '\uFF92', // メ => ﾒ
		'\u30E2': '\uFF93', // モ => ﾓ
		'\u30E4': '\uFF94', // ヤ => ﾔ
		'\u30E6': '\uFF95', // ユ => ﾕ
		'\u30E8': '\uFF96', // ヨ => ﾖ
		'\u30E9': '\uFF97', // ラ => ﾗ
		'\u30EA': '\uFF98', // リ => ﾘ
		'\u30EB': '\uFF99', // ル => ﾙ
		'\u30EC': '\uFF9A', // レ => ﾚ
		'\u30ED': '\uFF9B', // ロ => ﾛ
		'\u30F3': '\uFF9D', // ン => ﾝ
	}

	// h2zAnamayara 半角: ｱｲｴｵﾅﾆﾇﾈﾉﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾝ
	h2zAnamayara = reverse(z2hAnamayara)

	// z2hKasataha 全角　かさたは　行: カキクケコサシスセソタチツテトハヒフヘホウ
	z2hKasataha = map[rune]rune{
		'\u30AB': '\uFF76', // カ => ｶ
		'\u30AD': '\uFF77', // キ => ｷ
		'\u30AF': '\uFF78', // ク => ｸ
		'\u30B1': '\uFF79', // ケ => ｹ
		'\u30B3': '\uFF7A', // コ => ｺ
		'\u30B5': '\uFF7B', // サ => ｻ
		'\u30B7': '\uFF7C', // シ => ｼ
		'\u30B9': '\uFF7D', // ス => ｽ
		'\u30BB': '\uFF7E', // セ => ｾ
		'\u30BD': '\uFF7F', // ソ => ｿ
		'\u30BF': '\uFF80', // タ => ﾀ
		'\u30C1': '\uFF81', // チ => ﾁ
		'\u30C4': '\uFF82', // ツ => ﾂ
		'\u30C6': '\uFF83', // テ => ﾃ
		'\u30C8': '\uFF84', // ト => ﾄ
		'\u30CF': '\uFF8A', // ハ => ﾊ
		'\u30D2': '\uFF8B', // ヒ => ﾋ
		'\u30D5': '\uFF8C', // フ => ﾌ
		'\u30D8': '\uFF8D', // ヘ => ﾍ
		'\u30DB': '\uFF8E', // ホ => ﾎ
		'\u30A6': '\uFF73', // ウ => ｳ
	}

	// h2zKasataha 半角　かさたは　行: ｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾊﾋﾌﾍﾎｳ
	h2zKasataha = reverse(z2hKasataha)

	// z2hWaou 全角　わ　行: ワヲ
	z2hWaou = map[rune]rune{
		'\u30EF': '\uFF9C', // ワ => ﾜ
		'\u30F2': '\uFF66', // ヲ => ｦ
	}

	// h2zWaou 半角　わ　行: ﾜｦ
	h2zWaou = reverse(z2hWaou)

	// z2h 全角
	z2h = merge(z2hMark, z2hAyatu, z2hAnamayara, z2hKasataha, z2hWaou)

	// h2z 半角
	h2z = merge(h2zMark, h2zAyatu, h2zAnamayara, h2zKasataha, h2zWaou)

	// z2hDaku 全角　濁文字: ガギグゲゴザジズゼゾダヂヅデドバビブベボヴヷヸヹヺ
	z2hDaku = map[rune]rune{
		'\u30AC': '\uFF76', // ガ => ｶ
		'\u30AE': '\uFF77', // ギ => ｷ
		'\u30B0': '\uFF78', // グ => ｸ
		'\u30B2': '\uFF79', // ゲ => ｹ
		'\u30B4': '\uFF7A', // ゴ => ｺ
		'\u30B6': '\uFF7B', // ザ => ｻ
		'\u30B8': '\uFF7C', // ジ => ｼ
		'\u30BA': '\uFF7D', // ズ => ｽ
		'\u30BC': '\uFF7E', // ゼ => ｾ
		'\u30BE': '\uFF7F', // ゾ => ｿ
		'\u30C0': '\uFF80', // ダ => ﾀ
		'\u30C2': '\uFF81', // ヂ => ﾁ
		'\u30C5': '\uFF82', // ヅ => ﾂ
		'\u30C7': '\uFF83', // デ => ﾃ
		'\u30C9': '\uFF84', // ド => ﾄ
		'\u30D0': '\uFF8A', // バ => ﾊ
		'\u30D3': '\uFF8B', // ビ => ﾋ
		'\u30D6': '\uFF8C', // ブ => ﾌ
		'\u30D9': '\uFF8D', // ベ => ﾍ
		'\u30DC': '\uFF8E', // ボ => ﾎ
		'\u30F4': '\uFF73', // ヴ => ｳ
		'\u30F7': '\uFF9C', // ヷ => ﾜ
		'\u30F8': '\uFF72', // ヸ => ｲ
		'\u30F9': '\uFF74', // ヹ => ｴ
		'\u30FA': '\uFF66', // ヺ => ｦ
	}

	// h2zDaku 半角　濁文字: ｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾊﾋﾌﾍﾎｳﾜｦ
	h2zDaku = reverse(z2hDaku)

	// jaZenkakuHandakuRunes 全角　半濁文字: パピプペポ
	z2hHandaku = map[rune]rune{
		'\u30D1': '\uFF8A', // パ => ﾊ
		'\u30D4': '\uFF8B', // ピ => ﾋ
		'\u30D7': '\uFF8C', // プ => ﾌ
		'\u30DA': '\uFF8D', // ペ => ﾍ
		'\u30DD': '\uFF8E', // ポ => ﾎ
	}

	// h2zHandaku 半角　半濁文字: ﾊﾋﾌﾍﾎ
	h2zHandaku = reverse(z2hHandaku)
)

// ToZenkakuRune convert the rune c to zenkaku
func ToZenkakuRune(c rune) rune {
	if r, ok := h2z[c]; ok {
		return r
	}
	return ToFullRune(c)
}

// ToHankakuRune convert the rune c to hankaku
func ToHankakuRune(c rune) rune {
	if r, ok := z2h[c]; ok {
		return r
	}
	return ToASCIIRune(c)
}

// ToZenkakuDakuRune convert the rune c to zenkaku Daku
func ToZenkakuDakuRune(c rune) rune {
	if r, ok := h2zDaku[c]; ok {
		return r
	}
	return c
}

// ToHankakuDakuRune convert the rune c to hankaku Daku
func ToHankakuDakuRune(c rune) rune {
	if r, ok := z2hDaku[c]; ok {
		return r
	}
	return c
}

// ToZenkakuHandakuRune convert the rune c to zenkaku Handaku
func ToZenkakuHandakuRune(c rune) rune {
	if r, ok := h2zHandaku[c]; ok {
		return r
	}
	return c
}

// ToHankakuHandakuRune convert the rune c to hankaku Handaku
func ToHankakuHandakuRune(c rune) rune {
	if r, ok := z2hHandaku[c]; ok {
		return r
	}
	return c
}

// IsHankakuKatakanaRune return true if c is Hankaku Katakana rune
// http://charset.7jp.net/jis0201.html
func IsHankakuKatakanaRune(c rune) bool {
	return c >= '\uFF61' && c <= '\uFF9F'
}

// IsZenkakuKatakanaRune return true if c is Zenkaku Katakana rune
// https://ja.wikipedia.org/wiki/片仮名_(Unicodeのブロック)
func IsZenkakuKatakanaRune(c rune) bool {
	return c >= '\u30A1' && c <= '\u30FC'
}

// IsZenkakuHiraganaRune return true if c is Zenkaku Hiragana rune
// https://ja.wikipedia.org/wiki/平仮名_(Unicodeのブロック)
func IsZenkakuHiraganaRune(c rune) bool {
	return c >= '\u3041' && c <= '\u309F'
}

// IsHankakuRune return true if c is Hankaku rune
func IsHankakuRune(c rune) bool {
	if c < utf8.RuneSelf {
		return true
	}

	return IsHankakuKatakanaRune(c)
}

// IsZenkakuRune return true if c is Zenkaku
func IsZenkakuRune(c rune) bool {
	return !IsHankakuRune(c)
}

// IsHankakuKatakana return true if s is HankakuKatakana string
func IsHankakuKatakana(s string) bool {
	if s == "" {
		return false
	}

	for _, c := range s {
		if !IsHankakuKatakanaRune(c) {
			return false
		}
	}
	return true
}

// IsHankakuKatakanaSpace return true if s is HankakuKatakana or space string
func IsHankakuKatakanaSpace(s string) bool {
	if s == "" {
		return false
	}

	for _, c := range s {
		if c != ' ' && !IsHankakuKatakanaRune(c) {
			return false
		}
	}
	return true
}

// IsHankaku return true if s is Hankaku string
func IsHankaku(s string) bool {
	if s == "" {
		return false
	}

	for _, c := range s {
		if !IsHankakuRune(c) {
			return false
		}
	}
	return true
}

// IsZenkaku return true if s is Zenkaku string
func IsZenkaku(s string) bool {
	if s == "" {
		return false
	}

	for _, c := range s {
		if !IsZenkakuRune(c) {
			return false
		}
	}
	return true
}

// IsZenkakuKatakana return true if s is ZenkakuKatakana string
func IsZenkakuKatakana(s string) bool {
	if s == "" {
		return false
	}

	for _, c := range s {
		if !IsZenkakuKatakanaRune(c) {
			return false
		}
	}
	return true
}

// IsZenkakuKatakanaSpace return true if s is Zenkaku Katakana or Space string
func IsZenkakuKatakanaSpace(s string) bool {
	if s == "" {
		return false
	}

	for _, c := range s {
		if c != '　' && !IsZenkakuKatakanaRune(c) {
			return false
		}
	}
	return true
}

// IsZenkakuKatakanaSpaces return true if s is Zenkaku Katakana or Space (Zenkaku or Hankaku) string
func IsZenkakuKatakanaSpaces(s string) bool {
	if s == "" {
		return false
	}

	for _, c := range s {
		if c != ' ' && c != '　' && !IsZenkakuKatakanaRune(c) {
			return false
		}
	}
	return true
}

// IsZenkakuHiragana return true if s is ZenkakuHiragana string
func IsZenkakuHiragana(s string) bool {
	if s == "" {
		return false
	}

	for _, c := range s {
		if !IsZenkakuHiraganaRune(c) {
			return false
		}
	}
	return true
}

// IsZenkakuHiraganaSpace return true if s is Zenkaku Hiragana or Space string
func IsZenkakuHiraganaSpace(s string) bool {
	if s == "" {
		return false
	}

	for _, c := range s {
		if c != '　' && !IsZenkakuHiraganaRune(c) {
			return false
		}
	}
	return true
}

// IsZenkakuHiraganaSpaces return true if s is Zenkaku Hiragana or Space(Zenkaku or Hankaku) string
func IsZenkakuHiraganaSpaces(s string) bool {
	if s == "" {
		return false
	}

	for _, c := range s {
		if c != ' ' && c != '　' && !IsZenkakuHiraganaRune(c) {
			return false
		}
	}
	return true
}

// ToZenkaku convert the string from hankaku to Zenkaku
func ToZenkaku(s string) string {
	if s == "" {
		return s
	}

	len := len(s)

	sb := &strings.Builder{}
	sb.Grow(len)

	var c rune
	var w int
	for i := 0; i < len; {
		if w == 0 {
			c, w = utf8.DecodeRuneInString(s[i:])
		}
		if i+w < len {
			nc, nw := utf8.DecodeRuneInString(s[i+w:])
			if nc == '\uFF9E' { // ﾞ
				z := ToZenkakuDakuRune(c)
				if z != c {
					sb.WriteRune(z)
					i += w + nw
					w = 0
					continue
				}
			} else if nc == '\uFF9F' { // ﾟ
				z := ToZenkakuHandakuRune(c)
				if z != c {
					sb.WriteRune(z)
					i += w + nw
					w = 0
					continue
				}
			}

			c = ToZenkakuRune(c)
			sb.WriteRune(c)

			i += w
			c, w = nc, nw
			continue
		}

		c = ToZenkakuRune(c)
		sb.WriteRune(c)
		i += w
		w = 0
	}

	return sb.String()
}

// ToHankaku convert the string from zenkaku to hankaku
func ToHankaku(s string) string {
	if s == "" {
		return s
	}

	sb := &strings.Builder{}
	for _, c := range s {
		if IsHankakuRune(c) {
			sb.WriteRune(c)
			continue
		}

		r := ToHankakuDakuRune(c)
		if r != c {
			sb.WriteRune(r)
			sb.WriteRune('\uFF9E') // ﾞ
			continue
		}

		r = ToHankakuHandakuRune(c)
		if r != c {
			sb.WriteRune(r)
			sb.WriteRune('\uFF9F') // ﾟ
			continue
		}

		r = ToHankakuRune(c)
		sb.WriteRune(r)
	}

	return sb.String()
}
