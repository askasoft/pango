package sdk

import (
	"fmt"
	"time"

	"github.com/askasoft/pango/log"
)

type RateLimitedError struct {
	StatusCode int // http status code
	RetryAfter int // retry after seconds
}

func (e *RateLimitedError) Error() string {
	return fmt.Sprintf("%d Retry After %d seconds", e.StatusCode, e.RetryAfter)
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
				logger.Warnf("Sleep %d seconds for API Rate Limited", rle.RetryAfter)
			}
			time.Sleep(time.Duration(rle.RetryAfter) * time.Second)
			return true
		}
	}
	return false
}
