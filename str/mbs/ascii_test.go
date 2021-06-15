package mbs

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// halfNumber 半角数字: 0123456789
	halfNumber = ("0123456789")

	// fullNumber 全角数字: ０１２３４５６７８９
	fullNumber = ("０１２３４５６７８９")

	// halfLetter 半角英字: ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz
	halfLetter = ("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

	// fullLetter 全角英字: ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ
	fullLetter = ("ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ")

	// halfSymbol 半角記号: !"""#$%&'()*+,----./:;<=>?@[\]^_`{|}~~
	halfSymbol = (" !\"\"\"#$%&'()*+,----./:;<=>?@[\\]^_`{|}~~")

	// fullSymbol 全角記号: ！″“”＃＄％＆’（）＊＋，－ー‐−．／：；＜＝＞？＠［￥］＾＿｀｛｜｝～〜
	fullSymbol = ("　！″”“＃＄％＆’（）＊＋，－ー‐−．／：；＜＝＞？＠［￥］＾＿｀｛｜｝～〜")

	// halfASCII 半角: !"#$%&'()*+,-./0123456789:<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_`abcdefghijklmnopqrstuvwxyz{|}~~
	halfASCII      = halfNumber + halfLetter + halfSymbol
	halfASCIIRunes = []rune(halfASCII)

	// fullASCII 全角: ！“”＃＄％＆’（）＊＋，－．／０１２３４５６７８９：；＜＝＞？＠ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺ［￥］＾＿｀ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ｛｜｝～〜
	fullASCII      = fullNumber + fullLetter + fullSymbol
	fullASCIIRunes = []rune(fullASCII)
)

func testASCII(t *testing.T, name, zen, han string) {
	testPrint(t, name, zen, han, ToHalfRune, ToFullRune)
}

func specialZ(c rune) rune {
	switch c {
	case '“', '”':
		return '″'
	case '〜':
		return '～'
	case 'ー', '‐', '−':
		return '－'
	default:
		return c
	}
}

func testPrint(t *testing.T, name, zens, hans string, fz2h func(c rune) rune, fh2z func(c rune) rune) {
	zen := []rune(zens)
	han := []rune(hans)
	assert.Equal(t, len(zen), len(han), name)

	fmt.Println("// " + name)
	for i := 0; i < len(zen) && i < len(han); i++ {
		z := zen[i]
		h := han[i]

		if fh2z != nil && fz2h != nil {
			z2h := fz2h(z)
			h2z := fh2z(h)
			assert.Equal(t, h, z2h, string(h)+" <> "+string(z2h))
			sz := specialZ(z)
			assert.Equal(t, sz, h2z, string(sz)+" <> "+string(h2z))
		}
		fmt.Printf("'\\u%04X': '\\u%04X', // %s => %s\n", z, h, string(z), string(h))
	}
}

func TestASCIIPair(t *testing.T) {
	testASCII(t, "Number", fullNumber, halfNumber)
	testASCII(t, "Letter", fullLetter, halfLetter)
	testASCII(t, "Symbol", fullSymbol, halfSymbol)
	testASCII(t, "ASCII", fullASCII, halfASCII)
	testASCII(t, "ASCII", fullASCII, halfASCII)
}

func TestASCIIConvert(t *testing.T) {
	l := len(fullASCIIRunes)
	for i := 0; i < 100; i++ {
		src := &strings.Builder{}
		han := &strings.Builder{}
		zen := &strings.Builder{}
		for n := 0; n < 100; n++ {
			p := rand.Intn(l)
			h := halfASCIIRunes[p]
			z := fullASCIIRunes[p]
			src.WriteRune(h)
			src.WriteRune(z)
			han.WriteRune(h)
			han.WriteRune(h)
			sz := specialZ(z)
			zen.WriteRune(sz)
			zen.WriteRune(z)

			s := src.String()
			assert.Equal(t, han.String(), ToHalfWidth(s), s)
			assert.Equal(t, zen.String(), ToFullWidth(s), s)
		}
	}
}
