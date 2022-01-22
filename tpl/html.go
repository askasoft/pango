package tpl

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pandafw/pango/ars"
)

// HTMLTemplates html template engine
type HTMLTemplates struct {
	extensions []string // template extensions
	funcs      FuncMap  // template functions
	delims     Delims   // delimeters

	template *template.Template
}

// NewHTMLTemplates new template engine
func NewHTMLTemplates(extensions ...string) *HTMLTemplates {
	ht := &HTMLTemplates{
		delims: Delims{Left: "{{", Right: "}}"},
	}

	ht.Extensions(extensions...)
	return ht
}

// Extensions sets template entensions.
func (ht *HTMLTemplates) Extensions(extensions ...string) {
	if len(extensions) == 0 {
		extensions = []string{".html", ".gohtml"}
	}
	ht.extensions = extensions
}

// Delims sets template left and right delims and returns a Engine instance.
func (ht *HTMLTemplates) Delims(left, right string) {
	ht.delims = Delims{Left: left, Right: right}
}

// Funcs sets the FuncMap used for template.FuncMap.
func (ht *HTMLTemplates) Funcs(funcMap FuncMap) {
	ht.funcs = funcMap
}

// Load glob and parse template files under the root path
func (ht *HTMLTemplates) Load(root string) (err error) {
	tpl := template.New("")

	tpl.Delims(ht.delims.Left, ht.delims.Right)
	tpl.Funcs(template.FuncMap(ht.funcs))

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

		return ht.loadFile(tpl, nil, root, path)
	})
	if err != nil {
		return
	}

	ht.template = tpl
	return
}

// LoadFS glob and parse template files from FS
func (ht *HTMLTemplates) LoadFS(fsys fs.FS, root string) (err error) {
	tpl := template.New("")

	tpl.Delims(ht.delims.Left, ht.delims.Right)
	tpl.Funcs(template.FuncMap(ht.funcs))

	err = fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		return ht.loadFile(tpl, fsys, root, path)
	})
	if err != nil {
		return
	}

	ht.template = tpl
	return
}

// loadFile load template file
func (ht *HTMLTemplates) loadFile(tpl *template.Template, fsys fs.FS, root, path string) error {
	ext := filepath.Ext(path)
	if !ars.ContainsString(ht.extensions, ext) {
		return nil
	}

	text, err := readFile(fsys, path)
	if err != nil {
		return fmt.Errorf("HTMLTemplates load template %q error: %v", path, err)
	}

	path = toTemplateName(root, path, ext)

	tpl = tpl.New(path)
	_, err = tpl.Parse(text)
	if err != nil {
		return fmt.Errorf("HTMLTemplates parse template %q error: %v", path, err)
	}
	return nil
}

// Render render template with io.Writer
func (ht *HTMLTemplates) Render(w io.Writer, name string, data interface{}) error {
	err := ht.template.ExecuteTemplate(w, name, data)
	if err != nil {
		return fmt.Errorf("HTMLTemplates execute template %q error: %v", name, err)
	}

	return nil
}
