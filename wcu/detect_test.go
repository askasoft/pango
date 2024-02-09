package wcu

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/str"
)

var cs = []string{
	"euc-jp", "iso-2022-jp", "shift-jis", "utf-16be", "utf-8", "utf-8bom",
}

func testFilename(name string) string {
	return filepath.Join("testdata", name+".txt")
}

func testReadFile(t *testing.T, name string) []byte {
	fn := testFilename(name)
	bs, err := fsu.ReadFile(fn)
	if err != nil {
		t.Fatalf("Failed to read file %q: %v", fn, err)
	}
	return bs
}

func TestDetectCharsetBytes(t *testing.T) {
	for i, c := range cs {
		bs := testReadFile(t, c)

		a, err := DetectCharsetBytes(bs)
		if err != nil {
			t.Fatalf("[%d] Failed to DetectCharsetBytes(%q): %v", i, c, err)
		}

		a = str.ReplaceAll(a, "_", "-")
		if !str.EqualFold(a, str.Remove(c, "bom")) {
			t.Errorf("[%d] DetectCharsetBytes(%q) = %v, want %v", i, c, a, c)
		}
	}
}

func TestDetectCharsetReader(t *testing.T) {
	for i, c := range cs {
		bs := testReadFile(t, c)

		_, a, err := DetectCharsetReader(bytes.NewReader(bs))
		if err != nil {
			t.Fatalf("[%d] Failed to DetectCharsetReader(%q): %v", i, c, err)
		}

		a = str.ReplaceAll(a, "_", "-")
		if !str.EqualFold(a, str.Remove(c, "bom")) {
			t.Errorf("[%d] DetectCharsetReader(%q) = %v, want %v", i, c, a, c)
		}
	}
}

func TestDetectCharsetFile(t *testing.T) {
	for i, c := range cs {
		a, err := DetectCharsetFile(testFilename(c))
		if err != nil {
			t.Fatalf("[%d] Failed to DetectCharsetFile(%q): %v", i, c, err)
		}

		a = str.ReplaceAll(a, "_", "-")
		if !str.EqualFold(a, str.Remove(c, "bom")) {
			t.Errorf("[%d] DetectCharsetFile(%q) = %v, want %v", i, c, a, c)
		}
	}
}
