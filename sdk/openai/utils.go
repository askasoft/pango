package openai

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/askasoft/pango/doc/jsonx"
)

const (
	contentTypeJSON = `application/json; charset="utf-8"`
)

func toString(o any) string {
	return jsonx.Prettify(o)
}

func buildJsonRequest(a any) (io.Reader, string, error) {
	body, err := json.Marshal(a)
	if err != nil {
		return nil, "", err
	}

	buf := bytes.NewReader(body)
	return buf, contentTypeJSON, nil
}
