package tpl

import (
	"os"
	"strings"
	"testing"

	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/mag"
	"github.com/askasoft/pango/str"
)

var _ Templates = &TextTemplates{}

func testNewTextTpls() *TextTemplates {
	tt := NewTextTemplates()
	tt.Funcs(Functions())
	return tt
}

func testTextPage(t *testing.T, tt *TextTemplates, lang, page string, data map[string]any) {
	copy := make(map[string]any)
	mag.Copy(copy, data)
	data = copy

	fexp := "testdata/" + page + str.If(lang == "", "", "_"+lang) + ".txt.exp"
	fout := "testdata/" + page + str.If(lang == "", "", "_"+lang) + ".txt.out"

	os.Remove(fout)

	sb := &strings.Builder{}
	err := tt.Render(sb, lang, page, data)
	if err != nil {
		t.Errorf(`tt.Render(sb, %q, %q) = %v`, lang, page, err)
		return
	}

	out := sb.String()
	exp, _ := fsu.ReadString(fexp)

	if out != exp {
		fsu.WriteString(fout, out, 0666)
		t.Errorf("[%s] = \n%s\nwant:\n%s", page, out, exp)
	}
}

func testTextPages(t *testing.T, tt *TextTemplates, page string, data map[string]any) {
	for _, lang := range []string{"en", "ja", "ja-JP"} {
		testTextPage(t, tt, lang, page, data)
	}
}

func testTextLoad(t *testing.T, tt *TextTemplates) {
	testTextPage(t, tt, "", "index", map[string]any{
		"Title":   "Front Page",
		"Message": "Hello world!",
	})

	testTextPages(t, tt, "admin/admin", map[string]any{
		"Title":   "Admin Page",
		"Message": "Hello world!",
	})
}

func TestTextTemplatesLoad(t *testing.T) {
	roots := []string{"testdata", "./testdata"}

	for _, root := range roots {
		tt := testNewTextTpls()
		err := tt.Load(root)
		if err != nil {
			t.Fatalf(`ht.Load(%q) = %v`, root, err)
		}
		testTextLoad(t, tt)
	}
}

func TestTextTemplatesLoadFS(t *testing.T) {
	tt := testNewTextTpls()
	root := "testdata"

	err := tt.LoadFS(testdata, root)
	if err != nil {
		t.Errorf(`ht.LoadFS(%q) = %v`, root, err)
		return
	}

	testTextLoad(t, tt)
}
