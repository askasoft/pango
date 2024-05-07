package sdk

import (
	"time"

	"github.com/askasoft/pango/log"
)

type Retryable interface {
	GetRetryAfter() time.Duration
}

func getRetryAfter(err error) time.Duration {
	if re, ok := err.(Retryable); ok { //nolint: all
		return re.GetRetryAfter()
	}

	return 0
}

func NeverAbort() bool {
	return false
}

// RetryForError call api(), if api() returns a Retryable error,
// sleep until Retryable.GetRetryAfter() duration, retry call api().
// sleep 'sleep' duration, call 'abort()', if abort() returns true, returns the last error.
func RetryForError(api func() error, retries int, abort func() bool, sleep time.Duration, logger log.Logger) (err error) {
	for i := 1; ; i++ {
		err = api()
		if err == nil || i > retries {
			break
		}

		after := getRetryAfter(err)
		if after <= 0 {
			break
		}

		if logger != nil {
			logger.Warnf("Sleep %s for retry #%d: %s", after, i, err.Error())
		}

		if !SleepForRetry(sleep, after, abort) {
			break
		}
	}
	return
}

// SleepForRetry sleep until 'after' duration elapsed.
// call 'abort()' every 'sleep' interval, if abort() returns true, returns false.
// if 'sleep' <= 0, 'sleep' = time.Second.
// if 'sleep' > 'after', 'sleep' = 'after'.
func SleepForRetry(sleep, after time.Duration, abort func() bool) bool {
	if sleep <= 0 {
		sleep = time.Second
	}
	if sleep > after {
		sleep = after
	}

	for te := time.Now().Add(after); te.After(time.Now()); {
		if abort != nil && abort() {
			return false
		}
		time.Sleep(sleep)
	}

	return true
}
