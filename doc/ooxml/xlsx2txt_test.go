package ooxml

import (
	"os"
	"testing"

	"github.com/askasoft/pango/iox/fsu"
)

func TestXlsxFileTextifyString(t *testing.T) {
	cs := []string{"hello.xlsx"}

	for i, c := range cs {
		fn := testFilename(c)
		a, err := XlsxFileTextifyString(fn)
		if err != nil {
			t.Errorf("[%d] XlsxFileTextifyString(%s): %v", i, fn, err)
			continue
		}

		w := string(testReadFile(t, c+".txt"))
		if w != a {
			t.Errorf("[%d] XlsxFileTextifyString(%q):\nACTUAL: %q\n  WANT: %q\n", i, fn, a, w)
			fsu.WriteString(fn+".out", a, fsu.FileMode(0660))
		} else {
			os.Remove(fn + ".out")
		}
	}
}
