package openai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/askasoft/pango/iox"
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

func decodeResponse(res *http.Response, obj any, retryAfter time.Duration) error {
	defer iox.DrainAndClose(res.Body)

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode == http.StatusOK {
		if obj != nil {
			return decoder.Decode(obj)
		}
		return nil
	}

	er := &ErrorResult{
		Status:     res.Status,
		StatusCode: res.StatusCode,
	}
	_ = decoder.Decode(er)

	if er.StatusCode == http.StatusTooManyRequests || (er.StatusCode >= 500 && er.StatusCode <= 599) {
		er.RetryAfter = retryAfter
	}

	return er
}
