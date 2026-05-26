package cdt

import (
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

func TestDetector(t *testing.T) {
	type file_charset_language struct {
		File     string
		IsHtml   bool
		Charset  string
		Language string
	}
	var data = []file_charset_language{
		{"utf8.html", true, "UTF-8", ""},
		{"utf8_bom.html", true, "UTF-8", ""},
		{"8859_1_en.html", true, "ISO-8859-1", "en"},
		{"8859_1_da.html", true, "ISO-8859-1", "da"},
		{"8859_1_de.html", true, "ISO-8859-1", "de"},
		{"8859_1_es.html", true, "ISO-8859-1", "es"},
		{"8859_1_fr.html", true, "ISO-8859-1", "fr"},
		{"8859_1_pt.html", true, "windows-1252", "pt"},
		{"shift_jis.html", true, "Shift_JIS", "ja"},
		{"gb18030.html", true, "GB18030", "zh"},
		{"euc_jp.html", true, "EUC-JP", "ja"},
		{"euc_kr.html", true, "EUC-KR", "ko"},
		{"big5.html", true, "Big5", "zh"},
	}

	textDetector := NewTextDetector()
	htmlDetector := NewHtmlDetector()
	for i, d := range data {
		input := testReadFile(t, d.File)
		detector := textDetector
		if d.IsHtml {
			detector = htmlDetector
		}
		result, err := detector.DetectBest(input)
		if err != nil {
			t.Fatal(err)
		}
		if result.Charset != d.Charset {
			t.Errorf("[%d] %s: Expected charset %s, actual %s", i, d.File, d.Charset, result.Charset)
		}
		if result.Language != d.Language {
			t.Errorf("[%d] %s: Expected language %s, actual %s", i, d.File, d.Language, result.Language)
		}
	}
}
