package openai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type DetailError struct {
	Type    string `json:"type,omitempty"`
	Code    any    `json:"code,omitempty"`
	Param   any    `json:"param,omitempty"`
	Message string `json:"message,omitempty"`
}

func (de *DetailError) Error() string {
	var sb strings.Builder
	if de.Type != "" {
		sb.WriteString(de.Type)
	}
	if de.Code != nil {
		s := fmt.Sprint(de.Code)
		if s != "" {
			if sb.Len() > 0 {
				sb.WriteByte('/')
			}
			sb.WriteString(s)
		}
	}
	if de.Param != nil {
		s := fmt.Sprint(de.Param)
		if s != "" {
			if sb.Len() > 0 {
				sb.WriteByte('/')
			}
			sb.WriteString(s)
		}
	}
	if de.Message != "" {
		if sb.Len() > 0 {
			sb.WriteString(": ")
		}
		sb.WriteString(de.Message)
	}
	return sb.String()
}

type ResultError struct {
	StatusCode int          `json:"-"` // http status code
	Status     string       `json:"-"` // http status
	Detail     *DetailError `json:"error,omitempty"`
	RetryAfter time.Duration
}

func (re *ResultError) GetRetryAfter() time.Duration {
	return re.RetryAfter
}

func (re *ResultError) Error() string {
	es := re.Status

	if re.RetryAfter > 0 {
		es = fmt.Sprintf("%s (Retry After %s)", es, re.RetryAfter)
	}

	if re.Detail != nil {
		es = es + " - " + re.Detail.Error()
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
