package sdk

import (
	"fmt"
	"time"

	"github.com/askasoft/pango/log"
)

type RateLimitedError struct {
	StatusCode int           // http status code
	RetryAfter time.Duration // retry after time
}

func (e *RateLimitedError) Error() string {
	return fmt.Sprintf("%d Retry After %v", e.StatusCode, e.RetryAfter)
}

func RetryForRateLimited(api func() error, maxRetry int, logger log.Logger) (err error) {
	for i := 0; ; i++ {
		err = api()
		if i >= maxRetry {
			break
		}
		if !SleepForRateLimited(err, logger) {
			break
		}
	}
	return err
}

// SleepForRateLimited if err is RateLimitedError, sleep Retry-After and return true
func SleepForRateLimited(err error, logger log.Logger) bool {
	if err != nil {
		if rle, ok := err.(*RateLimitedError); ok { //nolint: errorlint
			if logger != nil {
				logger.Warnf("Sleep %v for API Rate Limited", rle.RetryAfter)
			}
			time.Sleep(rle.RetryAfter)
			return true
		}
	}
	return false
}
