package sdk

import (
	"net/http"
	"testing"
	"time"

	"github.com/askasoft/pango/log"
)

func TestRetryForError(t *testing.T) {
	w := "429 Too Many Requests (Retry After 2s): You exceeded your current quota, please check your plan and billing details."

	rte := &RateLimitedError{
		Status:     "429 Too Many Requests",
		StatusCode: http.StatusTooManyRequests,
		RetryAfter: 2 * time.Second,
		Message:    "You exceeded your current quota, please check your plan and billing details.",
	}

	called, aborted := 0, 0
	err := RetryForError(func() error {
		called++
		return rte
	}, 2, time.Millisecond*1500, func() bool {
		aborted++
		return false
	}, log.NewLog())

	if err.Error() != w {
		t.Errorf("Error(): %s, want %s", err.Error(), w)
	}
	if called != 3 {
		t.Errorf("called = %d, want %d", called, 3)
	}
	if aborted != 12 {
		t.Errorf("aborted = %d, want %d", aborted, 12)
	}
}
