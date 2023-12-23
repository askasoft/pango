package d2t

import (
	"testing"
)

func TestExtractTextFromXlsxFile(t *testing.T) {
	cs := []struct {
		f string
		w string
	}{
		{"hello.xlsx",
			"Hello world!\n" +
				"JA:\n今日はお元気ですか？\n" +
				"EN:\nHow do you do?\n" +
				"\u00a0\n" +
				"犬\n" +
				"ZH:\n今天\t好吗？\n",
		},
	}

	for i, c := range cs {
		lw := &LineWriter{}
		err := ExtractTextFromXlsxFile("testdata/"+c.f, lw)
		if err != nil {
			t.Errorf("[%d] ExtractTextFromXlsxFile(%s): %v", i, c.f, err)
			continue
		}

		a := lw.String()
		if c.w != a {
			t.Errorf("[%d] ExtractTextFromXlsxFile(%s):\nACTUAL: %q\n  WANT: %q\n", i, c.f, a, c.w)
		}
	}
}
