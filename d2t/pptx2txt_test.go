package d2t

import (
	"testing"
)

func TestExtractTextFromPptxFile(t *testing.T) {
	cs := []struct {
		f string
		w string
	}{
		{"hello.pptx",
			"Hello world!\n" +
				"EN: How do you do?\n" +
				"JA: 今日はお元気ですか？\n" +
				"ZH: 今天\t好吗？\n\n",
		},
		{"table.pptx",
			"Table\n\n" +
				"This is a cat.\n" +
				"1-1\n\u00a0\n" +
				"That is a dog.\n" +
				"1-2\n\u00a0\n" +
				"This is a fox.\n" +
				"1-3\n\u00a0\n" +
				"This is a bird.\n" +
				"2-1\n\u00a0\n" +
				"This is an apple.\n" +
				"2-2\n\u00a0\n" +
				"That is pear.\n" +
				"2-3\n\u00a0\n",
		},
	}

	for i, c := range cs {
		lw := &LineWriter{}
		err := ExtractTextFromPptxFile("testdata/"+c.f, lw)
		if err != nil {
			t.Errorf("[%d] ExtractTextFromPptxFile(%s): %v", i, c.f, err)
			continue
		}

		a := lw.String()
		if c.w != a {
			t.Errorf("[%d] ExtractTextFromPptxFile(%s):\nACTUAL: %q\n  WANT: %q\n", i, c.f, a, c.w)
		}
	}
}
