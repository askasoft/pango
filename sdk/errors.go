package sdk

import (
	"time"
)

type NetError struct {
	Err        error
	RetryAfter time.Duration
}

func NewNetError(err error, retryAfter ...time.Duration) *NetError {
	ne := &NetError{
		Err: err,
	}
	if len(retryAfter) > 0 {
		ne.RetryAfter = retryAfter[0]
	}
	return ne
}

func (ne *NetError) GetRetryAfter() time.Duration {
	return ne.RetryAfter
}

func (ne *NetError) Unwrap() error {
	return ne.Err
}

func (ne *NetError) Error() string {
	if ne == nil || ne.Err == nil {
		return "<nil>"
	}
	return ne.Err.Error()
}
