package ocrx

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
	path, err := exec.LookPath("tesseract")
	if path == "" || err != nil {
		t.Skip("Failed to find tesseract", path, err)
	}
}

func TestExtractTextFromPdfFile(t *testing.T) {
	testSkip(t)

	cs := []string{"jpn.jpg"}

	for i, c := range cs {
		ln := str.SubstrBeforeByte(c, '.')
		fn := testFilename(c)
		a, err := ExtractTextFromImgFile(context.Background(), fn, ln)
		if err != nil {
			fmt.Printf("[%d] ExtractTextFromImgFile(%s): %v\n", i, fn, err)
		} else {
			w := string(testReadFile(t, c+".txt"))

			a = str.RemoveByte(a, '\r')
			if w != a {
				t.Errorf("[%d] ExtractTextFromImgFile(%s):\nACTUAL: %q\n  WANT: %q\n", i, fn, a, w)
				fsu.WriteString(fn+".out", a, fsu.FileMode(0660))
			} else {
				os.Remove(fn + ".out")
			}
		}
	}
}

func TestExtractStringFromImgReader(t *testing.T) {
	testSkip(t)

	cs := []string{"jpn.jpg"}

	for i, c := range cs {
		ln := str.SubstrBeforeByte(c, '.')
		fn := testFilename(c)
		fr, err := os.Open(fn)
		if err != nil {
			t.Errorf("[%d] TestExtractStringFromImgReader(%s): %v\n", i, fn, err)
			continue
		}
		defer fr.Close()

		bw := &bytes.Buffer{}
		err = ExtractStringFromImgReader(context.Background(), bw, fr, ln)
		if err != nil {
			fmt.Printf("[%d] TestExtractStringFromImgReader(%s): %v\n", i, fn, err)
			continue
		}

		w := string(testReadFile(t, c+".txt"))
		a := str.RemoveByte(bw.String(), '\r')
		if w != a {
			t.Errorf("[%d] TestExtractStringFromImgReader(%s):\nACTUAL: %q\n  WANT: %q\n", i, fn, a, w)
			fsu.WriteString(fn+".out", a, fsu.FileMode(0660))
		} else {
			os.Remove(fn + ".out")
		}
	}
}
