package funcs

import (
	"encoding/json"
	"html/template"
)

// JSON returns a json marshal string.
func JSON(a any) (template.JS, error) {
	bs, err := json.Marshal(a)
	return template.JS(bs), err //nolint: gosec
}
