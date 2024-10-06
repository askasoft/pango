package fdk

import (
	"encoding/json"
	"testing"

	"github.com/askasoft/pango/sdk"
)

var _ sdk.Retryable = &ResultError{}

func TestResultErrorUnmarshall(t *testing.T) {
	body := `{"code":"xxx","message":"Bad message."}`
	want := `502 Bad Gateway - xxx/Bad message.`

	re := &ResultError{
		Status:     "502 Bad Gateway",
		StatusCode: 502,
	}
	err := json.Unmarshal([]byte(body), re)
	if err != nil {
		t.Fatal(err)
	}

	if re.Error() != want {
		t.Fatalf("want %s, but %s", want, re.Error())
	}
}
