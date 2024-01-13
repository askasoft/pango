package wcu

import (
	"testing"

	"github.com/askasoft/pango/str"
)

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
		bs, err := DetectAndReadFile(tf)
		if err != nil {
			t.Fatalf("[%d] Failed to DetectAndReadFile(%q): %v", i, tf, err)
		}
		as := str.SkipBOM(string(bs))

		if as != ws {
			t.Errorf("[%d] DetectAndReadFile(%q) = %v, want %v", i, c, as, ws)
		}
	}
}
