package xwj

import (
	"errors"
)

var (
	ErrJobAborted   = errors.New("job aborted")
	ErrJobCompleted = errors.New("job completed")
	ErrJobCheckout  = errors.New("job checkout failed: job running or missing")
	ErrJobMissing   = errors.New("job missing")
)
