package sdk

import (
	"net/http"
	"testing"
	"time"
)

func TestRateLimitedError(t *testing.T) {
	w := "429 Too Many Requests (Retry After 20s): You exceeded your current quota, please check your plan and billing details."

	err := &RateLimitedError{
		Status:     "429 Too Many Requests",
		StatusCode: http.StatusTooManyRequests,
		RetryAfter: 20 * time.Second,
		Message:    "You exceeded your current quota, please check your plan and billing details.",
	}

	if err.Error() != w {
		t.Errorf("Error(): %s, want %s", err.Error(), w)
	}
}
