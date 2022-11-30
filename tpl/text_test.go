package tpl

import (
	"strings"
	"testing"

	"github.com/pandafw/pango/fsu"
)

func testNewTextTpls() *TextTemplates {
	tt := NewTextTemplates()
	tt.Funcs(Functions())
	return tt
}

func testTextPage(t *testing.T, tt *TextTemplates, page string, ctx any) {
	fexp := "testdata/" + page + ".txt.exp"
	fout := "testdata/" + page + ".txt.out"

	sb := &strings.Builder{}
	err := tt.Render(sb, page, ctx)
	if err != nil {
		t.Errorf(`tt.Render(sb, "index", ctx) = %v`, err)
		return
	}

	out := sb.String()
	exp, _ := fsu.ReadString(fexp)

	if out != exp {
		fsu.WriteString(fout, out, fsu.FileMode(0666))
		t.Errorf("[%s] = %q, want %q", page, out, exp)
	}
}

func testTextLoad(t *testing.T, tt *TextTemplates) {
	testTextPage(t, tt, "index", map[string]any{
		"Title":   "Front Page",
		"Message": "Hello world!",
	})

	testTextPage(t, tt, "admin/admin", map[string]any{
		"Title":   "Admin Page",
		"Message": "Hello world!",
	})
}

func TestLoadText(t *testing.T) {
	tt := testNewTextTpls()
	root := "testdata"

	err := tt.Load(root)
	if err != nil {
		t.Errorf(`ht.Load(%q) = %v`, root, err)
		return
	}

	testTextLoad(t, tt)
}

func TestLoadText2(t *testing.T) {
	tt := testNewTextTpls()
	root := "./testdata"

	err := tt.Load(root)
	if err != nil {
		t.Errorf(`ht.Load(%q) = %v`, root, err)
		return
	}

	testTextLoad(t, tt)
}

func TestFSLoadText(t *testing.T) {
	tt := testNewTextTpls()
	root := "testdata"

	err := tt.LoadFS(testdata, root)
	if err != nil {
		t.Errorf(`ht.LoadFS(%q) = %v`, root, err)
		return
	}

	testTextLoad(t, tt)
}
