package tpl

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pandafw/pango/ars"
)

// TextTemplates text template engine
type TextTemplates struct {
	extensions []string // template extensions
	funcs      FuncMap  // template functions
	delims     Delims   // delimeters

	template *template.Template
}

// NewTextTemplates new templates instance
func NewTextTemplates(extensions ...string) *TextTemplates {
	tt := &TextTemplates{
		delims: Delims{Left: "{{", Right: "}}"},
	}
	tt.Extensions(extensions...)
	return tt
}

// Extensions sets template entensions.
func (tt *TextTemplates) Extensions(extensions ...string) {
	if len(extensions) == 0 {
		extensions = []string{".tpl", ".txt", ".gotxt"}
	}
	tt.extensions = extensions
}

// Delims sets template left and right delims and returns a Engine instance.
func (tt *TextTemplates) Delims(left, right string) {
	tt.delims = Delims{Left: left, Right: right}
}

// Funcs sets the FuncMap used for template.FuncMap.
func (tt *TextTemplates) Funcs(funcMap FuncMap) {
	tt.funcs = funcMap
}

func (tt *TextTemplates) init() {
	if tt.template == nil {
		tpl := template.New("")
		tpl.Delims(tt.delims.Left, tt.delims.Right)
		tpl.Funcs(template.FuncMap(tt.funcs))
		tt.template = tpl
	}
}

// Load glob and parse template files under the root path
func (tt *TextTemplates) Load(root string) (err error) {
	tt.init()

	root, err = filepath.Abs(root)
	if err != nil {
		return
	}

	root = filepath.ToSlash(root)
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return tt.loadFile(nil, root, path)
	})
	if err != nil {
		return
	}

	return
}

// LoadFS glob and parse template files from FS
func (tt *TextTemplates) LoadFS(fsys fs.FS, root string) (err error) {
	tt.init()

	err = fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		return tt.loadFile(fsys, root, path)
	})
	if err != nil {
		return
	}

	return
}

// loadFile load template file
func (tt *TextTemplates) loadFile(fsys fs.FS, root, path string) error {
	ext := filepath.Ext(path)
	if !ars.ContainsString(tt.extensions, ext) {
		return nil
	}

	text, err := readFile(fsys, path)
	if err != nil {
		return fmt.Errorf("TextTemplates load template %q error: %w", path, err)
	}

	path = toTemplateName(root, path, ext)

	tpl := tt.template.New(path)
	_, err = tpl.Parse(text)
	if err != nil {
		return fmt.Errorf("TextTemplates parse template %q error: %w", path, err)
	}
	return nil
}

// Render render template with io.Writer
func (tt *TextTemplates) Render(w io.Writer, name string, data any) error {
	err := tt.template.ExecuteTemplate(w, name, data)
	if err != nil {
		return fmt.Errorf("TextTemplates execute template %q error: %w", name, err)
	}

	return nil
}
