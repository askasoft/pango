package ret

import (
	"context"
	"time"

	"github.com/askasoft/pango/log"
)

type RetryError struct {
	Err        error
	RetryAfter time.Duration
}

func NewRetryError(err error, retryAfter ...time.Duration) *RetryError {
	re := &RetryError{
		Err: err,
	}
	if len(retryAfter) > 0 {
		re.RetryAfter = retryAfter[0]
	}
	return re
}

func (re *RetryError) GetRetryAfter() time.Duration {
	return re.RetryAfter
}

func (re *RetryError) Unwrap() error {
	return re.Err
}

func (re *RetryError) Error() string {
	if re == nil || re.Err == nil {
		return "<nil>"
	}
	return re.Err.Error()
}

type Retryable interface {
	GetRetryAfter() time.Duration
}

func GetRetryAfter(a any) time.Duration {
	if r, ok := a.(Retryable); ok {
		return r.GetRetryAfter()
	}
	return 0
}

// RetryForError loop max `retries` count to call api().
// returns the error if api() returns a non Retryable error.
func RetryForError(ctx context.Context, api func() error, retries int, logger log.Logger) error {
	err := api()
	if err == nil {
		return nil
	}

	count := 1
	after := GetRetryAfter(err)
	if after <= 0 || count > retries {
		return err
	}

	if logger != nil {
		logger.Warnf("Sleep %s for retry #%d: %s", after, count, err.Error())
	}

	timer := time.NewTimer(after)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			if err = api(); err == nil {
				return nil
			}

			if count++; count > retries {
				return err
			}

			if after = GetRetryAfter(err); after <= 0 {
				return err
			}

			if logger != nil {
				logger.Warnf("Sleep %s for retry #%d: %s", after, count, err.Error())
			}
			timer.Reset(after)
		}
	}
}
