package pdfx

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/askasoft/pango/fsu"
)

func testFilename(name string) string {
	return filepath.Join("testdata", name)
}

func testReadFile(t *testing.T, name string) []byte {
	fn := testFilename(name)
	bs, err := fsu.ReadFile(fn)
	if err != nil {
		t.Fatalf("Failed to read file %q: %v", fn, err)
	}
	return bs
}

func TestExtractTextFromPdfFile(t *testing.T) {
	cs := []string{"hello.pdf", "table.pdf"}

	for i, c := range cs {
		fn := testFilename(c)
		a, err := ExtractTextFromPdfFile(fn)
		if err != nil {
			fmt.Printf("[%d] ExtractTextFromPdfFile(%s): %v\n", i, fn, err)
		} else {
			w := string(testReadFile(t, c+".txt"))
			if w != a {
				t.Errorf("[%d] ExtractTextFromPdfFile(%s):\nACTUAL: %q\n  WANT: %q\n", i, fn, a, w)
				fsu.WriteString(fn+".out", a, fsu.FileMode(0660))
			} else {
				os.Remove(fn + ".out")
			}
		}
	}
}

func TestExtractStringFromPdfReader(t *testing.T) {
	cs := []string{"hello.pdf", "table.pdf"}

	for i, c := range cs {
		fn := testFilename(c)
		fr, err := os.Open(fn)
		if err != nil {
			t.Errorf("[%d] TestExtractStringFromPdfReader(%s): %v\n", i, fn, err)
			continue
		}
		defer fr.Close()

		bw := &bytes.Buffer{}
		err = ExtractStringFromPdfReader(fr, bw)
		if err != nil {
			fmt.Printf("[%d] ExtractTextFromPdfFile(%s): %v\n", i, fn, err)
			continue
		}

		w := string(testReadFile(t, c+".txt"))
		a := bw.String()
		if w != a {
			t.Errorf("[%d] ExtractTextFromPdfFile(%s):\nACTUAL: %q\n  WANT: %q\n", i, fn, a, w)
			fsu.WriteString(fn+".out", a, fsu.FileMode(0660))
		} else {
			os.Remove(fn + ".out")
		}
	}
}
