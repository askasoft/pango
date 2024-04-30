package xjm

import (
	"errors"
)

var (
	ErrJobAborted   = errors.New("job aborted")
	ErrJobCompleted = errors.New("job completed")
	ErrJobCheckout  = errors.New("job checkout failed: job running or missing")
	ErrJobExisting  = errors.New("job existing")
	ErrJobMissing   = errors.New("job missing")
	ErrJobPing      = errors.New("job ping failed")
)
