package tpl

import (
	"embed"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/pandafw/pango/str"
)

func htmlTestLoad(t *testing.T, ht *HTMLTemplates) {
	sb := &strings.Builder{}

	ctx := map[string]any{
		"Title":   "Front Page",
		"Message": "Hello world!",
		"Time":    time.Now(),
	}
	err := ht.Render(sb, "index", ctx)
	if err != nil {
		t.Errorf(`ht.Render(sb, "index", ctx) = %v`, err)
		return
	}

	htm := sb.String()
	if !str.Contains(htm, fmt.Sprintf("<title>%s</title>", ctx["Title"])) {
		t.Errorf("Incorrect Title\n%s", htm)
	}
	if !str.Contains(htm, fmt.Sprintf("<p>%s</p>", ctx["Message"])) {
		t.Errorf("Incorrect Message\n%s", htm)
	}
	if !str.Contains(htm, fmt.Sprintf("<p>Time: %s</p>", ctx["Time"].(time.Time).Format("2006/1/2 15:04:05"))) {
		t.Errorf("Incorrect Message\n%s", htm)
	}

	sb.Reset()
	ctx = map[string]any{
		"Title":   "Admin Page",
		"Message": "Hello world!",
		"Time":    time.Now(),
	}
	err = ht.Render(sb, "admin/admin", ctx)
	if err != nil {
		t.Errorf(`ht.Render(sb, "admin/admin", ctx) = %v`, err)
		return
	}

	htm = sb.String()
	if !str.Contains(htm, fmt.Sprintf("<title>%s</title>", ctx["Title"])) {
		t.Errorf("Incorrect Title\n%s", htm)
	}
	if !str.Contains(htm, fmt.Sprintf("<p>Admin: %s</p>", ctx["Message"])) {
		t.Errorf("Incorrect Message\n%s", htm)
	}
	if !str.Contains(htm, fmt.Sprintf("<p>Time: %s</p>", ctx["Time"].(time.Time).Format("2006/1/2 15:04:05"))) {
		t.Errorf("Incorrect Message\n%s", htm)
	}
}

func TestLoadHTML(t *testing.T) {
	ht := NewHTMLTemplates()
	root := "testdata"

	err := ht.Load(root)
	if err != nil {
		t.Errorf(`ht.Load(%q) = %v`, root, err)
		return
	}

	htmlTestLoad(t, ht)
}

func TestLoadHTML2(t *testing.T) {
	ht := NewHTMLTemplates()
	root := "./testdata"

	err := ht.Load(root)
	if err != nil {
		t.Errorf(`ht.Load(%q) = %v`, root, err)
		return
	}

	htmlTestLoad(t, ht)
}

//go:embed testdata
var testdata embed.FS

func TestFSLoadHTML(t *testing.T) {
	ht := NewHTMLTemplates()
	root := "testdata"

	err := ht.LoadFS(testdata, root)
	if err != nil {
		t.Errorf(`ht.LoadFS(%q) = %v`, root, err)
		return
	}

	htmlTestLoad(t, ht)
}
