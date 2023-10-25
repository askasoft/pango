package sdk

import (
	"errors"
	"net"
	"time"

	"github.com/askasoft/pango/log"
)

type RetryableError interface {
	RetryAfter() time.Duration
}

func RetryForError(api func() error, maxRetryCount int, maxRetryAfter time.Duration, abort func() bool, logger log.Logger) (err error) {
	for i := 0; ; i++ {
		err = api()
		if err == nil || i >= maxRetryCount {
			break
		}

		ra := getRetryAfter(err, maxRetryAfter)
		if ra <= 0 {
			break
		}

		if logger != nil {
			logger.Warnf("Sleep %s for %s", ra, err.Error())
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

func getRetryAfter(err error, maxRetryAfter time.Duration) time.Duration {
	if re, ok := err.(RetryableError); ok { //nolint: errorlint
		ra := re.RetryAfter()
		if ra <= 0 {
			return 0
		}
		if ra > maxRetryAfter {
			ra = maxRetryAfter
		}
		return ra
	}

	var noe *net.OpError
	if errors.As(err, &noe) {
		return maxRetryAfter
	}

	return 0
}
