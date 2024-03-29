package xjm

import (
	"encoding/base64"
	"encoding/json"

	"github.com/askasoft/pango/str"
)

func Encode(v any) string {
	if v == nil {
		return ""
	}

	if s, ok := v.(string); ok {
		return s
	}

	if bs, ok := v.([]byte); ok {
		return base64.StdEncoding.EncodeToString(bs)
	}

	bs, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return str.UnsafeString(bs)
}

func Decode(p string, v any) error {
	if p == "" {
		return nil
	}

	if ps, ok := v.(*string); ok {
		*ps = p
		return nil
	}

	if pb, ok := v.(*[]byte); ok {
		bs, err := base64.StdEncoding.DecodeString(p)
		if err != nil {
			return err
		}

		*pb = bs
		return nil
	}

	return json.Unmarshal(str.UnsafeBytes(p), v)
}

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
