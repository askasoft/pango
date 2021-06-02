package mbs

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	zenkakuDaku      = zenkakuGazadaba + zenkakuVavo + zenkakuVive
	zenkakuDakuRunes = []rune(zenkakuDaku)

	// hankakuHandaku 半角　半濁文字: ﾊﾋﾌﾍﾎ
	hankakuHandaku      = ("ﾊﾋﾌﾍﾎ")
	hankakuHandakuRunes = []rune(hankakuHandaku)

	// zenkakuHandaku 全角　半濁文字: パピプペポ
	zenkakuHandaku     = ("パピプペポ")
	zenkakuHadakuRunes = []rune(zenkakuHandaku)

	// hankaku 半角
	hankaku      = hankakuMark + hankakuAyatu + hankakuAnamayara + hankakuKasataha + hankakuWaou
	hankakuRunes = []rune(hankaku)

	// zenkaku 全角
	zenkaku      = zenkakuMark + zenkakuAyatu + zenkakuAnamayara + zenkakuKasataha + zenkakuWaou
	zenkakuRunes = []rune(zenkaku)
)

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
			assert.Equal(t, han.String(), ToHankaku(s), "h:"+s)
			assert.Equal(t, zen.String(), ToZenkaku(s), "z:"+s)
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
	assert.Equal(t, han, ToHankaku(zenkakuDaku))
	assert.Equal(t, zenkakuDaku, ToZenkaku(han))
}

func TestConvertHandaku(t *testing.T) {
	sb := &strings.Builder{}
	for _, c := range hankakuHandakuRunes {
		sb.WriteRune(c)
		sb.WriteRune('ﾟ')
	}
	han := sb.String()
	assert.Equal(t, han, ToHankaku(zenkakuHandaku))
	assert.Equal(t, zenkakuHandaku, ToZenkaku(han))
}

func TestIsHankakuKatakanaRune(t *testing.T) {
	fmt.Println(hankakuKatakana)
	for _, c := range hankakuKatakana {
		assert.True(t, IsHankakuKatakanaRune(c), fmt.Sprintf("%04X %s", c, string(c)))
	}
}

func TestIsHankakuKatakana(t *testing.T) {
	fmt.Println(hankakuKatakana)
	assert.True(t, IsHankakuKatakana(hankakuKatakana), hankakuKatakana)
}

func TestIsZenkakuKatakanaRune(t *testing.T) {
	fmt.Println(zenkakuKatakana)
	for _, c := range zenkakuKatakana {
		assert.True(t, IsZenkakuKatakanaRune(c), fmt.Sprintf("%04X %s", c, string(c)))
	}
}

func TestIsZenkakuKatakana(t *testing.T) {
	fmt.Println(zenkakuKatakana)
	assert.True(t, IsZenkakuKatakana(zenkakuKatakana), zenkakuKatakana)
}
