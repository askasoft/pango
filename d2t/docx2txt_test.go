package d2t

import (
	"testing"
)

func TestExtractTextFromDocxFile(t *testing.T) {
	cs := []struct {
		f string
		w string
	}{
		{"hello.docx",
			"Hello world!\n" +
				"EN:\nHow do you do?\n" +
				"JA:\n今日はお元気ですか？\n" +
				"ZH:\n今天\t\t好吗？\n",
		},
		{"history.docx",
			"History\n\n" +
				"Is this my bookmark?\n\n",
		},
		{"table.docx",
			"Table\n\n" +
				"This is a cat.\n" +
				"1-1\n\n" +
				"That is a dog.\n" +
				"1-2\n\n" +
				"This is a fox.\n" +
				"1-3\n\n" +
				"This is a bird.\n" +
				"2-1\n\n" +
				"This is an apple.\n" +
				"2-2\n\n" +
				"That is pear.\n" +
				"2-3\n\n\n\n",
		},
	}

	for i, c := range cs {
		lw := &LineWriter{}
		err := ExtractTextFromDocxFile("testdata/"+c.f, lw)
		if err != nil {
			t.Errorf("[%d] ExtractTextFromDocxFile(%s): %v", i, c.f, err)
			continue
		}

		a := lw.String()
		if c.w != a {
			t.Errorf("[%d] ExtractTextFromDocxFile(%s):\nACTUAL: %q\n  WANT: %q\n", i, c.f, a, c.w)
		}
	}
}
