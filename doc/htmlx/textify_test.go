package htmlx

import (
	"os"
	"testing"

	"github.com/askasoft/pango/fsu"
)

func TestExtractTextFromHTMLFile(t *testing.T) {
	cs := []string{"utf-8.html", "shift-jis.html"}

	w := string(testReadFile(t, "expect.txt"))

	for i, c := range cs {
		fn := testFilename(c)
		a, err := ExtractTextFromHTMLFile(fn)
		if err != nil {
			t.Fatalf("[%d] Failed to ExtractTextFromHTMLFile(%q): %v", i, c, err)
		}

		if w != a {
			t.Errorf("[%d] ExtractTextFromHTMLFile(%q):\n  GOT: %q\n WANT: %q\n", i, c, a, w)
			fsu.WriteString(fn+".out", a, fsu.FileMode(0660))
		} else {
			os.Remove(fn + ".out")
		}
	}
}

func TestExtractTextFromHTMLString(t *testing.T) {
	cs := []string{"utf-8.html"}

	w := string(testReadFile(t, "expect.txt"))

	for i, c := range cs {
		fn := testFilename(c)
		s := string(testReadFile(t, c))
		a, err := ExtractTextFromHTMLString(s)
		if err != nil {
			t.Fatalf("[%d] Failed to ExtractTextFromHTMLString(%q): %v", i, c, err)
		}

		if w != a {
			t.Errorf("[%d] ExtractTextFromHTMLString(%q):\n  GOT: %q\n WANT: %q\n", i, c, a, w)
			fsu.WriteString(fn+".out", a, fsu.FileMode(0660))
		} else {
			os.Remove(fn + ".out")
		}
	}
}
