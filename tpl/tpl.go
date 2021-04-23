package tpl

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// FuncMap is the type of the map defining the mapping from names to
// functions. Each function must have either a single return value, or two
// return values of which the second has type error. In that case, if the
// second (error) argument evaluates to non-nil during execution, execution
// terminates and Execute returns that error. FuncMap has the same base type
// as FuncMap in "text/template", copied here so clients need not import
// "text/template".
type FuncMap map[string]interface{}

// Delims delims for template
type Delims struct {
	Left  string
	Right string
}

// FileHandler file handler function to load file content
type FileHandler func(path string) (text string, err error)

// DefaultFileHandler new default file handler
func DefaultFileHandler(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Failed to read template %v, error: %v", path, err)
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
