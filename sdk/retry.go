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

func NeverAbort() error {
	return nil
}

// RetryForError loop max `retries` count to call api().
// returns the error if api() returns a non Retryable error.
// returns the error if abort() returns error.
// call SleepForRetry(sleep, after, abort), if api() returns a Retryable error.
func RetryForError(api func() error, retries int, abort func() error, sleep time.Duration, logger log.Logger) (err error) {
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

		if err = SleepForRetry(sleep, after, abort); err != nil {
			break
		}
	}
	return
}

// SleepForRetry loop to sleep(`sleep`) until 'after' duration elapsed.
// call 'abort()' every 'sleep' interval, if abort() returns error, returns it.
// if 'sleep' <= 0, 'sleep' = time.Second.
// if 'sleep' > 'after', 'sleep' = 'after'.
func SleepForRetry(sleep, after time.Duration, abort func() error) (err error) {
	if sleep <= 0 {
		sleep = time.Second
	}
	if sleep > after {
		sleep = after
	}

	for te := time.Now().Add(after); te.After(time.Now()); {
		if abort != nil {
			if err = abort(); err != nil {
				return
			}
		}
		time.Sleep(sleep)
	}
	return
}
