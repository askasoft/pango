package ret

import (
	"context"
	"time"

	"github.com/askasoft/pango/log"
)

// RetryForError call do() and retry max `retries` count if do() returns a error, and shouldRetry(err) returns a duration.
// returns the error if do() returns a non retryable error.
func RetryForError(ctx context.Context, do func() error, shouldRetry func(error) time.Duration, retries int, logger log.Logger) error {
	err := do()
	if err == nil {
		return nil
	}

	count := 1
	if count > retries {
		return err
	}

	after := shouldRetry(err)
	if after <= 0 {
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
			if err = do(); err == nil {
				return nil
			}

			if count++; count > retries {
				return err
			}

			if after = shouldRetry(err); after <= 0 {
				return err
			}

			if logger != nil {
				logger.Warnf("Sleep %s for retry #%d: %s", after, count, err.Error())
			}
			timer.Reset(after)
		}
	}
}

type Retryer struct {
	Logger      log.Logger
	MaxRetries  int
	ShouldRetry func(error) time.Duration
}

func (r *Retryer) Do(ctx context.Context, do func() error) error {
	return RetryForError(ctx, do, r.ShouldRetry, r.MaxRetries, r.Logger)
}
