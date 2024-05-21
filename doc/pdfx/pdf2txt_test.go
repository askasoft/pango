package pdfx

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/str"
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

func testSkip(t *testing.T) {
	path, err := exec.LookPath("pdftotext")
	if path == "" || err != nil {
		t.Skip("Failed to find pdftotext", path, err)
	}
}

func TestExtractTextFromPdfFile(t *testing.T) {
	testSkip(t)

	cs := []string{"hello.pdf", "table.pdf"}

	for i, c := range cs {
		fn := testFilename(c)
		a, err := ExtractTextFromPdfFile(context.Background(), fn)
		if err != nil {
			fmt.Printf("[%d] ExtractTextFromPdfFile(%s): %v\n", i, fn, err)
		} else {
			w := string(testReadFile(t, c+".txt"))

			a = str.RemoveByte(a, '\r')
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
	testSkip(t)

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
		err = ExtractStringFromPdfReader(context.Background(), fr, bw)
		if err != nil {
			fmt.Printf("[%d] TestExtractStringFromPdfReader(%s): %v\n", i, fn, err)
			continue
		}

		w := string(testReadFile(t, c+".txt"))
		a := str.RemoveByte(bw.String(), '\r')
		if w != a {
			t.Errorf("[%d] TestExtractStringFromPdfReader(%s):\nACTUAL: %q\n  WANT: %q\n", i, fn, a, w)
			fsu.WriteString(fn+".out", a, fsu.FileMode(0660))
		} else {
			os.Remove(fn + ".out")
		}
	}
}
