package sdk

import (
	"fmt"
	"time"
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
