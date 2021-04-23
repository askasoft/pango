package tpl

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
)

// HTMLTemplate html template engine
type HTMLTemplate struct {
	Extension   string      // template extension
	Funcs       FuncMap     // template functions
	Delims      Delims      // delimeters
	FileHandler FileHandler // file handler

	template *template.Template
}

// NewHTMLTemplate new template engine
func NewHTMLTemplate() *HTMLTemplate {
	return &HTMLTemplate{
		Extension:   ".html",
		Delims:      Delims{Left: "{{", Right: "}}"},
		FileHandler: DefaultFileHandler,
		template:    template.New(""),
	}
}

// Load glob and parse template files under root path
func (ht *HTMLTemplate) Load(root string) error {
	tpl := template.New("")
	tpl.Funcs(template.FuncMap(ht.Funcs))

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ht.Extension {
			return nil
		}

		text, err := ht.FileHandler(path)
		if err != nil {
			return fmt.Errorf("HTMLTemplate load template %q error: %v", path, err)
		}

		path = toTemplateName(root, path, ht.Extension)

		t2 := tpl.New(path)
		_, err = t2.Parse(text)
		if err != nil {
			return fmt.Errorf("HTMLTemplate parse template %q error: %v", path, err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	ht.template = tpl
	return nil
}

// Render render template with io.Writer
func (ht *HTMLTemplate) Render(w io.Writer, name string, data interface{}) error {
	err := ht.template.ExecuteTemplate(w, name, data)
	if err != nil {
		return fmt.Errorf("HTMLTemplate execute template %q error: %v", name, err)
	}

	return nil
}
