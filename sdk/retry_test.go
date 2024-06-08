package sdk

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/askasoft/pango/log"
)

var _ Retryable = NewNetError(nil)

type retryTestError struct {
	NetError

	status     string // http status
	statusCode int    // http status code
}

func (rte *retryTestError) Error() string {
	s := rte.status

	if rte.RetryAfter > 0 {
		s = fmt.Sprintf("%s (Retry After %s)", s, rte.RetryAfter)
	}

	return s
}

func TestRetryForErrorLoop(t *testing.T) {
	w := "429 Too Many Requests (Retry After 1s)"

	rte := &retryTestError{
		status:     "429 Too Many Requests",
		statusCode: http.StatusTooManyRequests,
	}
	rte.RetryAfter = time.Second

	called, aborted := 0, 0
	err := RetryForError(func() error {
		called++
		return rte
	}, 2, func() error {
		aborted++
		return nil
	}, time.Millisecond*250, log.NewLog())

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

func TestRetryForErrorAbort(t *testing.T) {
	w := errors.New("abort")

	rte := &retryTestError{
		status:     "429 Too Many Requests",
		statusCode: http.StatusTooManyRequests,
	}
	rte.RetryAfter = time.Second

	called, aborted := 0, 0
	err := RetryForError(func() error {
		called++
		return rte
	}, 2, func() error {
		aborted++
		if aborted == 2 {
			return w
		}
		return nil
	}, time.Millisecond*250, log.NewLog())

	if err != w {
		t.Errorf("Error(): %v, want %v", err, w)
	}
	if called != 1 {
		t.Errorf("called = %d, want %d", called, 1)
	}
	if aborted != 2 {
		t.Errorf("aborted = %d, want %d", aborted, 2)
	}
}
