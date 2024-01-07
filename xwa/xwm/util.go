package xwm

import (
	"encoding/json"

	"github.com/askasoft/pango/str"
)

func ToString(o any) string {
	bs, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(bs)
}

func ToMap(o string) (m map[string]any) {
	if o != "" {
		_ = json.Unmarshal(str.UnsafeBytes(o), &m)
	}
	return
}
