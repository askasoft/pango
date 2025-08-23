package render

import (
	"net/http"

	"github.com/askasoft/pango/tpl"
)

// HTMLRenderer is function to create a Render instance.
type HTMLRenderer func(locale, name string, data any) Render

// NewHTMLRenderer creates a HTMLRenderer function.
func NewHTMLRenderer(tpls tpl.Templates) HTMLRenderer {
	return func(locale, name string, data any) Render {
		return HTMLRender{
			Templates: tpls,
			Locale:    locale,
			Name:      name,
			Data:      data,
		}
	}
}

// HTMLRender is the HTML render implementation.
type HTMLRender struct {
	Templates tpl.Templates
	Locale    string
	Name      string
	Data      any
}

// Render renders data with html template.
func (r HTMLRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	return r.Templates.Render(w, r.Locale, r.Name, r.Data)
}

// WriteContentType write html content type
func (r HTMLRender) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, "text/html; charset=utf-8")
}
