package wcu

import (
	"testing"

	"github.com/askasoft/pango/str"
)

func TestDecodeBytes(t *testing.T) {
	ws := string(testReadFile(t, "utf-8"))

	for i, c := range cs {
		bs := testReadFile(t, c)

		abs, err := DecodeBytes(bs, str.Remove(c, "bom"))
		if err != nil {
			t.Fatalf("[%d] Failed to DecodeBytes(%q): %v", i, c, err)
		}
		as := str.SkipBOM(string(abs))

		if as != ws {
			t.Errorf("[%d] DecodeBytes(%q) = %v, want %v", i, c, as, ws)
		}
	}
}

func TestDetectAndDecodeBytes(t *testing.T) {
	ws := string(testReadFile(t, "utf-8"))

	for i, c := range cs {
		bs := testReadFile(t, c)

		abs, err := DetectAndDecodeBytes(bs)
		if err != nil {
			t.Fatalf("[%d] Failed to DetectAndDecodeBytes(%q): %v", i, c, err)
		}
		as := str.SkipBOM(string(abs))

		if as != ws {
			t.Errorf("[%d] DecodeBytes(%q) = %v, want %v", i, c, as, ws)
		}
	}
}
