package tpl

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/ref"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tpl/funcs"
)

// FuncMap is the type of the map defining the mapping from names to
// functions. Each function must have either a single return value, or two
// return values of which the second has type error. In that case, if the
// second (error) argument evaluates to non-nil during execution, execution
// terminates and Execute returns that error. FuncMap has the same base type
// as FuncMap in "text/template", copied here so clients need not import
// "text/template".
type FuncMap map[string]any

// Copy copy functions from other function maps `fms`
func (fm FuncMap) Copy(fms ...FuncMap) {
	for _, afm := range fms {
		for k, v := range afm {
			fm[k] = v
		}
	}
}

// Delims delims for template
type Delims struct {
	Left  string
	Right string
}

// Templates templates interface
type Templates interface {
	// Extensions sets template entensions.
	Extensions(extensions ...string)

	// Delims sets template left and right delims and returns a Engine instance.
	Delims(left, right string)

	// Funcs sets the FuncMap used for template.FuncMap.
	Funcs(funcMap FuncMap)

	// Load glob and parse template files under the root path
	Load(root string) error

	// LoadFS glob and parse template files from FS
	LoadFS(fsys fs.FS, root string) error

	// Render render template with io.Writer
	Render(w io.Writer, name string, data any) error
}

// readFile read file content to string
func readFile(fsys fs.FS, path string) (text string, err error) {
	var data []byte
	if fsys == nil {
		data, err = os.ReadFile(path)
	} else {
		data, err = fs.ReadFile(fsys, path)
	}

	if err != nil {
		return "", fmt.Errorf("Failed to read template %v, error: %w", path, err)
	}
	return string(data), nil
}

func toTemplateName(root, path, ext string) string {
	path = filepath.ToSlash(path)
	path = strings.TrimPrefix(path, root)
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, ext)
	return path
}

// Functions default utility functions for template
func Functions() FuncMap {
	return FuncMap{
		// unescape
		"HTML":     funcs.HTML,
		"HTMLAttr": funcs.HTMLAttr,
		"CSS":      funcs.CSS,
		"JS":       funcs.JS,
		"JSStr":    funcs.JSStr,
		"URL":      funcs.URL,
		"Srcset":   funcs.Srcset,

		// format
		"JSON":      funcs.JSON,
		"Comma":     funcs.Comma,
		"HumanSize": funcs.HumanSize,

		// math
		"Add":      funcs.Add,
		"Subtract": funcs.Subtract,
		"Multiply": funcs.Multiply,
		"Divide":   funcs.Divide,

		// string
		"StartsWith": str.StartsWith,
		"EndsWith":   str.EndsWith,
		"Contains":   str.Contains,
		"Split":      str.Split,
		"Fields":     str.Fields,
		"FieldsAny":  str.FieldsAny,

		// convert
		"Atoi": num.Atoi,
		"Atol": num.Atol,
		"Atof": num.Atof,

		// collection
		"Strings":  funcs.Strings,
		"Array":    funcs.Array,
		"ArrayGet": ref.ArrayGet,
		"ArraySet": ref.ArraySet,
		"Map":      funcs.Map,
		"MapGet":   ref.MapGet,
		"MapSet":   ref.MapSet,
	}
}
