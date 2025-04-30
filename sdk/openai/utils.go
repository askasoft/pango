package openai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

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

func shouldRetry(re *ResultError) bool {
	return re.StatusCode == http.StatusTooManyRequests || (re.StatusCode >= 500 && re.StatusCode <= 599)
}
