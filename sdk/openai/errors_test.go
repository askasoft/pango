package openai

import (
	"encoding/json"
	"testing"

	"github.com/askasoft/pango/sdk"
)

var _ sdk.Retryable = &ResultError{}

func TestResultErrorUnmarshall(t *testing.T) {
	body := `{"error":{"code":502,"message":"Bad gateway.","param":null,"type":"cf_bad_gateway"}}`
	want := `502 Bad Gateway - cf_bad_gateway/502: Bad gateway.`

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
