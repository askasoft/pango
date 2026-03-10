package ret

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/askasoft/pango/log"
)

type retryTestError struct {
	status     string // http status
	statusCode int    // http status code
}

func (rte *retryTestError) Error() string {
	return rte.status
}

func TestRetryForErrorLoop(t *testing.T) {
	rte := &retryTestError{
		status:     "429 Too Many Requests",
		statusCode: http.StatusTooManyRequests,
	}
	sr := func(_ error) time.Duration {
		return time.Millisecond * 100
	}

	w := rte.Error()

	called := 0
	do := func() error {
		called++
		return rte
	}

	ctx := context.Background()

	err := RetryForError(ctx, do, sr, 2, log.NewLog())

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
	sr := func(_ error) time.Duration {
		return time.Millisecond * 100
	}

	called := 0
	do := func() error {
		called++
		return rte
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Millisecond * 50)
		cancel()
	}()

	err := RetryForError(ctx, do, sr, 2, log.NewLog())

	if err != w {
		t.Errorf("Error(): %v, want %v", err, w)
	}
	if called != 1 {
		t.Errorf("called = %d, want %d", called, 1)
	}
}
