package xwj

import (
	"encoding/json"

	"github.com/askasoft/pango/str"
)

func toString(o any) string {
	bs, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(bs)
}

func toMap(o string) (m map[string]any) {
	if o != "" {
		_ = json.Unmarshal(str.UnsafeBytes(o), &m)
	}
	return
}
