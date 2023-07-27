package openai

import (
	"encoding/json"
	"net/http"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/sdk"
)

type RateLimitedError = sdk.RateLimitedError

func decodeResponse(res *http.Response, obj any) error {
	defer iox.DrainAndClose(res.Body)

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode == http.StatusOK {
		if obj != nil {
			return decoder.Decode(obj)
		}
		return nil
	}

	er := &ErrorResult{StatusCode: res.StatusCode, Status: res.Status}
	if res.StatusCode != http.StatusNotFound {
		_ = decoder.Decode(er)
	}
	return er
}
