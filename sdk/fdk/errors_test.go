package fdk

import (
	"encoding/json"
	"testing"

	"github.com/askasoft/pango/sdk"
)

var _ sdk.Retryable = &ErrorResult{}

func TestErrorResultUnmarshall(t *testing.T) {
	body := `{"code":"xxx","message":"Bad message."}`
	want := `502 Bad Gateway - xxx/Bad message.`

	er := &ErrorResult{
		Status:     "502 Bad Gateway",
		StatusCode: 502,
	}
	err := json.Unmarshal([]byte(body), er)
	if err != nil {
		t.Fatal(err)
	}

	if er.Error() != want {
		t.Fatalf("want %s, but %s", want, er.Error())
	}
}
