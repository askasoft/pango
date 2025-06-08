package render

import (
	"net/http/httptest"
	"testing"

	"github.com/askasoft/pango/test/assert"
)

func TestRenderHTMLTemplate(t *testing.T) {
	w := httptest.NewRecorder()

	ht := NewHTMLTemplates()
	ht.Load("testdata")
	instance := ht.Instance("hello", map[string]any{
		"name": "alexandernyquist",
	})

	err := instance.Render(w)

	assert.NoError(t, err)
	assert.Equal(t, "<h1>Hello alexandernyquist</h1>", w.Body.String())
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
}
