package jsonx

import (
	"encoding/json"

	"github.com/askasoft/pango/ref"
	"github.com/askasoft/pango/str"
)

// Prettify use json.MarshalIndent(o, "", "  ") to convert `o` to json string
// return "" if `o` is nil.
// return error string if json marshal raise error.
func Prettify(o any) string {
	if ref.IsNil(o) {
		return ""
	}

	bs, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return err.Error()
	}
	return str.UnsafeString(bs)
}

// Stringify use json.Marshal(o) to convert `o` to json string
// return "" if `o` is nil.
// return error string if json marshal raise error.
func Stringify(o any) string {
	if ref.IsNil(o) {
		return ""
	}

	bs, err := json.Marshal(o)
	if err != nil {
		return err.Error()
	}
	return str.UnsafeString(bs)
}
