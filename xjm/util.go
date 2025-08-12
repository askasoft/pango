package xjm

import (
	"encoding/base64"
	"encoding/json"

	"github.com/askasoft/pango/doc/jsonx"
	"github.com/askasoft/pango/str"
)

func toString(o any) string {
	return jsonx.Prettify(o)
}

func Encode(v any) (string, error) {
	if v == nil {
		return "", nil
	}

	if s, ok := v.(string); ok {
		return s, nil
	}

	if bs, ok := v.([]byte); ok {
		return base64.StdEncoding.EncodeToString(bs), nil
	}

	bs, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return str.UnsafeString(bs), nil
}

func MustEncode(v any) string {
	s, err := Encode(v)
	if err != nil {
		panic(err)
	}
	return s
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

func MustDecode(p string, v any) {
	err := Decode(p, v)
	if err != nil {
		panic(err)
	}
}
