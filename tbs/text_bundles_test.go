package tbs

import (
	"embed"
	"strings"
	"testing"

	"github.com/pandafw/pango/fsu"
	"github.com/pandafw/pango/iox"
)

//go:embed testdata
var testdata embed.FS

var testroot = "testdata"

func TestNewLoad(t *testing.T) {
	tbs := NewTextBundles()

	err := tbs.Load(testroot)
	if err != nil {
		t.Errorf(`tbs.Load(%q) = %v`, testroot, err)
		return
	}

	testFormat(t, func(locale, format string, args ...any) string {
		return tbs.Format(locale, format, args...)
	})

	testGetBundle(t, tbs)
}

func TestNewLoadFS(t *testing.T) {
	tbs := NewTextBundles()

	err := tbs.LoadFS(testdata, testroot)
	if err != nil {
		t.Errorf(`tbs.LoadFS(%q) = %v`, testroot, err)
		return
	}

	testFormat(t, func(locale, format string, args ...any) string {
		return tbs.Format(locale, format, args...)
	})

	testGetBundle(t, tbs)
}

func testFormat(t *testing.T, fmt func(locale, format string, args ...any) string) {
	cs := []struct {
		lang string
		name string
		args []any
		want string
	}{
		{"en", "title", nil, "hello world"},
		{"en", "label.welcome", []any{"home"}, "welcome home"},
		{"en", "label.new.hello", []any{"home"}, "hello home"},
		{"ja-JP", "title", nil, "こんにちは世界"},
		{"ja-JP", "label.welcome", []any{"ダーリン"}, "ようこそ ダーリン"},
		{"ja-JP", "label.new.hello", []any{"ダーリン"}, "ハロー ダーリン"},
	}

	for i, c := range cs {
		a := fmt(c.lang, c.name, c.args...)
		if a != c.want {
			t.Errorf("%d Foramt(%q, %q, %v) = %q, want %q", i, c.lang, c.name, c.args, a, c.want)
		}
	}
}

func testGetBundle(t *testing.T, tbs *TextBundles) {
	fexp := "testdata/bundles.ini"
	fout := "testdata/bundles.out"

	b := tbs.GetBundle("ja-JP")
	b.EOL = iox.CRLF

	sout := &strings.Builder{}
	if err := b.WriteData(sout); err != nil {
		t.Fatalf("WriteData(): %v", err)
	}

	sexp, _ := fsu.ReadString(fexp)
	if sexp != sout.String() {
		b.WriteFile(fout)
		t.Fatalf("GetBundle()\n actual: %v\n   want: %v", sout.String(), sexp)
	}
}
