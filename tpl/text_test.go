package tpl

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadText(t *testing.T) {
	tt := NewTextTemplate()
	root := "testdata"

	assert.Nil(t, tt.Load(root))

	sb := &strings.Builder{}

	ctx := map[string]interface{}{
		"Title":   "Front Page",
		"Message": "Hello world!",
		"Time":    time.Now(),
	}
	assert.Nil(t, tt.Render(sb, "index", ctx))
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println(sb.String())

	sb.Reset()
	ctx = map[string]interface{}{
		"Title":   "Admin Page",
		"Message": "Hello world!",
		"Time":    time.Now(),
	}
	assert.Nil(t, tt.Render(sb, "admin/admin", ctx))
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println(sb.String())
}
