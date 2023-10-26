package openai

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ErrorDetail struct {
	Type    string `json:"type,omitempty"`
	Code    any    `json:"code,omitempty"`
	Param   any    `json:"param,omitempty"`
	Message string `json:"message,omitempty"`
}

func (ed *ErrorDetail) String() string {
	var sb strings.Builder
	if ed.Type != "" {
		sb.WriteString(ed.Type)
	}
	if ed.Code != nil {
		s := fmt.Sprint(ed.Code)
		if s != "" {
			if sb.Len() > 0 {
				sb.WriteByte('/')
			}
			sb.WriteString(s)
		}
	}
	if ed.Param != nil {
		s := fmt.Sprint(ed.Param)
		if s != "" {
			if sb.Len() > 0 {
				sb.WriteByte('/')
			}
			sb.WriteString(s)
		}
	}
	if ed.Message != "" {
		if sb.Len() > 0 {
			sb.WriteString(": ")
		}
		sb.WriteString(ed.Message)
	}
	return sb.String()
}

type ErrorResult struct {
	StatusCode int          `json:"-"` // http status code
	Status     string       `json:"-"` // http status
	Detail     *ErrorDetail `json:"error,omitempty"`
}

func (er *ErrorResult) RetryAfter() time.Duration {
	switch er.StatusCode {
	case http.StatusTooManyRequests, http.StatusBadGateway, http.StatusServiceUnavailable:
		return time.Second * 20
	default:
		return 0
	}
}

func (er *ErrorResult) Error() string {
	es := er.Status

	if er.Detail != nil {
		es = es + " - " + er.Detail.String()
	}

	return es
}
