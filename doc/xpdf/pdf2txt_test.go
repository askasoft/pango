package xpdf

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/askasoft/pango/iox/fsu"
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

func TestPdfFileTextifyString(t *testing.T) {
	testSkip(t)

	cs := []string{"hello.pdf", "table.pdf"}

	for i, c := range cs {
		fn := testFilename(c)
		a, err := PdfFileTextifyString(context.Background(), fn, "-layout")
		if err != nil {
			fmt.Printf("[%d] PdfFileTextifyString(%s): %v\n", i, fn, err)
		} else {
			w := string(testReadFile(t, c+".txt"))

			a = str.RemoveByte(a, '\r')
			if w != a {
				t.Errorf("[%d] PdfFileTextifyString(%s):\nACTUAL: %q\n  WANT: %q\n", i, fn, a, w)
				fsu.WriteString(fn+".out", a, fsu.FileMode(0660))
			} else {
				os.Remove(fn + ".out")
			}
		}
	}
}

func TestPdfReaderTextify(t *testing.T) {
	testSkip(t)

	cs := []string{"hello.pdf", "table.pdf"}

	for i, c := range cs {
		fn := testFilename(c)
		fr, err := os.Open(fn)
		if err != nil {
			t.Errorf("[%d] PdfReaderTextify(%s): %v\n", i, fn, err)
			continue
		}
		defer fr.Close()

		bw := &bytes.Buffer{}
		err = PdfReaderTextify(context.Background(), bw, fr, "-layout")
		if err != nil {
			fmt.Printf("[%d] PdfReaderTextify(%s): %v\n", i, fn, err)
			continue
		}

		w := string(testReadFile(t, c+".txt"))
		a := str.RemoveByte(bw.String(), '\r')
		if w != a {
			t.Errorf("[%d] PdfReaderTextify(%s):\nACTUAL: %q\n  WANT: %q\n", i, fn, a, w)
			fsu.WriteString(fn+".out", a, fsu.FileMode(0660))
		} else {
			os.Remove(fn + ".out")
		}
	}
}
