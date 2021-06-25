package tpl

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
)

// TextTemplate text template engine
type TextTemplate struct {
	Extension string  // template extension
	Funcs     FuncMap // template functions
	Delims    Delims  // delimeters

	template *template.Template
}

// NewTextTemplate new template engine
func NewTextTemplate() *TextTemplate {
	return &TextTemplate{
		Extension: ".txt",
		Delims:    Delims{Left: "{{", Right: "}}"},
	}
}

// Load glob and parse template files under the root path
func (tt *TextTemplate) Load(root string) error {
	tpl := template.New("")
	tpl.Funcs(template.FuncMap(tt.Funcs))

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		return tt.loadFile(tpl, nil, root, path)
	})

	if err != nil {
		return err
	}

	tt.template = tpl
	return nil
}

// LoadFS glob and parse template files from FS
func (tt *TextTemplate) LoadFS(fsys fs.FS, root string) error {
	tpl := template.New("")
	tpl.Funcs(template.FuncMap(tt.Funcs))

	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		return tt.loadFile(tpl, fsys, root, path)
	})

	if err != nil {
		return err
	}

	tt.template = tpl
	return nil
}

// loadFile load template file
func (tt *TextTemplate) loadFile(tpl *template.Template, fsys fs.FS, root, path string) error {
	if filepath.Ext(path) != tt.Extension {
		return nil
	}

	text, err := readFile(fsys, path)
	if err != nil {
		return fmt.Errorf("TextTemplate load template %q error: %v", path, err)
	}

	path = toTemplateName(root, path, tt.Extension)

	tpl = tpl.New(path)
	_, err = tpl.Parse(text)
	if err != nil {
		return fmt.Errorf("TextTemplate parse template %q error: %v", path, err)
	}
	return nil
}

// Render render template with io.Writer
func (tt *TextTemplate) Render(w io.Writer, name string, data interface{}) error {
	err := tt.template.ExecuteTemplate(w, name, data)
	if err != nil {
		return fmt.Errorf("TextTemplate execute template %q error: %v", name, err)
	}

	return nil
}
