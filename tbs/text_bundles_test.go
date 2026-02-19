package tbs

import (
	"embed"
	"os"
	"strings"
	"testing"

	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/iox"
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

	testReplace(t, func(locale, format string, args ...any) string {
		return tbs.Replace(locale, format, args...)
	})

	testEvaluate(t, func(locale, format string, data any) string {
		return tbs.Evaluate(locale, format, data)
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

	testReplace(t, func(locale, format string, args ...any) string {
		return tbs.Replace(locale, format, args...)
	})

	testEvaluate(t, func(locale, format string, data any) string {
		return tbs.Evaluate(locale, format, data)
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
		{"en", "format.welcome", []any{"home"}, "welcome home"},
		{"en", "format.new.hello", []any{"home"}, "hello home"},
		{"ja-JP", "title", nil, "こんにちは世界"},
		{"ja-JP", "format.welcome", []any{"ダーリン"}, "ようこそ ダーリン"},
		{"ja-JP", "format.new.hello", []any{"ダーリン"}, "ハロー ダーリン"},
	}

	for i, c := range cs {
		a := fmt(c.lang, c.name, c.args...)
		if a != c.want {
			t.Errorf("%d Foramt(%q, %q, %v) = %q, want %q", i, c.lang, c.name, c.args, a, c.want)
		}
	}
}

func testReplace(t *testing.T, rep func(locale, format string, args ...any) string) {
	cs := []struct {
		lang string
		name string
		args []any
		want string
	}{
		{"en", "title", nil, "hello world"},
		{"en", "replace.welcome", []any{"{name}", "home"}, "welcome home"},
		{"en", "replace.new.hello", []any{"{name}", "home"}, "hello home"},
		{"ja-JP", "title2", []any{"こんにちは世界2"}, "こんにちは世界2"},
		{"ja-JP", "replace.welcome", []any{"{name}", "ダーリン"}, "ようこそ ダーリン"},
		{"ja-JP", "replace.new.hello", []any{"{name}", "ダーリン"}, "ハロー ダーリン"},
	}

	for i, c := range cs {
		a := rep(c.lang, c.name, c.args...)
		if a != c.want {
			t.Errorf("%d Replace(%q, %q, %v) = %q, want %q", i, c.lang, c.name, c.args, a, c.want)
		}
	}
}

func testEvaluate(t *testing.T, tra func(locale, format string, data any) string) {
	cs := []struct {
		lang string
		name string
		data map[string]any
		want string
	}{
		{"en", "evaluate.welcome", map[string]any{"some": "home"}, "welcome home"},
		{"en", "evaluate.new.hello", map[string]any{"some": "home"}, "hello home"},
		{"ja-JP", "evaluate.welcome", map[string]any{"some": "ダーリン"}, "ようこそ ダーリン"},
		{"ja-JP", "evaluate.new.hello", map[string]any{"some": "ダーリン"}, "ハロー ダーリン"},
	}

	for i, c := range cs {
		a := tra(c.lang, c.name, c.data)
		if a != c.want {
			t.Errorf("%d Evaluate(%q, %q, %v) = %q, want %q", i, c.lang, c.name, c.data, a, c.want)
		}
	}
}

func testGetBundle(t *testing.T, tbs *TextBundles) {
	fexp := "testdata/bundles.exp"
	fout := "testdata/bundles.out"

	os.Remove(fout)

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
