package tpl

import (
	"fmt"
	"io/fs"
	"os"
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

// readFile read file content to string
func readFile(fsys fs.FS, path string) (text string, err error) {
	var data []byte
	if fsys == nil {
		data, err = os.ReadFile(path)
	} else {
		data, err = fs.ReadFile(fsys, path)
	}

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
