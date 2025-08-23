package render

import (
	"net/http/httptest"
	"testing"

	"github.com/askasoft/pango/test/assert"
	"github.com/askasoft/pango/tpl"
)

func TestRenderHTMLTemplate(t *testing.T) {
	ht := tpl.NewHTMLTemplates()
	ht.Load("testdata")

	hr := NewHTMLRenderer(ht)
	for _, loc := range []string{"", "zh"} {
		t.Run("loc="+loc, func(t *testing.T) {
			w := httptest.NewRecorder()

			r := hr(loc, "hello", map[string]any{
				"name": "alexandernyquist",
			})

			err := r.Render(w)
			assert.NoError(t, err)
			assert.Equal(t, "<h1>Hello alexandernyquist</h1>", w.Body.String())
			assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
		})
	}

}

func TestRenderHTMLTemplateJA(t *testing.T) {
	ht := tpl.NewHTMLTemplates()
	ht.Load("testdata")

	hr := NewHTMLRenderer(ht)
	r := hr("ja", "hello", map[string]any{
		"name": "alexandernyquist",
	})

	w := httptest.NewRecorder()
	err := r.Render(w)
	assert.NoError(t, err)
	assert.Equal(t, "<h1>こんにちは、alexandernyquistさん！</h1>", w.Body.String())
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
}
