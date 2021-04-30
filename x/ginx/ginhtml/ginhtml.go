package ginhtml

import (
	"net/http"

	"github.com/gin-gonic/gin/render"
	"github.com/pandafw/pango/tpl"
)

// Engine html engine for gin
type Engine struct {
	*tpl.HTMLTemplate
}

// NewEngine create a html engine for gin
func NewEngine() *Engine {
	return &Engine{tpl.NewHTMLTemplate()}
}

// Render view render implement gin interface
type Render struct {
	ghe  *Engine
	name string
	data interface{}
}

// Instance implement gin interface
func (ghe *Engine) Instance(name string, data interface{}) render.Render {
	return Render{
		ghe:  ghe,
		name: name,
		data: data,
	}
}

// Render writes data with custom ContentType.
func (v Render) Render(w http.ResponseWriter) error {
	return v.ghe.Render(w, v.name, v.data)
}

// WriteContentType write html content type
func (v Render) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{"text/html"}
	}
}
