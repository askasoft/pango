package xmlrpc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// FaultError is feturned from the server when an invalid call is made
type FaultError struct {
	Method      string        `xmlrpc:"-"`
	StatusCode  int           `xmlrpc:"-"` // http status code
	RetryAfter  time.Duration `xmlrpc:"-"`
	FaultCode   int           `xmlrpc:"faultCode"`
	FaultString string        `xmlrpc:"faultString"`
}

func (fe *FaultError) GetRetryAfter() time.Duration {
	return fe.RetryAfter
}

// Error implements the error interface
func (fe FaultError) Error() string {
	return fmt.Sprintf("xmprpc: %s() %d fault(%d): %s", fe.Method, fe.StatusCode, fe.FaultCode, fe.FaultString)
}

func AsFaultError(err error) (fe *FaultError, ok bool) {
	ok = errors.As(err, &fe)
	return
}

func IsFaultError(err error) bool {
	_, ok := AsFaultError(err)
	return ok
}

func shouldRetry(err error) bool {
	if fe, ok := AsFaultError(err); ok {
		return fe.StatusCode == http.StatusTooManyRequests || (fe.StatusCode >= 500 && fe.StatusCode <= 599)
	}
	return !errors.Is(err, context.Canceled)
}
