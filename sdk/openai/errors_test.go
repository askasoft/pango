package openai

import (
	"encoding/json"
	"testing"
)

func TestErrorResultUnmarshall(t *testing.T) {
	body := `{"error":{"code":502,"message":"Bad gateway.","param":null,"type":"cf_bad_gateway"}}`
	want := `502 Bad Gateway - cf_bad_gateway/502: Bad gateway.`

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
