package jsonx

import (
	"encoding/json"

	"github.com/askasoft/pango/str"
)

func Prettify(o any) string {
	if o == nil {
		return ""
	}

	bs, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return err.Error()
	}
	return str.UnsafeString(bs)
}

func Stringify(o any) string {
	if o == nil {
		return ""
	}

	bs, err := json.Marshal(o)
	if err != nil {
		return err.Error()
	}
	return str.UnsafeString(bs)
}
