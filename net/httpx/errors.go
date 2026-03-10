package httpx

import (
	"errors"
	"fmt"
)

type HTTPError struct {
	URL        string
	Method     string
	StatusCode int
	Status     string
}

func (he *HTTPError) Error() string {
	return fmt.Sprintf("[%d] %s %s (%s)", he.StatusCode, he.Method, he.URL, he.Status)
}

func AsHTTPError(err error) (he *HTTPError, ok bool) {
	ok = errors.As(err, &he)
	return
}

func IsHTTPError(err error) bool {
	_, ok := AsHTTPError(err)
	return ok
}
