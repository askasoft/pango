package sdk

import (
	"fmt"
	"time"

	"github.com/askasoft/pango/log"
)

type RateLimitedError struct {
	Status     string        // http status
	StatusCode int           // http status code
	RetryAfter time.Duration // retry after time
	Message    string        // detail error message
}

func (rle *RateLimitedError) Error() string {
	s := rle.Status

	if rle.RetryAfter > 0 {
		s = fmt.Sprintf("%s (Retry After %s)", s, rle.RetryAfter)
	}

	if rle.Message != "" {
		s = fmt.Sprintf("%s: %s", s, rle.Message)
	}

	return s
}

func RetryForRateLimited(api func() error, maxRetry int, abort func() bool, logger log.Logger) (err error) {
	for i := 0; ; i++ {
		err = api()
		if i >= maxRetry {
			break
		}
		if !SleepForRateLimited(err, abort, logger) {
			break
		}
	}
	return err
}

// SleepForRateLimited if err is RateLimitedError, sleep Retry-After and return true
// return false if err is not ReteLimitedError or abort() returns true
func SleepForRateLimited(err error, abort func() bool, logger log.Logger) bool {
	if err != nil {
		if rle, ok := err.(*RateLimitedError); ok { //nolint: errorlint
			ra := rle.RetryAfter
			if ra <= 0 {
				ra = time.Second * 30 // default to 30s
			}

			if logger != nil {
				logger.Warnf("Sleep %s for API Rate Limited", ra)
			}

			for te := time.Now().Add(ra); te.After(time.Now()); {
				if abort != nil && abort() {
					return false
				}
				time.Sleep(time.Second)
			}

			return true
		}
	}
	return false
}
