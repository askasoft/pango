package htmlx

import (
	"os"
	"testing"

	"github.com/askasoft/pango/iox/fsu"
)

func TestHTMLFileTextifyString(t *testing.T) {
	cs := []string{
		"utf-8",
		"shift-jis",
	}

	for i, c := range cs {
		sfn := testFilename(c + ".html")
		ofn := testFilename(c + ".out")

		a, err := HTMLFileTextifyString(sfn, 1024)
		if err != nil {
			t.Fatalf("[%d] Failed to HTMLFileTextifyString(%q): %v", i, c, err)
		}

		w := string(testReadFile(t, c+".text"))
		if w != a {
			t.Errorf("[%d] HTMLFileTextifyString(%q):\n  GOT: %q\n WANT: %q\n", i, c, a, w)
			fsu.WriteString(ofn, a, fsu.FileMode(0660))
		} else {
			os.Remove(ofn)
		}
	}
}

func TestHTMLTextifyString(t *testing.T) {
	cs := []string{"utf-8"}

	for i, c := range cs {
		ofn := testFilename(c + ".out")

		s := string(testReadFile(t, c+".html"))
		a, err := HTMLTextifyString(s)
		if err != nil {
			t.Fatalf("[%d] Failed to HTMLTextifyString(%q): %v", i, c, err)
		}

		w := string(testReadFile(t, c+".text"))
		if w != a {
			t.Errorf("[%d] HTMLTextifyString(%q):\n  GOT: %q\n WANT: %q\n", i, c, a, w)
			fsu.WriteString(ofn, a, fsu.FileMode(0660))
		} else {
			os.Remove(ofn)
		}
	}
}
