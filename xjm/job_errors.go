package xjm

import (
	"errors"
)

var (
	ErrJobAborted  = errors.New("job aborted")  // indicates this job status is Aborted, returns by PingJob
	ErrJobComplete = errors.New("job complete") // indicates this job is finished, should update job status to Completed
	ErrJobCheckout = errors.New("job checkout failed: job running or missing")
	ErrJobExisting = errors.New("job existing")
	ErrJobMissing  = errors.New("job missing")
)
