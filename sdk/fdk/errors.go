package fdk

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/askasoft/pango/str"
)

type FieldError struct {
	Code    string `json:"code,omitempty"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

func (fe *FieldError) Error() string {
	return fmt.Sprintf("(%s: %s: %s)", fe.Code, fe.Field, fe.Message)
}

type ResultError struct {
	Method      string        `json:"-"` // http request method
	URL         *url.URL      `json:"-"` // http request URL
	StatusCode  int           `json:"-"` // http status code
	Status      string        `json:"-"` // http status
	Code        string        `json:"code,omitempty"`
	Message     string        `json:"message,omitempty"`
	Description string        `json:"description,omitempty"`
	Errors      []*FieldError `json:"errors,omitempty"`
	RetryAfter  time.Duration
}

func newResultError(res *http.Response) *ResultError {
	return &ResultError{
		Method:     res.Request.Method,
		URL:        res.Request.URL,
		StatusCode: res.StatusCode,
		Status:     res.Status,
	}
}

func (re *ResultError) GetRetryAfter() time.Duration {
	return re.RetryAfter
}

func (re *ResultError) Detail() string {
	var sb strings.Builder

	if re.Code != "" {
		sb.WriteString(re.Code)
	}
	if re.Message != "" {
		if sb.Len() > 0 {
			sb.WriteByte('/')
		}
		sb.WriteString(re.Message)
	}
	if re.Description != "" {
		if sb.Len() > 0 {
			sb.WriteByte('/')
		}
		sb.WriteString(re.Description)
	}
	for i, e := range re.Errors {
		sb.WriteString(str.If(i == 0, ": ", ", "))
		sb.WriteString(e.Error())
	}

	return sb.String()
}

func (re *ResultError) Error() string {
	es := re.Status

	if re.RetryAfter > 0 {
		es += " (Retry After " + re.RetryAfter.String() + ")"
	}

	es += " (" + re.Method + " " + re.URL.String() + ")"

	detail := re.Detail()
	if detail != "" {
		es += " - " + detail
	}

	return es
}

func shouldRetry(err error) bool {
	var re *ResultError
	if errors.As(err, &re) {
		return re.StatusCode == http.StatusTooManyRequests || (re.StatusCode >= 500 && re.StatusCode <= 599)
	}
	return !errors.Is(err, context.Canceled)
}
