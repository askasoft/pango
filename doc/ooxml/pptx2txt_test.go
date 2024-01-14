package ooxml

import (
	"os"
	"testing"

	"github.com/askasoft/pango/fsu"
)

func TestExtractTextFromPptxFile(t *testing.T) {
	cs := []string{"hello.pptx", "table.pptx"}

	for i, c := range cs {
		fn := testFilename(c)
		a, err := ExtractTextFromPptxFile(fn)
		if err != nil {
			t.Errorf("[%d] ExtractTextFromPptxFile(%q): %v", i, fn, err)
			continue
		}

		w := string(testReadFile(t, c+".txt"))
		if w != a {
			t.Errorf("[%d] ExtractTextFromPptxFile(%q):\nACTUAL: %q\n  WANT: %q\n", i, fn, a, w)
			fsu.WriteString(fn+".out", a, fsu.FileMode(0660))
		} else {
			os.Remove(fn + ".out")
		}
	}
}
