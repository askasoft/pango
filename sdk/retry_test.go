package sdk

import (
	"context"
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
	rte := &retryTestError{
		status:     "429 Too Many Requests",
		statusCode: http.StatusTooManyRequests,
	}
	rte.RetryAfter = time.Millisecond * 100

	w := rte.Error()

	called := 0
	api := func() error {
		called++
		return rte
	}

	ctx := context.Background()

	err := RetryForError(ctx, api, 2, log.NewLog())

	if err.Error() != w {
		t.Errorf("Error(): %s, want %s", err.Error(), w)
	}
	if called != 3 {
		t.Errorf("called = %d, want %d", called, 3)
	}
}

func TestRetryForErrorAbort(t *testing.T) {
	w := context.Canceled

	rte := &retryTestError{
		status:     "429 Too Many Requests",
		statusCode: http.StatusTooManyRequests,
	}
	rte.RetryAfter = time.Millisecond * 100

	called := 0
	api := func() error {
		called++
		return rte
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Millisecond * 50)
		cancel()
	}()

	err := RetryForError(ctx, api, 2, log.NewLog())

	if err != w {
		t.Errorf("Error(): %v, want %v", err, w)
	}
	if called != 1 {
		t.Errorf("called = %d, want %d", called, 1)
	}
}
