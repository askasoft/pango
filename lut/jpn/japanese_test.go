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
			a2 := ToHankaku(a)
			if a2 != a {
				t.Fatalf("ToHankaku(%q) = %q, want %q", a, a2, a)
			}

			e = zen.String()
			a = ToZenkaku(s)
			if a != e {
				t.Fatalf("ToZenkaku(%q) = %q, want %q", s, a, e)
			}
			a2 = ToZenkaku(a)
			if a2 != a {
				t.Fatalf("ToZenkaku(%q) = %q, want %q", a, a2, a)
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

func TestHiraganaToKatagana(t *testing.T) {
	cs := []struct {
		hira string
		kata string
	}{
		{`ぁあぃいぅうぇえぉおかがきぎく
ぐけげこごさざしじすずせぜそぞた
だちぢっつづてでとどなにぬねのは
ばぱひびぴふぶぷへべぺほぼぽまみ
むめもゃやゅゆょよらりるれろゎわ
ゐゑをんゔゕゖ`, `ァアィイゥウェエォオカガキギク
グケゲコゴサザシジスズセゼソゾタ
ダチヂッツヅテデトドナニヌネノハ
バパヒビピフブプヘベペホボポマミ
ムメモャヤュユョヨラリルレロヮワ
ヰヱヲンヴヵヶ`},
		{"", ""},
		{"abc", "abc"},
		{"aぁあぃいぅうぇえぉおかがきぎくz", "aァアィイゥウェエォオカガキギクz"},
	}

	for i, c := range cs {
		h2k := HiraganaToKatagana(c.hira)
		if h2k != c.kata {
			t.Errorf("#%d HiraganaToKatagana(%q) = %q, want %q", i, c.hira, h2k, c.kata)
		}

		k2h := KataganaToHiragana(c.kata)
		if k2h != c.hira {
			t.Errorf("#%d KataganaToHiragana(%q) = %q, want %q", i, c.kata, k2h, c.hira)
		}
	}
}

func TestCompareKana(t *testing.T) {
	tcs := []struct {
		s, t string
		out  int
	}{
		{"j", "J", 1},
		{"ぁあぃいぅうぇえぉお", "ァアィイゥウェエォオ", 0},
		{"abcdぁあぃいぅうぇえぉおefghij", "abcdァアィイゥウェエォオefghij", 0},
		{"abcdぁあぃいぅうぇえぉおefghijk", "abcdァアィイゥウェエォオefghij", 1},
		{"abcdぁあぃいぅうぇえぉおefghij", "abcdァアィイゥウェエォオefghiJ", 1},
	}

	for _, tt := range tcs {
		if out := CompareKana(tt.s, tt.t); out != tt.out {
			t.Errorf("CompareKana(%#q, %#q) = %v, want %v", tt.s, tt.t, out, tt.out)
		}
		if out := CompareKana(tt.t, tt.s); out != -tt.out {
			t.Errorf("CompareKana(%#q, %#q) = %v, want %v", tt.t, tt.s, out, -tt.out)
		}
	}
}

func TestCompareFoldKana(t *testing.T) {
	tcs := []struct {
		s, t string
		out  int
	}{
		{"ぁあぃいぅうぇえぉお", "ァアィイゥウェエォオ", 0},
		{"abcdぁあぃいぅうぇえぉおefghij", "aBcdァアィイゥウェエォオeFghij", 0},
		{"abc", "abc", 0},
		{"ABcd", "ABcd", 0},
		{"123abc", "123ABC", 0},
		{"abc", "xyz", -1},
		{"abc", "XYZ", -1},
		{"αβδ", "ΑΒΔ", 0},
		{"abcdefghijk", "abcdefghijX", -1},
		{"abcdefghijk", "abcdefghij\u212A", 0},
		{"abcdefghijK", "abcdefghij\u212A", 0},
		{"abcdefghijkz", "abcdefghij\u212Ay", 1},
		{"abcdefghijKz", "abcdefghij\u212Ay", 1},
		{"1", "2", -1},
		{"utf-8", "US-ASCII", 1},
	}

	for _, tt := range tcs {
		if out := CompareFoldKana(tt.s, tt.t); out != tt.out {
			t.Errorf("CompareFoldKana(%#q, %#q) = %v, want %v", tt.s, tt.t, out, tt.out)
		}
		if out := CompareFoldKana(tt.t, tt.s); out != -tt.out {
			t.Errorf("CompareFoldKana(%#q, %#q) = %v, want %v", tt.t, tt.s, out, -tt.out)
		}
	}
}
