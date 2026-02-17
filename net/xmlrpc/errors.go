package xmlrpc

import (
	"errors"
	"fmt"
)

type HTTPError struct {
	Method     string
	StatusCode int
	Status     string
}

func (he *HTTPError) Error() string {
	return fmt.Sprintf("xmlrpc: %s(): error %d - %s", he.Method, he.StatusCode, he.Status)
}

func AsHTTPError(err error) (he *HTTPError, ok bool) {
	ok = errors.As(err, &he)
	return
}

func IsHTTPError(err error) bool {
	_, ok := AsHTTPError(err)
	return ok
}

// FaultError is returned from the server when an invalid call is made
type FaultError struct {
	Method      string `xmlrpc:"-"`
	FaultCode   int    `xmlrpc:"faultCode"`
	FaultString string `xmlrpc:"faultString"`
}

func (fe *FaultError) Error() string {
	return fmt.Sprintf("xmlrpc: %s(): fault %d - %s", fe.Method, fe.FaultCode, fe.FaultString)
}

func AsFaultError(err error) (fe *FaultError, ok bool) {
	ok = errors.As(err, &fe)
	return
}

func IsFaultError(err error) bool {
	_, ok := AsFaultError(err)
	return ok
}
