package htmlx

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/askasoft/pango/fsu"
	"golang.org/x/net/html"
)

func testFilename(name string) string {
	return path.Join("testdata", name)
}

func testReadFile(t *testing.T, name string) []byte {
	fn := testFilename(name)
	bs, err := fsu.ReadFile(fn)
	if err != nil {
		t.Fatalf("Failed to read file %q: %v", fn, err)
	}
	return bs
}

func TestParseHTMLFile(t *testing.T) {
	fn := testFilename("utf-8.html")
	doc, err := ParseHTMLFile(fn)
	if err != nil {
		t.Fatalf("Failed to ParseHTMLFile(%q): %v", fn, err)
	}

	var f func(*html.Node, string)
	f = func(n *html.Node, p string) {
		switch n.Type {
		case html.DocumentNode:
			fmt.Printf("\n%s<doc>", p)
		case html.ElementNode:
			fmt.Printf("\n%s<%s>", p, n.Data)
		case html.CommentNode:
			fmt.Printf("\n%s<!--%s-->", p, n.Data)
		case html.TextNode:
			fmt.Printf("%q", n.Data)
		}

		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c, p+"  ")
			}
		}
	}

	f(doc, "")
}

func TestExtractTextFromHTMLFile(t *testing.T) {
	cs := []string{"utf-8.html", "shift-jis.html"}

	w := string(testReadFile(t, "expect.txt"))

	for i, c := range cs {
		fn := testFilename(c)
		a, err := ExtractTextFromHTMLFile(fn)
		if err != nil {
			t.Fatalf("[%d] Failed to ExtractTextFromHTMLFile(%q): %v", i, c, err)
		}

		if w != a {
			t.Errorf("[%d] ExtractTextFromHTMLFile(%q):\nACTUAL: %q\n  WANT: %q\n", i, c, a, w)
			fsu.WriteString(fn+".out", a, fsu.FileMode(0660))
		} else {
			os.Remove(fn + ".out")
		}
	}
}

func TestExtractTitleFromHTMLNode(t *testing.T) {
	cs := []string{"utf-8.html", "shift-jis.html"}

	w := "タイトル"

	for i, c := range cs {
		fn := testFilename(c)
		doc, err := ParseHTMLFile(fn)
		if err != nil {
			t.Fatalf("[%d] Failed to ParseHTMLFile(%q): %v", i, c, err)
		}

		a := ExtractTitleFromHTMLNode(doc)
		if w != a {
			t.Errorf("[%d] ExtractTitleFromHTMLNode(%q):\nACTUAL: %q\n  WANT: %q\n", i, c, a, w)
		}
	}
}

func TestExtractMetasFromHTMLNode(t *testing.T) {
	cs := []string{"utf-8.html", "shift-jis.html"}

	w := map[string]string{
		"keyword":     "metaキーワード",
		"description": "meta説明",
	}

	for i, c := range cs {
		fn := testFilename(c)
		doc, err := ParseHTMLFile(fn)
		if err != nil {
			t.Fatalf("[%d] Failed to ParseHTMLFile(%q): %v", i, c, err)
		}

		a := ExtractMetasFromHTMLNode(doc)
		if !reflect.DeepEqual(w, a) {
			t.Errorf("[%d] ExtractMetasFromHTMLNode(%q):\nACTUAL: %v\n  WANT: %v\n", i, c, a, w)
		}
	}
}
