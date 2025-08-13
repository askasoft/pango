package tpl

import (
	"embed"
	"os"
	"strings"
	"testing"

	"github.com/askasoft/pango/iox/fsu"
)

func testNewHtmlTpls() *HTMLTemplates {
	ht := NewHTMLTemplates()
	ht.Funcs(Functions())
	return ht
}

func testHtmlPage(t *testing.T, ht *HTMLTemplates, page string, ctx any) {
	fexp := "testdata/" + page + ".html.exp"
	fout := "testdata/" + page + ".html.out"

	os.Remove(fout)

	sb := &strings.Builder{}
	err := ht.Render(sb, page, ctx)
	if err != nil {
		t.Errorf(`ht.Render(sb, "index", ctx) = %v`, err)
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
	testHtmlPage(t, ht, "index", map[string]any{
		"Title":   "Front Page",
		"Message": "Hello world!",
	})

	testHtmlPage(t, ht, "admin/admin", map[string]any{
		"Title":   "Admin Page",
		"Message": "Hello world!",
	})
}

func TestLoadHTML(t *testing.T) {
	ht := testNewHtmlTpls()
	root := "testdata"

	err := ht.Load(root)
	if err != nil {
		t.Errorf(`ht.Load(%q) = %v`, root, err)
		return
	}

	testHtmlLoad(t, ht)
}

func TestLoadHTML2(t *testing.T) {
	ht := testNewHtmlTpls()
	root := "./testdata"

	err := ht.Load(root)
	if err != nil {
		t.Errorf(`ht.Load(%q) = %v`, root, err)
		return
	}

	testHtmlLoad(t, ht)
}

//go:embed testdata
var testdata embed.FS

func TestFSLoadHTML(t *testing.T) {
	ht := testNewHtmlTpls()
	root := "testdata"

	err := ht.LoadFS(testdata, root)
	if err != nil {
		t.Errorf(`ht.LoadFS(%q) = %v`, root, err)
		return
	}

	testHtmlLoad(t, ht)
}
