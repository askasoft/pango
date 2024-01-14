package ooxml

import (
	"os"
	"testing"

	"github.com/askasoft/pango/fsu"
)

func TestExtractTextFromXlsxFile(t *testing.T) {
	cs := []string{"hello.xlsx"}

	for i, c := range cs {
		fn := testFilename(c)
		a, err := ExtractTextFromXlsxFile(fn)
		if err != nil {
			t.Errorf("[%d] ExtractTextFromXlsxFile(%s): %v", i, fn, err)
			continue
		}

		w := string(testReadFile(t, c+".txt"))
		if w != a {
			t.Errorf("[%d] ExtractTextFromXlsxFile(%q):\nACTUAL: %q\n  WANT: %q\n", i, fn, a, w)
			fsu.WriteString(fn+".out", a, fsu.FileMode(0660))
		} else {
			os.Remove(fn + ".out")
		}
	}
}
