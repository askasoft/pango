package render

import (
	"net/http"

	"github.com/pandafw/pango/tpl"
)

// HTMLRender interface is to be implemented by HTMLProduction and HTMLDebug.
type HTMLRender interface {
	// Instance returns an HTML instance.
	Instance(string, interface{}) Render
}

// HTMLTemplate html engine for gin
type HTMLTemplate struct {
	*tpl.HTMLTemplate
}

// NewHTMLTemplate create a html engine for gin
func NewHTMLTemplate() *HTMLTemplate {
	return &HTMLTemplate{tpl.NewHTMLTemplate()}
}

// Instance implement gin interface
func (html *HTMLTemplate) Instance(name string, data interface{}) Render {
	return htmlRender{
		html: html,
		name: name,
		data: data,
	}
}

// HTMLRender view render implement gin interface
type htmlRender struct {
	html *HTMLTemplate
	name string
	data interface{}
}

// Render writes data with custom ContentType.
func (r htmlRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	return r.html.Render(w, r.name, r.data)
}

// WriteContentType write html content type
func (r htmlRender) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, []string{"text/html; charset=utf-8"})
}
