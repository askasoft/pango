package tpl

import (
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/str"
)

var HTMLTemplateExtensions = []string{".tpl", ".html", ".gohtml"}

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
		extensions = HTMLTemplateExtensions
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

func (ht *HTMLTemplates) init() {
	if ht.template == nil {
		tpl := template.New("")
		tpl.Delims(ht.delims.Left, ht.delims.Right)
		tpl.Funcs(template.FuncMap(ht.funcs))
		ht.template = tpl
	}
}

// Load glob and parse template files under the root path
func (ht *HTMLTemplates) Load(root string) (err error) {
	ht.init()

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

		return ht.loadFile(nil, root, path)
	})

	return
}

// LoadFS glob and parse template files from FS
func (ht *HTMLTemplates) LoadFS(fsys fs.FS, root string) (err error) {
	ht.init()

	err = fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		return ht.loadFile(fsys, root, path)
	})

	return
}

// loadFile load template file
func (ht *HTMLTemplates) loadFile(fsys fs.FS, root, path string) error {
	ext := filepath.Ext(path)
	if !asg.Contains(ht.extensions, ext) {
		return nil
	}

	text, err := readFile(fsys, path)
	if err != nil {
		return err
	}

	path = toTemplateName(root, path, ext)

	tpl := ht.template.New(path)
	_, err = tpl.Parse(text)
	return err
}

// Render render template with io.Writer
// If locale is not empty, it will try to load locale template first.
// For example, if locale is "zh-TW" and name is "hello", it will try to load template via the following order.
// 1. "hello.zh-TW.tpl"
// 2. "hello.zh.tpl"
// 3. "hello.tpl"
func (ht *HTMLTemplates) Render(w io.Writer, locale, name string, data any) error {
	for locale != "" {
		if t := ht.template.Lookup(name + "_" + locale); t != nil {
			return t.Execute(w, data)
		}

		if l2, _, ok := str.LastCutByte(locale, '-'); ok {
			locale = l2
		} else {
			break
		}
	}

	return ht.template.ExecuteTemplate(w, name, data)
}
