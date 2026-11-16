package httpx

import (
	"net/http"
	"testing"
)

func TestRedactHeader(t *testing.T) {
	cs := []struct {
		key string
		val string
	}{
		{
			key: "Authorization",
			val: "ZmRrLUFQSUtFWTpY",
		},
		{
			key: "Cookie",
			val: "SECRETSESSION123",
		},
		{
			key: "x-api-key",
			val: "sk-ant-SECRETKEY123",
		},
		{
			key: "x-assertion",
			val: "SECRETKEY123",
		},
		{
			key: "x-password",
			val: "SECRETKEY123",
		},
		{
			key: "x-secret",
			val: "SECRETKEY123",
		},
		{
			key: "x-subscription",
			val: "SECRETKEY123",
		},
		{
			key: "x-token",
			val: "SECRETKEY123",
		},
	}

	for _, c := range cs {
		hh := http.Header{}
		hh.Add(c.key, c.val)

		redactHeader(hh)

		got := hh.Get(c.key)
		if got != RedactedValue {
			t.Errorf("RedactHeader(%q) = %q", c.key, got)
		}
	}
}
