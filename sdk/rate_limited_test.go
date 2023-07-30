package sdk

import (
	"net/http"
	"testing"
	"time"
)

func TestRateLimitedError(t *testing.T) {
	w := "429 Retry After 20s"
	err := &RateLimitedError{StatusCode: http.StatusTooManyRequests, RetryAfter: 20 * time.Second}
	if err.Error() != w {
		t.Errorf("Error(): %s, want %s", err.Error(), w)
	}
}
