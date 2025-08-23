package tpl

import (
	"embed"
	"os"
	"strings"
	"testing"

	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/str"
)

var _ Templates = &HTMLTemplates{}

func testNewHtmlTpls() *HTMLTemplates {
	ht := NewHTMLTemplates()
	ht.Funcs(Functions())
	return ht
}

func testHtmlPages(t *testing.T, ht *HTMLTemplates, page string, data map[string]any) {
	for _, lang := range []string{"en", "ja", "ja-JP"} {
		testHtmlPage(t, ht, lang, page, data)
	}
}

func testHtmlPage(t *testing.T, ht *HTMLTemplates, lang, page string, data map[string]any) {
	copy := make(map[string]any)
	for k, v := range data {
		copy[k] = v
	}
	data = copy

	fexp := "testdata/" + page + str.If(lang == "", "", "_"+lang) + ".html.exp"
	fout := "testdata/" + page + str.If(lang == "", "", "_"+lang) + ".html.out"

	os.Remove(fout)

	sb := &strings.Builder{}
	err := ht.Render(sb, lang, page, data)
	if err != nil {
		t.Errorf(`ht.Render(sb, %q, %q) = %v`, lang, page, err)
		return
	}

	out := sb.String()
	exp, _ := fsu.ReadString(fexp)

	if out != exp {
		fsu.WriteString(fout, out, fsu.FileMode(0666))
		t.Errorf("[%s] = %q, want %q", page, out, exp)
	}
}

func testHtmlLoad(t *testing.T, ht *HTMLTemplates) {
	testHtmlPage(t, ht, "", "index", map[string]any{
		"Title":   "Front Page",
		"Message": "Hello world!",
	})

	testHtmlPages(t, ht, "admin/admin", map[string]any{
		"Title":   "Admin Page",
		"Message": "Hello world!",
	})
}

func TestHTMLTemplatesLoad(t *testing.T) {
	roots := []string{"testdata", "./testdata"}

	for _, root := range roots {
		ht := testNewHtmlTpls()
		err := ht.Load(root)
		if err != nil {
			t.Fatalf(`ht.Load(%q) = %v`, root, err)
		}
		testHtmlLoad(t, ht)
	}
}

//go:embed testdata
var testdata embed.FS

func TestHTMLTemplatesLoadFS(t *testing.T) {
	ht := testNewHtmlTpls()
	root := "testdata"

	err := ht.LoadFS(testdata, root)
	if err != nil {
		t.Errorf(`ht.LoadFS(%q) = %v`, root, err)
		return
	}

	testHtmlLoad(t, ht)
}
