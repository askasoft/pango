package jpn

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

var (
	// hankakuMark 半角: ｡｢｣､･ﾞﾟ
	hankakuMark = ("｡｢｣､･ﾞﾟ")

	// zenkakuMark 全角: 。「」、・゛゜
	zenkakuMark = ("。「」、・゛゜")

	// hankakuAyatu 半角: ｧｨｩｪｫｬｭｮｯｰ
	hankakuAyatu = ("ｧｨｩｪｫｬｭｮｯｰ")

	// zenkakuAyatu 全角: ァィゥェォャュョッー
	zenkakuAyatu = ("ァィゥェォャュョッー")

	// hankakuAnamayara 半角: ｱｲｴｵﾅﾆﾇﾈﾉﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾝ
	hankakuAnamayara = ("ｱｲｴｵﾅﾆﾇﾈﾉﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾝ")

	// zenkakuAnamayara 全角: アイエオナニヌネノマミムメモヤユヨラリルレロン
	zenkakuAnamayara = ("アイエオナニヌネノマミムメモヤユヨラリルレロン")

	// hankakuKasataha 半角　かさたは　行: ｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾊﾋﾌﾍﾎｳ
	hankakuKasataha = ("ｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾊﾋﾌﾍﾎｳ")

	// zenkakuKasataha 全角　かさたは　行: カキクケコサシスセソタチツテトハヒフヘホウ
	zenkakuKasataha = ("カキクケコサシスセソタチツテトハヒフヘホウ")

	// zenkakuGazadaba 全角　がざだば　行: ガギグゲゴザジズゼゾダヂヅデドバビブベボヴ
	zenkakuGazadaba = ("ガギグゲゴザジズゼゾダヂヅデドバビブベボヴ")

	// hankakuWaou 半角　わ　行: ﾜｦ
	hankakuWaou = ("ﾜｦ")

	// zenkakuWaou 全角　わ　行: ワヲ
	zenkakuWaou = ("ワヲ")

	// zenkakuVavo 全角　わ　行: ヷヺ
	zenkakuVavo = ("ヷヺ")

	// hankakuVive 半角　: ｲｴ
	hankakuVive = ("ｲｴ")

	// zenkakuVavo 全角　: ヸヹ
	zenkakuVive = ("ヸヹ")

	// hankakuKanakana 半角
	hankakuKatakana = hankakuAnamayara + hankakuAyatu + hankakuDaku + hankakuHandaku + hankakuKasataha + hankakuWaou + hankakuVive

	// zenkakuKanakana 全角
	zenkakuKatakana = zenkakuAnamayara + zenkakuAyatu + zenkakuDaku + zenkakuHandaku + zenkakuKasataha + zenkakuWaou + zenkakuVive

	// hankakuDaku 半角　濁文字: ｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾊﾋﾌﾍﾎｳﾜｦｲｴ
	hankakuDaku      = hankakuKasataha + hankakuWaou + hankakuVive
	hankakuDakuRunes = []rune(hankakuDaku)

	// zenkakuDaku 全角　濁文字: ガギグゲゴザジズゼゾダヂヅデドバビブベボヴヷヺヸヹ
	zenkakuDaku = zenkakuGazadaba + zenkakuVavo + zenkakuVive

	// hankakuHandaku 半角　半濁文字: ﾊﾋﾌﾍﾎ
	hankakuHandaku      = ("ﾊﾋﾌﾍﾎ")
	hankakuHandakuRunes = []rune(hankakuHandaku)

	// zenkakuHandaku 全角　半濁文字: パピプペポ
	zenkakuHandaku = ("パピプペポ")

	// hankaku 半角
	hankaku      = hankakuMark + hankakuAyatu + hankakuAnamayara + hankakuKasataha + hankakuWaou
	hankakuRunes = []rune(hankaku)

	// zenkaku 全角
	zenkaku      = zenkakuMark + zenkakuAyatu + zenkakuAnamayara + zenkakuKasataha + zenkakuWaou
	zenkakuRunes = []rune(zenkaku)
)

func specialZ(c rune) rune {
	switch c {
	case '“', '”', '″':
		return '＂'
	case '〜':
		return '～'
	case '’':
		return '＇'
	case '\u00a0':
		return '\u3000'
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

func testJapan(t *testing.T, name, zen, han string) {
	testPrint(t, name, zen, han, nil, nil)
}

func TestJapanPair(t *testing.T) {
	testJapan(t, "Mark:     ", zenkakuMark, hankakuMark)
	testJapan(t, "Ayatu:    ", zenkakuAyatu, hankakuAyatu)
	testJapan(t, "Anamayara:", zenkakuAnamayara, hankakuAnamayara)
	testJapan(t, "Kasataha: ", zenkakuKasataha, hankakuKasataha)
	testJapan(t, "Waou:     ", zenkakuWaou, hankakuWaou)

	testJapan(t, "Handaku:  ", zenkakuHandaku, hankakuHandaku)
	testJapan(t, "Daku:  ", zenkakuDaku, hankakuDaku)
}

func TestJapanConvert(t *testing.T) {
	l := len(zenkakuRunes)
	for i := 0; i < 100; i++ {
		src := &strings.Builder{}
		han := &strings.Builder{}
		zen := &strings.Builder{}
		for n := 0; n < 100; n++ {
			p := rand.Intn(l)
			h := hankakuRunes[p]
			z := zenkakuRunes[p]
			src.WriteRune(h)
			src.WriteRune(z)
			han.WriteRune(h)
			han.WriteRune(h)
			zen.WriteRune(z)
			zen.WriteRune(z)

			s := src.String()
			e := han.String()
			a := ToHankaku(s)
			if a != e {
				t.Fatalf("ToHankaku(%q) = %q, want %q", s, a, e)
			}

			e = zen.String()
			a = ToZenkaku(s)
			if a != e {
				t.Fatalf("ToZenkaku(%q) = %q, want %q", s, a, e)
			}
		}
	}
}

func TestConvertDaku(t *testing.T) {
	sb := &strings.Builder{}
	for _, c := range hankakuDakuRunes {
		sb.WriteRune(c)
		sb.WriteRune('ﾞ')
	}
	han := sb.String()

	s := zenkakuDaku
	a := ToHankaku(s)
	e := han
	if a != e {
		t.Errorf("ToHankaku(%q) = %q, want %q", s, a, e)
	}

	s = han
	a = ToZenkaku(s)
	e = zenkakuDaku
	if a != e {
		t.Errorf("ToZenkaku(%q) = %q, want %q", s, a, e)
	}
}

func TestConvertHandaku(t *testing.T) {
	sb := &strings.Builder{}
	for _, c := range hankakuHandakuRunes {
		sb.WriteRune(c)
		sb.WriteRune('ﾟ')
	}
	han := sb.String()

	e := han
	s := zenkakuHandaku
	a := ToHankaku(s)
	if a != e {
		t.Errorf("ToHankaku(%q) = %q, want %q", s, a, e)
	}

	s = han
	a = ToZenkaku(s)
	e = zenkakuHandaku
	if a != e {
		t.Errorf("ToZenkaku(%q) = %q, want %q", s, a, e)
	}
}

func TestIsHankakuKatakanaRune(t *testing.T) {
	fmt.Println(hankakuKatakana)
	for _, c := range hankakuKatakana {
		a := IsHankakuKatakanaRune(c)
		if !a {
			t.Errorf("IsHankakuKatakanaRune(%q) = %v, want %v", c, a, true)
		}
	}
}

func TestIsHankakuKatakana(t *testing.T) {
	fmt.Println(hankakuKatakana)
	a := IsHankakuKatakana(hankakuKatakana)
	if !a {
		t.Errorf("IsHankakuKatakana(%q) = %v, want %v", hankakuKatakana, a, true)
	}
}

func TestIsZenkakuKatakanaRune(t *testing.T) {
	fmt.Println(zenkakuKatakana)
	for _, c := range zenkakuKatakana {
		a := IsZenkakuKatakanaRune(c)
		if !a {
			t.Errorf("IsZenkakuKatakanaRune(%q) = %v, want %v", c, a, true)
		}
	}
}

func TestIsZenkakuKatakana(t *testing.T) {
	fmt.Println(zenkakuKatakana)
	a := IsZenkakuKatakana(zenkakuKatakana)
	if !a {
		t.Errorf("IsZenkakuKatakana(%q) = %v, want %v", zenkakuKatakana, a, true)
	}
}
