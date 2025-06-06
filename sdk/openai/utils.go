package openai

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/askasoft/pango/str"
)

const (
	contentTypeJSON = `application/json; charset="utf-8"`
)

func toJSONIndent(o any) string {
	if o == nil {
		return ""
	}

	bs, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return err.Error()
	}
	return str.UnsafeString(bs)
}

func buildJsonRequest(a any) (io.Reader, string, error) {
	body, err := json.Marshal(a)
	if err != nil {
		return nil, "", err
	}

	buf := bytes.NewReader(body)
	return buf, contentTypeJSON, nil
}
