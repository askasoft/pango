package sdk

import (
	"time"

	"github.com/askasoft/pango/log"
)

type RetryableError interface {
	RetryAfter() time.Duration
}

func RetryForError(api func() error, maxRetryCount int, abort func() bool, logger log.Logger) (err error) {
	for i := 1; ; i++ {
		err = api()
		if err == nil || i > maxRetryCount {
			break
		}

		ra := getRetryAfter(err)
		if ra <= 0 {
			break
		}

		if logger != nil {
			logger.Warnf("Sleep %s for retry [%d] %s", ra, i, err.Error())
		}

		if !sleepForRetry(ra, abort, logger) {
			break
		}
	}
	return
}

func sleepForRetry(ra time.Duration, abort func() bool, logger log.Logger) bool {
	for te := time.Now().Add(ra); te.After(time.Now()); {
		if abort != nil && abort() {
			return false
		}
		time.Sleep(time.Millisecond * 250)
	}

	return true
}

func getRetryAfter(err error) time.Duration {
	if re, ok := err.(RetryableError); ok { //nolint: all
		return re.RetryAfter()
	}

	return 0
}
