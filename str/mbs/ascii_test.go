package mbs

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
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

	// halfSymbol 半角記号: !""""#$%&'()*+,----./:;<=>?@[\\]^_`{|}~~
	halfSymbol = (" !\"\"\"\"#$%&'()*+,----./:;<=>?@[\\\\]^_`{|}~~")

	// fullSymbol 全角記号: ！″＂”“＃＄％＆’（）＊＋，－ー‐−．／：；＜＝＞？＠［＼￥］＾＿｀｛｜｝～〜
	fullSymbol = ("　！″＂”“＃＄％＆’（）＊＋，－ー‐−．／：；＜＝＞？＠［＼￥］＾＿｀｛｜｝～〜")

	// halfASCII 半角: !"#$%&'()*+,-./0123456789:<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_`abcdefghijklmnopqrstuvwxyz{|}~~
	halfASCII      = halfNumber + halfLetter + halfSymbol
	halfASCIIRunes = []rune(halfASCII)

	// fullASCII 全角: ！“”＃＄％＆’（）＊＋，－．／０１２３４５６７８９：；＜＝＞？＠ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺ［￥］＾＿｀ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ｛｜｝～〜
	fullASCII      = fullNumber + fullLetter + fullSymbol
	fullASCIIRunes = []rune(fullASCII)
)

func testASCII(t *testing.T, name, zen, han string) {
	testPrint(t, name, zen, han, ToASCIIRune, ToFullRune)
}

func specialZ(c rune) rune {
	switch c {
	case '“', '”', '″':
		return '＂'
	case '〜':
		return '～'
	case '’':
		return '＇'
	case '￥':
		return '＼'
	case 'ー', '‐', '−':
		return '－'
	default:
		return c
	}
}

func testPrint(t *testing.T, name, zens, hans string, fz2h func(c rune) rune, fh2z func(c rune) rune) {
	zen := []rune(zens)
	han := []rune(hans)

	if len(zen) != len(han) {
		t.Errorf("[%s] len(%q) != len(%q)", name, zen, han)
		return
	}

	fmt.Println("// " + name)
	for i := 0; i < len(zen) && i < len(han); i++ {
		z := zen[i]
		h := han[i]

		if fh2z != nil && fz2h != nil {
			z2h := fz2h(z)
			h2z := fh2z(h)

			if h != z2h {
				t.Errorf("[%s] %q != %q", name, h, z2h)
			}

			sz := specialZ(z)
			if sz != h2z {
				t.Errorf("[%s] %q != %q", name, sz, h2z)
			}
		}
		fmt.Printf("'\\u%04X': '\\u%04X', // %s => %s\n", z, h, string(z), string(h))
	}
}

func TestASCIIPair(t *testing.T) {
	testASCII(t, "Number", fullNumber, halfNumber)
	testASCII(t, "Letter", fullLetter, halfLetter)
	testASCII(t, "Symbol", fullSymbol, halfSymbol)
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

			e := han.String()
			a := ToASCII(s)
			if e != a {
				t.Errorf("ToHalfWidth(%q) = %q, want %q", s, a, e)
			}

			e = zen.String()
			a = ToFullWidth(s)
			if e != a {
				t.Errorf("ToFullWidth(%q) = %q, want %q", s, a, e)
			}
		}
	}
}

func TestHasFullWidth(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"", false},
		{"abc", false},
		{"abc123", false},
		{"!\"#$%&()<>/+=-_? ~^|.,@`{}[]", false},
		{"ひらがな・カタカナ、．漢字", true},
		{"３ー０　ａ＠ｃｏｍ", true},
		{"Ｆｶﾀｶﾅﾞﾬ", true},
		{"Good＝Parts", true},
		{"", false},
	}
	for i, c := range cs {
		a := HasFullWidth(c.s)
		if a != c.w {
			t.Errorf("[%d] HasFullWidth(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestHasHalfWidth(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"", false},
		{"あいうえお", false},
		{"００１１", false},
		{"!\"#$%&()<>/+=-_? ~^|.,@`{}[]", true},
		{"l-btn_02--active", true},
		{"abc123い", true},
		{"ｶﾀｶﾅﾞﾬ￩", true},
		{"", false},
	}
	for i, c := range cs {
		a := HasHalfWidth(c.s)
		if a != c.w {
			t.Errorf("[%d] HasHalfWidth(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsVariableWidth(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"", false},
		{"ひらがなカタカナ漢字ABCDE", true},
		{"３ー０123", true},
		{"Ｆｶﾀｶﾅﾞﾬ", true},
		{"", false},
		{"Good＝Parts", true},
		{"abc", false},
		{"abc123", false},
		{"!\"#$%&()<>/+=-_? ~^|.,@`{}[]", false},
		{"ひらがな・カタカナ、．漢字", false},
		{"１２３４５６", false},
		{"ｶﾀｶﾅﾞﾬ", false},
	}
	for i, c := range cs {
		a := IsVariableWidth(c.s)
		if a != c.w {
			t.Errorf("[%d] IsVariableWidth(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}
