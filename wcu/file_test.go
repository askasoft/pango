package wcu

import (
	"testing"

	"github.com/askasoft/pango/str"
)

func TestDetectCharsetFile(t *testing.T) {
	for i, c := range cs {
		a, err := DetectCharsetFile(testFilename(c), 1024)
		if err != nil {
			t.Fatalf("[%d] Failed to DetectCharsetFile(%q): %v", i, c, err)
		}

		a = str.ReplaceAll(a, "_", "-")
		if !str.EqualFold(a, str.Remove(c, "bom")) {
			t.Errorf("[%d] DetectCharsetFile(%q) = %v, want %v", i, c, a, c)
		}
	}
}

func TestReadFile(t *testing.T) {
	ws := string(testReadFile(t, "utf-8"))

	for i, c := range cs {
		tf := testFilename(c)
		bs, err := ReadFile(tf, str.Remove(c, "bom"))
		if err != nil {
			t.Fatalf("[%d] Failed to ReadFile(%q): %v", i, tf, err)
		}
		as := str.SkipBOM(string(bs))

		if as != ws {
			t.Errorf("[%d] ReadFile(%q) = %v, want %v", i, c, as, ws)
		}
	}
}

func TestDetectAndReadFile(t *testing.T) {
	ws := string(testReadFile(t, "utf-8"))

	for i, c := range cs {
		tf := testFilename(c)
		bs, enc, err := DetectAndReadFile(tf, 1024)
		if err != nil {
			t.Fatalf("[%d] Failed to DetectAndReadFile(%q): %v", i, tf, err)
		}

		// fmt.Printf("[%d] DetectAndReadFile(%q) = %v, want %v\n", i, c, enc, c)

		enc = str.ReplaceAll(enc, "_", "-")
		if !str.EqualFold(enc, str.Remove(c, "bom")) {
			t.Errorf("[%d] DetectAndReadFile(%q) = %v, want %v", i, c, enc, c)
		}

		as := str.SkipBOM(string(bs))
		if as != ws {
			t.Errorf("[%d] DetectAndReadFile(%q) = %q\nREAD: %v\nWANT: %v\n", i, c, enc, as, ws)
		}
	}
}
