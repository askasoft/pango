package sdk

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/askasoft/pango/log"
)

type retryTestError struct {
	NetError

	status     string // http status
	statusCode int    // http status code
}

func (rte *retryTestError) Error() string {
	s := rte.status

	if rte.retryAfter > 0 {
		s = fmt.Sprintf("%s (Retry After %s)", s, rte.retryAfter)
	}

	return s
}

func TestRetryForError(t *testing.T) {
	w := "429 Too Many Requests (Retry After 1s)"

	rte := &retryTestError{
		status:     "429 Too Many Requests",
		statusCode: http.StatusTooManyRequests,
	}
	rte.retryAfter = time.Second

	called, aborted := 0, 0
	err := RetryForError(func() error {
		called++
		return rte
	}, 2, func() bool {
		aborted++
		return false
	}, log.NewLog())

	if err.Error() != w {
		t.Errorf("Error(): %s, want %s", err.Error(), w)
	}
	if called != 3 {
		t.Errorf("called = %d, want %d", called, 3)
	}
	if aborted != 8 {
		t.Errorf("aborted = %d, want %d", aborted, 8)
	}
}
