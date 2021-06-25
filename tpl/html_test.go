package tpl

import (
	"embed"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func htmlTestLoad(t *testing.T, ht *HTMLTemplate) {
	sb := &strings.Builder{}

	ctx := map[string]interface{}{
		"Title":   "Front Page",
		"Message": "Hello world!",
		"Time":    time.Now(),
	}
	assert.Nil(t, ht.Render(sb, "index", ctx))
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println(sb.String())

	sb.Reset()
	ctx = map[string]interface{}{
		"Title":   "Admin Page",
		"Message": "Hello world!",
		"Time":    time.Now(),
	}
	assert.Nil(t, ht.Render(sb, "admin/admin", ctx))
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println(sb.String())
}

func TestLoadHTML(t *testing.T) {
	ht := NewHTMLTemplate()
	root := "testdata"

	assert.Nil(t, ht.Load(root))

	htmlTestLoad(t, ht)
}

//go:embed testdata
var testdata embed.FS

func TestFSLoadHTML(t *testing.T) {
	ht := NewHTMLTemplate()
	root := "testdata"

	assert.Nil(t, ht.LoadFS(testdata, root))

	htmlTestLoad(t, ht)
}
