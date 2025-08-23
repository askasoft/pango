package render

import (
	"net/http"

	"github.com/askasoft/pango/tpl"
)

// HTMLRender interface is to be implemented by HTMLProduction and HTMLDebug.
type HTMLRender interface {
	// Instance returns an HTML instance.
	Instance(locale, name string, data any) Render
}

// HTMLTemplates html templates interface for xin
type HTMLTemplates interface {
	tpl.Templates

	HTMLRender
}

// htmlTemplates html engine for xin
type htmlTemplates struct {
	*tpl.HTMLTemplates
}

// NewHTMLTemplates create a html templates instance for xin
func NewHTMLTemplates() HTMLTemplates {
	return &htmlTemplates{tpl.NewHTMLTemplates()}
}

// Instance implement xin interface
func (ht *htmlTemplates) Instance(locale, name string, data any) Render {
	return htmlRender{
		htpls:  ht,
		locale: locale,
		name:   name,
		data:   data,
	}
}

// HTMLRender view render implement xin interface
type htmlRender struct {
	htpls  *htmlTemplates
	locale string
	name   string
	data   any
}

// Render writes data with custom ContentType.
func (r htmlRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	return r.htpls.Render(w, r.locale, r.name, r.data)
}

// WriteContentType write html content type
func (r htmlRender) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, "text/html; charset=utf-8")
}
