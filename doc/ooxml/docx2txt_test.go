package ooxml

import (
	"os"
	"path"
	"testing"

	"github.com/askasoft/pango/fsu"
)

func testFilename(name string) string {
	return path.Join("testdata", name)
}

func testReadFile(t *testing.T, name string) []byte {
	fn := testFilename(name)
	bs, err := fsu.ReadFile(fn)
	if err != nil {
		t.Fatalf("Failed to read file %q: %v", fn, err)
	}
	return bs
}

func TestExtractTextFromDocxFile(t *testing.T) {
	cs := []string{"hello.docx", "history.docx", "table.docx"}

	for i, c := range cs {
		fn := testFilename(c)
		a, err := ExtractTextFromDocxFile(fn)
		if err != nil {
			t.Errorf("[%d] ExtractTextFromDocxFile(%q): %v", i, fn, err)
			continue
		}

		w := string(testReadFile(t, c+".txt"))
		if w != a {
			t.Errorf("[%d] ExtractTextFromDocxFile(%q):\nACTUAL: %q\n  WANT: %q\n", i, fn, a, w)
			fsu.WriteString(fn+".out", a, fsu.FileMode(0660))
		} else {
			os.Remove(fn + ".out")
		}
	}
}
